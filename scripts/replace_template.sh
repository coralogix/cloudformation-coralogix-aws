#!/bin/bash
# This code will take a template file and change it according to the requirements in the integration-definitions repo

file=$1
if grep -q "ParameterGroups" "$file"; then
    yq eval --inplace '.Metadata."AWS::CloudFormation::Interface".ParameterGroups[0].Parameters += "IntegrationId"' -i $file
fi
if grep -q "ParameterLabels" "$file"; then
    yq eval --inplace '.Metadata."AWS::CloudFormation::Interface".ParameterLabels += {"IntegrationId": {"default": "Integration ID"}}' $file
fi

yq eval --inplace '.Parameters += {"IntegrationId": {"Type": "String",  "Description": "The integration ID to register."}}' $file

echo "  # Used as a bridge because CF doesn't allow for conditional depends on clauses.
  NonNotifierResourcesAreReady:
    Type: AWS::CloudFormation::WaitConditionHandle
    Metadata:" >> "$file"

resources=$(cat "$file" | yq '.Resources' | grep -e '^[a-zA-Z]' | sed 's/:$//') # return the resources in the template
parameters=$(cat "$file" | yq '.Parameters' | grep -e '^[a-zA-Z]' | sed 's/:$//') # return the parameters in the template

no_condition_resource=()
while IFS= read -r resource; do
    if yq ".Resources[\"$resource\"] | has(\"Condition\")" "$file" | grep -q 'true'; then
        condition=$(yq ".Resources[\"$resource\"].Condition" "$file" | grep -oE '[^ ]+$')
        echo "      ${resource}Ready: !If [ $condition, !Ref ${resource}, \"\" ]" >> $file
    else
      no_condition_resource+=($resource)
    fi
done <<< "$resources"
if [[ $file == *"aws-shipper-lambda"* ]]; then
  yq eval --inplace '.Conditions += {"IsApiKeySecretArn": "condition"}' $file
  sed -i "s/IsApiKeySecretArn: condition/IsApiKeySecretArn: \!Not [\!Equals [\!Ref ApiKey , \!Select [0,\!Split [\":\" , \!Ref ApiKey]]]]/g" $file
fi

echo "  IntegrationStatusNotifier:
    Type: Custom::IntegrationsServiceNotifier
    DependsOn:" >> $file

for resource in "${no_condition_resource[@]}"; do
  echo "      - $resource" >> $file
done


echo "
    Properties:
      #      {{AWS_ACCOUNT_ID}} is replaced during the template synchronisation
      ServiceToken: !Sub \"arn:aws:lambda:\${AWS::Region}:{{AWS_ACCOUNT_ID}}:function:integrations-custom-resource-notifier\"
      IntegrationId: !Ref IntegrationId" >> $file
if [[ $file == *"aws-shipper-lambda"* ]]; then
  echo "      CoralogixDomain: !If
        - IsCustomDomain
        - !Ref CustomDomain
        - !FindInMap [CoralogixRegionMap, !Ref CoralogixRegion, Domain]
      CoralogixApiKey: !Ref ApiKey" >> $file
else
  echo "
      CoralogixDomain: !If
        - IsRegionCustomUrlEmpty
        - !Ref CustomDomain
        - !FindInMap [ CoralogixRegionMap, !Ref CoralogixRegion, LogUrl ]
      CoralogixApiKey: !Ref ApiKey" >> $file
fi
echo "
      # Parameters to track
      IntegrationNameField: !Ref \"AWS::StackName\"
      SubsystemField: !Ref SubsystemName
      ApplicationNameField: !Ref ApplicationName" >> $file

while IFS= read -r parameter; do
  if [[ $parameter != "ApiKey" ]] && [[ $parameter != "IntegrationId" ]] && [[ $parameter != "ApplicationName" ]] && [[ $parameter != "SubsystemName" ]]; then
    echo "      ${parameter}Field: !Ref $parameter" >> $file
  fi
done <<< "$parameters"
