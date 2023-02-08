### ECS Service and Task Definition

This template can be used to deploy an ECS Service and Task Definition for running the Open Telemetry agent on an ECS Cluster.

__Requires:__

- An existing ECS Cluster
- [aws-cli]() (_if deploying via CLI_)


__Parameters:__

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| ClusterName | The name of an __existing__ ECS Cluster |   | :heavy_check_mark: | 
| Image | The open telemtry collector container image.<br><br>ECR Images must be prefixed with the ECR image URI. For eg. `<AccountID>.dkr.ecr.<REGION>.amazonaws.com/image:tag` | coralogixrepo/otel-coralogix-ecs-wrapper | |
| Memory | The amount of memory to allocate to the Open Telemetry container.<br>_Assigning too much memory can lead to the ECS Service not being deployed. Make sure that values are within the range of what is available on your ECS Cluster_ | 256 | |
| CoralogixRegion | The region of your Coralogix Account | _Allowed Values:_<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US | :heavy_check_mark: |
| ApplicationName | You application name |  | :heavy_check_mark: |
| SubsystemName | You Subsystem name | AWS Account ID | __Required__ when using the default  _OTELConfig_ paramter. |
| PrivateKey | Your Coralogix Private Key | | __Required__ when using the default  _OTELConfig_ paramter. |
| OTELConfig | Base64 encoded open telemetry configuration yaml. This value is passed to the Docker container as an environment variable. The value is decoded by the container at runtime and applied to the OTel Agent.<br>Example configuration files can be found [here](https://github.com/coralogix/telemetry-shippers/blob/master/otel-agent/ecs-ec2/config.yaml) | Coralogix Default Configuration ||


__Deploy the Cloudformation template__:

```sh
aws cloudformation deploy --template-file cfn_template.yaml --stack-name cds-68 \
    --region <region> \
    --parameter-overrides \
        ApplicationName=<application name> \
        ClusterName=<ecs cluster name> \
        PrivateKey=<your-private-key> \
        OTELConfig=$(cat path/to/otelconfig.yaml | base64) \
        CoralogixRegion=<coralogix-region>
```