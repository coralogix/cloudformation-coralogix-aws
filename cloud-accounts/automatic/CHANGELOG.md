# Changelog

## [4.0.0] - 2026-10-29

### Added

Support for ingesting logs from S3 buckets.

- Shared `<RoleName>-sh` stock S3 log shipper execution role, with same-account object reads and S3-scoped KMS decrypt. Customer key policies must also grant access.
- Scoped S3 notification, shipper Lambda, owned SNS topic, same-account Lambda subscription, and shared-role pass permissions.

## [3.0.0] - 2026-09-08

### Added

Support for ingesting logs from CloudWatch log groups.

- Add `<RoleName>-logs` for CloudWatch Logs delivery to Coralogix-owned Firehose streams.
- Add `<RoleName>-lm` for least-privilege Lambda manager execution.
- Extend the Coralogix-assumed role to manage generated `cx-lc-*-lm` Lambda manager functions and `cx-lc-*-ebr` EventBridge rules, and to pass `<RoleName>-lm` to Lambda.

## [2.0.0] - 2026-07-03

### Changed

- Remove broad Firehose, S3, CloudWatch, and IAM managed policies from the Coralogix-assumed role.
- Add `<RoleName>-fh` and `<RoleName>-ms` service roles for Firehose-backed monitoring resources.

## [1.0.0] - 2026-06-18

### Added

- Add a CloudFormation template for creating the automatic cloud account IAM role with Coralogix region mapping, custom endpoint parameters, optional external ID support, and the default role name `coralogix-cloud-account`.
