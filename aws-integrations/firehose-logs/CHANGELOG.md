# Changelog

## firehose
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.0.1 / 25 Oct 2023
* [Update] Migrate transformation Lambda runtime

### 0.0.1 / 28 Aug 2023
* [Feature] Release of the Amazon Kinesis Data Firehose integration for logs and metrics

### 0.0.1 / 8 Sep 2023
* [Update] Remove DynamicMetadata from metrics integration, changed name from DynamicMetadata to DynamicMetadataLogs

### 0.0.1 / 30 Sep 2023
* [Update] Made naming & update changes to align Firehose template with Integration team's firehose-metricss one.

### 0.0.1 / 30 Oct 2023
* [Update] Fix for bucketname to be short and globally unique by stack ID

### 0.0.2 / 02 Nov 2023
* [Update] Changed CustomDomain to be proper domain instead of url. Added CloudWatch_Metrics_OpenTelemetry070_WithAggregations option.

### 0.0.2 / 02 Nov 2023
* [Update] Text change correction, test pr trigger script

### 0.0.2 / 06 Nov 2023
* [Update] Another text change correction

### 0.0.3 / 14 Nov 2023
* [Feature] Split Firehose Logs and Metrics

### 0.0.4 / 16 Nov 2023
* [Update] DynamicFromFields changes and documentation

### 0.0.5 / 02 Sept 2024
### 🛑 Breaking changes 🛑
* [Update] Update ingress domain from ingress-firehose.<domain.com>/firehose to ingress.<domain.com>/aws/firehose
* [Update] Update regions to follow [coralogix.com/docs/coralogix-domain] and added AP3 domain

### 0.0.6 / 03 April 2025
* [Update] Update buffer size to 1MB to be in line with documentation