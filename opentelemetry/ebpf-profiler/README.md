# Coralogix OpenTelemetry eBPF Profiler for ECS-EC2. CloudFormation template.

This CloudFormation template deploys an ECS daemon service for the OpenTelemetry eBPF profiler collector on an existing ECS EC2 cluster.

In supervised mode, the deployment includes default collector and Supervisor configurations in the template. Collector mode loads its configuration from S3. The deployment runs one profiler task per EC2 container instance.

## Image

This template supports two image modes:

- `collector` mode: `<EbpfProfilerImageRepository>:<EbpfProfilerImageVersion>`
- `supervised` mode: `<SupervisedImageRepository>:<SupervisedImageVersion>`

## Requirements

- AWS credentials with permissions for CloudFormation, ECS, EC2, IAM, and SSM, plus S3 when loading configurations from S3
- A Coralogix Send-Your-Data API key
- ECS container instances with Linux kernel `>= 5.2` (default Makefile AMI uses Amazon Linux 2023 ECS-optimized AMI)

## Parameters

| Parameter                     | Description                                                                | Default                                                               | Required |
|-------------------------------|----------------------------------------------------------------------------|-----------------------------------------------------------------------|----------|
| `ClusterName`                 | Name of the existing ECS cluster                                           |                                                                       | yes      |
| `CoralogixRegion`             | Coralogix region (`EU1`,`EU2`,`AP1`,`AP2`,`AP3`,`US1`,`US2`)               |                                                                       | yes      |
| `CoralogixApiKey`             | Coralogix Send-Your-Data API key                                           |                                                                       | yes      |
| `ProfilerImageMode`           | Image mode (`collector` or `supervised`)                                   | `collector`                                                           | no       |
| `EbpfProfilerImageRepository` | Collector image repository                                                 | `coralogixrepo/coralogix-otel-collector`                              | no       |
| `EbpfProfilerImageVersion`    | Image version/tag for collector mode                                       | `v0.5.8`                                                              | no       |
| `SupervisedImageRepository`   | Supervised image repository                                                | `cgx.jfrog.io/coralogix-docker-images/coralogix-otel-supervised-cdot` | no       |
| `SupervisedImageVersion`      | Supervised image version/tag                                               | `v0.10.0`                                                             | no       |
| `S3ConfigBucket`              | S3 bucket containing collector and optional Supervisor configurations      |                                                                       | collector mode |
| `S3ConfigKey`                 | S3 object key for the collector configuration                              |                                                                       | collector mode |
| `S3SupervisorConfigKey`       | S3 object key for the Supervisor configuration                             |                                                                       | no       |

## Deploy

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
  --region <aws_region> \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=<ecs_cluster_name> \
    CoralogixRegion=<coralogix_region> \
    CoralogixApiKey=<send_your_data_api_key> \
    ProfilerImageMode=collector \
    S3ConfigBucket=<s3_bucket_name> \
    S3ConfigKey=<path/to/collector-config.yaml> \
    EbpfProfilerImageRepository=coralogixrepo/coralogix-otel-collector \
    EbpfProfilerImageVersion=v0.5.8
```

For supervised mode:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
  --region <aws_region> \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=<ecs_cluster_name> \
    CoralogixRegion=<coralogix_region> \
    CoralogixApiKey=<send_your_data_api_key> \
    ProfilerImageMode=supervised \
    SupervisedImageRepository=cgx.jfrog.io/coralogix-docker-images/coralogix-otel-supervised-cdot \
    SupervisedImageVersion=v0.10.0
```

## S3 configurations

Collector mode requires `S3ConfigBucket` and `S3ConfigKey`. Supervised mode uses embedded configurations by default; set the bucket and relevant key to override either configuration from S3.

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
  --region <aws_region> \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=<ecs_cluster_name> \
    CoralogixRegion=<coralogix_region> \
    CoralogixApiKey=<send_your_data_api_key> \
    ProfilerImageMode=supervised \
    S3ConfigBucket=<s3_bucket_name> \
    S3ConfigKey=<path/to/collector-config.yaml> \
    S3SupervisorConfigKey=<path/to/supervisor-config.yaml>
```

When `S3SupervisorConfigKey` is set, the S3 file is used as-is.

## Reproducible smoke test

Example configurations are provided at:

- `examples/ebpf-profiler-config.yaml` (collector)
- `examples/ebpf-profiler-supervisor-config.yaml` (supervised)

Use the local `Makefile` to run the full flow (validate, create infra prerequisites, deploy, and basic runtime checks):

```sh
CLUSTER_NAME=<ecs_cluster_name> \
CORALOGIX_REGION=<coralogix_region> \
CORALOGIX_API_KEY=<send_your_data_api_key> \
AWS_REGION=<aws_region> \
S3_BUCKET=<s3_bucket_name> \
S3_CONFIG_KEY=<path/to/collector-config.yaml> \
make smoke
```

To run supervised mode:

```sh
make smoke-supervised
```

`deploy-supervised`, `smoke-supervised`, and `cleanup-supervised` reuse the corresponding main targets with `PROFILER_IMAGE_MODE=supervised` and all S3 configuration variables cleared. This allows the normal targets to keep S3 values configured in `.env` while the supervised targets always exercise embedded configurations.

To load configurations from S3 in the Makefile flow:

```sh
S3_BUCKET=<s3_bucket_name> \
S3_CONFIG_KEY=<path/to/collector-config.yaml> \
S3_SUPERVISOR_CONFIG_KEY=<path/to/supervisor-config.yaml> \
PROFILER_IMAGE_MODE=supervised \
make create-bucket upload-config deploy
```

`smoke` creates the configured S3 bucket and uploads each configured file before deployment. `cleanup` removes the configured bucket after deleting the stack and ECS resources. The S3 steps do nothing when `S3_BUCKET` is empty, so supervised mode can use embedded configurations without separate commands.

`make smoke` now also creates the ECS cluster and ECS EC2 container instance if they do not already exist.

Or use an env file for shorter commands:

```sh
cp .env.example .env
# edit .env and set real values
make smoke
```

Useful individual targets:

- `make validate`
- `make create-cluster`
- `make create-ec2-instance`
- `make wait-cluster-capacity`
- `make create-bucket`
- `make upload-config`
- `make delete-ec2-instance`
- `make deploy`
- `make deploy-supervised`
- `make smoke`
- `make smoke-supervised`
- `make status`
- `make tasks`
- `make delete`
- `make delete-cluster`
- `make delete-bucket`
- `make cleanup`
- `make cleanup-supervised`

## Notes

- The ECS service name is fixed to `coralogix-ebpf-profiler` and task family is fixed to `ebpf-profiler` for a minimal working setup.
- The task uses privileged mode with host network and host PID to support eBPF profiling.
- The template creates an ECS task execution role. It creates and attaches an S3 task role only when an S3 configuration is selected.
- `make delete-ec2-instance` only terminates EC2 instances created by this Makefile (`ManagedBy=ebpf-profiler-makefile`).
- If you previously created ECS instances with an older AMI/kernel, run `make delete-ec2-instance` and then `make create-ec2-instance` before retrying deploy.
