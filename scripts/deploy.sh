#!/bin/bash

set -e

echo "Starting local deployment process..."

required_vars=("AWS_ACCESS_KEY_ID" "AWS_SECRET_ACCESS_KEY" "AWS_KEY_NAME" "MY_IP_CIDR")
for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: $var is not set"
        exit 1
    fi
done

export AWS_REGION=${AWS_REGION:-us-west-2}
export TF_VAR_aws_region=$AWS_REGION
export TF_VAR_key_name=$AWS_KEY_NAME
export TF_VAR_my_ip_cidr=$MY_IP_CIDR

echo "Configuration:"
echo "  AWS Region: $AWS_REGION"
echo "  Key Name: $AWS_KEY_NAME"
echo "  Your IP: $MY_IP_CIDR"

echo "Building and testing application..."
go test ./...
docker build -t vm-chan:local .

echo "Provisioning infrastructure with Terraform..."
cd deployments/terraform
terraform init
terraform plan
terraform apply -auto-approve

export PUBLIC_IP=$(terraform output -raw public_ip)
echo "EC2 Instance created at: $PUBLIC_IP"

cd ../ansible
sed "s/\${PUBLIC_IP}/$PUBLIC_IP/g" inventory.tpl > inventory.ini

echo "Waiting for SSH access..."
for i in {1..30}; do
    if ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 -i ~/.ssh/id_rsa ubuntu@$PUBLIC_IP 'echo "SSH Ready"' 2>/dev/null; then
        break
    fi
    echo "  Attempt $i/30 - waiting..."
    sleep 10
done

echo "Installing k3s with Ansible..."
ansible-playbook -i inventory.ini site.yml

echo "Deployment completed successfully!"
echo ""
echo "Access your application:"
echo "  Health check: http://$PUBLIC_IP:30080/healthz"
echo "  API endpoint: http://$PUBLIC_IP:30080/analyze"
echo ""
echo "Management commands:"
echo "  SSH access: ssh -i ~/.ssh/id_rsa ubuntu@$PUBLIC_IP"
echo "  View pods: ssh -i ~/.ssh/id_rsa ubuntu@$PUBLIC_IP 'kubectl get pods -n vmchan'"
echo "  View logs: ssh -i ~/.ssh/id_rsa ubuntu@$PUBLIC_IP 'kubectl logs -n vmchan -l app=vm-chan'"
echo ""
echo "To cleanup: ./scripts/cleanup.sh"
