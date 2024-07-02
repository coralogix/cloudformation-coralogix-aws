# ECS Fargate APM integrations

## ECS Fargate Logs
Logs are collected using a sidecar deployment of [aws-for-fluent-bit](https://github.com/aws/aws-for-fluent-bit) as a firelens log_router.

Details of this integration can be found [here](https://github.com/coralogix/telemetry-shippers/tree/master/logs/fluent-bit/ecs-fargate)

## ECS Fargate Traces and Metrics
Traces and Metrics are collected using Opentelemetry Collector Contrib.

Details of this integration can be found [here](https://github.com/coralogix/telemetry-shippers/blob/master/otel-ecs-fargate/README.md)

## Example Cloudformation Template

The included Cloudformation Template is provided as a guide for building your own CF template for your Fargate Task Definitions, it is not intended to be used as is.

It does demonstrate which permissions are required, at a minimum, to deploy our integrations. These permissions may not be sufficient for your application needs, adjust as necessary.

**Requires:**

- An existing ECS Cluster
- [aws-cli]() (*if deploying via CLI*)

### Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| PrivateKey      | Your Coralogix Private Key                                                                                                                                                                                                           |                                                                              | :heavy_check_mark: |
| CoralogixRegion | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:*<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US<br>- US2 | :heavy_check_mark: |
| S3ConfigARN      | The S3 ARN for your uploaded Coralogix Fluent Bit configuration file (Explained in the ECS Fargate Logs integration documentation above)                                                                                                                                                                                                           |                                                                              | :heavy_check_mark: |


### Deploy the Cloudformation template:

```sh
aws cloudformation deploy --template-file ecs-fargate-cf.yaml \
    --stack-name <stack_name> \
    --region <aws region> \
    --capabilities "CAPABILITY_NAMED_IAM" \
    --parameter-overrides \
        PrivateKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        S3ConfigARN=<ARN of S3 Config>
```
