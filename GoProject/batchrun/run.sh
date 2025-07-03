# RUN apt install -y wget gnupg
# RUN apt install -y apt-transport-https && \
#   wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg && \
#   echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/9.x/apt stable main" | tee /etc/apt/sources.list.d/elastic-9.x.list && \
#   apt update



sed -i '' 's~/usr/share/elasticsearch~exec /sbin/runuser -u elasticsearch /usr/share/elasticsearch~g' artifacts/output/containers/mongod/entrypoint.sh
sed -i '' 's~kibana~kibana --allow-root~g' artifacts/output/containers/mongod/entrypoint.sh
