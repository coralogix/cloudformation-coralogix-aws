# Changelog

## ecs-fargate
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.0.7 / 2024-10-25
* [CHANGE] Specify grpc as appProtocol, needed to properly export traces.

### 0.0.6 / 2024-10-18
* [UPDATE] Update ecs-fargate integration cf to allow larger "Advanced" Parameter Store.
* [UPDATE] Adjust roles to true minimum requirements.

### 0.0.5 / 2024-09-11
### 🛑 Breaking changes 🛑
* [UPDATE] Update ecs-fargate integration cf template to OTEL only (remove fluentbit logrouter)

### 0.0.4 / 2024-09-02
* [UPDATE] Update domains to follow [coralogix-domain](coralogix.com/docs/coralogix-domain) and added AP3

### 0.0.3 / 2024-06-26
* [CHANGE] Migrate from ADOT to OTEL Collector Contrib

### 0.0.2 / 2023-11-23
* [CHANGE] Update the coralogix API
  
### 0.0.1 / 16 Aug 2023
* [Feature] Add example cloudformation template for ECS Fargate OTEL (ADOT) and fluentbit sidecar deployment
