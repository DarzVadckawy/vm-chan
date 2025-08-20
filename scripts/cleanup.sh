#!/bin/bash

set -e

echo "Starting infrastructure cleanup..."

cd deployments/terraform

echo "Destroying AWS infrastructure..."
terraform destroy -auto-approve \
  -var="key_name=${AWS_KEY_NAME}" \
  -var="my_ip_cidr=${MY_IP_CIDR}" \
  -var="aws_region=${AWS_REGION:-us-east-1}"

echo "Infrastructure cleanup completed!"
echo "All AWS resources have been destroyed."
