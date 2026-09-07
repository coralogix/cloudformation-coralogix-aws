# Changelog

## eventbridge
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.0.2 / 11 Aug 2026
* [Update] Migrate Coralogix endpoints to the regional domain format `<region>.coralogix.com` (`coralogix.com`, `coralogix.us`, `coralogix.in`, `coralogixsg.com`, `cx498.coralogix.com` -> `eu1`/`us1`/`ap1`/`ap2`/`us2.coralogix.com`).
* [Bug fix] `AP1` and `AP2` were mapped to each other's ingress endpoint (`AP1` -> Singapore, `AP2` -> Mumbai). `AP1` now resolves to `ingress.ap1.coralogix.com` (Mumbai) and `AP2` to `ingress.ap2.coralogix.com` (Singapore).
* [Feature] Add `US3` (`us3.coralogix.com`) as a supported `CoralogixRegion`.

### 0.0.1 / 02 Sept 2024
### 🛑 Breaking changes 🛑
* [Update] Update ingress domain from `aws-events.<domain.com>/aws/events` to `ingress.<domain.com>/aws/event-bridge`
* [Update] Update regions to follow [coralogix.com/docs/coralogix-domain] and added AP3 domain
