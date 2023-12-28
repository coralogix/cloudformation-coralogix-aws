# Prometheus on ECS-EC2 with Auto-Discovery

This cloudformation template deploys an Open Telemetry collector on an AWS ECS EC2 Cluster. The Open Telemetry collector is configured to automatically discover Prometheus scrape endpoints and store them in a a file using the [ecs_observer](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/observer/ecsobserver) extension. The collector is also configured to scrape the Prometheus endpoints it discovers and send the metrics to Coralogix.

Note that Open Telemetry handles, Auto Discovery, Scraping the detected Prometheus endpoints and sending the metrics to Coralogix. No other Prometheus components are required.

##### Requires:

In order for a Prometheus target container to be detected by the Cloudwatch Agent, it must be annotated with the following **Docker labels**:

| Label                        | Description                        | Default    |
|------------------------------|------------------------------------|------------|
| ECS_PROMETHEUS_EXPORTER_PORT | The Prometheus exporter port       | "9090"     |
| ECS_PROMETHEUS_METRICS_PATH  | The Prometheus expoert metric path | "/metrics" |

##### Parameters:

| Parameter           | Description                                                                                                       | Type   | Default                                | Required |
|---------------------|-------------------------------------------------------------------------------------------------------------------|--------|----------------------------------------|----------|
| ECSClusterName      | Enter the name of your ECS cluster from which you want to collect Prometheus metrics                              | String |                                        | ✔️        |
| CreateIAMRoles      | Whether to create new IAM roles or use an existing IAM roles for the ECS tasks                                    | String | `True`                                 |          |
| ECSNetworkMode      | ECS Network Mode for the Task                                                                                     | String | `bridge`                               |          |
| TaskRoleName        | Enter the CloudWatch Agent ECS task role name                                                                     | String | `ECSDiscoveryCWAgentTaskRoleName`      |          |
| ExecutionRoleName   | Enter the CloudWatch Agent ECS execution role name                                                                | String | `ECSDiscoveryCWAgentExecutionRoleName` |          |
| CoralogixRegion     | The Coralogix location region                                                                                     | String | `Europe`                               | ✔️        |
| CoralogixPrivateKey | The Coralogix Private Key                                                                                         | String | `NoEcho: true`                         | ✔️        |
| ImageTag            | The Coralogix Otel Collector image tag.<br>see [here](coralogixrepo/coralogix-otel-collector) for available tags: | String |                                        |          |

##### Open Telemetry Collector Configuration

