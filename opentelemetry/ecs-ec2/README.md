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

### Requires:

- An existing ECS Cluster
- [aws-cli]() (*if deploying via CLI*)

## Parameters:

| Parameter        | Description                                                                                                                                                                                                                          | Default Value                                                            | Required           |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|--------------------|
| ClusterName      | The name of an **existing** ECS Cluster                                                                                                                                                                                              |                                                                          | :heavy_check_mark: |
| CDOTImageVersion | The Coralogix OpenTelemetry Collector Image version/tag to use. See available tags [here](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector/tags)                                                                     |                                                                          |                    |
| Memory           | The amount of memory to allocate to the OpenTelemetry container.<br>*Assigning too much memory can lead to the ECS Service not being deployed. Make sure that values are within the range of what is available on your ECS Cluster* | 256                                                                      |                    |
| CoralogixRegion  | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:* <br>- EU1<br>- EU2<br>- AP1<br>- AP2<br>- AP3<br>- US1<br>- US2 | :heavy_check_mark: |
| ApplicationName  | Your application name                                                                                                                                                                                                                 |                                                                          | :heavy_check_mark: |
| SubsystemName    | Your Subsystem name |                                                            | :heavy_check_mark: |
| CoralogixApiKey  | The Send-Your-Data API key for your Coralogix account. See: https://coralogix.com/docs/send-your-data-api-key/                                                                                                                       |                                                                          | :heavy_check_mark: |
| CustomConfig       | The name of a Parameter Store to use as a custom configuration. Must be in the same region as your ECS cluster.            | none                                                                         |                    |
| TaskExecutionRoleARN       |       When using a Custom Configuration in Parameter Store, set to the ARN of a Task Execution Role that has access to the PS.            | Default                                                                        |                    |


## Deploy the Cloudformation template:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        DefaultApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region>
```

**When using a custom configuration for OpenTelemetry**

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        DefaultApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        CustomConfig=<Parameter Store Name> \
        TaskExecutionRoleARN=<task-execution-role-arn>
```

Note that these are just examples of how this could be deployed. You can also deploy this template using the AWS Console or any Cloudformation management tools.

### OpenTelemetry Configuration

An OpenTelemetry configuration is embedded in the cloudformation template by default, however, you do have the option of specifying your own configuration, stored in a Parameter Store, using the `CustomConfig` parameter. When using a custom configuration, you must ensure your EC2 host has access to the AWS Systems Manager API and the OTEL containers have read permission for the specified Parameter Store. You will need to create a task execution role with those permissions and align it using the `TaskExecutionRoleARN` parameter.

The default configuration will monitor container logs, listen for logs, metrics and traces on port `4317/4318`, and monitor container metrics using the `awsecscontainermetricsd` receiver. The `awsecscontainermetricsd` receiver is based on the [awsecscontainermetrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver) receiver, however, instead of monitoring a single task, metrics will be collected from all containers. If you do not wish to collect container metrics, comment out or delete the metrics/container-metrics pipeline from the configuration.

The embeded default OpenTelemetry configuration can be viewed [here](./template.yaml#L90).

### Health Check

The default config will expose a health check on port `13133` of the localhost via the path `/`. The health check is exposed using the [health_check](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/healthcheckextension) extension.

The healthy response should look like this:

```json
{
  "status": "Server available",
  "upSince": "2023-10-25T15:37:32.003837622Z",
  "uptime": "2m5.2610063s"
}
```

### Further info

See documentation: [AWS ECS-EC2 using OpenTelemetry](https://coralogix.com/docs/opentelemetry-using-ecs-ec2).
