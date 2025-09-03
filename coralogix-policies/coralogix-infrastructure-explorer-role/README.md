# Infrastructure explorer policy and role

The module will create a role with an inline policy to allow Coralogix integration to scrape AWS infrastructure metadata.

## Fields

| Parameter | Description | Default Value | Required |
|-----------|-------------|---------------|----------|
| CoralogixRegion | The Alias for the Coralogix region, possible options are [US1, US2, EU1, EU2, AP1, AP2, AP3, dev, staging, custom] | EU1 | :heavy_check_mark: |
| RoleName | The name of the role that template will create in your AWS account | n\a | :heavy_check_mark: |
| CustomAWSAccountId | In case you want to use a custom coralogix account, enter the aws account id that you want to use. | n\a | |
| CoralogixCompanyId | Your coralogix account company ID, will be used for security validation. | n\a | :heavy_check_mark: |
| ExternalId | "sts:ExternalId" this id is used for increased security, the value of the ExternalId will be `ExternalIdSecret@CoralogixCompanyId`. | n\a | :heavy_check_mark: |

Run the following command to deploy the integration:

```sh
aws cloudformation deploy --capabilities CAPABILITY_IAM  CAPABILITY_NAMED_IAM --template-file template.yaml --stack-name <the name of the stack that will be deploy in aws> --parameter-overrides CoralogixRegion=<coralogix account region> RoleName=<name of the role> ExternalIdSecret=<external id secret part> CoralogixCompanyId=<coralogix company id>
```
