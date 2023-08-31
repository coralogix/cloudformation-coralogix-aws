# AWS Kinesis Firehose Integration to Coralogix

This template can be used to deploy an AWS Kinesis Firehose Integration to send resource logs and metrics to Coralogix.

For a more detailed description of the settigns and architecture of this AWS Kinesis Data Firehose setup, please refer to the Coralogix documentation on [AWS Kinesis Data Firehose – Logs](https://coralogix.com/docs/aws-firehose/) and [AWS Kinesis Data Firehose – Metrics](https://coralogix.com/docs/amazon-kinesis-data-firehose-metrics/).

## Prerequisites
* AWS account.
* Coralogix account.

## Main Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CoralogixRegion | The region of your Coralogix Account | _Allowed Values:_<br>- ireland<br>- stockholm<br>- india<br>- singapore<br>- us<br>- us2<br>_Default: ireland_ | :heavy_check_mark: |
| ApiKey | Your Coralogix Private Key | |  :heavy_check_mark: |
| ApplicationName | Your Coralogix Application name | | |
| SubsystemName | Your Coralogix Subsystem name | | |
| CustomUrl | The custom url domain. If set, will be the url used to send telemetry. | | |

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
| OutputFormat | The output format of the cloudwatch metric stream | _Allowed Values:_<br>- opentelemetry0.7<br>- json<br> _Default_: opentelemetry0.7 | |
| IncludeNamespaces | A string comma-delimited list of namespaces to include to the metric stream <br>e.g. `AWS/EC2,AWS/EKS,AWS/ELB,AWS/Logs,AWS/S3` | | |
| IncludeNamespacesMetricNames | A string json list of namespaces and metric_names to include to the metric stream. JSON stringify the input to avoid format errors. <br>e.g. {"AWS/EC2":["CPUUtilization","NetworkOut"],"AWS/S3":["BucketSizeBytes"]} | | |
| AddtionalStatisticsConfigurations | A json list of additional statistics to include to the metric stream following [MetricStream StatisticsConfiguration](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-metricstream-metricstreamstatisticsconfiguration.html). <br>JSON stringify the input to avoid format errors. | "p50","p75","p95","p99" of the following <br>- AWS/EBS:[VolumeTotalReadTime,VolumeTotalWriteTime]<br>- AWS/ELB:[Latency,Duration], <br>- AWS/Lambda:[PostRuntimeExtensionsDuration]<br>- AWS/S3:[FirstByteLatency,TotalRequestLatency] | |

## Optional Parameters
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CloudwatchRetentionDays | Enable logs streaming to Coralogix | false | |
| DynamicMetadata | When set to true, it fetches the applicationName / subsystemName dynamically | false | |

## Notes:

* If you want to use the Kinesis Stream as a source for logs, you must create the Kinesis Stream before deploying the Cloudformation template and set the KinesisStreamAsSourceARN parameter to the ARN of the Kinesis Stream.
* If DynamicMetadata is set to `dynamicMetadata`: `true`, the applicationName and subsystemName for logs will be based on the selected IntegrationTypeLogs and follow the below Dynamic values table:

| Type | Dynamic applicationName | Dynamic subsystemName | Notes |
| --- | --- | --- | --- |
| CloudWatch_JSON | the cloudwatch log group | none | supplied by aws |
| CloudWatch_CloudTrail | the cloudwatch log group | none | supplied by aws |
| Default | ‘applicationName’ field	| ‘subsystemName’ field	| need to be supplied in the log to be used |
| EksFargate | ‘kubernetes.namespace_name’ field | ‘kubernetes.container_name’ field | supplied by the default configuration |
| WAF | The web acl name | none | supplied by aws |

* `LambdaMetricsTagsProcessors` lambda function code is deployed to the following s3 regions: [ _us-east-1 us-east-2 us-west-1 us-west-2 ap-south-1 ap-northeast-2 ap-southeast-1 ap-southeast-2 ap-northeast-1 ca-central-1 eu-central-1 eu-west-1 eu-west-2 eu-west-3 eu-north-1 eu-south-1 sa-east-1_ ]. If you are using a different region, please contact Coralogix support.

## Deploy the Cloudformation template using aws cli

With the aws cli installed and configured, run the following command:

```sh
aws cloudformation create-stack --stack-name <stack_name> --template-body template.yaml --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_IAM CAPABILITY_NAMED_IAM --parameter-overrides ApiKey=<coralogix_api_key> CoralogixRegion=<region> ApplicationName=<application_name> SubsystemName=<subsystem_name> EnableLogsStream=<true/false> EnableMetricsStream=<true/false> 
```

or with a parameters json file example:

```sh
aws cloudformation create-stack --stack-name <stack_name> --template-body template.yaml --parameters parameters.json --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_IAM CAPABILITY_NAMED_IAM
```
