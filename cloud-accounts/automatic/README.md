# Coralogix Cloud Account Automatic Integration

This CloudFormation template creates named IAM roles that allow Coralogix to collect AWS account metrics and manage CloudWatch Logs ingestion through the Coralogix Cloud Accounts feature.

The published template URL is:

```text
https://cgx-cloudformation-templates.s3.amazonaws.com/cloud-accounts/automatic/template.yaml
```

## Prerequisites

* AWS account permissions to create IAM roles, attach AWS managed policies, and add inline IAM policies.
* The Coralogix company ID for your account.
* `CAPABILITY_NAMED_IAM` when deploying, because the stack creates named IAM roles.
* An active CloudTrail trail that delivers write management events is required for automatic discovery of `STANDARD` log groups created after a log collection configuration is deployed. Existing `STANDARD` log groups are scanned when the configuration is deployed. The template does not create, inspect, or manage CloudTrail resources.

## IAM Permissions

This template creates five IAM roles:

* `coralogix-cloud-account` by default, or the custom `RoleName` value. Coralogix assumes this role.
* `<RoleName>-fh`, a service role trusted by Kinesis Data Firehose.
* `<RoleName>-ms`, a service role trusted by CloudWatch Metric Streams.
* `<RoleName>-logs`, a service role trusted by CloudWatch Logs that can write only to customer `cx-*` Firehose delivery streams.
* `<RoleName>-lm`, an execution role for generated Coralogix `cx-lc-*-lm` Lambda manager functions. It can write only to their Lambda log groups; inspect, create, and delete subscriptions on customer log groups; tag and untag those log groups; query the Resource Groups Tagging API for the manager's tag index; and pass only `<RoleName>-logs` to CloudWatch Logs. It cannot modify its own Lambda configuration or permissions.

The Coralogix-assumed role keeps AWS `ReadOnlyAccess` and `AWSResourceExplorerFullAccess`. An inline policy allows only the actions needed to manage Coralogix-owned `cx-*` Firehose delivery streams, CloudWatch Metric Streams, backup S3 buckets/objects, generated `cx-lc-*-lm` Lambda manager functions, and generated `cx-lc-*-ebr` EventBridge rules/targets. Its Lambda permissions include the configuration reads used by AWS SDK waiters and setting reserved concurrency on generated managers. `lambda:AddPermission` is restricted to generated manager functions and requests whose principal is `events.amazonaws.com`; AWS provides no IAM condition keys for constraining that API's requested action, source ARN, or statement ID. It can pass the Firehose, Metric Streams, and Lambda manager roles only to their intended AWS services. Because it creates the manager Lambda and its `cx-*` Firehose stream, its effective access includes the `<RoleName>-lm` and `<RoleName>-logs` permissions, including account-wide subscription management and configuration of the stream's HTTP endpoint. It has no CloudTrail management or direct CloudWatch Logs subscription actions.

## Parameters

| Parameter | Description | Default value | Required |
|---|---|---|---|
| `CoralogixRegion` | Coralogix region for your account. Allowed values are `staging`, `EU1`, `EU2`, `AP1`, `AP2`, `AP3`, `US1`, `US2`, `gov1`, and `CustomEndpoint`. | `EU1` | :heavy_check_mark: |
| `CoralogixCompanyId` | Numeric Coralogix company ID. Used to build `<ExternalIdSecret>@<CoralogixCompanyId>` when `ExternalIdSecret` is provided. | | :heavy_check_mark: |
| `RoleName` | Name of the main IAM role to create. Static service roles use this name with `-fh`, `-ms`, `-logs`, and `-lm` suffixes. | `coralogix-cloud-account` | |
| `ExternalIdSecret` | Optional external ID secret for `sts:AssumeRole`. When blank, the trust policy has no external ID condition and the `ExternalId` output is omitted. | | |
| `CustomAWSAccountId` | Custom 12-digit Coralogix AWS account ID. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |
| `CustomRoleSuffix` | Custom suffix for the trusted `coralogix-ingestion-<suffix>` role. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |

For BYOC (Bring Your Own Cloud) Coralogix environments, use `CoralogixRegion=CustomEndpoint` and provide both `CustomAWSAccountId` and `CustomRoleSuffix`.

## Outputs

| Output | Description |
|---|---|
| `CoralogixAwsMetricsRoleArn` | ARN of the IAM role created for Coralogix. |
| `ExternalId` | `<ExternalIdSecret>@<CoralogixCompanyId>`. Present only when `ExternalIdSecret` is not blank. |
