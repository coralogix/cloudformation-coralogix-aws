# Fluentbit ECS Fargate

This template provides a starting point for deploying Fluentbit as a sidecar container in an ECS Fargate task. The sidecar approach to log collection works by running a separate container in the same ecs task as the application container. The sidecar container is responsible for collecting logs from the application container and sending them to a destination. This approach is useful when you want to collect logs from an application container that does not support logging to a file or stdout.

### How to add Fluentbit as a sidecar container in an ECS Fargate task

ECS Fargate allows us to add a fluentbit container as a sidecar, the term sidecar here simply means that the container is running in the same ECS task as the application container.

AWS supports adding Fluentbit as a sidecar using a FireLens configuration.

In order to configure Fluentbit as a sidecar container in an ECS Fargate task, you need to perform the steps below:

### Fluentbit Configuration

When using the Fluentbit for AWS container image, you need to embed your configuration file in the Image itself. This is because the Fluentbit for AWS image already has a default configuration that is required for the integration to work. When we embed our configuration file in the image, it runs alongside the default configuration.

Once your fluentbit image is created it needs to be pushed to ECR or any other publically accessible container registry.

Create a configuration file. For example, extra.conf:

```
[SERVICE]
    Flush 1
    Grace 30
    log_level debug
    Parsers_File /fluent-bit/parsers/parsers.conf

[FILTER]
    Name        nest
    Match       *
    Operation   nest
    Wildcard    *
    Nest_under  json

[FILTER]
    Name    modify
    Match   *
    Add    applicationName ${APP_NAME}

# retrieve subsystem name from ENV if container_id is not defined.
# using in-line lua script
[FILTER]
    Name        lua
    Match       *
    call        subsystemNameFromEnv
    code        function subsystemNameFromEnv(tag, timestamp, record) record["subsystemName"] = record["json"]["container_id"] or os.getenv("SUB_SYSTEM") return 2, timestamp, record end
```

The configuration above will capture the ecs container_id and set it as the Coralogix subsystemName. If no `container_id` is found it will use the SUB_SYSTEM environment variable.

- Create a Dockerfile

```Dockerfile
FROM amazon/aws-for-fluent-bit:latest
COPY extra.conf /extra.conf
```

- Build the image and push it to ECR

```sh
docker build -t <your image + tag> .
ecs-cli push <your image + tag> --region <region>
```

Note that we're using ecs-cli to push the image to ECR, this is one of the simplest ways to push to ECR however, you can use any other method you prefer.

### ECS Task Definition (JSON)

Next we need to add Log Configuration to our application task definition. The following example is a JSON representation of a simple task definition with a single container plus the fluentbit container.

```json
{
    "Family": "my-fluentbit-task",
    "RequiresCompatibilities": ["FARGATE"],
    "NetworkMode": "awsvpc",
    "Cpu": "256",
    "Memory": "512",
    "ContainerDefinitions": [
      {
        "Name": "<YOUR APPLICATION NAME>",
        "Image": "<YOUR APPLICATION IMAGE>",
        "Essential": true,
        // Log Configuration for the application container
        "LogConfiguration": {
          "LogDriver": "awsfirelens",
          "Options": {
            "Format": "json_lines",
            "Header": "private_key <YOUR PRIVATE KEY>",
            "Retry_Limit": "10",
            "compress": "gzip",
            "Port": "443",

            // your coralogix domain: 
            "Host": "ingress.coralogix.com",
            "TLS": "On",
            "URI": "/logs/rest/singles",
            "Name": "http"
          }
        }
      },

    // Fluentbit container
      {
        "Name": "log_router",
        "Image": "<YOUR FLUENTBIT IMAGE>",
        "Essential": true,
        "User": "0",
        "Environment": [
          {
            "Name": "APP_NAME",
            "Value": "testfluent"
          },
          {
            "Name": "SUB_SYSTEM",
            "Value": "aaaa"
          }
        ],
        "FirelensConfiguration": {
          "Type": "fluentbit",
          "Options": {
            "config-file-type": "file",
            // path to the extra.conf file inside container
            "config-file-value": "/extra.conf" 
          }
        },
        "LogConfiguration": {
          "LogDriver": "awslogs",
          "Options": {
            "awslogs-create-group": "true",
            "awslogs-group": "/ecs/fluentbit-agent",
            "awslogs-region": "eu-west-1",
            "awslogs-stream-prefix": "ecs"
          }
        }
      }
    ]
  }
```

The Task Definition will require a Task Execution Role with the following permissions:

- "logs:CreateLogGroup"
- "logs:CreateLogStream"
- "logs:PutLogEvents"
- "logs:DescribeLogStreams"


Once the task definition is defined, you can save it and create service from it. ECS will create a task with two containers, one for your application and one for fluentbit. The fluentbit container will collect logs from the application container and send them to Coralogix.

The JSON example above can be used to deploy via the AWS Console or API. Note that the [Cloudformation](template.yaml) format for defining a task definition is different.

### Cloudformation

When using this Cloudformation template you will need to modify it to suit your needs. You should start by replacing the [app container definition](template.yaml#72-#87) in the task to your own.

Once the template is modified, you can deploy it using the AWS Console or CLI.
