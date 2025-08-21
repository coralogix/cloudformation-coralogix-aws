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
- **For S3 Configuration**: S3 bucket with OpenTelemetry configuration files (See examples in examples/ folder)
- **For Existing Clusters**: ECS instance role with required permissions (see IAM Requirements below)

### IAM Requirements

#### For New Clusters (Created by Templates)
The templates automatically create an ECS instance role with the following permissions:
- `AmazonEC2ContainerServiceforEC2Role` - ECS container service permissions
- `AmazonSSMManagedInstanceCore` - SSM access for instance management
- `AWSCloudMapFullAccess` - Cloud Map service discovery permissions
- **S3 Permissions** (when using S3 configuration):
  - `s3:GetObject` and `s3:GetObjectVersion` on your S3 bucket

#### For Existing Clusters
If using an existing ECS cluster, ensure the instance role has these permissions:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecs:*",
        "ec2:*",
        "elasticloadbalancing:*",
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "logs:*"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "servicediscovery:*"
      ],
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:GetObjectVersion"
      ],
      "Resource": "arn:aws:s3:::your-config-bucket/*"
    }
  ]
}
```

## Deployment Options

### Option 1: Complete Cluster Deployment (Recommended)

Deploy the entire tail sampling stack including CloudMap, Autoscaling group, Task Definitions, Gateways and Receivers in one step.

```bash
# Using built-in template configuration
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
    ImageVersion=x.x.x \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ApplicationName=my-app \
    SubsystemName=production \
    EnableSpanMetrics=true \
    LoadBalancerTaskCount=1 \
    SamplingTaskCount=1

# Using S3 configuration files
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
    ImageVersion=x.x.x \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ApplicationName=my-app \
    SubsystemName=production \
    ConfigSource=s3 \
    S3ConfigBucket=your-config-bucket \
    LbS3ConfigKey=configs/lb-config-with-spanmetrics.yaml \
    SamplingS3ConfigKey=configs/sampling-config.yaml \
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
  --template-file ecs-cluster.yaml \
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
# Using built-in template configuration
aws cloudformation deploy \
  --stack-name otel-load-balancer \
  --template-file load-balancer-agents.yaml \
  --parameter-overrides \
    ClusterName=otel-cluster \
    NamespaceId=$NAMESPACE_ID \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ImageVersion=x.x.x \
    Memory=1024 \
    LoadBalancerTaskCount=1

# Using S3 configuration
aws cloudformation deploy \
  --stack-name otel-load-balancer \
  --template-file load-balancer-agents.yaml \
  --parameter-overrides \
    ClusterName=otel-cluster \
    NamespaceId=$NAMESPACE_ID \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ImageVersion=x.x.x \
    Memory=1024 \
    ConfigSource=s3 \
    S3ConfigBucket=your-config-bucket \
    S3ConfigKey=configs/lb-config-with-spanmetrics.yaml \
    LoadBalancerTaskCount=1
```

#### Step 4: Deploy Sampling Agents

```bash
# Use NAMESPACE_ID from step 1
# Using built-in template configuration
aws cloudformation deploy \
  --stack-name otel-sampling-agents \
  --template-file sampling-agents.yaml \
  --parameter-overrides \
    ClusterName=otel-cluster \
    NamespaceId=$NAMESPACE_ID \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    ImageVersion=x.x.x \
    Memory=1024 \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ApplicationName=my-app \
    SubsystemName=sampling \
    SamplingTaskCount=1

# Using S3 configuration
aws cloudformation deploy \
  --stack-name otel-sampling-agents \
  --template-file sampling-agents.yaml \
  --parameter-overrides \
    ClusterName=otel-cluster \
    NamespaceId=$NAMESPACE_ID \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    ImageVersion=x.x.x \
    Memory=1024 \
    CoralogixDomain=coralogix.com \
    CoralogixPrivateKey=your-private-key \
    ApplicationName=my-app \
    SubsystemName=sampling \
    ConfigSource=s3 \
    S3ConfigBucket=your-config-bucket \
    S3ConfigKey=configs/sampling-config.yaml \
    SamplingTaskCount=1
