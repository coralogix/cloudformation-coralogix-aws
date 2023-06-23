# Prometheus on ECS-EC2 with Auto-Discovery

This cloudformation template deploys the CWAgent and Prometheus Node Exporter on an ECS EC2 Cluster. The Cloudwatch Agent is configured with auto-discover to detect Prometheus scrape endpoints and store them in a a file. The Promotheus node exporter is configured to scrape the endpoints discovered by the Cloudwatch Agent from the file.

```
+---------------------------------------------+
|   ECS Cluster                               |
|                                             |
|                        Scan cluster for     |
| +-------------------+  Prometheus targets   |
| |                   |---------------------+ |
| |   CloudWatch      |                       |
| |     Agents        |<--------------------+ |
| |                   |                       |
| +-------------------+                       |
|             |                               |
|             |  Store targets                |
|             V                               |
| +----------------------------------+        |
| | Prometheus Scrape Endpoints File |        |
| +----------------------------------+        |
|             ^                               |
|             | Uses                          |
| +------------------+                        |
| |                  |                        |
| | Prometheus Node  |   -------------------------> [Coralogix]
| |     Exporter     |                        |
| |                  |                        |
| +------------------+                        |
|             |                               |
|             | Scrapes                       |
|             V                               |
| +--------------------+                      |
| | Prometheus Targets |                      |
| +--------------------+                      |
|                                             |
+---------------------------------------------+


```

##### Requires:

In order for a Prometheus target container to be detected by the Cloudwatch Agent, it must be annotated with the following __Docker labels__:

| Label | Description | Default |
| --- | --- | --- |
| ECS_PROMETHEUS_EXPORTER_PORT | The Prometheus exporter port | "9090" |
| ECS_PROMETHEUS_METRICS_PATH | The Prometheus expoert metric path | "/metrics" |




##### Parameters:

| Parameter | Description | Type | Default | Required |
|---|---|---|---|---|
| ECSClusterName | Enter the name of your ECS cluster from which you want to collect Prometheus metrics | String | | ✔️ |
| CreateIAMRoles | Whether to create new IAM roles or use an existing IAM roles for the ECS tasks | String | `True` | |
| ECSNetworkMode | ECS Network Mode for the Task | String | `bridge` | |
| TaskRoleName | Enter the CloudWatch Agent ECS task role name | String | `ECSDiscoveryCWAgentTaskRoleName` | |
| ExecutionRoleName | Enter the CloudWatch Agent ECS execution role name | String | `ECSDiscoveryCWAgentExecutionRoleName` | |
| CoralogixRegion | The Coralogix location region | String | `Europe` | ✔️ |
| CoralogixPrivateKey | The Coralogix Private Key | String | `NoEcho: true` | ✔️ |
| PrometheusNodeExporterImage | The Prometheus Node Exporter Image | String | `prom/node-exporter:v1.0.1` | |


##### Deploy

```
aws cloudformation deploy --template-file template.yaml \
    --stack-name prometheus-ecs-ec2-autodiscovery \
    --region <region> \
    --parameter-overrides \
        ECSClusterName=<your cluster name> \
        CoralogixPrivateKey=<your-private-key> \
        CoralogixRegion=<your-region> \
    --capabilities CAPABILITY_NAMED_IAM
```


##### Validation & Troubleshooting

- The Cloudwatch Agent will store the Prometheus scrape endpoints in the file: `/tmp/cwagent_ecs_auto_sd.yaml`. The file will be mounted to the container from  the host. You can check the contents of the file to verify that the Prometheus scrape endpoints are being discovered.

The contents should look similar to the following:

```yaml
- targets:
  - 10.0.0.127:32768
  labels:
    __metrics_path__: /metrics
    ECS_PROMETHEUS_EXPORTER_PORT: "9090"
    ECS_PROMETHEUS_METRICS_PATH: /metrics
    InstanceType: t2.large
    LaunchType: EC2
    StartedBy: ecs-svc/7562722944147020298
    SubnetId: subnet-076f1018de94369dd
    TaskClusterName: cds-305
    TaskDefinitionFamily: prometheus-task-definition
    TaskGroup: service:prometheus-service
    TaskId: 07f62e85057342cc916c27e7232b9546
    TaskRevision: "10"
    VpcId: vpc-02bbde9df35cf5d95
    container_name: prometheus
- targets:
  - 10.0.1.77:32902
  labels:
    __metrics_path__: /metrics
    ECS_PROMETHEUS_EXPORTER_PORT: "9090"
    ECS_PROMETHEUS_METRICS_PATH: /metrics
    InstanceType: t2.large
    LaunchType: EC2
    StartedBy: ecs-svc/7562722944147020298
    SubnetId: subnet-0445660746d6d06b4
    TaskClusterName: cds-305
    TaskDefinitionFamily: prometheus-task-definition
    TaskGroup: service:prometheus-service
    TaskId: 43479e93fa184e24bbf0ddf0e42c6cd2
    TaskRevision: "10"
    VpcId: vpc-02bbde9df35cf5d95
    container_name: prometheus
```

- By default the Prometheus Node Export is configured to send to Coralogix. You can also verify that the Prometheus scrape endpoints are being scraped by checking the Coralogix console.
<br>
- If data is not being received, check to make sure that there is network connectivity between the Scraper and all the targets. In order for targets to be scraped from within the ECS Cluster, the relevant security groups need to allow sufficient inbound/outbound access to the Prometheus Node Exporter and the Prometheus targets.
<br>
- If targets are not being detected, check to make sure that labels are defined correctly. Also note that if the port or path specific are invalid or do not represent a prometheus compatible endpoint, the target will not be detected. For eg. If we set the `ECS_PROMETHEUS_EXPORTER_PORT` to `8080` and the `ECS_PROMETHEUS_METRICS_PATH` to `/metrics`, the target will not be detected because if the Prometheus Node Exporter does not expose metrics on port `8080` and the `/metrics` path is not valid for the Prometheus Node Exporter.

---

This template is meant to be a starting point for leveraging the CW Agent for Prometheus Auto Discovery in ECS-EC2 using Docker Labels, it is also possible to detect prometheus endpoints using Task Definition ARN patterns.

The detection mechanism can also detect any prometheus compatible endpoint, not just the Prometheus Node Exporters. For eg. Cadvisor, Prometheus Push Gateway, etc.

For more information on additional configurations and features, please reference the [AWS Docuementation](https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContainerInsights-Prometheus-Setup-autodiscovery-ecs.html).