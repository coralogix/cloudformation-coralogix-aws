SHELL := /bin/bash

ENV_FILE ?= .env

ifneq (,$(wildcard $(ENV_FILE)))
include $(ENV_FILE)
export
endif

TEMPLATE_FILE ?= template.yaml
CONFIG_FILE ?= examples/ebpf-profiler-config.yaml
SUPERVISOR_CONFIG_FILE ?= examples/ebpf-profiler-supervisor-config.yaml
STACK_NAME ?= ebpf-profiler-test
AWS_REGION ?= eu-north-1
S3_CONFIG_KEY ?= configs/ebpf-profiler-config.yaml
S3_SUPERVISOR_CONFIG_KEY ?= configs/ebpf-profiler-supervisor-config.yaml
SERVICE_NAME ?= coralogix-ebpf-profiler
PROFILER_IMAGE_MODE ?= collector
EBPF_PROFILER_IMAGE_REPOSITORY ?= coralogixrepo/coralogix-otel-collector
EBPF_PROFILER_IMAGE_VERSION ?= v0.5.8
SUPERVISED_IMAGE_REPOSITORY ?= cgx.jfrog.io/coralogix-docker-images/coralogix-otel-supervised-cdot
SUPERVISED_IMAGE_VERSION ?= v0.0.1

ECS_INSTANCE_PROFILE ?= ecsInstanceRole
ECS_EC2_AMI_SSM_PARAM ?= /aws/service/ecs/optimized-ami/amazon-linux-2023/recommended/image_id
ECS_EC2_AMI_ID ?=
EC2_INSTANCE_TYPE ?= t3.medium
EC2_SUBNET_ID ?=
EC2_SECURITY_GROUP_ID ?=
EC2_KEY_NAME ?=

# Required for deploy/smoke:
# CLUSTER_NAME, CORALOGIX_REGION, CORALOGIX_API_KEY, S3_BUCKET

.PHONY: help validate create-cluster create-instance-profile create-ec2-instance wait-cluster-capacity create-bucket upload-config prepare-stack deploy status tasks smoke delete delete-ec2-instance delete-cluster delete-bucket cleanup

help:
	@printf "eBPF profiler ECS CFN helper targets\n\n"
	@printf "Required vars: CLUSTER_NAME CORALOGIX_REGION CORALOGIX_API_KEY S3_BUCKET\n"
	@printf "Optional vars: STACK_NAME AWS_REGION S3_CONFIG_KEY S3_SUPERVISOR_CONFIG_KEY TEMPLATE_FILE CONFIG_FILE SUPERVISOR_CONFIG_FILE PROFILER_IMAGE_MODE EBPF_PROFILER_IMAGE_REPOSITORY EBPF_PROFILER_IMAGE_VERSION SUPERVISED_IMAGE_REPOSITORY SUPERVISED_IMAGE_VERSION ECS_INSTANCE_PROFILE ECS_EC2_AMI_ID ECS_EC2_AMI_SSM_PARAM EC2_INSTANCE_TYPE EC2_SUBNET_ID EC2_SECURITY_GROUP_ID EC2_KEY_NAME\n\n"
	@printf "If $(ENV_FILE) exists, values are auto-loaded from it.\n\n"
	@printf "Targets:\n"
	@printf "  make validate            Validate CloudFormation template\n"
	@printf "  make create-cluster      Create ECS cluster if missing\n"
	@printf "  make create-ec2-instance Create ECS EC2 instance if missing\n"
	@printf "  make create-bucket       Create S3 bucket if missing\n"
	@printf "  make upload-config       Upload collector config to S3\n"
	@printf "  make deploy              Deploy stack (with rollback-state recovery)\n"
	@printf "  make status              Show ECS service status\n"
	@printf "  make tasks               List ECS tasks for service\n"
	@printf "  make smoke               Full flow: validate + infra + deploy + checks\n"
	@printf "  make delete              Delete CloudFormation stack\n"
	@printf "  make delete-ec2-instance Delete EC2 instances created by this Makefile\n"
	@printf "  make delete-cluster      Delete ECS cluster if empty\n"
	@printf "  make delete-bucket       Delete S3 bucket (idempotent)\n"
	@printf "  make cleanup             delete-bucket + delete + delete-ec2-instance + delete-cluster\n"

