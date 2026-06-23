# Coralogix Cloud Account Automatic Integration

This CloudFormation template creates a named IAM role that allows Coralogix to collect AWS account metrics and is meant to be used via Cloud Accounts feature in Coralogix.

The published template URL is:

```text
https://cgx-cloudformation-templates.s3.amazonaws.com/cloud-accounts/automatic/template.yaml
```

## Prerequisites

* AWS account permissions to create IAM roles and attach AWS managed policies.
* The Coralogix company ID for your account.
* `CAPABILITY_NAMED_IAM` when deploying, because the stack creates a named IAM role.

## IAM Permissions

This template intentionally attaches broad AWS managed policies, including IAM, S3, Kinesis Data Firehose, CloudWatch, Resource Explorer, and read-only account access. The automatic Cloud Accounts role is used by Coralogix to create supporting IAM roles and AWS resources required for Firehose metric streaming, not only to read existing metrics. A future template will provide a narrower permission model for deployments that do not require Coralogix to create those supporting resources.

## Parameters

| Parameter | Description | Default value | Required |
|---|---|---|---|
| `CoralogixRegion` | Coralogix region for your account. Allowed values are `staging`, `EU1`, `EU2`, `AP1`, `AP2`, `AP3`, `US1`, `US2`, `gov1`, and `CustomEndpoint`. | `EU1` | :heavy_check_mark: |
| `CoralogixCompanyId` | Numeric Coralogix company ID. Used to build `<ExternalIdSecret>@<CoralogixCompanyId>` when `ExternalIdSecret` is provided. | | :heavy_check_mark: |
| `RoleName` | Name of the IAM role to create. | `coralogix-cloud-account` | |
| `ExternalIdSecret` | Optional external ID secret for `sts:AssumeRole`. When blank, the trust policy has no external ID condition and the `ExternalId` output is omitted. | | |
| `CustomAWSAccountId` | Custom 12-digit Coralogix AWS account ID. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |
| `CustomRoleSuffix` | Custom suffix for the trusted `coralogix-ingestion-<suffix>` role. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |

For BYOC (Bring Your Own Cloud) Coralogix environments, use `CoralogixRegion=CustomEndpoint` and provide both `CustomAWSAccountId` and `CustomRoleSuffix`.

## Outputs

| Output | Description |
|---|---|
| `CoralogixAwsMetricsRoleArn` | ARN of the IAM role created for Coralogix. |
| `ExternalId` | `<ExternalIdSecret>@<CoralogixCompanyId>`. Present only when `ExternalIdSecret` is not blank. |
