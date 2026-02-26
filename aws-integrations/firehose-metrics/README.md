# AWS Kinesis Firehose Integration to Coralogix

This template can be used to deploy an AWS Kinesis Firehose Integration to send resource metrics to Coralogix.

For a more detailed description of the settigns and architecture of this AWS Kinesis Data Firehose setup, please refer to the Coralogix documentation on [AWS Kinesis Data Firehose – Metrics](https://coralogix.com/docs/amazon-kinesis-data-firehose-metrics/).

## Prerequisites

* AWS account.
* Coralogix account.

## Main Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CoralogixRegion | The region of your Coralogix Account. If set to Custom, you must provide a CustomDomain otherwise url will be invalid. | _Allowed Values:_<br>- Custom<br>- EU1<br>- EU2<br>- AP1<br>- AP2<br>- AP3<br>- US1<br>- US2<br>_Default_: Custom | :heavy_check_mark: |
| CustomDomain | The Custom Coralogix domain. If set, will be the domain to send telemetry. | | |
| ApiKey | Your Coralogix Private Key | |  :heavy_check_mark: |
| ApplicationName | Your Coralogix Application name | | |
| SubsystemName | Your Coralogix Subsystem name | | |

## Metrics Stream Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| IntegrationTypeMetrics | The data structure of the Firehose delivery stream for metrics. For `CloudWatch_Metrics_OpenTelemetry070_WithAggregations`, additional aggregations here are `_min`, `_max`, `_avg` recorded as gauges. See https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/CloudWatch-metric-streams-formats-opentelemetry-translation.html | _Allowed Values:_<br>- CloudWatch_Metrics_OpenTelemetry070<br>- CloudWatch_Metrics_OpenTelemetry070_WithAggregations<br> _Default_: CloudWatch_Metrics_OpenTelemetry070_WithAggregations | |
| OutputFormat | The output format of the cloudwatch metric stream | _Allowed Values:_<br>- opentelemetry0.7<br>- json<br> _Default_: opentelemetry0.7 | |
| IncludeNamespaces | A string comma-delimited list of namespaces to include to the metric stream <br>e.g. `AWS/EC2,AWS/EKS,AWS/ELB,AWS/Logs,AWS/S3` | | |
| IncludeNamespacesMetricNames | A string json list of namespaces and metric_names to include to the metric stream. JSON stringify the input to avoid format errors. <br>e.g. {"AWS/EC2":["CPUUtilization","NetworkOut"],"AWS/S3":["BucketSizeBytes"]} | | |
| AddtionalStatisticsConfigurations | A json list of additional statistics to include to the metric stream following [MetricStream StatisticsConfiguration](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudwatch-metricstream-metricstreamstatisticsconfiguration.html). <br>JSON stringify the input to avoid format errors. | "p50","p75","p95","p99" of the following <br>- AWS/EBS:[VolumeTotalReadTime,VolumeTotalWriteTime]<br>- AWS/ELB:[Latency,Duration], <br>- AWS/Lambda:[PostRuntimeExtensionsDuration]<br>- AWS/S3:[FirstByteLatency,TotalRequestLatency] | |
| EnableMetricsTagsProcessors | Enable the lambda metrics tags processor function. Set to false to remove the processor | true | |
| IncludeLinkedAccountsMetrics | Enable cross-account observability to include metrics from linked source accounts (requires CloudWatch OAM setup between monitoring and source accounts) | false | |
| CrossAccountEnabled | Enable Lambda cross-account tag enrichment. When true, Lambda assumes per-account roles from `CrossAccountRoles`. | false | |
| CrossAccountRoles | JSON map of source account IDs to role ARNs used for tag enrichment. Example: `{"597078901540":"arn:aws:iam::597078901540:role/CoralogixMetricsReader"}` | {} | |
| FileCacheEnabled | Enable local file cache for resource discovery in Lambda. | true | |
| FileCacheExpiration | Cache expiration in Go duration format (for example `1h`, `30m`). | 1h | |
| FileCachePath | Existing directory path for Lambda cache files. Use `/tmp` in Lambda. | /tmp | |

## Optional Parameters
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CloudwatchRetentionDays | Days of retention in Cloudwatch retention days | 1 | |

## Notes:

* For cross-account enrichment, each linked account role (for example `CoralogixMetricsReader`) must trust the monitoring account Lambda role and allow `sts:AssumeRole`:
  - Principal: `arn:aws:iam::<monitoring-account-id>:role/<stack-name>-lambda`
  - Action: `sts:AssumeRole`

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
