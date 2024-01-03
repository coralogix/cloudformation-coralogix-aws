# EventBridge Policy

The module will create a role with an inline policy to allow Coralogix to send events to an EventBridge event bus.

## Fields

| Parameter              | Description                                                                                        | Default Value | Required           |
|------------------------|----------------------------------------------------------------------------------------------------|---------------|--------------------|
| CoralogixRegionAlias   | The Alias for the Coralogix region, possible options are [us1, us2, eu1, eu2, ap1, ap2, custom]    | n\a           | :heavy_check_mark: |
| CustomCoralogixAccount | In case you want to use a custom coralogix account, enter the aws account id that you want to use. | n\a           |                    |
| CustomCoralogixRole    | In case you want to use a custom coralogix role, enter the role name that you want to use.         | n\a           |                    |
| EventBusArn            | The ARN corresponding to the Event Bus that will receive events via the PutEvents method.          | n\a           | :heavy_check_mark: |
