# AWS Kinesis Firehose Integration to Coralogix

This template can be used to deploy an AWS Kinesis Firehose Integration to Coralogix.

## Prerequisites
* AWS account.
* Coralogix account.

## Main Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CoralogixRegion | The region of your Coralogix Account | _Allowed Values:_<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US<br>_Default: Europe_ | :heavy_check_mark: |
| CoralogixApiKey | Your Coralogix Private Key | |  :heavy_check_mark: |
| ApplicationName | Your Coralogix Application name | | |
| SubsystemName | Your Coralogix Subsystem name | | |

## Log Stream Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| EnableLogsStream | Enable logs streaming to Coralogix | false | |
| IntegrationTypeLogs | The data structure of the Firehose delivery stream for logs | _Allowed Values:_<br>- CloudWatch_JSON<br>- WAF<br>- CloudWatch_CloudTrail<br>- EksFargate<br>- Default<br>- RawText | |
| KinesisStreamAsSourceARN | If KinesisStreamAsSource for logs is desired, input the ARN of the Kinesis stream |  | |

## Metrics Stream Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| EnableMetricsStream | Enable metrics streaming to Coralogix | true | |
| IntegrationTypeMetrics | The data structure of the Firehose delivery stream for metrics | _Allowed Values:_<br>- opentelemetry0.7<br>- CloudWatch_Metrics_JSON<br> _Default_: CloudWatch_Metrics_OpenTelemetry070 | |
| OutputFormat | The output format of the cloudwatch metric stream | _Allowed Values:_<br>- opentelemetry0.7<br>- json<br> _Default_: **opentelemetry0.7 | |
| IncludeNamespaces | A string comma-delimited list of namespaces to include to the metric stream <br>e.g. `AWS/EC2,AWS/EKS,AWS/ELB,AWS/Logs,AWS/S3` | | |
| IncludeNamespacesMetricNames | A string json list of namespaces and metric_names to include to the metric stream. JSON stringify the input to avoid format errors. <br>e.g. {"AWS/EC2":["CPUUtilization","NetworkOut"],"AWS/S3":["BucketSizeBytes"]} | | |
| AddtionalStatisticsConfigurations | A json list of additional statistics to include to the metric stream following [MetricStream StatisticsConfiguration](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-metricstream-metricstreamstatisticsconfiguration.html). <br>JSON stringify the input to avoid format errors. | | |

## Optional Parameters
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CloudwatchRetentionDays | Enable logs streaming to Coralogix | false | |
| DynamicMetadata | When set to true, it fetches the applicationName / subsystemName dynamically | false | |



## Deploy the Cloudformation template

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> --capabilities CAPABILITY_NAMED_IAM  --parameter-overrides ApplicationName=<application name> SubsystemName=<subsystem name> EventbridgeStream=<EventBridge delivery stream name> RoleName=<EventBridge Role> PrivateKey=<your-private-key> CoralogixRegion=<coralogix-region> CustomUrl=<Custom Coralogix url>
```
