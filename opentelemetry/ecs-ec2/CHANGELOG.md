# Changelog

## opentelemetry ecs-ec2
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.0.5 / 2023-10-25
* Added Healthcheck to default ecs-ec2 configuration


### 0.0.5 / 2023-10-23
* Updated default ecs-ec2 default templates 
    - removed  unnecessary OTEL_RESOURCE_ATTRIBUTES from default configuration  
    - updated default `ecsattributes` config to include `docker.name`
    - added resourcedetection for otel-collector metrics
    - removed unnecssary differences between default and metric configurations


### 0.0.4 / 2023-10-04
* Removed ecsattributes filters from default configuration

### 0.0.3 / 2023-09-29
* Add Otel Collector metrics to default configuration embedded in template.

### 0.0.2 / 2023-09-20
* Update cdot image description to advise users that a tag must be selected, latest is not supported

### 0.0.2 / 2023-09-11
* Added EC2 ECS Windows Example for metrics collection

### 0.0.1 / 2023--08-02
* Updated Otel ECS-EC2 cloudformation template. Added embedded support for logs, metrics and traces
* Updated the default Otel container image used to coralogixrepo/coralogix-otel-collector

### 0.0.2 / 2023-08-16
* Updated Otel ECS-EC2 to support US2 region.

### 0.0.2 / 2023-10-17
* Removed log.file.path from coralogix exporter subsystem_name_attributes, too many permutations.