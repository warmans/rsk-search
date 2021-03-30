# setup home dir
cd
ln -s /mnt/volume_fra1_01/postgres-data ./postgres-data
ln -s /mnt/volume_fra1_01/postgres-backups ./postgres-backups

# packages
apt-get update && sudo apt install -y docker.io docker-compose ufw

# firewall
ufw allow 22
ufw allow 80
ufw enable
