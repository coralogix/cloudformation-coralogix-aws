# Changelog

## firehose
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.0.12 / Feb 27 2026 
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
