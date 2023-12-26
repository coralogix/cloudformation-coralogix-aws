# EventBridge Policy

The module will create a policy on a given EventBridge Event Bus so that Coralogix can send data to it.

## Fields

| Parameter          | Description                                                                                   | Default Value | Required |
|--------------------|-----------------------------------------------------------------------------------------------|---------------|----------|
| CustomCoralogixArn | In case you want to use a custom coralogix arn enter the aws account id that you want to use. | n\a           |          |
| EventBusArn        | The ARN corresponding to the Event Bus that will receive events via the PutEvents method.     | n\a           |          |