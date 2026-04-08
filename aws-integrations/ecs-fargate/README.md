# ECS Fargate APM integrations

## ECS Fargate Logs, Metrics and Traces
Logs, Metrics and Traces are collected using Opentelemetry Collector Contrib.

Details of this integration can be found on [github](https://github.com/coralogix/telemetry-shippers/blob/master/otel-ecs-fargate/README.md) and at [coralogix.com](https://coralogix.com/docs/integrations/aws/opentelemetry-ecs-fargate/)

## Example Cloudformation Template

This Cloudformation Template is provided as a guide for building your own CF template for your Fargate Task Definitions, it is not intended to be used as is.
It does demonstrate which permissions are required, at a minimum, to deploy our integration. These permissions may not be sufficient for your application needs, adjust as necessary.

**Requires:**

- [aws-cli]() (*if deploying via CLI*)

## Parameters:

| Parameter        | Description                                                                                                                                                                                                                          | Default Value                                                            | Required           |
|------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------|--------------------|
| Image            | The OpenTelemetry Collector Image to use.                                                                      | otel/opentelemetry-collector-contrib:0.147.0                                                                     |                    |
| Memory           | The amount of memory to allocate to the OpenTelemetry container.<br>*Assigning too much memory can lead to the ECS Service not being deployed. Make sure that values are within the range of what is available on your ECS Cluster* | 2048                                                                     |                    |
| AwsPlatform      | Commercial AWS or AWS GovCloud. Must match the partition where you deploy this stack (the template validates the stack partition). GovCloud uses us-gov-west-1 / us-gov-east-1.<br>*Allowed Values:* Standard, AWSGovCloud | Standard |                    |
| CoralogixRegion  | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:* <br>- EU1<br>- EU2<br>- AP1<br>- AP2<br>- AP3<br>- US1<br>- US2 | :heavy_check_mark: |
| CoralogixApiKey  | The Send-Your-Data API key for your Coralogix account. See: https://coralogix.com/docs/send-your-data-api-key/                                                                                                                       |                                                                          | :heavy_check_mark: |
| S3ConfigBucket   | S3 bucket name containing the configuration file. |                                                                          | :heavy_check_mark: |
| S3ConfigKey      | S3 object key (file path) for the configuration file. |                                                                          | :heavy_check_mark: |
| HealthCheckEnabled | Enable ECS container health check for the OTEL agent container. | false                                                                    |                    |
| HealthCheckInterval | Health check interval (seconds)                                                                                                                                                                                                     | 30                                                                       |                    |
| HealthCheckTimeout | Health check timeout (seconds)                                                                                                                                                                                                      | 5                                                                        |                    |
| HealthCheckRetries | Health check retries                                                                                                                                                                                                               | 3                                                                        |                    |
| HealthCheckStartPeriod | Health check start period (seconds)                                                                                                                                                                                               | 10                                                                       |                    |

### Deploy the Cloudformation template:

```sh
aws cloudformation deploy --template-file ecs-fargate-cf.yaml \
    --stack-name <stack_name> \
    --region <aws region> \
    --capabilities "CAPABILITY_NAMED_IAM" \
    --parameter-overrides \
        CoralogixApiKey=<your-private-key> \
        CoralogixRegion=<coralogix-region> \
        S3ConfigBucket=<example-s3-bucket> \
        S3ConfigKey=<example-s3-config-object>
```

For AWS GovCloud, include `AwsPlatform=AWSGovCloud` in `--parameter-overrides`. For commercial regions, `Standard` is the default and does not need to be set.
