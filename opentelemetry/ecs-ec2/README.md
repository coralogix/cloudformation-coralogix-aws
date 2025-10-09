# Coralogix OpenTelemetry Agent for ECS-EC2. Cloudformation template.

This CloudFormation template deploys an ECS Service and Task Definition for running the OpenTelemetry Collector agent on an ECS Cluster. This deployment is able to collect Logs, Metrics and Traces. The template will deploy a daemonset which runs an instance of OpenTelemetry Collector on each node in a cluster.

CloudFormation template to launch the Coralogix Distribution for OpenTelemetry ("CDOT") into an existing ECS Cluster. This CDOT deployment is able to collect Logs, Metrics and Traces. CDOT is deployed in the OTEL [Agent deployment](https://opentelemetry.io/docs/collector/deployment/agent/) pattern, as an ECS Daemon Service type, which runs an instance of the OpenTelemetry Collector agent on each node in a cluster.

## Image

This solution uses the coralogixrepo/coralogix-otel-collector image which is a custom distribution of OpenTelemetry containing custom components developed by Coralogix. The image is available on [Docker Hub](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector). The ECS components are described [here](./components.md)

The OTEL Collector/agent/daemon image used is the [Coralogix Distribution for OpenTelemetry](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector) Docker Hub image. It is deployed as a [*Daemon*](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html#service_scheduler_daemon) ECS Task, i.e. one OTEL Collector agent container on each EC2 instance (i.e. ECS container instance) across the cluster.

CDOT extends upon the main [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) project, adding features specifically to enhance integration with AWS ECS, among other improvements.

The OTEL agent is deployed as a Daemon ECS Task and connected using [`host` network mode](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/networking-networkmode-host.html). OTEL-instrumented application containers that need to send telemetry to the local OTEL agent can lookup the IP address of the CDOT container [using a number of methods](https://coralogix.com/docs/opentelemetry-using-ecs-ec2/#otel-agent-network-service-discovery), making it easier for Application Tasks using `awsvpc` and `bridge` network modes to connect with the OTEL agent. OTEL-instrumented application containers should also consider which resource attributes to use as telemetry identifiers.

The CDOT OTEL agent also features enhancements specific to ECS integration. These improvements are proprietary to the Coralogix Distribution for OpenTelemetry.

### Logs

The OTEL agent uses a [filelog receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filereceiver) to read the docker logs of all containers on the EC2 host. OTLP is also accepted. Coralogix provides the `awsecscontainermetricsd` receiver which enables metrics collection of all tasks on the same host. The [coralogix exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/coralogixexporter) forwards telemetry to your configured Coralogix endpoint.

Logs are collected from all containers that log to `/var/lib/docker/containers/*/*.log`. The container requires privileges to mount the host read-only file path `/var/lib/docker/`.

### Container Metrics

Container metrics are collected from all containers running on the ECS Cluster. The metrics are collected using the [awsecscontainermetricsd](./components.md#awsecscontainermetricsd) receiver. If you do not wish to collect container metrics, comment out or delete the `metrics/containermetrics` pipeline from the configuration.

### OpenTelemetry Collector Metrics

The default configuration exposes OpenTelemetry Collector metrics on port `8888` via the path `/metrics`. The metrics are collected using a [prometheus](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/prometheusreceiver) scrape job. These are performance metrics of the OpenTelemetry Collector containers. Records received and processed, submission faults and more.

### Traces

A GRPC(*4317*) and HTTP(*4318*) endpoint is exposed for sending traces to the local OTLP endpoint.

By default, traces are sampled at 10% rate using head sampling. Head sampling is a feature that allows you to sample traces at the collection point before any processing occurs. When enabled, it creates a separate pipeline for sampled traces using probabilistic sampling. This helps reduce the volume of traces while maintaining a representative sample.

The sampling configuration can be adjusted using the following parameters:
- `EnableHeadSampler`: Enable/disable head sampling
- `SamplerMode`: Choose between proportional, equalizing, or hash_seed sampling modes
- `SamplingPercentage`: Set the desired sampling rate (0-100%)

### Span Metrics

When enabled, the spanmetrics connector generates metrics from traces, providing insights into trace performance and patterns. This feature creates additional metrics pipelines that convert span data into metrics for monitoring and alerting purposes.

### Database Traces

When enabled, database operation traces are processed separately with dedicated metrics generation. This feature provides specialized monitoring for database operations with optimized bucket configurations and filtering.

### Requires:

- An existing ECS Cluster
- [aws-cli]() (*if deploying via CLI*)

## Template Versions

This repository contains two CloudFormation templates:

- **`template.yaml`** - **S3-only configuration**  Uses S3 to store OpenTelemetry configuration files.
- **`template-legacy.yaml`** - **Full-featured template** with multiple configuration options (template, S3, Parameter Store) for direct user deployment.

## Parameters:

| Parameter        | Description                                                                                                                                                                                                                          | Default Value                                                            | Required           |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|--------------------|
| ConfigSource     | Select the configuration source for OpenTelemetry Collector:<br>- `template`: Use built-in template configuration<br>- `s3`: Use configuration file from S3 (via S3 URI)<br>- `parameter-store`: Use configuration from AWS Systems Manager Parameter Store | template                                                                |                    |
| ClusterName      | The name of an **existing** ECS Cluster                                                                                                                                                                                              |                                                                          | :heavy_check_mark: |
| CDOTImageVersion | The Coralogix OpenTelemetry Collector Image version/tag to use. See available tags [here](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector/tags)                                                                     |                                                                          |                    |
| Image            | The OpenTelemetry Collector Image to use. If specified, this value will override the CDOTImageVersion parameter and the coralogix otel collector image.                                                                              | none                                                                     |                    |
| Memory           | The amount of memory to allocate to the OpenTelemetry container.<br>*Assigning too much memory can lead to the ECS Service not being deployed. Make sure that values are within the range of what is available on your ECS Cluster* | 2048                                                                     |                    |
| CoralogixRegion  | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:* <br>- EU1<br>- EU2<br>- AP1<br>- AP2<br>- AP3<br>- US1<br>- US2 | :heavy_check_mark: |
| DefaultApplicationName | Your application name                                                                                                                                                                                                                 | OTEL                                                                     |                    |
| DefaultSubsystemName | Your Subsystem name                                                                                                                                                                                                                 | ECS-EC2                                                                  |                    |
| CoralogixApiKey  | The Send-Your-Data API key for your Coralogix account. See: https://coralogix.com/docs/send-your-data-api-key/                                                                                                                       |                                                                          | :heavy_check_mark: |
| S3ConfigBucket   | S3 bucket name containing the configuration file. Required when ConfigSource is 's3'.                                                                                                                                                |                                                                          |                    |
| S3ConfigKey      | S3 object key (file path) for the configuration file. Required when ConfigSource is 's3'.                                                                                                                                           |                                                                          |                    |
| CustomConfig     | The name of an AWS Systems Manager Parameter Store parameter to use as a custom configuration. Required when ConfigSource is 'parameter-store'.                                                                                                | none                                                                     |                    |
| EnableHeadSampler | Enable or disable head sampling for traces. When enabled, sampling decisions are made at the collection point before any processing occurs.                                                                                        | true                                                                     |                    |
| EnableSpanMetrics | Enable or disable the spanmetrics connector and pipeline. When enabled (default), span metrics will be generated from traces.                                                                                                      | true                                                                     |                    |
| EnableTracesDB   | Enable or disable the traces/db pipeline for database operation metrics. When enabled, database operation metrics will be generated. Note: This feature requires spanmetrics to be enabled.                                           | false                                                                    |                    |
| SamplerMode      | The sampling mode to use:<br>**proportional**: Maintains the relative proportion of traces across services.<br>**equalizing**: Attempts to sample equal numbers of traces from each service.<br>**hash_seed**: Uses consistent hashing to ensure the same traces are sampled across restarts. | proportional                                                             |                    |
| SamplingPercentage | The percentage of traces to sample (0-100). A value of 100 means all traces will be sampled, while 0 means no traces will be sampled.                                                                                            | 10                                                                       |                    |
| HealthCheckEnabled | Enable ECS container health check for the OTEL agent container. Requires OTEL collector image version v0.4.2 or later.                                                                                                            | false                                                                    |                    |
| HealthCheckInterval | Health check interval (seconds)                                                                                                                                                                                                     | 30                                                                       |                    |
| HealthCheckTimeout | Health check timeout (seconds)                                                                                                                                                                                                      | 5                                                                        |                    |
| HealthCheckRetries | Health check retries                                                                                                                                                                                                               | 3                                                                        |                    |
| HealthCheckStartPeriod | Health check start period (seconds)                                                                                                                                                                                               | 10                                                                       |                    |

## Deploy the Cloudformation template:

### Default Template Configuration

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        DefaultApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        HealthCheckEnabled=true
```

### S3 Configuration

For large configurations that exceed SSM Parameter Store limits, you can store your configuration in S3:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        ConfigSource=s3 \
        S3ConfigBucket=<your-s3-bucket> \
        S3ConfigKey=<path/to/config.yaml> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        HealthCheckEnabled=true
```

### AWS Systems Manager Parameter Store Configuration

For custom configurations stored in AWS Systems Manager Parameter Store:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        ConfigSource=parameter-store \
        CustomConfig=<Parameter Store Name> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        HealthCheckEnabled=true
```

Note that these are just examples of how this could be deployed. You can also deploy this template using the AWS Console or any Cloudformation management tools.

### OpenTelemetry Configuration

The template supports three configuration sources:

1. **Template Configuration (Default)**: An OpenTelemetry configuration is embedded in the cloudformation template by default.

2. **S3 Configuration**: For large configurations that exceed AWS Systems Manager Parameter Store limits, you can store your configuration file in S3. The template will automatically create the necessary IAM roles for S3 access.

3. **AWS Systems Manager Parameter Store**: You can specify your own configuration stored in AWS Systems Manager Parameter Store using the `CustomConfig` parameter. When using Parameter Store configuration (ConfigSource=parameter-store), you must ensure your EC2 host has access to the AWS Systems Manager API and the OTEL containers have read permission for the specified parameter.

The default configuration will monitor container logs, listen for logs, metrics and traces on port `4317/4318`, and monitor container metrics using the `awsecscontainermetricsd` receiver. The `awsecscontainermetricsd` receiver is based on the [awsecscontainermetrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver) receiver, however, instead of monitoring a single task, metrics will be collected from all containers. If you do not wish to collect container metrics, comment out or delete the metrics/container-metrics pipeline from the configuration.

The embeded default OpenTelemetry configuration can be viewed [here](./template.yaml#L90).

### Using Complete Example Configuration with S3

For customers who need to customize their OpenTelemetry configuration beyond the template defaults, you can use the complete example configuration from the `examples` folder as a starting point for S3 deployment.

#### Example Configuration File

A comprehensive example configuration is provided in [`examples/comprehensive-config.yaml`](./examples/comprehensive-config.yaml) that includes:

- **Receivers**: OTLP (gRPC/HTTP), AWS ECS Container Metrics, Prometheus, File Log, Host Metrics
- **Processors**: ECS Attributes, Transform, Filter, Batch, Resource Detection, Probabilistic Sampler
- **Exporters**: Coralogix (logs, metrics, traces), Coralogix Resource Catalog
- **Connectors**: Forward, Span Metrics, Database Span Metrics
- **Pipelines**: Multiple pipeline configurations for different sampling and metrics scenarios
- **Telemetry**: Proper metrics configuration with Prometheus endpoint

#### Configuration Scenarios

The example configuration demonstrates different deployment scenarios:

1. **ALL ENABLED** (EnableSpanMetrics=true, EnableTracesDB=true, EnableSampling=true) - *Currently active*
2. **ALL DISABLED** (EnableSpanMetrics=false, EnableTracesDB=false, EnableSampling=false)
3. **SPAN METRICS ONLY** (EnableSpanMetrics=true, EnableTracesDB=false, EnableSampling=true)
4. **SPAN METRICS NO SAMPLING** (EnableSpanMetrics=true, EnableTracesDB=false, EnableSampling=false)

#### Steps to Use with S3:

1. **Download the Example**: Copy the configuration from [`examples/comprehensive-config.yaml`](./examples/comprehensive-config.yaml)
2. **Modify as Needed**: 
   - Choose the appropriate scenario by uncommenting/commenting the relevant pipeline sections
   - Add your custom receivers, processors, or exporters
3. **Upload to S3**: Save the modified configuration as a YAML file and upload to your S3 bucket
4. **Deploy with S3**: Use the S3 deployment method with your custom configuration

This approach allows you to leverage the full power of OpenTelemetry while maintaining the convenience of CloudFormation deployment.

### Health Check

Requires OTEL collector image version v0.4.2 or later.
The default config will expose a health check on port `13133` of the localhost via the path `/`. The health check is exposed using the [health_check](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/healthcheckextension) extension.

The healthy response should look like this:

```json
{
  "status": "Server available",
  "upSince": "2023-10-25T15:37:32.003837622Z",
  "uptime": "2m5.2610063s"
}
```

### ECS Container Health Check

You can customize the health check settings using the following parameters:
- `HealthCheckInterval` (default: 30)
- `HealthCheckTimeout` (default: 5)
- `HealthCheckRetries` (default: 3)
- `HealthCheckStartPeriod` (default: 10)

Example deployment with custom health check settings:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        DefaultApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        HealthCheckEnabled=true \
        HealthCheckInterval=60 \
        HealthCheckTimeout=10 \
        HealthCheckRetries=5 \
        HealthCheckStartPeriod=20
```

### Further info

See documentation: [AWS ECS-EC2 using OpenTelemetry](https://coralogix.com/docs/opentelemetry-using-ecs-ec2).
