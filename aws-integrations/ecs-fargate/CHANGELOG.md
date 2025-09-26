# Changelog

## ecs-fargate
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->
### 0.0.10 / 14 Aug 2025
* [FIX] Changed telemetry.metrics to new syntax
* [UPDATE] spanmetric added to the config

### 0.0.9 / 7 Oct 2025
* [UPDATE] Change coralogix domains and endpoints to new format `<coralogix_region>.coralogix.com`.

### 0.0.8 / 25 Mar 2025
* [UPDATE] Added head sampling to reduce the initial data ingestion.

### 0.0.7 / 2 Dec 2024
* [FIX] Rename Parameter Store to prevent deployment failure.
* [UPDATE] Update and flatten configuration to reduce size.
* [UPDATE] Add Resource Catalog configs to parameter store.
* [UPDATE] Changed default parameter store name.

### 0.0.6 / 18 Oct 2024
* [UPDATE] Update ecs-fargate integration cf to allow larger "Advanced" Parameter Store.
* [UPDATE] Adjust roles to true minimum requirements.

### 0.0.5 / 11 Sep 2024
### 🛑 Breaking changes 🛑
* [UPDATE] Update ecs-fargate integration cf template to OTEL only (remove fluentbit logrouter)

### 0.0.4 / 2 Sep 2024
* [UPDATE] Update domains to follow [coralogix-domain](coralogix.com/docs/coralogix-domain) and added AP3

### 0.0.3 / 26 Jun 2024
* [CHANGE] Migrate from ADOT to OTEL Collector Contrib

### 0.0.2 / 23 Nov 2023
* [CHANGE] Update the coralogix API
  
### 0.0.1 / 16 Aug 2023
* [Feature] Add example cloudformation template for ECS Fargate OTEL (ADOT) and fluentbit sidecar deployment
