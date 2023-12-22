# EventBridge Policy

The module will create a policy on a given EventBridge Event Bus so that Coralogix can send data to it. 

## Fields

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CustomCoraloigxArn | In case you ship logs from custom coralogix account specify the aws id of this account. | n\a | |
| ByPassRegion | Use only with approval from our CS team, use to bypass region restrictions | false | |
| LogsBucketName | The name of the S3 bucket to create for the logs and traces archive (Leave empty if not needed), Note: bucket name must follow [AWS naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html) | n\a | |
| MetricsBucketName | The name of the S3 bucket to create for the metrics archive (Leave empty if not needed), Note: bucket name must follow [AWS naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html) | n\a | |
| MetricsKmsArn | The arn of your kms for the metrics bucket , Note: make sure that the kms is in the same region as your bucket | n\a | |
| LogsBucketName | The arn of your kms for the logs and traces bucket , Note: make sure that the kms is in the same region as your bucket | n\a | |
