# EventBridge Policy

The module will create a role with an inline policy to allow Coralogix to send events to an EventBridge event bus.

## Fields

| Parameter              | Description                                                                                                                                        | Default Value | Required           |
|------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------|---------------|--------------------|
| CoralogixRegionAlias   | The Alias for the Coralogix region, possible options are [us1, us2, eu1, eu2, ap1, ap2, ap3,  custom]                                                    | n\a           | :heavy_check_mark: |
| Role Name              | Don't change it! It needs to match the one that was input on the Coralogix form. Corresponds to the name of the AWS IAM role that will be created. | n\a           | :heavy_check_mark: |
| CustomCoralogixAccount | In case you want to use a custom coralogix account, enter the aws account id that you want to use.                                                 | n\a           |                    |
| CustomCoralogixRole    | In case you want to use a custom coralogix role, enter the role name that you want to use.                                                         | n\a           |                    |
| EventBusArn            | The ARN corresponding to the Event Bus that will receive events via the PutEvents method.                                                          | n\a           | :heavy_check_mark: |

Run the following command to deploy the integration:

```sh
aws cloudformation deploy --capabilities CAPABILITY_IAM  CAPABILITY_NAMED_IAM --template-file template.yaml --stack-name <the name of the stack that will be deploy in aws> --parameter-overrides CoralogixRegionAlias=<coralogix account region> EventBusArn=<EventBusArn>
```