require-deploy-vars:
	@test -n "$(CLUSTER_NAME)" || (echo "Missing required variable: CLUSTER_NAME" && exit 1)
	@test -n "$(CORALOGIX_REGION)" || (echo "Missing required variable: CORALOGIX_REGION" && exit 1)
	@test -n "$(CORALOGIX_API_KEY)" || (echo "Missing required variable: CORALOGIX_API_KEY" && exit 1)
	@test -n "$(S3_BUCKET)" || (echo "Missing required variable: S3_BUCKET" && exit 1)

require-cluster:
	@test -n "$(CLUSTER_NAME)" || (echo "Missing required variable: CLUSTER_NAME" && exit 1)

require-bucket:
	@test -n "$(S3_BUCKET)" || (echo "Missing required variable: S3_BUCKET" && exit 1)

validate:
	aws cloudformation validate-template --region "$(AWS_REGION)" --template-body "file://$(TEMPLATE_FILE)"

create-cluster: require-cluster
	@if aws ecs describe-clusters --region "$(AWS_REGION)" --clusters "$(CLUSTER_NAME)" --query "clusters[0].status" --output text | grep -q "^ACTIVE$$"; then \
		echo "ECS cluster already exists: $(CLUSTER_NAME)"; \
	else \
		echo "Creating ECS cluster: $(CLUSTER_NAME)"; \
		aws ecs create-cluster --region "$(AWS_REGION)" --cluster-name "$(CLUSTER_NAME)" >/dev/null; \
		echo "ECS cluster created: $(CLUSTER_NAME)"; \
	fi

create-instance-profile:
	@if aws iam get-instance-profile --instance-profile-name "$(ECS_INSTANCE_PROFILE)" >/dev/null 2>&1; then \
		echo "Instance profile already exists: $(ECS_INSTANCE_PROFILE)"; \
	else \
		echo "Creating IAM role/instance profile: $(ECS_INSTANCE_PROFILE)"; \
		if ! aws iam get-role --role-name "$(ECS_INSTANCE_PROFILE)" >/dev/null 2>&1; then \
			aws iam create-role --role-name "$(ECS_INSTANCE_PROFILE)" --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"ec2.amazonaws.com"},"Action":"sts:AssumeRole"}]}' >/dev/null; \
			aws iam attach-role-policy --role-name "$(ECS_INSTANCE_PROFILE)" --policy-arn "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role" >/dev/null; \
		fi; \
		aws iam create-instance-profile --instance-profile-name "$(ECS_INSTANCE_PROFILE)" >/dev/null; \
		aws iam add-role-to-instance-profile --instance-profile-name "$(ECS_INSTANCE_PROFILE)" --role-name "$(ECS_INSTANCE_PROFILE)" >/dev/null; \
		echo "Waiting for IAM propagation"; \
		sleep 10; \
	fi