The Open Telemetry configuration used for this template is [embedded in the template.yaml](./template.yaml#L91-L142) file. You can update this as required for your use case.

```yaml
        extensions:
          ecs_observer:
            result_file: /tmp/ecs_sd_targets.yaml
            cluster_region: 'eu-west-1'
            cluster_name: cds-305
            services:
              - name_pattern: ^.*$
            docker_labels:
              - port_label: ECS_PROMETHEUS_EXPORTER_PORT
                metrics_path_label: ECS_PROMETHEUS_METRICS_PATH

        receivers:
          prometheus:
            config:
              scrape_configs:
                - job_name: "ecs-task"
                  file_sd_configs:
                    - files:
                        - '/tmp/ecs_sd_targets.yaml' # MUST match the file name in ecs_observer.result_file

        processors:
          batch:

        exporters:
          coralogix:
            domain: "${DOMAIN}"
            private_key: "${CX_TOKEN}"
            application_name: "otel-discovery-collector"
            subsystem_name: "ecs"
            application_name_attributes:
            - "aws.ecs.container.name"
            - "docker.name"
            - "APP_NAME"
            subsystem_name_attributes:
            - "ecs.task.definition.family"
            - "log.file.name"
            - "service.name"
            - "SUB_SYS"
            timeout: 30s

        service:
          pipelines:
            metrics:
              receivers:
                - prometheus
              processors:
                - batch
              exporters:
                - coralogix

          extensions:
            - ecs_observer
```

##### Deploy

```
aws cloudformation deploy --template-file template.yaml \
    --stack-name otel-prometheus-ecs-ec2-autodiscovery \
    --region <region> \
    --parameter-overrides \
        ECSClusterName=<your cluster name> \
        CoralogixPrivateKey=<your-private-key> \
        CoralogixRegion=<your-region> \
    --capabilities CAPABILITY_NAMED_IAM
```

##### Validation & Troubleshooting

- The Open Telemetry collector will store the Prometheus scrape endpoints in the file: `/tmp/ecs_sd_targets.yaml`. The file will be mounted to the container from the host. You can check the contents of the file to verify that the Prometheus scrape endpoints are being discovered.

The contents should look similar to the following:

```yaml
- targets:
  - 10.0.0.127:32776
  labels:
    __meta_ecs_cluster_name: cds-305
    __meta_ecs_container_labels_ECS_PROMETHEUS_EXPORTER_PORT: "8080"
    __meta_ecs_container_labels_ECS_PROMETHEUS_METRICS_PATH: /metrics
    __meta_ecs_container_name: cadvisor
    __meta_ecs_ec2_instance_id: i-09421dc29f5304f5d
    __meta_ecs_ec2_instance_type: t2.large
    __meta_ecs_ec2_private_ip: 10.0.0.127
    __meta_ecs_ec2_public_ip: 3.252.44.108
    __meta_ecs_ec2_subnet_id: subnet-076f1018de94369dd
    __meta_ecs_ec2_tags_Name: ECS Instance - amazon-ecs-cli-setup-cds-305
    __meta_ecs_ec2_tags_aws_autoscaling_groupName: amazon-ecs-cli-setup-cds-305-EcsInstanceAsg-MJGW9NNJ81I4
    __meta_ecs_ec2_tags_aws_cloudformation_logical_id: EcsInstanceAsg
    __meta_ecs_ec2_tags_aws_cloudformation_stack_id: arn:aws:cloudformation:eu-west-1:035955823196:stack/amazon-ecs-cli-setup-cds-305/4eeb94c0-f4a0-11ed-862b-0691223038f1
    __meta_ecs_ec2_tags_aws_cloudformation_stack_name: amazon-ecs-cli-setup-cds-305
    __meta_ecs_ec2_vpc_id: vpc-02bbde9df35cf5d95
    __meta_ecs_health_status: UNKNOWN
    __meta_ecs_service_name: cadvisor-service
    __meta_ecs_source: arn:aws:ecs:eu-west-1:035955823196:task/cds-305/998e8700260a40afba5007358682bd14
    __meta_ecs_task_definition_family: cadvisor-task-definition
    __meta_ecs_task_definition_revision: "7"
    __meta_ecs_task_group: service:cadvisor-service
    __meta_ecs_task_launch_type: EC2
    __meta_ecs_task_started_by: ecs-svc/6176391408914978728
    __metrics_path__: /metrics
- targets:
  - 10.0.1.77:32911
  labels:
    __meta_ecs_cluster_name: cds-305
    __meta_ecs_container_labels_ECS_PROMETHEUS_EXPORTER_PORT: "8080"
    __meta_ecs_container_labels_ECS_PROMETHEUS_METRICS_PATH: /metrics
    __meta_ecs_container_name: cadvisor
    __meta_ecs_ec2_instance_id: i-0d1032d189ea93453
    __meta_ecs_ec2_instance_type: t2.large
    __meta_ecs_ec2_private_ip: 10.0.1.77
    __meta_ecs_ec2_public_ip: 3.248.181.116
    __meta_ecs_ec2_subnet_id: subnet-0445660746d6d06b4
    __meta_ecs_ec2_tags_Name: ECS Instance - amazon-ecs-cli-setup-cds-305
    __meta_ecs_ec2_tags_aws_autoscaling_groupName: amazon-ecs-cli-setup-cds-305-EcsInstanceAsg-MJGW9NNJ81I4
    __meta_ecs_ec2_tags_aws_cloudformation_logical_id: EcsInstanceAsg
    __meta_ecs_ec2_tags_aws_cloudformation_stack_id: arn:aws:cloudformation:eu-west-1:035955823196:stack/amazon-ecs-cli-setup-cds-305/4eeb94c0-f4a0-11ed-862b-0691223038f1
    __meta_ecs_ec2_tags_aws_cloudformation_stack_name: amazon-ecs-cli-setup-cds-305
    __meta_ecs_ec2_vpc_id: vpc-02bbde9df35cf5d95
    __meta_ecs_health_status: UNKNOWN
    __meta_ecs_service_name: cadvisor-service
    __meta_ecs_source: arn:aws:ecs:eu-west-1:035955823196:task/cds-305/ec7ff82b7a3a44a5bbbe9bcf11daee33
    __meta_ecs_task_definition_family: cadvisor-task-definition
    __meta_ecs_task_definition_revision: "7"
    __meta_ecs_task_group: service:cadvisor-service
    __meta_ecs_task_launch_type: EC2
    __meta_ecs_task_started_by: ecs-svc/6176391408914978728
    __metrics_path__: /metrics
```

- By default the Open Telemetry collector is configured to send to Coralogix. You can also verify that the Prometheus scrape endpoints are being scraped by checking the Coralogix console. <br>
- If data is not being received, check to make sure that there is network connectivity between the Scraper and all the targets. In order for targets to be scraped from within the ECS Cluster, the relevant security groups need to allow sufficient inbound/outbound access to the Open Telemetry collector and the Prometheus targets. <br>
- If targets are not being detected, check to make sure that labels are defined correctly. Also note that if the port or path specific are invalid or do not represent a prometheus compatible endpoint, the target will not be detected. For eg. If we set the `ECS_PROMETHEUS_EXPORTER_PORT` to `8080` and the `ECS_PROMETHEUS_METRICS_PATH` to `/metrics`, the target will not be detected if the Prometheus endpoint does not expose metrics on port `8080` with the path `/metrics` path.

---

This template is meant to be a starting point for leveraging Open Telemetry for Prometheus Auto Discovery in ECS-EC2 using Docker Labels, it is also possible to detect prometheus endpoints using other patterns, such as Service Name.

The detection mechanism can also detect any prometheus compatible endpoint, not just the Prometheus Node Exporters. For eg. Cadvisor, Prometheus Push Gateway, etc.

For more information on additional configurations and features, please reference the [ecs_observer](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/extension/observer/ecsobserver) extension documentation or for more informatio on configuring Open Telemetry visit [here](https://opentelemetry.io/docs/)
