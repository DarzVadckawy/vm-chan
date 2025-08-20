[k3s_cluster]
k3s-master ansible_host=${public_ip} ansible_user=ubuntu ansible_ssh_private_key_file=~/.ssh/id_rsa

[k3s_cluster:vars]
ansible_ssh_common_args='-o StrictHostKeyChecking=no'