create-ec2-instance: create-cluster create-instance-profile require-cluster
	@registered=$$(aws ecs describe-clusters --region "$(AWS_REGION)" --clusters "$(CLUSTER_NAME)" --query "clusters[0].registeredContainerInstancesCount" --output text); \
	if [ "$$registered" != "None" ] && [ "$$registered" != "0" ]; then \
		echo "ECS cluster already has $$registered container instance(s), skipping EC2 creation."; \
		exit 0; \
	fi; \
	if [ -n "$(ECS_EC2_AMI_ID)" ]; then \
		ami_id="$(ECS_EC2_AMI_ID)"; \
	else \
		ami_id=$$(aws ssm get-parameter --region "$(AWS_REGION)" --name "$(ECS_EC2_AMI_SSM_PARAM)" --query "Parameter.Value" --output text); \
	fi; \
	if [ -n "$(EC2_SUBNET_ID)" ]; then \
		subnet_id="$(EC2_SUBNET_ID)"; \
		vpc_id=$$(aws ec2 describe-subnets --region "$(AWS_REGION)" --subnet-ids "$$subnet_id" --query "Subnets[0].VpcId" --output text); \
	else \
		vpc_id=$$(aws ec2 describe-vpcs --region "$(AWS_REGION)" --filters Name=isDefault,Values=true --query "Vpcs[0].VpcId" --output text); \
		subnet_id=$$(aws ec2 describe-subnets --region "$(AWS_REGION)" --filters Name=vpc-id,Values=$$vpc_id Name=default-for-az,Values=true --query "Subnets[0].SubnetId" --output text); \
	fi; \
	if [ "$$subnet_id" = "None" ] || [ -z "$$subnet_id" ]; then \
		echo "Could not resolve subnet. Set EC2_SUBNET_ID."; \
		exit 1; \
	fi; \
	if [ -n "$(EC2_SECURITY_GROUP_ID)" ]; then \
		sg_id="$(EC2_SECURITY_GROUP_ID)"; \
	else \
		sg_id=$$(aws ec2 describe-security-groups --region "$(AWS_REGION)" --filters Name=vpc-id,Values=$$vpc_id Name=group-name,Values=default --query "SecurityGroups[0].GroupId" --output text); \
	fi; \
	if [ "$$sg_id" = "None" ] || [ -z "$$sg_id" ]; then \
		echo "Could not resolve security group. Set EC2_SECURITY_GROUP_ID."; \
		exit 1; \
	fi; \
	user_data=$$(printf '#!/bin/bash\ncat <<EOF >> /etc/ecs/ecs.config\nECS_CLUSTER=%s\nECS_ENABLE_TASK_IAM_ROLE=true\nECS_ENABLE_TASK_IAM_ROLE_NETWORK_HOST=true\nEOF\n' "$(CLUSTER_NAME)"); \
	echo "Launching ECS EC2 instance into cluster $(CLUSTER_NAME)"; \
	if [ -n "$(EC2_KEY_NAME)" ]; then \
		instance_id=$$(aws ec2 run-instances --region "$(AWS_REGION)" --image-id "$$ami_id" --instance-type "$(EC2_INSTANCE_TYPE)" --iam-instance-profile Name="$(ECS_INSTANCE_PROFILE)" --subnet-id "$$subnet_id" --security-group-ids "$$sg_id" --key-name "$(EC2_KEY_NAME)" --user-data "$$user_data" --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$(CLUSTER_NAME)-ecs-container-instance},{Key=ManagedBy,Value=ebpf-profiler-makefile},{Key=ClusterName,Value=$(CLUSTER_NAME)}]" --query "Instances[0].InstanceId" --output text); \
	else \
		instance_id=$$(aws ec2 run-instances --region "$(AWS_REGION)" --image-id "$$ami_id" --instance-type "$(EC2_INSTANCE_TYPE)" --iam-instance-profile Name="$(ECS_INSTANCE_PROFILE)" --subnet-id "$$subnet_id" --security-group-ids "$$sg_id" --user-data "$$user_data" --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$(CLUSTER_NAME)-ecs-container-instance},{Key=ManagedBy,Value=ebpf-profiler-makefile},{Key=ClusterName,Value=$(CLUSTER_NAME)}]" --query "Instances[0].InstanceId" --output text); \
	fi; \
	echo "Created EC2 instance: $$instance_id"; \
	aws ec2 wait instance-running --region "$(AWS_REGION)" --instance-ids "$$instance_id"

wait-cluster-capacity: require-cluster
	@echo "Waiting for ECS container instance registration in cluster $(CLUSTER_NAME)"; \
	for i in $$(seq 1 30); do \
		registered=$$(aws ecs describe-clusters --region "$(AWS_REGION)" --clusters "$(CLUSTER_NAME)" --query "clusters[0].registeredContainerInstancesCount" --output text); \
		if [ "$$registered" != "None" ] && [ "$$registered" != "0" ]; then \
			echo "ECS container instances registered: $$registered"; \
			exit 0; \
		fi; \
		sleep 10; \
	done; \
	echo "Timed out waiting for ECS container instance registration."; \
	exit 1

create-bucket: require-bucket
	@if aws s3api head-bucket --bucket "$(S3_BUCKET)" >/dev/null 2>&1; then \
		echo "S3 bucket already exists: $(S3_BUCKET)"; \
	else \
		echo "Creating S3 bucket: $(S3_BUCKET) in region $(AWS_REGION)"; \
		if [ "$(AWS_REGION)" = "us-east-1" ]; then \
			aws s3api create-bucket --bucket "$(S3_BUCKET)" --region "$(AWS_REGION)"; \
		else \
			aws s3api create-bucket --bucket "$(S3_BUCKET)" --region "$(AWS_REGION)" --create-bucket-configuration LocationConstraint="$(AWS_REGION)"; \
		fi; \
	fi

