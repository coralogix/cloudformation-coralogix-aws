# coralogix-aws-shipper (Beta)

[![license](https://img.shields.io/github/license/coralogix/coralogix-aws-shipper.svg)](https://raw.githubusercontent.com/coralogix/coralogix-aws-shipper/master/LICENSE)

![publish workflow](https://github.com/coralogix/coralogix-aws-shipper/actions/workflows/publish.yaml/badge.svg)

![Dynamic TOML Badge](https://img.shields.io/badge/dynamic/toml?url=https%3A%2F%2Fraw.githubusercontent.com%2Fcoralogix%2Fcoralogix-aws-shipper%2Fmaster%2FCargo.toml%3Ftoken%3DGHSAT0AAAAAACJIQT3CA3OFRU7Z5NU4T6YKZLPLLSQ&query=%24.package.version&label=version)

![Static Badge](https://img.shields.io/badge/status-beta-purple)


## Overview
Coralogix provides a predefined AWS Lambda function to easily forward your logs to the Coralogix platform.

The `coralogix-aws-shipper` supports forwarding of logs for the following AWS Services:

* [Amazon CloudWatch](https://docs.aws.amazon.com/cloudwatch/)
* [AWS CloudTrail](https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-log-file-examples.html)
* [Amazon VPC Flow logs](https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs-s3.html)
* AWS Elastic Load Balancing access logs ([ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-access-logs.html), [NLB](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-access-logs.html) and [ELB](https://docs.aws.amazon.com/elasticloadbalancing/latest/classic/access-log-collection.html))
* [Amazon CloudFront](https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/AccessLogs.html)
* [AWS Network Firewall](https://docs.aws.amazon.com/network-firewall/latest/developerguide/logging-s3.html)
* [Amazon Redshift](https://docs.aws.amazon.com/redshift/latest/mgmt/db-auditing.html#db-auditing-manage-log-files)
* [Amazon S3 access logs](https://docs.aws.amazon.com/AmazonS3/latest/userguide/ServerLogs.html)
* [Amazon VPC DNS query logs](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/resolver-query-logs.html)
* [AWS WAF](https://docs.aws.amazon.com/waf/latest/developerguide/logging-s3.html)

Additionally, you can ingest any generic text, JSON and csv logs stored in your S3 bucket

## Prerequisites

* AWS account (Your AWS user should have permissions to create lambdas and IAM roles)
* Coralogix account
* The application should be installed in the same AWS region as your resource are (i.e the S3 bucket you want to send the logs from)

## Deployment instructions

### AWS Serverless Application

The lambda can be deployed by clicking the link below and signing into your AWS account:
[Deployment link](https://serverlessrepo.aws.amazon.com/applications/eu-central-1/597078901540/Coralogix-aws-shipper)
Please make sure you selecet the AWS region before you deploy

### Coralogix In Product Integration

Link To Coralogix Document TBD

### AWS CloudFormation Application

Log into your AWS account and deploy the CloudFormation Stack with the button below
[![Launch Stack](https://s3.amazonaws.com/cloudformation-examples/cloudformation-launch-stack.png)](https://console.aws.amazon.com/cloudformation/home#/stacks/create/review?stackName=coralogix-aws-shipper&templateURL=https://cgx-cloudformation-templates.s3.amazonaws.com/aws-integrations/aws-shipper-lambda/template.yaml)

### Terraform

Link To coralogix Module TBD

## Paramaters 

### Coralogix configuration
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| Application name | The stack name of this application created via AWS CloudFormation. |   | :heavy_check_mark: |
| IntegrationType | The integration type. Can be one of: s3, cloudtrail, vpcflow, cloudwatch, s3_csv' |  s3 | :heavy_check_mark: | 
| CoralogixRegion | The Coralogix location region, possible options are [Custom, Europe, Europe2, India, Singapore, US, US2] If this value is set to Custom you must specify the Custom Domain to use via the CustomDomain parameter |  Custom | :heavy_check_mark: | 
| CustomDomain | The Custom Domain. If set, will be the domain used to send telemetry (e.g. cx123.coralogix.com) |   |   |
| ApplicationName | The [name](https://coralogix.com/docs/application-and-subsystem-names/) of your application. for dynamically value from the log you should use $.my_log.field |   | :heavy_check_mark: | 
| SubsystemName | The [name](https://coralogix.com/docs/application-and-subsystem-names/) of your subsystem. for dynamically value from the log you should use $.my_log.field . for cloudwatch loggroup leave empty |   |   |
| ApiKey | Your Coralogix Send Your Data - [API Key](https://coralogix.com/docs/send-your-data-api-key/) which is used to validate your authenticity, This value can be a Coralogix API Key or an AWS Secret Manager ARN that holds the API Key |   | :heavy_check_mark: |
| StoreAPIKeyInSecretsManager | Store the API key in AWS Secrets Manager.  If this option is set to false, the ApiKey will apeear in plain text as an environment variable in the lambda function console. | True  | :heavy_check_mark: |

### Integration S3/Cloudtrail/vpcflow/s3_csv configuration
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| S3BucketName | The name of the AWS S3 bucket to watch |   | :heavy_check_mark: |
| S3KeyPrefix | The AWS S3 path prefix to watch | cloudtrail 'AWSLogs/' |   |
| S3KeySuffix | The AWS S3 path suffix to watch | cloudtrail/vpcflow '.json.gz' |   |
| NewlinePattern | Regular expression to detect a new log line for multiline logs from S3 source, e.g., use expression \n(?:\r\n\|\r\|\n) | \n(?:\r\n\|\r\|\n) |   |
| SNSTopicArn | The ARN of SNS topic that will contain the SNS subscription for retrieving logs from S3 |   |   |
| CsvDelimiter | Single Character for using as a Delimiter when ingesting CSV (This value is applied when the s3_csv integration type  is selected), e.g. "," or " " |   |   |

### Integration Cloudwatch configuration
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CloudWatchLogGroupName | A comma separated list of CloudWatch log groups names to watch  e.g, (log-group1,log-group2,log-group3) |   | :heavy_check_mark: | 

### Integration Generic Config (Optional)
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| NotificationEmail | Failure notification email address |   |   | 
| BlockingPattern | Regular expression to detect lines that should be excluded from sent to Coralogix |   |   | 
| SamplingRate | Send messages with specific rate (1 out of N) e.g., put the value 10 if you want to send every 10th log | 1 | :heavy_check_mark: | 


### Lambda configuration (Optional)
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| FunctionMemorySize | Memory size for lambda function in mb | 1024 | :heavy_check_mark: | 
| FunctionTimeout | Timeout for the lambda function in sec | 300 | :heavy_check_mark: | 
| LogLevel | Log level for the Lambda function. Can be one of: INFO, WARNING, ERROR, DEBUG | INFO | :heavy_check_mark: | 
| LambdaLogRetention | CloudWatch log retention days for logs generated by the Lambda function | 5 | :heavy_check_mark: | 

### VPC configuration (Optional)
| Parameter | Description | Default Value | Required |
|---|---|---|---|
| LambdaSubnetID | ID of Subnet into which to deploy the integration |   | :heavy_check_mark: | 
| LambdaSecurityGroupID | ID of the SecurityGroup into which to deploy the integration |   | :heavy_check_mark: | 
| UsePrivateLink | Will you be using our PrivateLink? | false | :heavy_check_mark: | 

## Advanced

### AWS PrivateLink
To use privatelink please forllow the instruction in this [link](https://coralogix.com/docs/coralogix-amazon-web-services-aws-privatelink-endpoints/)

## Troubleshooting
TBD
