# Changelog

## [5.0.0] - 2026-11-20

### Added

Support for delivering metrics over AWS PrivateLink.

- `<RoleName>-plsh` execution role for PrivateLink metrics shippers, with Lambda log and VPC network interface permissions.
- Tag-scoped VPC, subnet, security group, and interface VPC endpoint management, including cross-region endpoints, for resources tagged `coralogix:cloud-account-owner`.
- `route53:AssociateVPCWithHostedZone`, needed to enable private DNS on the interface VPC endpoint.
- PrivateLink shipper Lambda management and `<RoleName>-plsh` pass permissions; `<RoleName>-fh` can invoke PrivateLink shippers.
- Interface VPC endpoints can be created only to services owned by Amazon or by the Coralogix account of the selected region.
- Shipper Lambda network interfaces can be created only in subnets and security groups tagged `coralogix:cloud-account-owner`, and never by the function code.
- `logs:DeleteLogGroup` on the log groups of Coralogix Lambdas, so they are removed with the function.
- S3 log shippers can run in the PrivateLink network: `<RoleName>-shipper` manages the Lambda VPC network interfaces, and tag-scoped route table management lets the network reach S3 through a gateway VPC endpoint.

## [4.0.0] - 2026-11-20

### Added

Support for ingesting logs from S3 buckets.

- Shared `<RoleName>-shipper` stock S3 log shipper execution role, with same-account object reads and S3-scoped KMS decrypt. Customer key policies must also grant access.
- Scoped S3 notification, shipper Lambda, owned SNS topic, same-account Lambda subscription, and shared-role pass permissions.

### Changed

- Allow the partition-appropriate CloudWatch Logs principal alongside S3/SNS for same-account shipper AddPermission management, preparing a future additional destination without implementing its transport.

## [3.1.0] - 2026-11-05

### Added

- Add `US3` to `CoralogixRegion`, trusting `coralogix-ingestion-us3` in AWS account `826308361216`.

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
