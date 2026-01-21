# SIEM service Integration - AWS Management Role 
This cloud formation template creates an organization-wide role that the CSPM agent can assume in order to scan all the accounts in your organization. It only has one parameter, that is the principal that can assume the created role. This should be set to the AWS account id of the account where the CSPM agent is going to run.


## Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| Principal      |  Account ID of the principal                                                                                                                                                                                                             |                                                                              | :heavy_check_mark: |


