# Changelog

## [2.0.0] - 2026-07-03

### Changed

- Remove broad Firehose, S3, CloudWatch, and IAM managed policies from the Coralogix-assumed role.
- Add `<RoleName>-fh` and `<RoleName>-ms` service roles for Firehose-backed monitoring resources.

## [1.0.0] - 2026-06-18

### Added

- Add a CloudFormation template for creating the automatic cloud account IAM role with Coralogix region mapping, custom endpoint parameters, optional external ID support, and the default role name `coralogix-cloud-account`.
