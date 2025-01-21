# Changelog

## AwsMetrics

### 22.1.2025 - Breaking changes to be able to tie role to specific company:

- rename `ExternalId` to `ExternalIdSecret`
- Add validation to `ExternalIdSecret`, must be a valid pattern of `[\w+=,.:\/-]*`
- Add validation to `CustomerAccountId`, must be a valid pattern of `[0-9]*`
- rename `CustomerAccountId` to `CustomAccountId`

### 12.12.2024 New permissions, that would allow integration to get data from Amazon ElastiCache API

- Added permissions to the policy:

    ```cloudformation
    - elasticache:DescribeCacheClusters
    - elasticache:ListTagsForResource
    ```

### 27.11.2024 Fix CustomAccountId field name for custom account and align template description

### 18.10.2024 Fix rds:ListTagsForResource permission

### 30.9.2024 Make the policy more secure by adding more specific permissions and not using `*` in the policy

### 25.9.2024 New permission to be able to get RDS and EC2 instance metadata like allocated memory, CPU, etc

### 6.9.2024 New permission for ECS enhanced monitoring

### 2.9.2024 Enable ap3 as allowed value in AWSAccount

### 19.8.2024 New permission for RDS enhanced monitoring & add support for new environment ap3

### 7.8.2024 Add output to the role that the module will create

### 10.7.2024 Depend on specific role to make it more secure & pass wiz security check

### 2.7.2024 New permissions for future convenience, we would be adding calls that need these permissions in near future

- Add permissions to the policy:

    ```cloudformation
    - apigateway:GET
    - autoscaling:DescribeAutoScalingGroups
    - aps:ListWorkspaces
    - dms:DescribeReplicationInstances
    - dms:DescribeReplicationTasks
    - ec2:DescribeTransitGatewayAttachments
    - ec2:DescribeSpotFleetRequests
    - storagegateway:ListGateways
    - storagegateway:ListTagsForResource
    ```

### 23.4.2024 New permission needed

- Add permission `cloudwatch:GetMetricData` to the policy

### 11.4.2024 🚀 New components 🚀

- Create metrics integration module
