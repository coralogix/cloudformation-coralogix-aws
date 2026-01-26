# Changelog

## opentelemetry ecs-ec2-tail-sampling
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->

### 1.1.0 / 2026-01-25
* [Security] Added TaskRoleArn parameter to separate execution and task IAM roles, following principle of least privilege
* [Security] Auto-create minimal task role with S3 read and CloudMap discovery permissions when TaskRoleArn is not provided
* [Update] Standardized all templates and examples to use CORALOGIX_PRIVATE_KEY environment variable instead of PRIVATE_KEY
* [Update] Added CloudMap discovery permissions (servicediscovery:DiscoverInstances, servicediscovery:ListInstances, servicediscovery:ListServices, servicediscovery:ListNamespaces) to auto-created task roles for tail-sampling deployments

### 1.0.0 / 2025-08-21
* [Added] Initial release of ECS EC2 Tail Sampling CloudFormation templates
* [Added] Example folder with example config files
