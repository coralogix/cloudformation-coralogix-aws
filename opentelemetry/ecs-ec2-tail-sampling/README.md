# Tail Sampling with OpenTelemetry using AWS ECS EC2

This tutorial demonstrates how to configure an AWS ECS EC2 cluster, deploy OpenTelemetry to collect logs, metrics, and traces, and enable intelligent trace sampling. We will cover four deployment approaches:

1. **[ECS Cluster with Tail Sampling](#ecs-cluster-with-tail-sampling)**: Deploy agents and gateway directly to your application cluster
2. **[Central Collector Cluster for Tail Sampling](#central-collector-cluster-for-tail-sampling)**: Deploy a dedicated telemetry cluster separate from your applications
3. **[Verification using Telemetrygen](#verification-using-telemetrygen)**: Test and verify your deployment with sample traces
4. **[Individual Templates](#individual-templates)**: Deploy components individually for granular control

## ECS Cluster with Tail Sampling

### How it Works

The AWS ECS EC2 OpenTelemetry Integration consists of the following components:

**OpenTelemetry Agent.** The Agent is deployed as a daemon service on each ECS instance within the cluster and collects telemetry data from the applications running on that instance. The agent is configured to send the logs and metrics to Coralogix and Traces to OpenTelemetry Gateway. The agent ensures that traces with the same ID are sent to the same gateway. This allows tail sampling to be performed on the traces correctly, even if they span multiple applications and instances.

**OpenTelemetry Gateway.** The Gateway is responsible for receiving telemetry data from the agents and forwarding it to the Coralogix backend. The Gateway is also responsible for performing tail sampling decisions and load balancing the telemetry data to the Coralogix backend.

## Prerequisites

- **Existing AWS ECS Cluster**: EC2-based cluster with instances running
- **S3 Bucket**: For storing OpenTelemetry configuration files
- **Send-Your-Data API Key**: Your Coralogix [Send-Your-Data API key](../../../user-guides/account-management/api-keys/send-your-data-api-key/index.md)
- **Coralogix Domain**: Your [Coralogix domain](../../../user-guides/account-management/account-settings/coralogix-domain/index.md) (e.g., `eu2.coralogix.com`, `us2.coralogix.com`)
- **Application and Subsystem Names**: [Application and subsystem names](../../../user-guides/account-management/account-settings/application-and-subsystem-names/index.md) for organizing your data in Coralogix

## Deployment Options

### Option 1: ECS Cluster with Tail Sampling

Deploy OpenTelemetry agents and gateway directly to your application cluster. This approach is suitable when you want to keep telemetry processing close to your applications.

**Quick Start:**

#### 1. Prepare Configuration Files

Upload your OpenTelemetry configuration files to S3:

```bash
# Agent configuration (for daemon)
aws s3 cp examples/agent-config.yaml s3://your-bucket/configs/agent-config.yaml

# Sampling configuration (for sampling agents)
aws s3 cp examples/sampling-config.yaml s3://your-bucket/configs/sampling-config.yaml
```

#### 2. Deploy using CloudFormation template 

```bash
aws cloudformation deploy \
  --stack-name otel-tail-sampling \
  --template-file otel-tail-sampling.yaml \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=your-existing-cluster \
    VpcId=vpc-xxxxxxxxx \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    CoralogixRegion=EU2 \
    CoralogixApiKey=your-api-key \
    CDOTImageVersion=image-version\
    S3ConfigBucket=your-bucket \
    AgentS3ConfigKey=configs/agent-config.yaml \
    SamplingS3ConfigKey=configs/sampling-config.yaml \
    ApplicationName=my-app \
    SubsystemName=production \
    SamplingTaskCount=2 \
    TaskExecutionRoleArn=if-not-provided-template-creates-it
```

### Option 2: Central Collector Cluster for Tail Sampling

Deploy a dedicated telemetry cluster separate from your application workloads. This approach provides centralized telemetry processing and is ideal for multi-cluster environments.

**See the [Central Collector Cluster for Tail Sampling](#central-collector-cluster-for-tail-sampling) section below for detailed instructions.**

## Sampling Policies

For detailed information about available sampling policies and configuration options, see the [OpenTelemetry Tail Sampling Processor documentation](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/processor/tailsamplingprocessor).

## Service Discovery

The solution uses AWS CloudMap for service discovery:

- **Namespace**: `cx-otel`
- **Service**: `grpc-gateway`
- **DNS**: `grpc-gateway.cx-otel`

The agent automatically discovers sampling agents using the loadbalancing exporter.

### Scaling

Adjust the number of sampling agents based on your workload:

```bash
aws ecs update-service \
  --cluster your-cluster \
  --service coralogix-otel-sampling \
  --desired-count 4
```

### External IAM Roles

Use existing IAM roles instead of creating new ones. The role must have the following permissions:

**Required Policies:**
- `AmazonECSTaskExecutionRolePolicy` (AWS managed policy)
- Custom policy for S3 read access to your configuration bucket
- Custom policy for CloudMap service discovery

**Example Role Policy:**
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:GetObjectVersion"
            ],
            "Resource": "arn:aws:s3:::your-config-bucket/*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "servicediscovery:DiscoverInstances",
                "servicediscovery:ListInstances",
                "servicediscovery:ListServices",
                "servicediscovery:ListNamespaces"
            ],
            "Resource": "*"
        }
    ]
}
```

**Usage:**
```bash
--parameter-overrides \
  TaskExecutionRoleArn=arn:aws:iam::ACCOUNT:role/existing-role
