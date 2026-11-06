# Coralogix Cloud Account Manual Integration

This CloudFormation template creates the one-time bootstrap IAM role that allows Coralogix to collect AWS account metrics in Manual mode. Terraform or OpenTofu manages the monitoring resources after onboarding.

The published template URL is:

```text
https://cgx-cloudformation-templates.s3.amazonaws.com/cloud-accounts/manual/template.yaml
```

## Prerequisites

* AWS account permissions to create a named IAM role and attach AWS managed policies.
* The Coralogix company ID for your account.
* `CAPABILITY_NAMED_IAM` when deploying, because the stack creates a named IAM role.

## IAM Permissions

The template creates only `coralogix-cloud-account` by default, or the custom `RoleName` value. Coralogix assumes this role, which has these AWS managed policies:

* `ReadOnlyAccess`
* `AWSResourceExplorerFullAccess`

The Coralogix-assumed role does not receive the automatic template's resource-management or `iam:PassRole` inline policy. It also does not create the Firehose or Metric Stream service roles.

The customer-controlled identity that runs the generated Terraform or OpenTofu configuration is separate from the Coralogix-assumed role. That execution identity must already be authorized to create, update, and delete the generated IAM, S3, Kinesis Data Firehose, and CloudWatch resources. It must also be authorized to pass the two Terraform-created service roles only to their intended services:

* `<RoleName>-fh` to `firehose.amazonaws.com`
* `<RoleName>-ms` to `streams.metrics.cloudwatch.amazonaws.com`

Terraform or OpenTofu owns those service roles and all generated Firehose monitoring resources after onboarding.

### S3 log collection

The v1 bootstrap stays read-only. Generated account-level Terraform owns the shared `<RoleName>-sh` shipper role, stock shipper Lambda, invoke permission and selected routing. The execution identity needs Lambda/SNS lifecycle and subscription permissions, S3 notification replacement, and pass-role on `<RoleName>-sh` to `lambda.amazonaws.com`. The shared role's same-account S3/KMS read envelope is broader than selected buckets; customer key policies must also grant decrypt access. Terraform state/plans and Lambda configuration can contain the account SendData key: use encrypted restricted state and protect sensitive generated files.

**Direct and provisioned-SNS delivery transfer the entire bucket-notification document to one native Terraform resource per bucket.** The generated root aggregates its selected rules, rejects unrelated observed notifications during initial handover, and clears the document on destroy. Independently managed routes are unsupported. Review the complete plan and avoid new notifications between generation and apply. Existing-SNS delivery manages only the shipper/permission/matching subscription; bucket routing, topic/key policies and upstream format/filter correctness remain customer-owned. These ownership modes cannot be mixed on a bucket. Switching ownership requires staged deletion, apply/confirmation and customer handover before creating the replacement configuration. Direct/provisioned-SNS switches retain ownership and are supported.

Apply all same-bucket changes sequentially. Best-effort owned-topic sharing inspection occurs while preparing deletion/update, not continuously in Terraform; do not repurpose topics between generation and apply. Existing customer topics and source buckets are never deleted by the S3 module. Configurations remain pending until exact-generation confirmation; confirmation does not verify ingestion. There is no DLQ, and finite AWS retries, duplicate/unordered delivery and transition gaps remain delivery limitations.

## Parameters

| Parameter | Description | Default value | Required |
|---|---|---|---|
| `CoralogixRegion` | Coralogix region for your account. Allowed values are `staging`, `EU1`, `EU2`, `AP1`, `AP2`, `AP3`, `US1`, `US2`, `gov1`, and `CustomEndpoint`. | `EU1` | :heavy_check_mark: |
| `CoralogixCompanyId` | Numeric Coralogix company ID. Used to build `<ExternalIdSecret>@<CoralogixCompanyId>` when `ExternalIdSecret` is provided. | | :heavy_check_mark: |
| `RoleName` | Name of the IAM role to create for Coralogix. | `coralogix-cloud-account` | |
| `ExternalIdSecret` | Optional external ID secret for `sts:AssumeRole`. When blank, the trust policy has no external ID condition and the `ExternalId` output is omitted. | | |
| `CustomAWSAccountId` | Custom 12-digit Coralogix AWS account ID. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |
| `CustomRoleSuffix` | Custom suffix for the trusted `coralogix-ingestion-<suffix>` role. Required only when `CoralogixRegion` is `CustomEndpoint`. | | |

For BYOC (Bring Your Own Cloud) Coralogix environments, use `CoralogixRegion=CustomEndpoint` and provide both `CustomAWSAccountId` and `CustomRoleSuffix`.

## Outputs

| Output | Description |
|---|---|
| `CoralogixAwsMetricsRoleArn` | ARN of the IAM role created for Coralogix. |
| `ExternalId` | `<ExternalIdSecret>@<CoralogixCompanyId>`. Present only when `ExternalIdSecret` is not blank. |
