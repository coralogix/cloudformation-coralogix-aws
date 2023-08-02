### ECS Service and Task Definition

This template can be used to deploy an ECS Service and Task Definition for running the Open Telemetry agent on an ECS Cluster. This deployment is able to collect Logs, Metrics and Traces. The template will deploy a daemonset which runs an instance open telemetry on each node in a cluster.

**Logs**

> Logs are collected from all containers that log to `/var/lib/docker/containers/*/*.log` 

**Metrics**

> Metrics are collected from all containers running on the ECS Cluster. The metrics are collected using the [awsecscontainermetricsd](./components.md#awsecscontainermetricsd) receiver.

**Traces**

> A GRPC(_4317_) and HTTP(_4318_) endpoint is exposed for sending traces to the local OTLP endpoint.

**Requires:**

- An existing ECS Cluster
- [aws-cli]() (*if deploying via CLI*)

### Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| ClusterName     | The name of an **existing** ECS Cluster                                                                                                                                                                                              |                                                                              | :heavy_check_mark: |
| Image           | The open telemtry collector container image.<br><br>ECR Images must be prefixed with the ECR image URI. For eg. `<AccountID>.dkr.ecr.<REGION>.amazonaws.com/image:tag`                                                               | coralogixrepo/otel-coralogix-ecs-ec2                                         |                    |
| Memory          | The amount of memory to allocate to the Open Telemetry container.<br>*Assigning too much memory can lead to the ECS Service not being deployed. Make sure that values are within the range of what is available on your ECS Cluster* | 256                                                                          |                    |
| CoralogixRegion | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:*<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US | :heavy_check_mark: |
| ApplicationName | You application name                                                                                                                                                                                                                 |                                                                              | :heavy_check_mark: |
| SubsystemName   | You Subsystem name                                                                                                                                                                                                                   | AWS Account ID                                                               | :heavy_check_mark: |
| PrivateKey      | Your Coralogix Private Key                                                                                                                                                                                                           |                                                                              | :heavy_check_mark: |
| Metrics         | Enable/Disable Metrics                                                                                                                                                                                                               | disable                                                                      |                    |

### Deploy the Cloudformation template:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        ApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        PrivateKey=<your-private-key> \
        CoralogixRegion=<coralogix-region>
```

### Open Telemetry Configuration

The Open Telemetry configuration is embedded in this cloudformation template.

The default configuration will monitor container logs and listen for traces on port `4317/4318`. To enable metrics, set the `Metrics` parameter to `enabled`. This will add the `awsecscontainermetricsd` receiver which is based on the [awsecscontainermetrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver) receiver, however, instead of monitoring a single task, metrics will be collected from all containers.

The default Open Telemetry configuration can be view [here](./template.yaml#L75-L161) and the metrics [here](./template.yaml#L164-L262).


### Image

This example uses the coralogixrepo/coralogix-otel-collector image which is a custom distribution of Open Telemetry containing custom components developed by Coralogix. The image is available on [Docker Hub](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector). The ECS components are described [here](./components.md)



