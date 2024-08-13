# Changelog

## AwsMetrics

### 13.8.2024
### New permission for RDS enhanced monitoring

### 7.8.2024
### Add output to the role that the module will create

### 10.7.2024
### Depend on specific role to make it more secure & pass wiz security check

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

