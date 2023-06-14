# AWS EventBridge Integration to Coralogix

This template can be used to deploy an AWS EventBridge Integration to Coralogix.

## Prerequisites
* AWS account.
* Coralogix account.


## Parameters

| Parameter | Description | Default Value | Required |
|---|---|---|---|
| ApplicationName | Your Coralogix application name |  | :heavy_check_mark: |
| SubsystemName | Your Coralogix Subsystem name | | :heavy_check_mark: |
| EventbridgeStream | AWS EventBridge delivery stream name |  | :heavy_check_mark: |
| RoleName | The name of the EventBridge Role |  | :heavy_check_mark: |
| PrivateKey | Your Coralogix Private Key | |  :heavy_check_mark: |
| CoralogixRegion | The region of your Coralogix Account | _Allowed Values:_<br>- ireland<br>- stockholm<br>- india<br>- singapore<br>- us | :heavy_check_mark: |
| CustomUrl | Custom Coralogix url (Endpoint) |  |  |


## Deploy the Cloudformation template

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> --capabilities CAPABILITY_NAMED_IAM  --parameter-overrides ApplicationName=<application name> SubsystemName=<subsystem name> EventbridgeStream=<EventBridge delivery stream name> RoleName=<EventBridge Role> PrivateKey=<your-private-key> CoralogixRegion=<coralogix-region> CustomUrl=<Custom Coralogix url>
```
