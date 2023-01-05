# setup home dir
cd
mkdir -p /mnt/volume_fra1_01/postgres-data && ln -s /mnt/volume_fra1_01/postgres-data ./postgres-data
mkdir -p /mnt/volume_fra1_01/postgres-backups && ln -s /mnt/volume_fra1_01/postgres-backups ./postgres-backups
mkdir -p /mnt/volume_fra1_01/swag-config && ln -s /mnt/volume_fra1_01/swag-config ./swag-config
mkdir -p /mnt/volume_fra1_01/grafana/data && chown -R 472:472 mnt/volume_fra1_01/grafana && ln -s /mnt/volume_fra1_01/grafana ./grafana
mkdir -p /mnt/volume_fra1_01/prometheus && chown -R nobody:nobody mnt/volume_fra1_01/prometheus && ln -s /mnt/volume_fra1_01/prometheus ./prometheus

# fix permission on audio dir to allow writing from docker container
chmod -R nobody:nogroup /mnt/audio

# packages
apt-get update && sudo apt install -y docker.io docker-compose ufw

# firewall
ufw allow 22
ufw allow 80
ufw allow 443
ufw enable
