# Changelog

## AwsMetrics

### 2.7.2024
### New permissions for future convenience, we would be adding calls that need these permissions in near future
- Add permissions:
    ```
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

### 23.4.2024
### New permission needed
- Add permission `cloudwatch:GetMetricData` to the policy

### 11.4.2024
### 🚀 New components 🚀
- Create metrics integration module

