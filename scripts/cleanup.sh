#!/bin/bash

# Cleanup script for vm-chan infrastructure
set -e

echo "🗑️  Starting infrastructure cleanup..."

# Change to terraform directory
cd deployments/terraform

# Destroy infrastructure
echo "Destroying AWS infrastructure..."
terraform destroy -auto-approve \
  -var="key_name=${AWS_KEY_NAME}" \
  -var="my_ip_cidr=${MY_IP_CIDR}" \
  -var="aws_region=${AWS_REGION:-us-west-2}"

echo "✅ Infrastructure cleanup completed!"
echo "💡 All AWS resources have been destroyed."
