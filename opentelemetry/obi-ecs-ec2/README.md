# Coralogix OBI (eBPF Instrumentation) — ECS EC2

Deploy [OBI](https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation)
(OpenTelemetry eBPF Instrumentation) as an ECS **DAEMON** service. OBI captures
application spans & metrics via eBPF (zero-code) for **every process on the host** —
including services with no OpenTelemetry SDK — and exports OTLP to the node-local
Coralogix OpenTelemetry Collector, which forwards to Coralogix.

This mirrors the Kubernetes OBI integration (the `opentelemetry-ebpf-instrumentation`
Helm chart) and the `telemetry-shippers` `otel-ecs-ec2` integration, packaged as a
CloudFormation template alongside the existing `opentelemetry/ecs-ec2` collector and
`opentelemetry/ebpf-profiler`.

## Requirements

- **ECS EC2 launch type** — OBI **cannot run on Fargate** (it needs `privileged`
  mode, host PID namespace, and eBPF host mounts). One task runs per container
  instance (`SchedulingStrategy: DAEMON`).
- **Linux kernel >= 5.8** on the container instances.
- **A collector on the cluster.** OBI exports to the node-local collector, not
  directly to Coralogix. Deploy [`opentelemetry/ecs-ec2`](../ecs-ec2/) first (it
  runs as a host-network DAEMON listening on `4317`/`4318`), or point
  `OtlpTracesEndpoint` / `OtlpMetricsEndpoint` at your own collector.

## How it works

- A privileged `coralogix-obi` container runs the OBI image with host PID/network
  and the eBPF host mounts (`/sys/fs/cgroup`, `/sys/kernel/debug`,
  `/sys/kernel/tracing`, `/sys/fs/bpf`), reading its config from
  `OTEL_EBPF_CONFIG_PATH`.
- A `config-loader` sidecar writes the OBI config to a shared volume — either the
  template's embedded default or a file pulled from S3 (`S3ConfigBucket` /
  `S3ConfigKey`).
- The default config uses host-process discovery (`open_ports` / `exe_path`),
  GenAI + DB payload extraction, and exports to `localhost:4317` (traces) /
  `localhost:4318` (metrics).

## Parameters

| Parameter | Default | Description |
| --- | --- | --- |
| `ClusterName` | — (required) | Existing ECS cluster (EC2 launch type). |
| `ObiImageRepository` | `ghcr.io/open-telemetry/opentelemetry-ebpf-instrumentation/ebpf-instrument` | OBI image repository. |
| `ObiImageVersion` | `v0.10.0` | OBI image tag. |
| `ContextPropagationMode` | `headers,tcp` | eBPF distributed context propagation; `disabled` to turn off. |
| `OtlpTracesEndpoint` | `http://localhost:4317` | OTLP gRPC endpoint for traces (the node-local collector). |
| `OtlpMetricsEndpoint` | `http://localhost:4318` | OTLP HTTP endpoint for metrics (the node-local collector). |
| `S3ConfigBucket` | `""` | Optional S3 bucket to load a custom OBI config from. |
| `S3ConfigKey` | `""` | Optional S3 key for the OBI config. |

## Deploy

Using the AWS CLI:

```bash
aws cloudformation deploy \
  --template-file template.yaml \
  --stack-name coralogix-obi \
  --region <region> \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides ClusterName=<your-ecs-cluster>
```

Or with the helper Makefile (auto-loads `.env`; see `.env.example`):

```bash
cp .env.example .env   # set CLUSTER_NAME
make deploy            # deploy against an existing cluster + collector
make status            # show the ECS service
make smoke             # full flow: create cluster + EC2 instance + deploy + checks
make cleanup           # tear everything down
```

## Scoping what gets instrumented

Edit `examples/obi-config.yaml` (or the embedded default in `template.yaml`) —
narrow `discovery.instrument` to specific ports or `exe_path` globs — then upload it
to S3 and set `S3ConfigBucket` / `S3ConfigKey`.
