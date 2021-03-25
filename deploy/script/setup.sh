cd
ln -s /mnt/volume_fra1_01 ./postgres-data
apt-get update && sudo apt install -y docker.io docker-compose ufw
ufw allow 22
ufw allow 80
ufw enable
