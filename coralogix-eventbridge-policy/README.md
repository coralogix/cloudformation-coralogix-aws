# EventBridge Policy

The module will create a policy on a given EventBridge Event Bus so that Coralogix can send data to it.

## Fields

| Parameter              | Description                                                                                        | Default Value | Required |
|------------------------|----------------------------------------------------------------------------------------------------|---------------|----------|
| CoralogixRegion        | The Coralogix region, possible options are [Europe, Europe2, India, Singapore, US, US2, Custom]    |---------------| :heavy_check_mark: |
| CustomCoralogixAccount | In case you want to use a custom coralogix account, enter the aws account id that you want to use. | n\a           |          |
| EventBusArn            | The ARN corresponding to the Event Bus that will receive events via the PutEvents method.          | n\a           | :heavy_check_mark: |