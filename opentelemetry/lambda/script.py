import boto3
import cfnresponse
import json

def handler(event, context):
    print("received event: " + json.dumps(event, indent=2))

    try:
        status = cfnresponse.SUCCESS

        if event['RequestType'] in ['Create', 'Update']:
            lambda_arn = event['ResourceProperties']['LambdaArn']
            lambda_client = boto3.client('lambda')

            # Get the current environment variables of the Lambda function
            function_configuration = lambda_client.get_function_configuration(FunctionName=lambda_arn)
            current_env = function_configuration['Environment']['Variables']

            # Merge the current environment variables with any new ones specified in the CloudFormation event
            # new_environment = current_environment.copy()
            current_env.update(event['ResourceProperties']['env'])

            # Get the list of layers to add to the Lambda function
            layers = event['ResourceProperties'].get('LayerArns', [])

            # Add the specified layers and environment variables to the Lambda function
            lambda_client.update_function_configuration(
                FunctionName=lambda_arn, 
                Environment={'Variables': current_env},
                Layers=layers)
        
        if event['RequestType'] == 'Delete':
            # TODO: implement delete logic
            pass
        

    except Exception as e:
        print(f"custom resource failed to update layers and env vars: {e}")
        status = cfnresponse.FAILED
    
    finally:
        # Send a status response to CloudFormation
        response_data = {}
        cfnresponse.send(event, context, status, response_data)
        return
