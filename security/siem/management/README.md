# SIEM service Integration - AWS Management Role 
This cloud formation template creates an organization-wide role that the SIEM service can assume in order to enrich your logs with data coming from all the accounts in your organization. It only has one parameter, that is the principal that can assume the created role. This should be set to the Coralogix AWS account id you've been provided.


## Parameters:

| Parameter       | Description                                                                                                                                                                                                                          | Default Value                                                                | Required           |
|-----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|--------------------|
| Principal      |  Account ID of the principal                                                                                                                                                                                                             |                                                                              | :heavy_check_mark: |


