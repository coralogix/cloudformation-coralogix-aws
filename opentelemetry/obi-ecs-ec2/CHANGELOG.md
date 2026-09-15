# Changelog

## opentelemetry obi-ecs-ec2
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.1.0 / 2026-08-16
* [Feature] Initial OpenTelemetry eBPF Instrumentation (OBI) ECS EC2 daemon template. Runs OBI (`ghcr.io/open-telemetry/opentelemetry-ebpf-instrumentation/ebpf-instrument:v0.10.0`) as a privileged, host-PID, host-network `DAEMON` service that instruments every process on the host via eBPF and exports OTLP to the node-local Coralogix collector. Config is embedded by default or loaded from S3.
