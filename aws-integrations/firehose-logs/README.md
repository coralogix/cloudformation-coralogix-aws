# AWS Kinesis Firehose Integration to Coralogix

This template can be used to deploy an AWS Kinesis Firehose Integration to send resource logs Coralogix.

For a more detailed description of the settigns and architecture of this AWS Kinesis Data Firehose setup, please refer to the Coralogix documentation on [AWS Kinesis Data Firehose – Logs](https://coralogix.com/docs/aws-firehose/).

## Prerequisites
* AWS account.
* Coralogix account.

## Main Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CoralogixRegion | The region of your Coralogix Account. If set to Custom, you must provide a CustomDomain otherwise url will be invalid. | _Allowed Values:_<br>- Custom<br>- Europe<br>- Europe2<br>- India<br>- Singapore<br>- US<br>- US2<br>_Default_: Custom | :heavy_check_mark: |
| CustomDomain | The Custom Coralogix domain. If set, will be the domain to send telemetry. | | |
| ApiKey | Your Coralogix Private Key | |  :heavy_check_mark: |
| ApplicationName | Your Coralogix Application name | | |
| SubsystemName | Your Coralogix Subsystem name | | |

## Log Stream Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| IntegrationTypeLogs | The data structure of the Firehose delivery stream for logs | _Allowed Values:_<br>- CloudWatch_JSON<br>- WAF<br>- CloudWatch_CloudTrail<br>- EksFargate<br>- Default<br>- RawText | |
| DynamicMetadataLogs | When set to true, it fetches the applicationName / subsystemName dynamically for logs | false | |
| KinesisStreamAsSourceARN | If KinesisStreamAsSource for logs is desired, input the ARN of the Kinesis stream |  | |

## Optional Parameters
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CloudwatchRetentionDays | Days of retention in Cloudwatch retention days | 1 | |

## Notes:

* If you want to use the Kinesis Stream as a source for logs, you must create the Kinesis Stream before deploying the Cloudformation template and set the KinesisStreamAsSourceARN parameter to the ARN of the Kinesis Stream.

## Dynamic Values Table for Logs

For `ApplicationName` and/or `SubsystemName` to be set dynamically in relation to their `integrationType` resource fields (e.g. CloudWatch_JSON's loggroup name, EksFargate's k8s namespace). The source's `var` has to be mapped as a string literal to the `integrationType`'s as a DyanamicFromFrield with pre-defined values:

| Field | Source `var` | Expected String Literal | Integration Type | Notes |
|-------|--------------|-------------------------|------------------|-------|
| `applicationName` field in logs | applicationName | `${applicationName}` | Default | need to be supplied in the log to be used |
| `subsystemName` field in logs | subsystemName | `${subsystemName}` | Default |  need to be supplied in the log to be used |
| CloudWatch LogGroup name | logGroup | `${logGroup}` | CloudWatch_JSON <br/> CloudWatch_CloudTrail | supplied by aws |
| `kubernetes.namespace_name` field | kubernetesNamespaceName | `${kubernetesNamespaceName}` | EksFargate | supplied by the default configuration |
| `kubernetes.container_name` field | kubernetesContainerName | `${kubernetesContainerName}` | EksFargate | supplied by the default configuration |
| name part of the `log.webaclId` field | webAclName | `${webAclName}` | WAF | supplied by aws |

For more information - visit [Kinesis Data Firehose - Logs](https://coralogix.com/docs/aws-firehose/).

Note: `RawText` integrationType does not support dynamic values.

## Deploy the Cloudformation template using aws cli

With the aws cli installed and configured, run the following command:

```sh
aws cloudformation create-stack --stack-name <stack_name> --template-body template.yaml --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_IAM CAPABILITY_NAMED_IAM --parameter-overrides CoralogixDomain=<domain> ApiKey=<coralogix_api_key> ApplicationName=<application_name> SubsystemName=<subsystem_name> 
```

or with a parameters json file example:

```sh
aws cloudformation create-stack --stack-name <stack_name> --template-body template.yaml --parameters parameters.json --capabilities CAPABILITY_AUTO_EXPAND CAPABILITY_IAM CAPABILITY_NAMED_IAM
```
