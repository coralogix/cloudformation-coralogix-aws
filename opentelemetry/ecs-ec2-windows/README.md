# Coralogix OpenTelemetry Agent for ECS EC2 (Windows). CloudFormation template.

This CloudFormation template deploys the Coralogix Distribution for OpenTelemetry (CDOT) as a **Daemon** ECS service on an **existing** AWS ECS cluster with **Windows EC2** container instances. The integration aligns with the [Terraform ecs-ec2-windows module](https://github.com/coralogix/terraform-coralogix-aws/tree/main/modules/ecs-ec2-windows) and the [ecs-ec2-windows telemetry shipper](https://github.com/coralogix/telemetry-shippers/tree/main/otel-ecs-ec2-windows): Windows-optimized OTEL config (ECS Task Metadata, no Docker API), same collector image and feature gates.

The agent runs one task per Windows EC2 instance, uses **awsvpc** network mode, mounts `C:\` and `C:\ProgramData\Amazon\ECS` for ECS metadata, and sends logs to CloudWatch via **awslogs**. The built-in config supports OTLP, Jaeger, Zipkin, StatsD, Prometheus scrape, ECS container metrics (sidecar mode), and the Coralogix exporter.

## Requirements

- Existing ECS cluster with **Windows** EC2 capacity (e.g. `WINDOWS_SERVER_2022_CORE`).
- Subnets and security groups where the Daemon service will run (private subnets recommended; outbound allowed for Coralogix and optional S3/Secrets).
- [AWS CLI](https://aws.amazon.com/cli/) if deploying via CLI.

## Comparison: ecs-ec2 (Linux) vs ecs-ec2-windows

| Aspect | ecs-ec2 (Linux) | ecs-ec2-windows |
|--------|------------------|------------------|
| **OS / cluster** | Amazon Linux 2 (EC2 ECS-optimized) | Windows Server 2022 Core (EC2 ECS-optimized) |
| **Network mode** | `host` (agent shares instance network) | `awsvpc` (agent gets its own ENI) |
| **Subnets / security groups** | Not required (host mode) | Required (`SubnetIds`, `SecurityGroupIds`) |
| **Agent task** | Privileged; host mounts (`/var/lib/docker`, `/var/run/docker.sock`) | Not privileged; mounts `C:\`, `C:\ProgramData\Amazon\ECS` |
| **Agent image** | Linux tags (e.g. `v0.5.0`) | Windows tags (e.g. `v0.5.10-windowsserver-2022`) |
| **Service discovery** | — | Optional `ServiceDiscoveryRegistryArn` so other tasks reach agent via DNS (e.g. `agent.otel.local:4317`) |
| **Logging** | `json-file` (host) | `awslogs` (CloudWatch); template can create log group |

## Parameters

| Parameter | Description | Default | Required |
|-----------|-------------|---------|:--------:|
| ClusterName | Name of the existing Windows ECS cluster | — | ✓ |
| SubnetIds | Comma-separated subnet IDs for the ECS service (awsvpc) | — | ✓ |
| SecurityGroupIds | Comma-separated security group IDs for the ECS service | — | ✓ |
| ConfigSource | Config source: `template`, `s3`, `parameter-store` | `template` | |
| CDOTImageVersion | OTEL Collector image tag (use Windows tag, e.g. v0.5.10-windowsserver-2022) | `v0.5.10-windowsserver-2022` | |
| Image | Override image repository (empty = coralogixrepo/coralogix-otel-collector) | `""` | |
| CoralogixRegion | Coralogix region [EU1\|EU2\|AP1\|AP2\|AP3\|US1\|US2\|custom] | — | ✓ |
| CustomDomain | Coralogix custom domain (e.g. private link). Required when region is custom | `""` | |
| CoralogixApiKey | Send-Your-Data API key | — | ✓ * |
| UseApiKeySecret | Use API key from Secrets Manager | `false` | |
| ApiKeySecretArn | ARN of the secret (required if UseApiKeySecret is true) | `""` | |
| TaskExecutionRoleArn | Task execution role (ECR, logs, optional Secrets/SSM). If empty, a role with ECR + CloudWatch Logs is created | `""` | |
| TaskRoleArn | Task role for runtime (e.g. S3 config). If empty and ConfigSource=s3, a minimal S3 read role is created | `""` | |
| S3ConfigBucket | S3 bucket for config (required when ConfigSource=s3) | `""` | |
| S3ConfigKey | S3 key for config (required when ConfigSource=s3) | `""` | |
| CustomConfigParameterStoreName | SSM Parameter Store name for config (required when ConfigSource=parameter-store) | `""` | |
| DefaultApplicationName | Default Coralogix application name | `otel` | |
| DefaultSubsystemName | Default Coralogix subsystem name | `ecs-ec2` | |
| Cpu | Task CPU units (1024 = 1 vCPU) | `1024` | |
| Memory | Task memory (MiB) | `2048` | |
| CloudWatchLogGroupName | CloudWatch log group name; if empty, one is created | `""` | |
| CloudwatchLogRetentionDays | Retention for the created log group | `7` | |
| ServiceDiscoveryRegistryArn | Cloud Map service ARN so other tasks can reach agent via DNS | `""` | |
| HealthCheckEnabled | Enable container health check (Windows: CMD /C exit 0) | `false` | |
| HealthCheckInterval | Health check interval (seconds) | `30` | |
| HealthCheckTimeout | Health check timeout (seconds) | `5` | |
| HealthCheckRetries | Health check retries | `3` | |
| HealthCheckStartPeriod | Health check start period (seconds) | `10` | |
| EnableHeadSampler | Enable head sampling (template config) | `true` | |
| SamplingPercentage | Sampling percentage 0–100 (template config) | `10` | |
| SamplerMode | Sampler mode: proportional, equalizing, hash_seed (template config) | `proportional` | |
| EnableSpanMetrics | Enable span metrics (template config) | `true` | |
| EnableTracesDB | Enable traces/db pipeline (template config) | `false` | |

\* CoralogixApiKey is required unless UseApiKeySecret is true.

## Configuration sources

- **template** (default): Uses the built-in Windows OTEL config in the template (domain, application name, subsystem from parameters; API key via env or secret).
- **s3**: Load config from S3 at runtime. Provide `S3ConfigBucket` and `S3ConfigKey`. The template can create a task role with S3 read, or you can supply `TaskRoleArn`.
- **parameter-store**: Load config from SSM Parameter Store. Provide `CustomConfigParameterStoreName` and a `TaskExecutionRoleArn` with Parameter Store read access.

## IAM roles

- **CloudFormation execution role**: The IAM role that runs CloudFormation (e.g. when using `aws cloudformation deploy`) must be allowed to create/update ECS resources. If you see `AccessDeniedException` for `ecs:RegisterTaskDefinition`, add ECS permissions to that role (e.g. `ecs:RegisterTaskDefinition`, `ecs:DeregisterTaskDefinition`, `ecs:DescribeTaskDefinition`, and any other ECS actions the stack uses).
- **Execution role**: Used by ECS for pulling images and writing CloudWatch logs. If `TaskExecutionRoleArn` is empty, the template creates a role with `AmazonECSTaskExecutionRolePolicy`. When using API Key Secret or Parameter Store config, you must provide an execution role with the required secrets/SSM permissions.
- **Task role**: Used by the container at runtime (e.g. reading S3 config). If `ConfigSource` is `s3` and `TaskRoleArn` is empty, the template creates a minimal role with S3 read on the config bucket.

## Outputs

| Output | Description |
|--------|-------------|
| CoralogixOtelAgentServiceId | ECS service ID of the OTEL agent Daemon |
| CoralogixOtelAgentTaskDefinitionArn | Task definition ARN of the OTEL agent |
| CloudWatchLogGroupName | CloudWatch log group name used by the agent |

## Service discovery

To have other tasks in the same VPC reach the agent via DNS (e.g. `agent.otel.local:4317`), set `ServiceDiscoveryRegistryArn` to the ARN of the AWS Cloud Map service (e.g. from the [telemetry-shippers otel-ecs-ec2-windows](https://github.com/coralogix/telemetry-shippers/tree/main/otel-ecs-ec2-windows) stack). Ensure the agent runs in subnets and security groups that allow TCP 4317 from those tasks.

## Example deployment

```bash
aws cloudformation deploy --template-file template.yaml --stack-name coralogix-otel-windows \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=my-windows-ecs-cluster \
    "SubnetIds=subnet-xxx,subnet-yyy" \
    "SecurityGroupIds=sg-xxx" \
    CDOTImageVersion=v0.5.10-windowsserver-2022 \
    CoralogixRegion=EU2 \
    CoralogixApiKey=your-send-your-data-api-key
```

With API key from Secrets Manager and custom execution role:

```bash
aws cloudformation deploy --template-file template.yaml --stack-name coralogix-otel-windows \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=my-windows-ecs-cluster \
    "SubnetIds=subnet-xxx,subnet-yyy" \
    "SecurityGroupIds=sg-xxx" \
    CDOTImageVersion=v0.5.10-windowsserver-2022 \
    CoralogixRegion=EU2 \
    UseApiKeySecret=true \
    ApiKeySecretArn=arn:aws:secretsmanager:region:account:secret:name \
    TaskExecutionRoleArn=arn:aws:iam::account:role/your-execution-role
```

## Image

Use a Windows Server image tag from [Docker Hub](https://hub.docker.com/r/coralogixrepo/coralogix-otel-collector/tags), e.g. `v0.5.10-windowsserver-2022`. Linux tags are not compatible with this template.

## Notes

- The template does **not** create the ECS cluster, launch template, or ASG. Use an existing Windows ECS cluster.
- Health check on Windows uses `CMD /C exit 0` (no `/healthcheck` binary).
- Built-in template config matches the telemetry-shippers Windows integration (awsecscontainermetricsd sidecar, resourcedetection env+ec2, no Docker/opamp on Windows).
