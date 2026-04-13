# CSPM Agent Integration
This cloud formation template creates a role with the permissions that the CSPM agent needs in order to function. It only has one parameter, that is the principal that can assume the created role.

There are two possible ways of using this template:
1. You do not have an AWS organization. In this case you need to run this template against all the accounts that you wish to scan with the CSPM agent, specifying the account id from which the CSPM agent is running as the principal. You're then going to have to provide the agent with the ARNs of the created roles.

2. You have an AWS organization. In this case you're going to have to run the management template against your organization, setting the account from which the CSPM agent is running as the principal. This will create an organization-wide role that the CSPM agent will be able to assume. You're then going to have to run the execution template against all the accounts in the organization, specifying the management role that you creted in the previous step as the principal (you can also use StackSets to accomplish this https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/stacksets-concepts.html). 


## Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| Principal      |  Account ID/Role ARN of the principal                                                                                                                                                                                                             |                                                                              | :heavy_check_mark: |

