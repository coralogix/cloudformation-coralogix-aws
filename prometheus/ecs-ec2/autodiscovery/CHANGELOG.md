# Changelog

## prometheus ecs-ec2 autodiscovery

- [Update] Migrate Coralogix endpoints to the regional domain format `<region>.coralogix.com` (`coralogix.com`, `coralogix.us`, `coralogix.in`, `coralogixsg.com` -> `eu1`/`us1`/`ap1`/`ap2.coralogix.com`).
- [Update] `CoralogixRegion` now accepts the region codes `EU1`, `EU2`, `AP1`, `AP2`, `AP3`, `US1`, `US2`, `US3`, adding `AP3`, `US2` and `US3` support. The former names (`Europe`, `Europe2`, `India`, `Singapore`, `US`) are still accepted and map to the same endpoints.
- Removed default image parameter from autodiscovery workflow template
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->