```

## Central Collector Cluster for Tail Sampling

This section describes how to deploy OpenTelemetry services to an existing ECS cluster for centralized tail sampling. This approach separates the telemetry collection infrastructure from your application workloads, providing centralized telemetry processing and intelligent sampling capabilities.

### Architecture

The central collector deployment consists of:

- **Receiver Services**: Collect telemetry data from applications and route traces to the gateway
- **Gateway Services**: Perform tail sampling and forward to Coralogix
- **CloudMap**: Service discovery for dynamic load balancing

**Traces Flow**: Applications → Receiver → Gateway → Coralogix

### Prerequisites

- **Existing AWS ECS Cluster**: EC2-based cluster with instances running
- **S3 Bucket**: For storing OpenTelemetry configuration files
- **Send-Your-Data API Key**: Your Coralogix [Send-Your-Data API key](../../../user-guides/account-management/api-keys/send-your-data-api-key/index.md)
- **Coralogix Domain**: Your [Coralogix domain](../../../user-guides/account-management/account-settings/coralogix-domain/index.md) (e.g., `eu2.coralogix.com`, `us2.coralogix.com`)
- **Application and Subsystem Names**: [Application and subsystem names](../../../user-guides/account-management/account-settings/application-and-subsystem-names/index.md) for organizing your data in Coralogix

### Important Configuration Note

**⚠️ Spanmetrics Configuration**: When using the central collector approach, spanmetrics should be configured on either the agents/receivers OR the gateway, but not both. Configuring spanmetrics on both components can lead to duplicate metrics and incorrect data.

- **Option A**: Configure spanmetrics on agents/receivers (recommended for most use cases)
- **Option B**: Configure spanmetrics on the gateway (use when you need centralized span processing)

### Deployment Steps

#### 1. Prepare Configuration Files

Upload your OpenTelemetry configuration files to S3:

```bash
# Receiver configuration
aws s3 cp examples/receiver-config.yaml s3://your-bucket/configs/receiver-config.yaml

# Gateway configuration (with tail sampling)
aws s3 cp examples/gateway-config.yaml s3://your-bucket/configs/gateway-config.yaml
```

#### 2. Deploy to Existing Cluster

Deploy OpenTelemetry services to your existing ECS cluster:

```bash
aws cloudformation deploy \
  --stack-name otel-complete-cluster \
  --template-file otel-complete-cluster.yaml \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides \
    ClusterName=otel-cluster \
    VpcId=vpc-xxxxxxxxx \
    SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" \
    SecurityGroupId=sg-xxxxxxxxx \
    KeyName=your-key-pair \
    ECSAMI=ami-xxxxxxxxx \
    CoralogixRegion=EU2 \
    CoralogixApiKey=your-api-key \
    CDOTImageVersion=Image-version \
    S3ConfigBucket=your-bucket \
    ReceiverS3ConfigKey=configs/receiver-config.yaml \
    GatewayS3ConfigKey=configs/gateway-config.yaml \
    ApplicationName=my-app \
    SubsystemName=production \
    ReceiverTaskCount=2 \
    GatewayTaskCount=1 \
    TaskExecutionRoleArn=if-not-provided-template-creates-it
