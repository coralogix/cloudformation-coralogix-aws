# Changelog

## s3-archive

### 0.0.7 / 10.7.2025
* [fix] Fix bug with missing dependencies for resource `aws_role_region`

### 0.0.6 / 7.7.2025
* [update] Change the metrics role to be the same as the logs role: `arn:aws:iam::${aws_account_id}:role/coralogix-archive-${aws_role_region}`

### 0.0.5 / 2.9.2024
* [update] Add option to run module in AP3 region

### 0.0.4 / 19.5.2024
* [update] 
 - Add delete permissions to the buckets rule
 - replaced ap1 with ap2 in the mapping

### 0.0.3 / 17.1.2024
* [Bug Fix] Changed the role for metrics bucket

### 0.0.2 / 7.1.2024
* [update] Update the lambda role premission

### 0.0.1 / 21.8.2023
* [update] Add support to US2 and add an option to use CustomCoralogixArn and ByPassRegion
<!-- To add a new entry write: -->
<!-- ### version / full date -->
<!-- * [Update/Bug fix] message that describes the changes that you apply -->