upload-config: require-bucket
	aws s3 cp "$(CONFIG_FILE)" "s3://$(S3_BUCKET)/$(S3_CONFIG_KEY)" --region "$(AWS_REGION)"
	aws s3 cp "$(SUPERVISOR_CONFIG_FILE)" "s3://$(S3_BUCKET)/$(S3_SUPERVISOR_CONFIG_KEY)" --region "$(AWS_REGION)"

prepare-stack:
	@status=$$(aws cloudformation describe-stacks --region "$(AWS_REGION)" --stack-name "$(STACK_NAME)" --query "Stacks[0].StackStatus" --output text 2>/dev/null || echo "NOT_FOUND"); \
	case "$$status" in \
		NOT_FOUND) echo "Stack not found, will create: $(STACK_NAME)" ;; \
		ROLLBACK_COMPLETE|ROLLBACK_FAILED|CREATE_FAILED) \
			echo "Stack $(STACK_NAME) is in $$status, deleting before deploy"; \
			aws cloudformation delete-stack --region "$(AWS_REGION)" --stack-name "$(STACK_NAME)"; \
			aws cloudformation wait stack-delete-complete --region "$(AWS_REGION)" --stack-name "$(STACK_NAME)" ;; \
		DELETE_IN_PROGRESS) \
			echo "Stack deletion in progress, waiting: $(STACK_NAME)"; \
			aws cloudformation wait stack-delete-complete --region "$(AWS_REGION)" --stack-name "$(STACK_NAME)" ;; \
		CREATE_IN_PROGRESS|UPDATE_IN_PROGRESS|UPDATE_COMPLETE_CLEANUP_IN_PROGRESS|UPDATE_ROLLBACK_IN_PROGRESS|UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS|REVIEW_IN_PROGRESS) \
			echo "Stack is currently busy ($$status). Wait for completion and retry."; \
			exit 1 ;; \
		*) echo "Stack status is $$status, proceeding with deploy" ;; \
	esac

deploy: require-deploy-vars prepare-stack
	aws cloudformation deploy \
		--template-file "$(TEMPLATE_FILE)" \
		--stack-name "$(STACK_NAME)" \
		--region "$(AWS_REGION)" \
		--capabilities CAPABILITY_NAMED_IAM \
		--parameter-overrides \
			ClusterName="$(CLUSTER_NAME)" \
			CoralogixRegion="$(CORALOGIX_REGION)" \
			CoralogixApiKey="$(CORALOGIX_API_KEY)" \
			S3ConfigBucket="$(S3_BUCKET)" \
			S3ConfigKey="$(S3_CONFIG_KEY)" \
			S3SupervisorConfigKey="$(S3_SUPERVISOR_CONFIG_KEY)" \
			ProfilerImageMode="$(PROFILER_IMAGE_MODE)" \
			EbpfProfilerImageRepository="$(EBPF_PROFILER_IMAGE_REPOSITORY)" \
			EbpfProfilerImageVersion="$(EBPF_PROFILER_IMAGE_VERSION)" \
			SupervisedImageRepository="$(SUPERVISED_IMAGE_REPOSITORY)" \
			SupervisedImageVersion="$(SUPERVISED_IMAGE_VERSION)"

status:
	aws ecs describe-services \
		--region "$(AWS_REGION)" \
		--cluster "$(CLUSTER_NAME)" \
		--services "$(SERVICE_NAME)" \
		--query "services[0].{status:status,desired:desiredCount,running:runningCount,pending:pendingCount,events:events[0:5].message}"

tasks:
	aws ecs list-tasks \
		--region "$(AWS_REGION)" \
		--cluster "$(CLUSTER_NAME)" \
		--service-name "$(SERVICE_NAME)"

smoke: validate create-cluster create-ec2-instance wait-cluster-capacity create-bucket upload-config deploy status tasks

