# Changelog

## [1.1.0] - 2026-11-05

### Added

- Add `US3` to `CoralogixRegion`, trusting `coralogix-ingestion-us3` in AWS account `826308361216`.

## [1.0.0] - 2026-07-15

### Added

- Add the one-time Manual cloud account bootstrap role with AWS `ReadOnlyAccess`, `AWSResourceExplorerFullAccess`, Coralogix region mapping, custom endpoint parameters, and optional external ID support.
- Tag the role with `CoralogixCloudAccountTemplateVersion` value `1` for Manual onboarding.

Manual and Automatic template versions are independent release lines. Both templates use the `CoralogixCloudAccountTemplateVersion` IAM tag, and integrations-service interprets its integer value together with the account permission mode.
