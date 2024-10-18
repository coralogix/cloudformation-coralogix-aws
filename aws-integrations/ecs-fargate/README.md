# ECS Fargate APM integrations

## Note: This ECS Fargate integration used to require the use of fluentbit logrouter for logs processing. This version does everything within OTEL.

## ECS Fargate Logs, Metrics and Traces
Logs, Metrics and Traces are collected using Opentelemetry Collector Contrib.

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
| CoralogixRegion | The region of your Coralogix Account                                                                                                                                                                                                 | *Allowed Values:*<br>- EU1<br>- EU2<br>- AP1<br>- AP2<br>- AP3<br>- US1<br>- US2 | :heavy_check_mark: |
| ParameterName | The name of the Parameter Store to create and use | /CX_OTEL/config.yaml ||
| StorageType | Storage type to use for configuration | *Allowed Values:*<br>- ParameterStore (default)<br>- ParameterStoreAdvanced ||
|

### Deploy the Cloudformation template:

```sh
aws cloudformation deploy --template-file ecs-fargate-cf.yaml \
    --stack-name <stack_name> \
    --region <aws region> \
    --capabilities "CAPABILITY_NAMED_IAM" \
    --parameter-overrides \
        PrivateKey=<your-private-key> \
        CoralogixRegion=<coralogix-region>
```
