# Coralogix ECS Otel Components:

- [ecsattributes](#The-ecsattributes-Processor)
- [awsecscontainermetricsd](#The-awsecscontainermetricsd-Receiver)
---


# AWS ECS Container Metrics Daemonset Receiver

## Overview

AWS ECS Container Metrics Daemoonset Receiver (`awsecscontainermetricsd`) is based on the [awsecscontainermetrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver) receiver. It uses the Docker API to identify all the available [Amazon ECS Task Metadata Endpoints](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint.html) from all containers and reads task metadata and [docker stats](https://docs.docker.com/engine/api/v1.30/#operation/ContainerStats), and generates resource usage metrics (such as CPU, memory, network, and disk) from them. To get the full list of metrics, see the [Available Metrics](#available-metrics) section below.

This receiver works only for [ECS Task Metadata Endpoint V4](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint-v4.html). Amazon ECS tasks on Amazon EC2 that are running at least version 1.39.0 of the Amazon ECS container agent can utilize this receiver. For more information, see [Amazon ECS Container Agent Versions](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-agent-versions.html).

### Pre-requisites
- Privileged mode must be enabled for the container running the collector
- The `docker.sock` must be mounted to the container running the collector at `/var/run/docker.sock`

## Configuration

Example:

```yaml
receivers:
  awsecscontainermetricsd:
    collection_interval: 20s
    sidecar: true
```

#### collection_interval:

This receiver collects task metadata and container stats at a fixed interval and emits metrics to the next consumer of OpenTelemetry pipeline. `collection_interval` will determine the frequency at which metrics are collected and emitted by this receiver.

default: `20s`

#### sidecar:

If true, the receiver will run in sidecar mode. In this mode the receiver will collect metrics from all containers within the same task definition in the same way the [awsecscontainermetrics](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/receiver/awsecscontainermetricsreceiver) receiver does. This mode is required for ECS EC2 Windows Server environments.

default: `false`

## Enabling the AWS ECS Container Metrics Receiver

To enable the awsecscontainermetricsd receiver, add the name under receiver section in the OpenTelemetry config file. By default, the receiver scrapes the ECS task metadata endpoints every 20s and collects all metrics (For the full list of metrics, see [Available Metrics](#available-metrics)).

The following configuration collects AWS ECS resource usage metrics by using `awsecscontainermetricsd` receiver and sends them to the Coralogix exporter.

```yaml
receivers:
  awsecscontainermetricsd:
exporters:
  coralogix:
    domain: coralogix.com
    private_key: '<your-product-key>'
    application_name: "opentelemetry"
    subsystem_name: "ecs"
    application_name_attributes:
    - "APP_NAME"
    subsystem_name_attributes:
    - "SUB_SYS"
    - "aws.ecs.task.family"
    timeout: 30s

service:
  pipelines:
      metrics:
          receivers: [awsecscontainermetricsd]
          exporters: [coralogix]
```

## Set Metrics Collection Interval

Customers can configure `collection_interval` under `awsecscontainermetricsd` receiver to scrape and gather metrics at a specific interval. The following example configuration will collect metrics every 40 seconds.

```yaml
receivers:
  awsecscontainermetricsd:
      collection_interval: 40s
exporters:
  logging:
    verbosity: detailed

service:
  pipelines:
      metrics:
          receivers: [awsecscontainermetrics]
          exporters: [logging]
```

## Collect specific metrics and update metric names

The previous configurations collect all the metrics and sends them to Amazon CloudWatch using default names. Customers can use `filter` and `metrictransform` processors to send specific metrics and rename them respectively.

The following configuration example collects only the `ecs.task.memory.utilized` metric and renames it to `MemoryUtilized` before sending to CloudWatch.

```yaml
receivers:
  awsecscontainermetricsd:
exporters:
  logging:
    verbosity: detailed
processors:
  filter:
    metrics:
      include:
        match_type: strict
        metric_names:
          - ecs.task.memory.utilized

  metricstransform:
    transforms:
      - include: ecs.task.memory.utilized
        action: update
        new_name: MemoryUtilized

service:
  pipelines:
      metrics:
          receivers: [awsecscontainermetricsd]
          processors: [filter, metricstransform]
          exporters: [logging]
```

## Available Metrics

Following is the full list of metrics emitted by this receiver. Note that all metrics are support on Linux, however there are some that will not work on Windows. See the **Supported in Windows** column in the table below.

| Task Level Metrics                   | Container Level Metrics               | Unit         | Supported in Windows |
|--------------------------------------|---------------------------------------|--------------|----------------------|
| ecs.task.memory.usage                | container.memory.usage                | Bytes        | :heavy_check_mark:   |
| ecs.task.memory.usage.max            | container.memory.usage.max            | Bytes        | :heavy_check_mark:   |
| ecs.task.memory.usage.limit          | container.memory.usage.limit          | Bytes        |                      |
| ecs.task.memory.reserved             | container.memory.reserved             | Megabytes    |                      |
| ecs.task.memory.utilized             | container.memory.utilized             | Megabytes    | :heavy_check_mark:   |
| ecs.task.cpu.usage.total             | container.cpu.usage.total             | Nanoseconds  | :heavy_check_mark:   |
| ecs.task.cpu.usage.kernelmode        | container.cpu.usage.kernelmode        | Nanoseconds  | :heavy_check_mark:   |
| ecs.task.cpu.usage.usermode          | container.cpu.usage.usermode          | Nanoseconds  | :heavy_check_mark:   |
| ecs.task.cpu.usage.system            | container.cpu.usage.system            | Nanoseconds  |                      |
| ecs.task.cpu.usage.vcpu              | container.cpu.usage.vcpu              | vCPU         | :heavy_check_mark:   |
| ecs.task.cpu.cores                   | container.cpu.cores                   | Count        |                      |
| ecs.task.cpu.onlines                 | container.cpu.onlines                 | Count        |                      |
| ecs.task.cpu.reserved                | container.cpu.reserved                | vCPU         | :heavy_check_mark:   |
| ecs.task.cpu.utilized                | container.cpu.utilized                | Percent      | :heavy_check_mark:   |
| ecs.task.network.rate.rx             | container.network.rate.rx             | Bytes/Second | :heavy_check_mark:   |
| ecs.task.network.rate.tx             | container.network.rate.tx             | Bytes/Second | :heavy_check_mark:   |
| ecs.task.network.io.usage.rx_bytes   | container.network.io.usage.rx_bytes   | Bytes        | :heavy_check_mark:   |
| ecs.task.network.io.usage.rx_packets | container.network.io.usage.rx_packets | Count        | :heavy_check_mark:   |
| ecs.task.network.io.usage.rx_errors  | container.network.io.usage.rx_errors  | Count        | :heavy_check_mark:   |
| ecs.task.network.io.usage.rx_dropped | container.network.io.usage.rx_dropped | Count        | :heavy_check_mark:   |
| ecs.task.network.io.usage.tx_bytes   | container.network.io.usage.tx_bytes   | Bytes        | :heavy_check_mark:   |
| ecs.task.network.io.usage.tx_packets | container.network.io.usage.tx_packets | Count        | :heavy_check_mark:   |
| ecs.task.network.io.usage.tx_errors  | container.network.io.usage.tx_errors  | Count        | :heavy_check_mark:   |
| ecs.task.network.io.usage.tx_dropped | container.network.io.usage.tx_dropped | Count        | :heavy_check_mark:   |
| ecs.task.storage.read_bytes          | container.storage.read_bytes          | Bytes        |                      |
| ecs.task.storage.write_bytes         | container.storage.write_bytes         | Bytes        |                      |

## Resource Attributes and Metrics Labels

Metrics emitted by this receiver comes with a set of resource attributes. These resource attributes can be converted to metrics labels using appropriate processors/exporters (See `Full Configuration Examples` section below). Finally, these metrics labels can be set as metrics dimensions while exporting to desired destinations. Check the following table to see available resource attributes for Task and Container level metrics. Container level metrics have three additional attributes than task level metrics.

| Resource Attributes for Task Level Metrics | Resource Attributes for Container Level Metrics |
|--------------------------------------------|-------------------------------------------------|
| aws.ecs.cluster.name                       | aws.ecs.cluster.name                            |
| aws.ecs.task.family                        | aws.ecs.task.family                             |
| aws.ecs.task.arn                           | aws.ecs.task.arn                                |
| aws.ecs.task.id                            | aws.ecs.task.id                                 |
| aws.ecs.task.revision                      | aws.ecs.task.revision                           |
| aws.ecs.service.name                       | aws.ecs.service.name                            |
| cloud.availability_zone                    | cloud.availability_zone                         |
| cloud.account.id                           | cloud.account.id                                |
| cloud.region                               | cloud.region                                    |
| aws.ecs.task.pull_started_at               | aws.ecs.container.started_at                    |
| aws.ecs.task.pull_stopped_at               | aws.ecs.container.finished_at                   |
| aws.ecs.task.known_status                  | aws.ecs.container.know_status                   |
| aws.ecs.launch_type                        | aws.ecs.launch_type                             |
| &nbsp;                                     | aws.ecs.container.created_at                    |
| &nbsp;                                     | container.name                                  |
| &nbsp;                                     | container.id                                    |
| &nbsp;                                     | aws.ecs.docker.name                             |
| &nbsp;                                     | container.image.tag                             |
| &nbsp;                                     | aws.ecs.container.image.id                      |
| &nbsp;                                     | aws.ecs.container.exit_code                     |

## Full Configuration Examples

This receiver emits 52 unique metrics. Customer may not want to send all of them to destinations. Following sections will show full configuration file for filtering and transforming existing metrics with different processors/exporters.

The example shows a full configuration to get most useful task level metrics. It uses `awsecscontainermetricsd` receiver to collect all the resource usage metrics from ECS task metadata endpoint. It applies `filter` processor to select only 8 task-level metrics and update metric names using `metricstransform` processor. It also renames the resource attributes using `resource` processor which will be used as metric dimensions in the Coralogix.

**Note:** AWS OpenTelemetry Collector has a [default configuration](https://github.com/aws-observability/aws-otel-collector/blob/main/config/ecs/container-insights/otel-task-metrics-config.yaml) backed into it for Container Insights experience which is smiliar to this one. Follow our [setup](https://aws-otel.github.io/docs/setup/ecs) doc to check how to use that default config.

```yaml
receivers:
  awsecscontainermetricsd: # collect 52 metrics

processors:
  filter: # filter metrics
    metrics:
      include:
        match_type: strict
        metric_names: # select only 8 task level metrics out of 52
          - ecs.task.memory.reserved
          - ecs.task.memory.utilized
          - ecs.task.cpu.reserved
          - ecs.task.cpu.utilized
          - ecs.task.network.rate.rx
          - ecs.task.network.rate.tx
          - ecs.task.storage.read_bytes
          - ecs.task.storage.write_bytes
  metricstransform: # update metric names
    transforms:
      - include: ecs.task.memory.utilized
        action: update
        new_name: MemoryUtilized
      - include: ecs.task.memory.reserved
        action: update
        new_name: MemoryReserved
      - include: ecs.task.cpu.utilized
        action: update
        new_name: CpuUtilized
      - include: ecs.task.cpu.reserved
        action: update
        new_name: CpuReserved
      - include: ecs.task.network.rate.rx
        action: update
        new_name: NetworkRxBytes
      - include: ecs.task.network.rate.tx
        action: update
        new_name: NetworkTxBytes
      - include: ecs.task.storage.read_bytes
        action: update
        new_name: StorageReadBytes
      - include: ecs.task.storage.write_bytes
        action: update
        new_name: StorageWriteBytes
  resource:
    attributes: # rename resource attributes which will be used as dimensions
      - key: ClusterName
        from_attribute: aws.ecs.cluster.name
        action: insert
      - key: aws.ecs.cluster.name
        action: delete
      - key: ServiceName
        from_attribute: aws.ecs.service.name
        action: insert
      - key: aws.ecs.service.name
        action: delete
      - key: TaskId
        from_attribute: aws.ecs.task.id
        action: insert
      - key: aws.ecs.task.id
        action: delete
      - key: TaskDefinitionFamily
        from_attribute: aws.ecs.task.family
        action: insert
      - key: aws.ecs.task.family
        action: delete
exporters:
  coralogix:
    domain: coralogix.com
    private_key: '<your-product-key>'
    application_name: "opentelemetry"
    subsystem_name: "ecs"
    application_name_attributes:
    - "ServiceName"
    subsystem_name_attributes:
    - "TaskId"
    timeout: 30s
service:
  pipelines:
    metrics:
      receivers: [awsecscontainermetricsd ]
      processors: [filter, metricstransform, resource]
      exporters: [ coralogix ]
```

## Reference
1. [Setup OpenTelemetry Collector on Amazon ECS](https://aws-otel.github.io/docs/setup/ecs)
2. [Getting Started with ECS Container Metrics Receiver in the OpenTelemetry Collector](https://aws-otel.github.io/docs/components/ecs-metrics-receiver) g


---

# The ecsattributes Processor

---

| Status    |             |                       |
|-----------|-------------|-----------------------|
| Stability | alpha: logs | WEP: metrics & traces |

The coralogixrepo/otel-coralogix-ecs-ec2 docker image includes an Open Telemetry distribution with a dedicated processor designed to handle metadata enrichment for logs collected at the Host level. This processor enables the collector to discover metadata endpoints for all active containers on an instance, utilizing container IDs to indentify metadata endpoints to enrich logs and establish correlations. It's important to note that the default resourcedetection processor does not offer this specific functionality.

### Pre-requisites
- Privileged mode must be enabled for the container running the collector
- The `docker.sock` must be mounted to the container running the collector at `/var/run/docker.sock`
- This processor uses the container ID to identify the correct metadata endpoint for each container. The processor checks for the container ID in **resource** attribute(s) specified during configuration. If no container ID can be determined, no metadata will be added.

### Attributes

| Attribute                       | Value                                                                               | Default |
|---------------------------------|-------------------------------------------------------------------------------------|---------|
| aws.ecs.task.definition.family  | The ECS task defintion family                                                       | ✔️       |
| aws.ecs.task.definition.version | The ECS task defintion version                                                      | ✔️       |
| image                           | The container image                                                                 | ✔️       |
| aws.ecs.container.name          | The name of the running container. The name given to the container by the ECS Agent | ✔️       |
| aws.ecs.container.arn           | The ECS instance ARN                                                                | ✔️       |
| aws.ecs.cluster                 | The ECS cluster name                                                                | ✔️       |
| aws.ecs.task.arn                | The ECS task ARN                                                                    | ✔️       |
| image.id                        | The image ID of the running container                                               | ✔️       |
| docker.name                     | The name of the running container. The is name you will see if you run `docker ps`  | ✔️       |
| docker.id                       | The docker container ID                                                             | ✔️       |
| name                            | Same as `ecs.container.name`                                                        | ✔️       |
| limits.cpu                      | The CPU limit of the container                                                      |         |
| limits.memory                   | The memory limit of the container                                                   |         |
| type                            | The type of the container                                                           |         |
| aws.ecs.known.status            | The lifecycle state of the container                                                |         |
| created.at                      | The time the container was created                                                  |         |
| `networks.*.ipv4.addresses.*`   | An expression that matches the IP address(s) assigned to a container                |         |
| `networks.*.network.mode`       | An expression that matches the network mode(s) associated with the container        |         |
| labels.*                        | An expression that matches the docker labels associated with the container          |         |

Only containers with a valid ECS metadata endpoint will have attributes assigned, all others will be ignored.

To verify your container has a valid ECS metadata endpoint, you can check for the following environment variables in the your running container:

- ECS_CONTAINER_METADATA_URI
- ECS_CONTAINER_METADATA_URI_V4

Atleast one must be present.

### Configuration

The ecsattributes processor is enabled by adding the keyword `ecsattributes` to the `processors` section of the configuration file. The processor can be configured using the following options:

| Config               | Description                                                                                                                                                                      |
|----------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| attributes           | A list of regex patterns that match specific or multiple attribute keys.                                                                                                         |
| container_id.sources | The **resource** attribute key that contains the container ID. Defaults to `container.id`. If multiple attribute keys are provided, the first none-empty value will be selected. |

Note, given a `log.file.name=<container.id>-json.log`, the `ecsattributesprocessor` will automatically remove the `-json.log` suffix from the container ID when correlating metadata.

The following config, will collect all the [default attributes](#attributes).

```yaml
processors:
  ecsattributes:

  # check for container id in the following attributes:
  container_id:
    sources:
      - "container.id"
      - "log.file.name"
```

You can specify which attributes should be collected by using the `attributes` option which represents a list of regex patterns that match specific or multiple attribute keys.

```yaml
processors:
  ecsattributes:
    attributes:
      - '^aws.ecs.*' # all attributes that start with ecs
      - '^docker.*' # all attributes that start with docker
      - '^image.*|^network.*' # all attributes that start with image or network
```

### Important

- The `ecsattributesprocessor` uses the Docker API to collect metadata for each container. It refreshes the metadata every 60 seconds, everytime a new container is detected via Docker events and each time a log is detected without metadata. If the processor is unable to collect metadata for a container or if there are errors during the refresh process, the processor will log the error and continue processing the next log record. **It will not halt/crash the open telemetry process**. If you notice metadata not being added to your logs, please check the logs for the collector for any errors related to the `ecsattributesprocessor`.

- If logs are received with no attributes, it is possible that these logs are from the ECS Agent; the agent responsible for managing containers on an ECS Node. This container does not have a metadata endpoint. Also, logs from rogue containers that are run on ECS outside the control of the ECS Agent will not be assigned a metadata endpoint and will not have attributes added.

**TODO:**
- Implement a configuration option that allows the user to specify what action to take on error. For example, `continue` or `halt`.
