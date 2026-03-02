# Coralogix OpenTelemetry eBPF Profiler for ECS-EC2. CloudFormation template.

This CloudFormation template deploys an ECS daemon service for the OpenTelemetry eBPF profiler collector on an existing ECS EC2 cluster.

The deployment uses an S3-hosted collector configuration file and runs one profiler task per EC2 container instance.

## Image

This template supports two image modes:

- `collector` mode: `<EbpfProfilerImageRepository>:<EbpfProfilerImageVersion>`
- `supervised` mode: `<SupervisedImageRepository>:<SupervisedImageVersion>`

## Requirements

- AWS credentials with permissions for CloudFormation, ECS, EC2, IAM, S3, and SSM
- A Coralogix Send-Your-Data API key
- ECS container instances with Linux kernel `>= 5.2` (default Makefile AMI uses Amazon Linux 2023 ECS-optimized AMI)

## Parameters

| Parameter                     | Description                                                  | Default                                                               | Required |
|-------------------------------|--------------------------------------------------------------|-----------------------------------------------------------------------|----------|
| `ClusterName`                 | Name of the existing ECS cluster                             |                                                                       | yes      |
| `CoralogixRegion`             | Coralogix region (`EU1`,`EU2`,`AP1`,`AP2`,`AP3`,`US1`,`US2`) |                                                                       | yes      |
| `CoralogixApiKey`             | Coralogix Send-Your-Data API key                             |                                                                       | yes      |
| `S3ConfigBucket`              | S3 bucket containing collector configuration                 |                                                                       | yes      |
| `S3ConfigKey`                 | S3 object key for collector configuration                    |                                                                       | yes      |
| `S3SupervisorConfigKey`       | S3 object key for supervisor configuration                   | `configs/ebpf-profiler-supervisor-config.yaml`                        | no       |
| `ProfilerImageMode`           | Image mode (`collector` or `supervised`)                     | `collector`                                                           | no       |
| `EbpfProfilerImageRepository` | Collector image repository                                   | `coralogixrepo/coralogix-otel-collector`                              | no       |
| `EbpfProfilerImageVersion`    | Image version/tag for collector mode                         | `v0.5.8`                                                              | no       |
| `SupervisedImageRepository`   | Supervised image repository                                  | `cgx.jfrog.io/coralogix-docker-images/coralogix-otel-supervised-cdot` | no       |
| `SupervisedImageVersion`      | Supervised image version/tag                                 | `v0.0.1`                                                              | no       |

## Deploy

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
  --region <aws_region> \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=<ecs_cluster_name> \
    CoralogixRegion=<coralogix_region> \
    CoralogixApiKey=<send_your_data_api_key> \
    S3ConfigBucket=<s3_bucket_name> \
    S3ConfigKey=<path/to/collector-config.yaml> \
    S3SupervisorConfigKey=<path/to/supervisor-config.yaml> \
    ProfilerImageMode=collector \
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
    S3ConfigBucket=<s3_bucket_name> \
    S3ConfigKey=<path/to/collector-config.yaml> \
    S3SupervisorConfigKey=<path/to/supervisor-config.yaml> \
    ProfilerImageMode=supervised \
    SupervisedImageRepository=cgx.jfrog.io/coralogix-docker-images/coralogix-otel-supervised-cdot \
    SupervisedImageVersion=v0.0.1
```

## Reproducible smoke test

Example configurations are provided at:

- `examples/ebpf-profiler-config.yaml` (collector)
- `examples/ebpf-profiler-supervisor-config.yaml` (supervised)

Use the local `Makefile` to run the full flow (validate, create infra prerequisites, upload config, deploy, and basic runtime checks):

```sh
CLUSTER_NAME=<ecs_cluster_name> \
CORALOGIX_REGION=<coralogix_region> \
CORALOGIX_API_KEY=<send_your_data_api_key> \
S3_BUCKET=<s3_bucket_name> \
AWS_REGION=<aws_region> \
make smoke
```

To run supervised mode:

```sh
PROFILER_IMAGE_MODE=supervised make smoke
```

`make smoke` now also creates the ECS cluster, ECS EC2 container instance, and S3 bucket if they do not already exist.

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
- `make delete-ec2-instance`
- `make create-bucket`
- `make upload-config`
- `make deploy`
- `make status`
- `make tasks`
- `make delete`
- `make delete-cluster`
- `make delete-bucket`
- `make cleanup`

## Notes

- The ECS service name is fixed to `coralogix-ebpf-profiler` and task family is fixed to `ebpf-profiler` for a minimal working setup.
- The task uses privileged mode with host network and host PID to support eBPF profiling.
- IAM is split between execution role and task role (both created by the template).
- `make delete-ec2-instance` only terminates EC2 instances created by this Makefile (`ManagedBy=ebpf-profiler-makefile`).
- If you previously created ECS instances with an older AMI/kernel, run `make delete-ec2-instance` and then `make create-ec2-instance` before retrying deploy.
