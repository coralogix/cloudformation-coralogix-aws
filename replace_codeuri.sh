#!/bin/bash
# This script will change the codeUri in a given template, so that the template will be validate
# for cloudformation.

file=$1
puckage_name=$2
bucket_string="!Sub coralogix-serverless-repo-\${AWS::Region}"

if grep -q "LambdaFunction" "$file"; then
    yq eval '.Resources.LambdaFunction.Properties.CodeUri = {"Bucket": "'"$bucket_string"'", "S3Key": "'"$puckage_name"'.zip"}' -i $file
    sed -i "s/'!Sub coralogix-serverless-repo-\${AWS::Region}/!Sub 'coralogix-serverless-repo-\${AWS::Region}/g" $file
fi
if grep -q "LambdaFunctionSSM" "$file"; then
    yq eval '.Resources.LambdaFunctionSSM.Properties.CodeUri = {"Bucket": "'"$bucket_string"'", "S3Key": "'"$puckage_name"'.zip"}' -i $file
    sed -i "s/'!Sub coralogix-serverless-repo-\${AWS::Region}/!Sub 'coralogix-serverless-repo-\${AWS::Region}/g" $file
fi
if grep -q "CustomResourceLambdaTriggerFunction" "$file"; then
    yq eval '.Resources.CustomResourceLambdaTriggerFunction.Properties.CodeUri = {"Bucket": "'"$bucket_string"'", "S3Key": "helper.zip"}' -i $file
    sed -i "s/'!Sub coralogix-serverless-repo-\${AWS::Region}/!Sub 'coralogix-serverless-repo-\${AWS::Region}/g" $file
fi
