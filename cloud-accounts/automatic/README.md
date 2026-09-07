# Coralogix Cloud Account Automatic Integration

This CloudFormation template creates a named IAM role that allows Coralogix to collect AWS account metrics and is meant to be used via Cloud Accounts feature in Coralogix.

The published template URL is:

```text
https://cgx-cloudformation-templates.s3.amazonaws.com/cloud-accounts/automatic/template.yaml
```

## Prerequisites

* AWS account permissions to create IAM roles, attach AWS managed policies, and add inline IAM policies.
* The Coralogix company ID for your account.
* `CAPABILITY_NAMED_IAM` when deploying, because the stack creates a named IAM role.

## IAM Permissions

This template creates three IAM roles:

* `coralogix-cloud-account` by default, or the custom `RoleName` value. Coralogix assumes this role.
* `<RoleName>-fh`, a service role trusted by Kinesis Data Firehose.
* `<RoleName>-ms`, a service role trusted by CloudWatch Metric Streams.

The Coralogix-assumed role keeps AWS `ReadOnlyAccess` and `AWSResourceExplorerFullAccess`. An inline policy allows only the actions needed to manage Coralogix-owned `cx-*` Firehose delivery streams, CloudWatch Metric Streams, backup S3 buckets/objects, and to pass the two service roles to their intended AWS services.

## Parameters

| Parameter | Description | Default value | Required |
|---|---|---|---|
| `CoralogixRegion` | Coralogix region for your account. Allowed values are `staging`, `EU1`, `EU2`, `AP1`, `AP2`, `AP3`, `US1`, `US2`, `gov1`, and `CustomEndpoint`. | `EU1` | :heavy_check_mark: |
| `CoralogixCompanyId` | Numeric Coralogix company ID. Used to build `<ExternalIdSecret>@<CoralogixCompanyId>` when `ExternalIdSecret` is provided. | | :heavy_check_mark: |
| `RoleName` | Name of the main IAM role to create. Static Firehose service roles use this name with `-fh` and `-ms` suffixes. | `coralogix-cloud-account` | |
| `ExternalIdSecret` | Optional external ID secret for `sts:AssumeRole`. When blank, the trust policy has no external ID condition and the `ExternalId` output is omitted. | | |
| `CustomAWSAccountId` | Custom 12-digit Coralogix AWS account ID. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |
| `CustomRoleSuffix` | Custom suffix for the trusted `coralogix-ingestion-<suffix>` role. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |

For BYOC (Bring Your Own Cloud) Coralogix environments, use `CoralogixRegion=CustomEndpoint` and provide both `CustomAWSAccountId` and `CustomRoleSuffix`.

## Outputs

| Output | Description |
|---|---|
| `CoralogixAwsMetricsRoleArn` | ARN of the IAM role created for Coralogix. |
| `ExternalId` | `<ExternalIdSecret>@<CoralogixCompanyId>`. Present only when `ExternalIdSecret` is not blank. |