delete:
	@status=$$(aws cloudformation describe-stacks --region "$(AWS_REGION)" --stack-name "$(STACK_NAME)" --query "Stacks[0].StackStatus" --output text 2>/dev/null || echo "NOT_FOUND"); \
	if [ "$$status" = "NOT_FOUND" ]; then \
		echo "Stack does not exist, skipping: $(STACK_NAME)"; \
	else \
		echo "Deleting stack: $(STACK_NAME)"; \
		aws cloudformation delete-stack --stack-name "$(STACK_NAME)" --region "$(AWS_REGION)"; \
		echo "Waiting for stack deletion: $(STACK_NAME)"; \
		aws cloudformation wait stack-delete-complete --stack-name "$(STACK_NAME)" --region "$(AWS_REGION)"; \
		echo "Stack deleted: $(STACK_NAME)"; \
	fi

delete-ec2-instance: require-cluster
	@instance_ids=$$(aws ec2 describe-instances --region "$(AWS_REGION)" --filters Name=tag:ManagedBy,Values=ebpf-profiler-makefile Name=tag:ClusterName,Values="$(CLUSTER_NAME)" Name=instance-state-name,Values=pending,running,stopping,stopped --query "Reservations[].Instances[].InstanceId" --output text); \
	if [ -z "$$instance_ids" ] || [ "$$instance_ids" = "None" ]; then \
		echo "No managed EC2 instances found for cluster $(CLUSTER_NAME)."; \
		exit 0; \
	fi; \
	echo "Terminating EC2 instance(s): $$instance_ids"; \
	aws ec2 terminate-instances --region "$(AWS_REGION)" --instance-ids $$instance_ids >/dev/null; \
	aws ec2 wait instance-terminated --region "$(AWS_REGION)" --instance-ids $$instance_ids

delete-cluster: require-cluster
	@status=$$(aws ecs describe-clusters --region "$(AWS_REGION)" --clusters "$(CLUSTER_NAME)" --query "clusters[0].status" --output text); \
	if [ "$$status" != "ACTIVE" ]; then \
		echo "ECS cluster does not exist or is not active, skipping: $(CLUSTER_NAME)"; \
		exit 0; \
	fi; \
	services=$$(aws ecs list-services --region "$(AWS_REGION)" --cluster "$(CLUSTER_NAME)" --query "length(serviceArns)" --output text); \
	tasks=$$(aws ecs list-tasks --region "$(AWS_REGION)" --cluster "$(CLUSTER_NAME)" --query "length(taskArns)" --output text); \
	instances=$$(aws ecs list-container-instances --region "$(AWS_REGION)" --cluster "$(CLUSTER_NAME)" --query "length(containerInstanceArns)" --output text); \
	if [ "$$services" != "0" ] || [ "$$tasks" != "0" ] || [ "$$instances" != "0" ]; then \
		echo "ECS cluster is not empty. services=$$services tasks=$$tasks instances=$$instances"; \
		exit 1; \
	fi; \
	echo "Deleting ECS cluster: $(CLUSTER_NAME)"; \
	aws ecs delete-cluster --region "$(AWS_REGION)" --cluster "$(CLUSTER_NAME)" >/dev/null

delete-bucket: require-bucket
	@if aws s3api head-bucket --bucket "$(S3_BUCKET)" >/dev/null 2>&1; then \
		echo "Removing objects from s3://$(S3_BUCKET)"; \
		aws s3 rm "s3://$(S3_BUCKET)" --recursive --region "$(AWS_REGION)" >/dev/null 2>&1 || true; \
		echo "Deleting S3 bucket $(S3_BUCKET)"; \
		if aws s3api delete-bucket --bucket "$(S3_BUCKET)" --region "$(AWS_REGION)" >/dev/null 2>&1; then \
			echo "S3 bucket deleted: $(S3_BUCKET)"; \
		else \
			if aws s3api head-bucket --bucket "$(S3_BUCKET)" >/dev/null 2>&1; then \
				echo "Failed to delete S3 bucket: $(S3_BUCKET)"; \
				exit 1; \
			else \
				echo "S3 bucket already deleted: $(S3_BUCKET)"; \
			fi; \
		fi; \
	else \
		echo "S3 bucket does not exist or is not accessible, skipping: $(S3_BUCKET)"; \
	fi

cleanup: delete delete-ec2-instance delete-cluster delete-bucket
	@echo "Done"
