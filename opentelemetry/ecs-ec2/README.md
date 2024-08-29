# Coralogix OpenTelemetry Agent for ECS-EC2. Cloudformation template.

This CloudFormation template deploys an ECS Service and Task Definition for running the Open Telemetry agent on an ECS Cluster. This deployment is able to collect Logs, Metrics and Traces. The template will deploy a daemonset which runs an instance open telemetry on each node in a cluster.

CloudFormation template to launch the Coralogix Distribution for Open Telemetry ("CDOT") into an existing ECS Cluster. This CDOT deployment is able to collect Logs, Metrics and Traces. CDOT is deployed in the OTEL [Agent deployment](https://opentelemetry.io/docs/collector/deployment/agent/) pattern, as an ECS Daemon Service type, which runs an instance of the Open Telemetry collector agent on each node in a cluster.

## Image

<!--
'solution' vs 'example'. We are as a team already supporting this repo. What should be our actual support posture? If this repo is positioned as "supported" by CX? Then:'This solution', is more accurate than:'This example'.
-->

This example uses the coralogixrepo/coralogix-otel-collector image which is a custom distribution of Open Telemetry containing custom components developed by Coralogix. The image is available on [Docker Hub](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector). The ECS components are described [here](./components.md)

The OTEL collector/agent/daemon image used is the [Coralogix Distribution for Open Telemetry](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector) Docker Hub image. It is deployed as a [*Daemon*](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html#service_scheduler_daemon) ECS Task, i.e. one OTEL collector agent container on each EC2 instance (i.e. ECS container instance) across the cluster.

CDOT extends upon the main [Open Telemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) project, adding features specifically to enhance integration with AWS ECS, among other improvements.

The OTEL agent is deployed as a Daemon ECS Task and connected using [`host` network mode](https://docs.aws.amazon.com/AmazonECS/latest/bestpracticesguide/networking-networkmode-host.html). OTEL-instrumented application containers that need to send telemetry to the local OTEL agent can lookup the IP address of the CDOT container [using a number of methods](https://coralogix.com/docs/opentelemetry-using-ecs-ec2/#otel-agent-network-service-discovery), making it easier for Application Tasks using `awsvpc` and `bridge` network modes to connect with the OTEL agent. OTEL-instrumented application containers should also consider which resource attributes to use as telemetry identifiers.

The CDOT OTEL agent also features enhancements specific to ECS integration. These improvements are proprietary to the Coralogix Distribution for Open Telemetry.

### Logs

The OTEL agent uses a [filelog receiver](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/filereceiver) to read the docker logs of all containers on the EC2 host. OTLP is also accepted. Coralogix provides the `awsecscontainermetricsd` receiver which enables metrics collection of all tasks on the same host. The [coralogix exporter](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/coralogixexporter) forwards telemetry to your configured Coralogix endpoint.

Logs are collected from all containers that log to `/var/lib/docker/containers/*/*.log`. The container requires privileges to mount the host read-only file path `/var/lib/docker/`.

### Metrics

Metrics are collected from all containers running on the ECS Cluster. The metrics are collected using the [awsecscontainermetricsd](./components.md#awsecscontainermetricsd) receiver.

### Traces

A GRPC(*4317*) and HTTP(*4318*) endpoint is exposed for sending traces to the local OTLP endpoint.

### Requires:

- An existing ECS Cluster
- [aws-cli]() (*if deploying via CLI*)

## Parameters:

| Parameter        | Description                                                                                                                                                                                                                          | Default Value                                                            | Required           |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|--------------------|
| ClusterName      | The name of an **existing** ECS Cluster                                                                                                                                                                                              |                                                                          | :heavy_check_mark: |
| CDOTImageVersion | The Coralogix Open Telemetry Collector Image version/tag to use. See available tags [here](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector/tags)                                                                     |                                                                          |                    |
| Memory           | The amount of memory to allocate to the Open Telemetry container.<br>*Assigning too much memory can lead to the ECS Service not being deployed. Make sure that values are within the range of what is available on your ECS Cluster* | 256                                                                      |                    |
| CoralogixRegion  | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:* <br>- EU1<br>- EU2<br>- AP1<br>- AP2<br>- US1<br>- US2 | :heavy_check_mark: |
| ApplicationName  | You application name                                                                                                                                                                                                                 |                                                                          | :heavy_check_mark: |
| SubsystemName    | You Subsystem name                                                                                                                                                                                                                   | AWS Account ID                                                           | :heavy_check_mark: |
| CoralogixApiKey  | The Send-Your-Data API key for your Coralogix account. See: https://coralogix.com/docs/send-your-data-api-key/                                                                                                                       |                                                                          | :heavy_check_mark: |
| Metrics          | Enable/Disable Metrics                                                                                                                                                                                                               | disable                                                                  |                    |
| OtelConfig       | The Open Telemetry Configuration Yaml string. This will be used instead of the embedded configuration if specified.<br>**Note** that as of image `v0.3.0` decoding base64 encoded env variables, is no longer supporter.             |                                                                          |                    |

## Deploy the Cloudformation template:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        DefaultApplicationName=<application name> \
        CDOTImageVersion=<image tag> \
        ClusterName=<ecs cluster name> \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region>
```

**When using a custom configuration for Open Telemetry**

```sh
aws cloudformation deploy --template-file cfn_template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        DefaultApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        CDOTImageVersion=<image tag> \
        CoralogixApiKey=<your-private-key> \
        OtelConfig=$(cat path/to/otelconfig.yaml) \
        CoralogixRegion=<coralogix-region>
```

Note that these are just examples of how this could be deployed. You can also deploy this template using the AWS Console or any Cloudformation management tools.

### Open Telemetry Configuration

The Open Telemetry configuration is embedded in this cloudformation template by default, however, you do have the option of specifying your own configuration using the `OtelConfig` parameter. The configuration must be base64 encoded.

The default configuration will monitor container logs and listen for traces on port `4317/4318`. To enable metrics, set the `Metrics` parameter to `enabled`. This will add the `awsecscontainermetricsd` receiver which is based on the [awsecscontainermetrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver) receiver, however, instead of monitoring a single task, metrics will be collected from all containers.

The default Open Telemetry configuration can be view [here](./template.yaml#L75-L161) and the metrics [here](./template.yaml#L164-L262).

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

### Open Telemetry Collector Metrics

The default Open Telemetry configuration embedded in this cloudformation template exposes metrics about the collector on port `8888` via the path `/metrics`. The metrics are collected using a [prometheus](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/prometheusreceiver) scrape job.

### Further info

See documentation: [AWS ECS-EC2 using OpenTelemetry](https://coralogix.com/docs/opentelemetry-using-ecs-ec2).
