# Changelog

## opentelemetry ecs-ec2-windows
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->
### 2.0.0 / 2026-03-18
* [Update] Aligned template with Terraform module ecs-ec2-windows: Daemon ECS service only (removed sidecar app example)
* [Update] Switched to awsvpc network mode with required SubnetIds and SecurityGroupIds; Windows Server 2022 Core, X86_64
* [Update] Added config sources: template (default), s3, parameter-store; optional Service Discovery registry ARN
* [Update] Task definition: volumes C:\ and C:\ProgramData\Amazon\ECS; container command uses --feature-gates=service.profilesSupport
* [Update] Optional CloudWatch log group creation; optional task execution role and task role (S3 read when config from S3)
* [Update] Coralogix regions EU1, EU2, AP1, AP2, AP3, US1, US2, custom (CustomDomain); API key from parameter or Secrets Manager
* [Update] Health check uses Windows command CMD /C exit 0; added sampling and span-metrics parameters for template config

### 1.1.0 / 2026-01-25
* [Security] Added TaskRoleArn parameter to separate execution and task IAM roles, following principle of least privilege
* [Update] Standardized to use CORALOGIX_PRIVATE_KEY environment variable instead of PRIVATE_KEY

### 1.0.3 / 2025-05-20
- Added healthcheck to ECS task.

### 0.0.3 / 2025-3-31
* Add command  `--config env:OTEL_CONFIG` to the windows otel collector container example

### 0.0.2 / 2025-3-15
* Added Collector metric collection to ECS-EC2-Windows example
* Updated doc to remove deprecated S3 config and Base64 env var features

### 0.0.1 / 2023-09-11
* Added EC2 ECS Windows Example for metrics collection