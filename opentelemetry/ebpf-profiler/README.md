# Coralogix OpenTelemetry eBPF Profiler for ECS-EC2. CloudFormation template.

This CloudFormation template deploys an ECS daemon service for the OpenTelemetry eBPF profiler collector on an existing ECS EC2 cluster.

The deployment includes the collector and supervisor configurations in the template and runs one profiler task per EC2 container instance.

## Image

This template supports two image modes:

- `collector` mode: `<EbpfProfilerImageRepository>:<EbpfProfilerImageVersion>`
- `supervised` mode: `<SupervisedImageRepository>:<SupervisedImageVersion>`

## Requirements

- AWS credentials with permissions for CloudFormation, ECS, EC2, IAM, and SSM
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
| `SupervisedImageVersion`      | Supervised image version/tag                                               | `v0.0.1`                                                              | no       |
| `InitialFallbackConfigs`      | Comma-delimited S3 URIs for `agent.initial_fallback_configs` in supervised mode |                                                               | no       |

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
    SupervisedImageVersion=v0.0.1
```

To enable initial fallback configs in supervised mode, pass a comma-delimited list of S3 URIs:

```sh
aws cloudformation deploy --template-file template.yaml --stack-name <stack_name> \
  --region <aws_region> \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=<ecs_cluster_name> \
    CoralogixRegion=<coralogix_region> \
    CoralogixApiKey=<send_your_data_api_key> \
    ProfilerImageMode=supervised \
    InitialFallbackConfigs=s3://<BUCKET>.s3.<REGION>.amazonaws.com/<ACCOUNT_ID>/<GROUP_NAME>/<COLLECTOR_VERSION>/<REMOTE_CONFIG_NAME>/config.yaml,s3://<BUCKET>.s3.<REGION>.amazonaws.com/<ACCOUNT_ID>/<GROUP_NAME>/EMPTY_VERSION/<REMOTE_CONFIG_NAME>/config.yaml
```

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
make smoke
```

To run supervised mode:

```sh
PROFILER_IMAGE_MODE=supervised make smoke
```

To enable initial fallback configs in the Makefile flow, pass a comma-delimited list of S3 URIs:

```sh
INITIAL_FALLBACK_CONFIGS=s3://<BUCKET>.s3.<REGION>.amazonaws.com/<ACCOUNT_ID>/<GROUP_NAME>/<COLLECTOR_VERSION>/<REMOTE_CONFIG_NAME>/config.yaml,s3://<BUCKET>.s3.<REGION>.amazonaws.com/<ACCOUNT_ID>/<GROUP_NAME>/EMPTY_VERSION/<REMOTE_CONFIG_NAME>/config.yaml make deploy
```

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
- `make delete-ec2-instance`
- `make deploy`
- `make status`
- `make tasks`
- `make delete`
- `make delete-cluster`
- `make cleanup`

## Notes

- The ECS service name is fixed to `coralogix-ebpf-profiler` and task family is fixed to `ebpf-profiler` for a minimal working setup.
- The task uses privileged mode with host network and host PID to support eBPF profiling.
- The template creates an ECS task execution role. No S3 bucket or task role is required for configuration loading.
- `make delete-ec2-instance` only terminates EC2 instances created by this Makefile (`ManagedBy=ebpf-profiler-makefile`).
- If you previously created ECS instances with an older AMI/kernel, run `make delete-ec2-instance` and then `make create-ec2-instance` before retrying deploy.
