# ECS EC2 Windows (Metrics)


This template section provides an example template for deploying the Open Telemetry collector as a sidecar for collecting Metrics on ECS EC2 Windows. This template is not meant to be used in production as is, it is intended as a demonstration/example.



**Requires:**

- An existing Windows ECS EC2 Cluster
- [aws-cli]() (*if deploying via CLI*)

### Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| ClusterName     | The name of an **existing** ECS Cluster                                                                                                                                                                                              |                                                                              | :heavy_check_mark: |
| OTelImage           | The open telemtry collector container image.<br><br>ECR Images must be prefixed with the ECR image URI. For eg. `<AccountID>.dkr.ecr.<REGION>.amazonaws.com/image:tag`                                                               | coralogixrepo/otel-coralogix-ecs-ec2                                         |                    |
| AppImage          | Your windows application container image |                                                                          |                    |
| CoralogixRegion | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:*<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US | :heavy_check_mark: |
| ApplicationName | You application name                                                                                                                                                                                                                 |                                                                              | :heavy_check_mark: |
| SubsystemName   | You Subsystem name                                                                                                                                                                                                                   | AWS Account ID                                                               | :heavy_check_mark: |
| PrivateKey      | Your Coralogix Private Key                                                                                                                                                                                                           |                                                                              | :heavy_check_mark: |
| TaskRoleArn     | Optional ARN of the task role (IAM role) that the container can assume. If not provided, the task will run without a task role (null). This is separate from the execution role which is used by ECS to pull images and retrieve secrets. | "" | |
| HealthCheckEnabled     | Enable ECS container health check for the OTEL agent container. | false | |
| HealthCheckInterval    | Health check interval (seconds)             | 30      |          |
| HealthCheckTimeout     | Health check timeout (seconds)              | 5       |          |
| HealthCheckRetries     | Health check retries                        | 3       |          |
| HealthCheckStartPeriod | Health check start period (seconds)         | 60      |          |



### IAM Role Management

This template separates execution and task IAM roles for security:

- **Execution Role** (`ECSTaskExecutionRole`): Used by the ECS agent for infrastructure operations (pulling images, sending logs, retrieving secrets). This role has broader permissions needed for ECS operations.
- **Task Role** (`TaskRoleArn`): Optional IAM role used by the container at runtime. If not provided, the task runs without a task role (null). This follows the principle of least privilege - the container doesn't need AWS permissions for this deployment since it only sends metrics to Coralogix.

**Note:** Since this template doesn't require AWS service access at runtime (no S3 config, no CloudMap discovery), the task role can be omitted. If you need to add AWS permissions for your specific use case, provide a custom `TaskRoleArn` with minimal required permissions.

### How it works

This deployment uses the the [awsecscontainermetricsd](../ecs-ec2/components.md#aws-ecs-container-metrics-daemonset-receiver) receiver by Coralogix, to collect metrics from Windows application. For windows deployments the receiver must be deployed in Sidecar mode.

```yaml
receivers:
  awsecscontainermetricsd:
    collection_interval: 20s
    sidecar: true
```

For Windows, this receiver does not support being run as a daemonset, as such, each intances of the collector must be added as a sidecar/container within the Task Definition of the container(s) you wish to collect metrics from.


### Open Telemetry Configuration

The Open Telemetry configuration is embedded in this CloudFormation template by default. However, you have the option of specifying your own configuration by modifying the template. The Coralogix Open Telemetry distribution supports reading configuration over HTTP as well as an Environment Variable. Note that Environment Variables must be raw strings. Base64 encoding is not supported.

### ECS Container Health Check

You can enable and customize the ECS container health check for the OTEL agent container using the following parameters:
- `HealthCheckEnabled` (default: false)
- `HealthCheckInterval` (default: 30)
- `HealthCheckTimeout` (default: 5)
- `HealthCheckRetries` (default: 3)
- `HealthCheckStartPeriod` (default: 60)

Example deployment with custom health check settings:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
    --region <region> \
    --parameter-overrides \
        ClusterName=<ecs cluster name> \
        AppImage=<your-app-image> \
        OtelImage=<otel-image> \
        PrivateKey=<your-private-key> \
        TaskRoleArn=<optional-task-role-arn> \
        ApplicationName=<application name> \
        SubsystemName=<subsystem name> \
        CoralogixRegion=<coralogix-region> \
        HealthCheckEnabled=true \
        HealthCheckInterval=60 \
        HealthCheckTimeout=10 \
        HealthCheckRetries=5 \
        HealthCheckStartPeriod=60
```

### Image

This example uses the [coralogixrepo/coralogix-otel-collector:0.4.1-windowsserver-1809](https://hub.docker.com/layers/coralogixrepo/coralogix-otel-collector/0.4.1-windowsserver-1809/images/sha256-c436b2b29501592449e2b72a2393f4825e7216bdd62d90cb5e14463a46fafd95?context=explore) image which is a custom distribution of Open Telemetry containing custom components developed by Coralogix. The image is available on [Docker Hub](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector). The ECS components are described [here](../ecs-ec2/components.md)
