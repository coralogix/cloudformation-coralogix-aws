# Coralogix OpenTelemetry Tail Sampling for ECS-EC2 - Example Deployment

This CloudFormation template provides an example of how to deploy a complete OpenTelemetry tail sampling solution on Amazon ECS with EC2 launch type. The solution enables intelligent trace sampling based on trace characteristics and provides a scalable architecture for processing high-volume telemetry data.

## Overview

The ECS OpenTelemetry Tail Sampling solution consists of the following components:

- **Load Balancer (Gateway)**: Receives traces from ECS applications via OTLP protocol and distributes them across sampling agents using trace ID routing for consistent sampling decisions
- **Sampling Agents**: Perform tail sampling decisions and forward selected traces to Coralogix
- **AWS Cloud Map**: Provides service discovery for the Load Balancer to dynamically discover and route traces to Sampling Agents, enabling automatic scaling and failover of sampling instances. This is the mechanism that allows traces from ECS applications to reach the Load Balancer endpoint (`grpc-lb.cx-otel:4317`)

## Prerequisites

- Amazon ECS cluster with EC2 launch type
- VPC with public and private subnets
- Coralogix Send-Your-Data API key

## Deployment Options

### Option 1: Complete Cluster Deployment (Recommended)

Deploy the entire tail sampling stack including CloudMap, Autoscaling group, Task Definitions, Gateways and Receivers in one step.

```bash
aws cloudformation deploy \
  --stack-name otel-tail-sampling \
  --template-file otel-complete-cluster.yaml \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    VpcId=vpc-xxxxxxxxx \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    KeyName=your-key-pair \
    ECSAMI=ami-xxxxxxxxx \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ApplicationName=my-app \
    SubsystemName=production \
    EnableSpanMetrics=true \
    LoadBalancerTaskCount=1 \
    SamplingTaskCount=1
```

### Option 2: Step-by-Step Deployment

For more control over the deployment process, you can deploy each component individually:

#### Step 1: Create Cloud Map Namespace

```bash
aws servicediscovery create-private-dns-namespace \
  --name cx-otel \
  --vpc vpc-xxxxxxxxx \
  --description "Cloud Map namespace for OpenTelemetry services"
```

**Note**: Save the returned `Id` - you'll need it for subsequent steps.

#### Step 2: Deploy ECS Cluster Infrastructure

```bash
aws cloudformation deploy \
  --stack-name otel-ecs-cluster \
  --template-file step2-ecs-cluster.yaml \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    ClusterName=otel-cluster \
    InstanceType=t3.medium \
    KeyName=your-key-pair \
    ECSAMI=ami-xxxxxxxxx \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx
```

#### Step 3: Deploy Load Balancer Service

```bash
# Use NAMESPACE_ID from step 1
aws cloudformation deploy \
  --stack-name otel-load-balancer \
  --template-file step3-load-balancer.yaml \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    ClusterName=otel-cluster \
    NamespaceId=$NAMESPACE_ID \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ImageVersion=latest \
    Memory=1024
```

#### Step 4: Deploy Sampling Agents

```bash
# Use NAMESPACE_ID from step 1
aws cloudformation deploy \
  --stack-name otel-sampling-agents \
  --template-file step4-sampling-agents.yaml \
  --capabilities CAPABILITY_IAM \
  --parameter-overrides \
    ClusterName=otel-cluster \
    NamespaceId=$NAMESPACE_ID \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    ImageVersion=latest \
    Memory=1024 \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ApplicationName=my-app \
    SubsystemName=sampling
```

## Tail Sampling Configuration

The sampling agents are configured with tail sampling policies. Here is the default Basic Probabilistic Sampling configuration used in this template:

```yaml
processors:
  tail_sampling:
    decision_wait: 10s
    num_traces: 5000
    policies:
      [
        {
          name: sampling-policy,
          type: probabilistic,
          probabilistic: { sampling_percentage: 25 },
        }
      ]
```

### Sampling Policies

You can customize the sampling policies based on your requirements:

- **Probabilistic Sampling**: Random sampling based on percentage
- **Error-based Sampling**: Sample traces with errors
- **Latency-based Sampling**: Sample traces exceeding latency thresholds
- **Composite Policies**: Combine multiple sampling criteria

## Application Configuration

### OTLP Endpoint Configuration

Configure your applications to send traces to the load balancer endpoint:

```bash
# For applications running in the same VPC
OTEL_EXPORTER_OTLP_ENDPOINT=http://grpc-lb.cx-otel:4317

# For applications outside the VPC (if NLB is configured)
OTEL_EXPORTER_OTLP_ENDPOINT=http://your-nlb-endpoint:4317
```

### Span Metrics Configuration

**Important**: Choose one approach for span metrics to avoid duplicate metrics generation.

#### Option 1: Span Metrics in Load Balancer (Recommended for centralized processing)
- Set `EnableSpanMetrics=true` in the CloudFormation parameters
- **Disable span metrics in your ECS OTEL collector configuration**
- All span metrics are generated centrally in the Load Balancer
- Benefits: Single point of processing, consistent aggregation, centralized configuration

#### Option 2: Span Metrics in ECS Collectors (Recommended for distributed processing)
- Set `EnableSpanMetrics=false` in the CloudFormation parameters
- **Enable span metrics in your ECS OTEL collector configuration**
- Each ECS instance generates its own span metrics
- Benefits: Local processing, better fault tolerance, lower latency

**Note**: Both approaches will preserve the same resource attributes and span metadata. The choice depends on your operational preferences and architecture requirements.


## Parameters

### Complete Cluster Template Parameters

| Parameter | Description | Default | Required |
|-----------|-------------|---------|----------|
| VpcId | VPC ID where resources will be deployed | | Yes |
| SubnetIds | Comma-separated list of subnet IDs | | Yes |
| SecurityGroupId | Security group ID for ECS tasks | | Yes |
| KeyName | EC2 key pair name | | Yes |
| ECSAMI | ECS-optimized AMI ID | | Yes |
| CoralogixDomain | Coralogix domain (e.g., coralogix.com) | | Yes |
| CoralogixPrivateKey | Coralogix Send-Your-Data API key | | Yes |
| ApplicationName | Application name for telemetry | my-app | No |
| SubsystemName | Subsystem name for telemetry | production | No |
| InstanceType | EC2 instance type | t3.medium | No |
| DesiredCapacity | Desired number of EC2 instances | 2 | No |
| MaxSize | Maximum number of EC2 instances | 4 | No |
| MinSize | Minimum number of EC2 instances | 1 | No |
| EnableSpanMetrics | Enable span metrics generation in Load Balancer | false | No |
| LoadBalancerTaskCount | Number of Load Balancer tasks to run | 1 | No |
| SamplingTaskCount | Number of Sampling Agent tasks to run | 1 | No |

### Step-by-Step Template Parameters

Each step template has its own set of parameters. Refer to the individual template files for detailed parameter descriptions.
