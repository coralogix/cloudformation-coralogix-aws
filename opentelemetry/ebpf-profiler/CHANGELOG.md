# Changelog

## opentelemetry ebpf-profiler
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.1.1 / 2026-07-09
- [Update] Embedded the eBPF profiler and Supervisor configurations in the CloudFormation template instead of loading them from S3.
- [Update] Added the `InitialFallbackConfigs` CloudFormation parameter for Supervisor initial fallback config sources, written as a comma-separated list.

### 0.1.0 / 2026-03-04
- [Feature] Added initial ECS EC2 CloudFormation template for eBPF profiler with S3-based configuration.
