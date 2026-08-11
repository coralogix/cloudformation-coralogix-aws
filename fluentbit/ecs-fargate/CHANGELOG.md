# Changelog

## fluentbit ecs-fargate
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 0.0.2 / 11 Aug 2026
* [Update] Migrate the hardcoded example ingress host to the regional domain format: `ingress.coralogix.com` -> `ingress.eu1.coralogix.com`.
* [Update] Move the example off the deprecated log-ingestion contract (end-of-life 30 September 2026): `URI` `/logs/rest/singles` -> `/logs/v1/singles`, and the auth header `private_key <key>` -> `Authorization Bearer <key>`.

### 0.0.1 / 16 Aug 2023
[Deprecate] Integration is deprecated. README.md updated accordingly.
[FIX] Updated Otel ECS-Fargate cloudformation template to remove unnecessary resources.