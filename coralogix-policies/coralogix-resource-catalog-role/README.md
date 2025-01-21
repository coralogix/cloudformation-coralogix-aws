# Resource catalog role

The module will create a role with an inline policy to allow Coralogix to describe ec2 instances.

## Fields

| Parameter | Description | Default Value | Required |
|-----------|-------------|---------------|----------|
| CoralogixRegion | The Alias for the Coralogix region, possible options are [us1, us2, eu1, eu2, ap1, ap2, ap3, custom] | n\a | :heavy_check_mark: |
| RoleName | The name of the role that will get created in your AWS account | n\a | :heavy_check_mark: |

Run the following command to deploy the integration:

```sh
aws cloudformation deploy --capabilities CAPABILITY_IAM  CAPABILITY_NAMED_IAM --template-file template.yaml --stack-name <the name of the stack that will be deploy in aws> --parameter-overrides CoralogixRegion=<coralogix account region> RoleName=<RoleName>
```
