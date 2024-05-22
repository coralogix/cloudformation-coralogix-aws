# ECS EC2 Windows (Metrics)


This template section provides an example template for deployin the Open Telemetry collector as a sidecar for collecting Metrics on ECS EC2 Windows. This template is not meant to be used in production as is, it is intended as a demonstration/example.



**Requires:**

- An existing Windows ECS EC2 Cluster
- [aws-cli]() (*if deploying via CLI*)

### Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| ClusterName     | The name of an **existing** ECS Cluster                                                                                                                                                                                              |                                                                              | :heavy_check_mark: |
| OTelImage           | The open telemtry collector container image.<br><br>ECR Images must be prefixfd with the ECR image URI. For eg. `<AccountID>.dkr.ecr.<REGION>.amazonaws.com/image:tag`                                                               | coralogixrepo/otel-coralogix-ecs-ec2                                         |                    |
| AppImage          | Your windows application container image |                                                                          |                    |
| CoralogixRegion | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:*<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US | :heavy_check_mark: |
| ApplicationName | You application name                                                                                                                                                                                                                 |                                                                              | :heavy_check_mark: |
| SubsystemName   | You Subsystem name                                                                                                                                                                                                                   | AWS Account ID                                                               | :heavy_check_mark: |
| PrivateKey      | Your Coralogix Private Key                                                                                                                                                                                                           |                                                                              | :heavy_check_mark: |



### How it works

This deployment uses the [awsecscontainermetricsd](../ecs-ec2/components.md#aws-ecs-container-metrics-daemonset-receiver) receiver by Coralogix, to collect metrics from Windows applications. For windows deployments the receiver must be deployed in Sidecar mode.

```yaml
receivers:
  awsecscontainermetricsd:
    collection_interval: 20s
    sidecar: true
```

For Windows, this receiver does not support being run as a daemonset, as such, each intances of the collector must be added as a sidecar/container within the Task Definition of the container(s) you wish to collect metrics from.


### Open Telemetry Configuration

The Open Telemetry configuration is embedded in this CloudFormation template by default. However, you have the option of specifying your own configuration by modifying the template. The Coralogix Open Telemetry distribution supports reading configuration from S3 as well as from an Environment Variable. Note that environment variables can be raw strings or Base64 encoded. Configuration from S3 must be passed to the collector using the S3 URL of the object. For example: cdot --config `s3://{your-bucket}.s3.{region}.amazonaws.com/{your-object-key}`. When using this feature in ECS, the host or task must have sufficient permissions to read S3 objects.

### Image

This example uses the _windows-server-1809_ tagged version of the [coralogixrepo/coralogix-otel-collector](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector/tags) image which is a custom distribution of Open Telemetry containing custom components developed by Coralogix. The image is available on [Docker Hub](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector). The ECS components are described [here](../ecs-ec2/components.md)

