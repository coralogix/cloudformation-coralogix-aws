# Changelog

## ecs-fargate
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->
### 0.0.10 / 2025-08-14
* [FIX] Changed telemetry.metrics to new syntax
* [UPDATE] spanmetric added to the config

### 0.0.9 / 2025-10-07
* [UPDATE] Change coralogix domains and endpoints to new format `<coralogix_region>.coralogix.com`.

### 0.0.8 / 2025-25-03
* [UPDATE] Added head sampling to reduce the initial data ingestion.

### 0.0.7 / 2024-12-02
* [FIX] Rename Parameter Store to prevent deployment failure.
* [UPDATE] Update and flatten configuration to reduce size.
* [UPDATE] Add Resource Catalog configs to parameter store.
* [UPDATE] Changed default parameter store name.

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
