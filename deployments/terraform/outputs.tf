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

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh -i ~/.ssh/id_rsa ubuntu@${aws_instance.k3s.public_ip}"
}
