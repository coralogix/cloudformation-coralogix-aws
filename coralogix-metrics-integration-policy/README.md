# Aws Metrics integration Role

The module will create a role to be used with AwsMetrics integration

## Fields

| Parameter | Description | Default Value | Required |
|-----------|-------------|---------------|----------|
| AWSAccount | The Alias for the Coralogix region, possible options are [US1, US2, EU1, EU2, AP1, AP2, dev, staging, custom] | EU1 | :heavy_check_mark: |
| RoleName | The name of the rule that template will create in your AWS account | n\a | :heavy_check_mark: |
| CustomeAccountId | In case you want to use a custom coralogix account, enter the aws account id that you want to use.| n\a  | |