```

## Configuration Options

### Configuration Source

The templates support two configuration sources:

#### Option 1: Built-in Template Configuration (Default)
- Configuration is embedded directly in the CloudFormation templates
- Environment variables are substituted at runtime
- No external dependencies required

#### Option 2: S3 Configuration Files
- Configuration files stored in S3 bucket
- Supports dynamic configuration updates without template changes
- Requires S3 bucket with proper permissions

**Example Configuration Files:**
- [`examples/lb-config-with-spanmetrics.yaml`](examples/lb-config-with-spanmetrics.yaml) - Load Balancer configuration with span metrics enabled
- [`examples/lb-config-no-spanmetrics.yaml`](examples/lb-config-no-spanmetrics.yaml) - Load Balancer configuration without span metrics
- [`examples/sampling-config.yaml`](examples/sampling-config.yaml) - Sampling agents configuration with tail sampling policies

### Tail Sampling Configuration

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


### Modifying the inline OpenTelemetry Configuration

If you need to customize the OpenTelemetry Collector configuration, it's recommended to use S3 configuration files for easier management. However, you can also modify the embedded configuration directly in the CloudFormation templates:

#### **Complete Cluster Template (`otel-complete-cluster.yaml`)**

**Load Balancer Configuration:**
- **Location**: Lines ~200-300 (search for `OTEL_CONFIG`)
- **What to modify**: Receivers, processors, exporters, connectors, and pipelines for the load balancer service

**Sampling Agents Configuration:**
- **Location**: Lines ~400-500 (search for `OTEL_CONFIG`)
- **What to modify**: Tail sampling policies, processors, and exporters for the sampling agents

#### **Step-by-Step Templates**

**Load Balancer (`load-balancer-agents.yaml`):**
- **Location**: Lines ~80-200 (search for `OTEL_CONFIG`)
- **What to modify**: Load balancer configuration including spanmetrics connector

**Sampling Agents (`sampling-agents.yaml`):**
- **Location**: Lines ~80-150 (search for `OTEL_CONFIG`)
- **What to modify**: Tail sampling policies and processing configuration

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
| **Configuration Source Parameters** |
| ConfigSource | Configuration source: 'template' or 's3' | template | No |
| S3ConfigBucket | S3 bucket name for configuration files | | Required if ConfigSource=s3 |
| LbS3ConfigKey | S3 key for Load Balancer configuration file | | Required if ConfigSource=s3 |
| SamplingS3ConfigKey | S3 key for Sampling configuration file | | Required if ConfigSource=s3 |
| ECSInstanceRoleArn | ARN of existing ECS instance role | | No |

### Docker Image Tags

**Important**: The `ImageVersion` parameter is required and must be explicitly provided. We recommend using a specific version tag rather than `latest` for production deployments.

- **Recommended**: Use the official Coralogix supported version from the [Helm chart](https://github.com/coralogix/opentelemetry-helm-charts/blob/main/charts/opentelemetry-collector/Chart.yaml#L15)
- **All Available Tags**: Browse all available Docker image tags at [Docker Hub](https://hub.docker.com/r/otel/opentelemetry-collector-contrib/tags)


### Step-by-Step Template Parameters

Each step template has its own set of parameters. Refer to the individual template files for detailed parameter descriptions.

## Using Existing ECS Clusters

The individual service templates (`load-balancer-agents.yaml` and `sampling-agents.yaml`) can be deployed to existing ECS clusters. Ensure your existing cluster meets these requirements:

### Prerequisites for Existing Clusters

1. **ECS Instance Role Permissions**: The cluster's instance role must have the permissions listed in the IAM Requirements section above
2. **Cloud Map Namespace**: Access to the Cloud Map namespace specified in `NamespaceId`
3. **Network Access**: Subnets and security groups must allow required traffic
4. **ECS Capacity**: Sufficient ECS capacity to run the services

### Deployment to Existing Cluster

```bash
# Deploy Load Balancer to existing cluster
aws cloudformation deploy \
  --stack-name coralogix-lb-service \
  --template-file load-balancer-agents.yaml \
  --parameter-overrides \
    ClusterName=your-existing-cluster \
    NamespaceId=ns-xxxxx \
    SubnetIds="subnet-xxx,subnet-yyy" \
    SecurityGroupId=sg-xxx \
    ImageVersion=x.x.x \
    CoralogixDomain=eu2.coralogix.com \
    CoralogixPrivateKey=your-key \
    ConfigSource=s3 \
    S3ConfigBucket=your-bucket \
    S3ConfigKey=configs/lb-config.yaml

# Deploy Sampling to existing cluster  
aws cloudformation deploy \
  --stack-name coralogix-sampling-service \
  --template-file sampling-agents.yaml \
  --parameter-overrides \
    ClusterName=your-existing-cluster \
    NamespaceId=ns-xxxxx \
    SubnetIds="subnet-xxx,subnet-yyy" \
    SecurityGroupId=sg-xxx \
    ImageVersion=x.x.x \
    CoralogixDomain=eu2.coralogix.com \
    CoralogixPrivateKey=your-key \
    ConfigSource=s3 \
    S3ConfigBucket=your-bucket \
    S3ConfigKey=configs/sampling-config.yaml
```