```

**Key Features:**
- **Load Balancing**: Routes traces to multiple gateway instances
- **Trace ID Routing**: Ensures traces with the same ID go to the same gateway
- **Service Discovery**: Automatically discovers gateway instances via CloudMap

### Service Discovery

The solution uses AWS CloudMap for service discovery:

- **Namespace**: `cx-otel`
- **Receiver Service**: `grpc-receiver.cx-otel`
- **Gateway Service**: `grpc-gateway.cx-otel`

### Connecting Applications

To connect your applications to the central collector cluster:

1. **Update Application Configuration**: Point your application's OTLP exporter to the receiver service
2. **Use Service Discovery**: Applications can discover receivers using the CloudMap DNS name
3. **Load Balancing**: Multiple receivers provide high availability and load distribution

### Enabling Other Clusters to Send Data

To enable other ECS clusters or applications to send telemetry data to the central collector cluster:

#### Method 1: Cross-Cluster Service Discovery

1. **Share CloudMap Namespace**: Ensure both clusters can access the same CloudMap namespace
2. **Update Application Configuration**: Point applications to the receiver service:
   ```bash
   --otlp-endpoint=grpc-receiver.cx-otel:4317
   ```
3. **Network Connectivity**: Ensure VPC peering or transit gateway for cross-cluster communication

#### Method 2: Load Balancer Endpoint

1. **Create Application Load Balancer**: Expose the receiver service via ALB
2. **Update Security Groups**: Allow traffic from other clusters
3. **Use ALB DNS Name**: Point applications to the load balancer endpoint

#### Method 3: VPC Endpoints

1. **Create VPC Endpoints**: For ECS and CloudMap services
2. **Cross-Account Access**: If clusters are in different AWS accounts
3. **IAM Permissions**: Ensure proper permissions for cross-cluster access

#### Configuration Example

For applications in other clusters, update their OpenTelemetry configuration:

```yaml
exporters:
  otlp:
    endpoint: grpc-receiver.cx-otel:4317
    tls:
      insecure: true
```

## Verification using Telemetrygen

To verify that your OpenTelemetry deployment is working correctly and traces are reaching Coralogix, you can deploy a telemetrygen task that generates test traces.

### Deploying Telemetrygen Task

#### 1. Create Task Definition

Create a task definition for telemetrygen that connects to your receiver service:

```json
{
    "compatibilities": [
        "EXTERNAL",
        "EC2"
    ],
    "containerDefinitions": [
        {
            "command": [
                "traces",
                "--otlp-endpoint=grpc-receiver.cx-otel:4317",
                "--otlp-insecure",
                "--rate=10",
                "--duration=1h"
            ],
            "cpu": 256,
            "environment": [],
            "essential": true,
            "image": "ghcr.io/open-telemetry/opentelemetry-collector-contrib/telemetrygen:latest",
            "logConfiguration": {
                "logDriver": "awslogs",
                "options": {
                    "awslogs-group": "/ecs/telemetrygen",
                    "awslogs-region": "us-east-1",
                    "awslogs-stream-prefix": "ecs"
                }
            },
            "memory": 512,
            "mountPoints": [],
            "name": "telemetrygen",
            "portMappings": [
                {
                    "containerPort": 8080,
                    "hostPort": 0,
                    "protocol": "tcp"
                }
            ],
            "systemControls": [],
            "volumesFrom": []
        }
    ],
    "cpu": "256",
    "executionRoleArn": "arn:aws:iam::ACCOUNT:role/ecsTaskExecutionRole",
    "family": "telemetrygen-task",
    "memory": "512",
    "networkMode": "bridge",
    "placementConstraints": [],
    "requiresCompatibilities": [
        "EC2"
    ],
    "volumes": []
}
```

#### 2. Register Task Definition

```bash
aws ecs register-task-definition --cli-input-json file://telemetrygen-task-definition.json
```

#### 3. Run Telemetrygen Task

```bash
aws ecs run-task \
  --cluster your-cluster-name \
  --task-definition telemetrygen-task \
  --launch-type EC2 \
  --count 1
