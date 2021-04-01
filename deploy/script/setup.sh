# setup home dir
cd
mkdir -p /mnt/volume_fra1_01/postgres-data && ln -s /mnt/volume_fra1_01/postgres-data ./postgres-data
mkdir -p /mnt/volume_fra1_01/postgres-backups && ln -s /mnt/volume_fra1_01/postgres-backups ./postgres-backups
mkdir -p /mnt/volume_fra1_01/swag-config && ln -s /mnt/volume_fra1_01/swag-config ./swag-config

# packages
apt-get update && sudo apt install -y docker.io docker-compose ufw

# firewall
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable
