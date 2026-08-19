# Changelog

## [3.0.0] - 2026-08-11

### Added

- Add `<RoleName>-logs` for CloudWatch Logs delivery to Coralogix-owned Firehose streams.
- Add `<RoleName>-lm` for least-privilege Lambda manager execution.

### Changed

- Allow the Coralogix-assumed role to manage only generated `cx-lc-*-lm` Lambda managers and `cx-lc-*-lg` EventBridge rules/targets, including the Lambda configuration reads required by AWS SDK waiters, and to pass only the Lambda manager execution role to Lambda.
- Scope Lambda manager execution to generated manager log groups and configuration, and use log-group ARNs for customer subscription management.
- Document the CloudTrail prerequisite for discovering newly created log groups.

## [2.0.0] - 2026-07-03

### Changed

- Remove broad Firehose, S3, CloudWatch, and IAM managed policies from the Coralogix-assumed role.
- Add `<RoleName>-fh` and `<RoleName>-ms` service roles for Firehose-backed monitoring resources.

## [1.0.0] - 2026-06-18

### Added

- Add a CloudFormation template for creating the automatic cloud account IAM role with Coralogix region mapping, custom endpoint parameters, optional external ID support, and the default role name `coralogix-cloud-account`.
