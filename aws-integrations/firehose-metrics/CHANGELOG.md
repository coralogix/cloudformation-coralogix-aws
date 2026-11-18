# Changelog

## firehose
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.1.6 / 18 Nov 2026
* [Update] AWS GovCloud: the processor is loaded from the Coralogix distribution bucket `coralogix-public-artifacts-<region>`, which holds the FIPS build from 0.1.19 on. `LambdaProcessorVersion` now defaults to empty, which deploys 0.1.18 in commercial regions, as before, and 0.1.19 in GovCloud. In GovCloud, 0.1.18 is rejected. Commercial stacks are unchanged: updating to this template with the previous `LambdaProcessorVersion`, or with it cleared, makes no changes. Existing GovCloud stacks keep their stored `LambdaProcessorVersion` (0.1.18 if created from 0.1.5) on update, and 0.1.18 is rejected: when updating a GovCloud stack to this template, clear `LambdaProcessorVersion` or set it to 0.1.19 or later. Its processor then loads from the distribution bucket. Stacks with `EnableMetricsTagsProcessors` set to false are not affected.

### 0.1.5 / 16 Nov 2026
* [Feature] Add `LambdaProcessorVersion` (default `0.1.18`): the Lambda processor is deployed from the versioned, write-once `firehose-metrics-transformer-<version>.zip` instead of the unversioned `firehose-metrics-transformer.zip`. Update the stack with a newer version to upgrade. Updating an existing stack to this template version updates the function to 0.1.18.

### 0.1.4 / 15 Nov 2026
* [Feature] Document optional `drop_values` in `DerivedLabels` rules: data points where any final label with the `target` key matches a value (case-insensitive) are dropped before they reach Coralogix. Requires processor v0.1.17 or later; existing Lambda functions are not updated automatically. The template passes `DerivedLabels` through unchanged.

### 0.1.3 / 5 Oct 2026
* [Feature] Add `DerivedLabels` parameter, passed to the Lambda processor as `DERIVED_LABELS`, to add labels copied from existing resource tags.

### 0.1.2 / 7 Sep 2026
* [Fix] Add AWS GovCloud (`aws-us-gov` partition) compatibility by replacing hardcoded `arn:aws:` prefixes with `arn:${AWS::Partition}:` across all IAM policy ARNs.
* [Update] Bump the CloudWatch Metric Stream custom resource Lambda runtime from EOL `python3.8` to `python3.12`.

### 0.1.1 / 11 Aug 2026
* [Update] Refresh the `CustomDomain` example domains to the regional format (`eu1.coralogix.com`, `eu2.coralogix.com`, `ap2.coralogix.com`).
* [Feature] Add `US3` (`us3.coralogix.com`) as a supported `CoralogixRegion`.

### 0.1.0 / 4 Aug 2026
* [Feature] Add support for retrieving the Coralogix API key from AWS Secrets Manager, including a rule to allow an empty API key parameter.
* [Feature] Add customer-managed key (CMK/KMS) encryption support for the Firehose delivery stream and associated resources.
* [Update] Set an explicit `DeliveryStreamName` only for the non-Secrets-Manager path, and gate the `IntegrationStatusNotifier` on the new `IsNotSecretDeploy` condition so it is omitted on the Secrets Manager migration path.

### 0.0.13 / 28 Feb 2026
* [Fix] Scope `sts:AssumeRole` IAM permission to `arn:aws:iam::*:role/*` instead of `*` when `CrossAccountEnabled` is true.

### 0.0.12 / 27 Feb 2026
* [Feature] Add cross-account tag enrichment support. New `CrossAccountEnabled` and `CrossAccountRoles` parameters allow the Lambda processor to assume roles in linked accounts and enrich metrics with resource tags via AWS OAM.

### 0.0.11 / 26 Feb 2026
* [Update] Change Lambda runtime from `provided.al2` to `provided.al2023` in firehose metrics template to align with AWS deprecation timeline.

### 0.0.10 / 20 Nov 2025
* [Update] Migrate Lambda ZIP Package to the common serverless repo `coralogix-serverless-repo`

### 0.0.9 / 3 Nov 2025
* [Feature] Support for AWS CloudWatch Cross-Account Observability. Added `IncludeLinkedAccountsMetrics` parameter to enable centralized metrics monitoring from multiple AWS accounts via CloudWatch Observability Access Manager (OAM)

### 0.0.8 / 10 Jul 2025
* [UPDATE] Change coralogix domains and endpoints to new format `<coralogix_region>.coralogix.com`.

### 0.0.7 / 3 Apr 2025
* [Update] Update buffer size to 1MB to be in line with documentation

### 0.0.6 / 2 Sep 2024
### 🛑 Breaking changes 🛑
* [Update] Update ingress domain from ingress-firehose.<domain.com>/firehose to ingress.<domain.com>/aws/firehose
* [Update] Update regions to follow [coralogix.com/docs/coralogix-domain] and added AP3 domain

### 0.0.5 / 7 May 2024
* [Update] default CloudWatch_Metrics_OpenTelemetry070_WithAggregations metrics integrationType option

### 0.0.4 / 7 Dec 2023
* [Update] Remove CloudWatch_Metrics_JSON metrics integrationType option

### 0.0.3 / 14 Nov 2023
* [Feature] Split Firehose Logs and Metrics

### 0.0.2 / 6 Nov 2023
* [Update] Another text change correction

### 0.0.2 / 2 Nov 2023
* [Update] Text change correction, test pr trigger script

### 0.0.2 / 2 Nov 2023
* [Update] Changed CustomDomain to be proper domain instead of url. Added CloudWatch_Metrics_OpenTelemetry070_WithAggregations option.

### 0.0.1 / 30 Oct 2023
* [Update] Fix for bucketname to be short and globally unique by stack ID

### 0.0.1 / 25 Oct 2023
* [Update] Migrate transformation Lambda runtime

### 0.0.1 / 30 Sep 2023
* [Update] Made naming & update changes to align Firehose template with Integration team's firehose-metricss one.

### 0.0.1 / 8 Sep 2023
* [Update] Remove DynamicMetadata from metrics integration, changed name from DynamicMetadata to DynamicMetadataLogs

### 0.0.1 / 28 Aug 2023
* [Feature] Release of the Amazon Kinesis Data Firehose integration for logs and metrics
