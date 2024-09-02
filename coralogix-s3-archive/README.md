# s3-archive

The module s3-archive will create s3 buckets to archive your coralogix logs, traces and metrics

The module can run only on the following regions eu-west-1,eu-north-1,ap-southeast-1,ap-south-1,ap-southeast-3,us-east-2,us-west-2.

## Fields

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| CustomCoraloigxArn | In case you ship logs from custom coralogix account specify the aws id of this account. | n\a | |
| ByPassRegion | Use only with approval from our CS team, use to bypass region restrictions | false | |
| LogsBucketName | The name of the S3 bucket to create for the logs and traces archive (Leave empty if not needed), Note: bucket name must follow [AWS naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html) | n\a | |
| MetricsBucketName | The name of the S3 bucket to create for the metrics archive (Leave empty if not needed), Note: bucket name must follow [AWS naming rules](https://docs.aws.amazon.com/AmazonS3/latest/userguide/bucketnamingrules.html) | n\a | |
| MetricsKmsArn | The arn of your kms for the metrics bucket , Note: make sure that the kms is in the same region as your bucket | n\a | |
| LogsBucketName | The arn of your kms for the logs and traces bucket , Note: make sure that the kms is in the same region as your bucket | n\a | |