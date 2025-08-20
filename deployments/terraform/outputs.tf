output "public_ip" {
  description = "Public IP address of the k3s instance"
  value       = aws_instance.k3s.public_ip
}

output "instance_id" {
  description = "EC2 instance ID"
  value       = aws_instance.k3s.id
}

output "security_group_id" {
  description = "Security group ID"
  value       = aws_security_group.k3s.id
}

output "key_name" {
  description = "EC2 key pair name"
  value       = aws_key_pair.k3s_key.key_name
}

output "private_key" {
  description = "Private key for SSH access"
  value       = tls_private_key.k3s_key.private_key_pem
  sensitive   = true
}

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh -i ~/.ssh/${aws_key_pair.k3s_key.key_name}.pem ubuntu@${aws_instance.k3s.public_ip}"
}

output "ansible_inventory" {
  description = "Ansible inventory content"
  value = templatefile("${path.module}/inventory.tpl", {
    PUBLIC_IP = aws_instance.k3s.public_ip
  })
}