```

### Configuration Options

#### For Complete Cluster Deployment

When using the complete cluster template (`otel-complete-cluster.yaml`), use the receiver service:

```bash
--otlp-endpoint=grpc-receiver.cx-otel:4317
```

#### For Tail Sampling Deployment

When using the tail sampling template (`otel-complete-tail-sampling.yaml`), use the agent service:

```bash
--otlp-endpoint=grpc-sampling.cx-otel:4317
```


**Verify in Coralogix**: Check your Coralogix dashboard for incoming traces
   - Look for traces with service name matching your configuration
   - Verify trace sampling is working as expected

## Individual Templates

For more granular control over your deployment, you can use individual CloudFormation templates located in the `individual-templates/` directory:

### Available Templates:

- **`cloudmap-namespace.yaml`**: Creates the CloudMap namespace for service discovery
- **`otel-daemon-template.yaml`**: Deploys the OpenTelemetry agent as a daemon service
- **`load-balancer-agents.yaml`**: Deploys the receiver service (for central collector approach)
- **`sampling-agents.yaml`**: Deploys the gateway service for tail sampling

### Example Deployment:

These commands will deploy a complete telemetry infrastructure with agents sending logs and metrics directly to Coralogix, and traces to the Receiver for load balancing, then to the Gateway for tail sampling decisions.

1. **CloudMap Namespace** (if not using complete template):
   ```bash
   aws cloudformation deploy --stack-name otel-cloudmap --template-file individual-templates/cloudmap-namespace.yaml --capabilities CAPABILITY_NAMED_IAM --parameter-overrides VpcId=vpc-xxxxxxxxx
   ```

2. **Daemon Agent**:
   ```bash
   aws cloudformation deploy --stack-name otel-daemon --template-file individual-templates/otel-daemon-template.yaml --capabilities CAPABILITY_NAMED_IAM --parameter-overrides ClusterName=your-cluster S3ConfigBucket=your-bucket S3ConfigKey=configs/agent-config-to-receiver.yaml CoralogixRegion=EU2 CoralogixApiKey=your-api-key CDOTImageVersion=v0.5.0
   ```

3. **Receiver** (for central collector approach):
   ```bash
   aws cloudformation deploy --stack-name otel-receiver --template-file individual-templates/load-balancer-agents.yaml --capabilities CAPABILITY_NAMED_IAM --parameter-overrides ClusterName=your-cluster NamespaceId=ns-xxxxxxxxx SubnetIds="subnet-xxxxxxxxx,subnet-yyyyyyyyy" SecurityGroupId=sg-xxxxxxxxx CDOTImageVersion=v0.5.0 CoralogixRegion=EU2 CoralogixApiKey=your-api-key S3ConfigBucket=your-bucket ReceiverS3ConfigKey=configs/receiver-config-with-span-metrics.yaml ReceiverTaskCount=2
   ```

4. **Gateway**:
   ```bash
   aws cloudformation deploy --stack-name otel-gateway --template-file individual-templates/sampling-agents.yaml --capabilities CAPABILITY_NAMED_IAM --parameter-overrides ClusterName=your-cluster NamespaceId=ns-xxxxxxxxx Subnets="subnet-xxxxxxxxx,subnet-yyyyyyyyy" SecurityGroupId=sg-xxxxxxxxx CDOTImageVersion=v0.5.0 CoralogixRegion=EU2 CoralogixApiKey=your-api-key S3ConfigBucket=your-bucket GatewayS3ConfigKey=configs/gateway-config.yaml GatewayTaskCount=2
   ```

### Template Features:

All individual templates support:
- **External IAM Roles**: Use existing roles via `TaskExecutionRoleArn` parameter
- **Conditional Resource Creation**: IAM roles are only created when external roles are not provided
- **Flexible Configuration**: Each component can be deployed independently
- **Service Discovery**: Automatic CloudMap integration for load balancing

## Support

For issues and questions:

- **Documentation**: [Coralogix OpenTelemetry Documentation](https://coralogix.com/docs/opentelemetry/)
- **Community**: [Coralogix Community](https://community.coralogix.com/)

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
