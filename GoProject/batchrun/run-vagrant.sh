#!/bin/bash

# Go to GoProject root directory
cd ..
make artifactgenerator
make builddocker
make rundocker
make export
make removedocker

bin/removedocker nginx 2> /dev/null 1> /dev/null && echo 'nginx container stopped and removed (if there was one)'
bin/removedocker apache 2> /dev/null 1> /dev/null && echo 'apache container stopped and removed (if there was one)'
bin/removedocker mariadb 2> /dev/null 1> /dev/null && echo 'mariadb container stopped and removed (if there was one)'
bin/removedocker redis-server 2> /dev/null 1> /dev/null && echo 'redis container stopped and removed (if there was one)'
bin/removedocker elk 2> /dev/null 1> /dev/null && echo 'elk container stopped and removed (if there was one)'
bin/removedocker mongod 2> /dev/null 1> /dev/null && echo 'mongo container stopped and removed (if there was one)'
bin/removedocker grafana-server 2> /dev/null 1> /dev/null && echo 'grafana container stopped and removed (if there was one)'
bin/removedocker prometheus 2> /dev/null 1> /dev/null && echo 'prometheus container stopped and removed (if there was one)'

bashscripts/rmfromknownhosts.sh localhost

# # NGINX:
# bin/artifactgenerator yamlfiles/vagrant/nginx.yaml
# sed -i '' 's/daemon on/daemon off/g' artifacts/output/containers/nginx/entrypoint.sh
# bin/builddockerimage nginx
# bin/rundocker nginx
# bin/exportmetrics Nginx
#
# bashscripts/rmfromknownhosts.sh localhost
#
# APACHE:
bin/artifactgenerator yamlfiles/vagrant/apache.yaml
sed -i '' 's/apachectl start/apachectl -DFOREGROUND/g' artifacts/output/containers/apache/entrypoint.sh
bin/builddockerimage apache
bin/rundocker apache
bin/exportmetrics Apache
#
# bashscripts/rmfromknownhosts.sh localhost
#
# # MARIADB:
# bin/artifactgenerator yamlfiles/vagrant/mariadb.yaml
# sed -i '' 's/mariadbd/mariadbd --console/g' artifacts/output/containers/mariadb/entrypoint.sh
# bin/builddockerimage mariadb
# bin/rundocker mariadb
# bin/exportmetrics Mariadb
#
# bashscripts/rmfromknownhosts.sh localhost
#
# # REDIS:
# bin/artifactgenerator yamlfiles/vagrant/redis.yaml
# sed -i '' 's/daemonize yes/daemonize no/g' artifacts/output/containers/redis-server/entrypoint.sh
# bin/builddockerimage redis-server
# bin/rundocker redis-server
# bin/exportmetrics Redis
#
# # bashscripts/rmfromknownhosts.sh localhost
#
# # # # ELASTICSEARCH&KIBANA:
# bin/artifactgenerator yamlfiles/vagrant/elk.yaml
# sed -i '' 's~/usr/share/elasticsearch~exec /sbin/runuser -u elasticsearch /usr/share/elasticsearch~g' artifacts/output/containers/mongod/entrypoint.sh
# sed -i '' 's~/usr/share/kibana/bin/kibana~/usr/share/kibana/bin/kibana --allow-root~g' artifacts/output/containers/mongod/entrypoint.sh
#
# # Add to Dockerfile after first apt update:
# # RUN apt install -y wget gnupg
# # RUN apt install -y apt-transport-https && \
# #   wget -qO - https://artifacts.elastic.co/GPG-KEY-elasticsearch | gpg --dearmor -o /usr/share/keyrings/elasticsearch-keyring.gpg && \
# #   echo "deb [signed-by=/usr/share/keyrings/elasticsearch-keyring.gpg] https://artifacts.elastic.co/packages/9.x/apt stable main" | tee /etc/apt/sources.list.d/elastic-9.x.list && \
# #   apt update
# bin/builddockerimage elk
# bin/rundocker elk
# bin/exportmetrics ElasticsearchKibana
#
# #
# # MONGO:
# bin/artifactgenerator yamlfiles/vagrant/mongo.yaml
# sed -i '' 's/mongod --fork/mongod/g' artifacts/output/containers/mongod/entrypoint.sh
# echo 'change dockerfile of mongod!!!'
#
# # Add to Dockerfile after first apt update:
# #
# # RUN apt install -y wget gnupg
# # RUN wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | apt-key add -
# # RUN apt install -y mongodb-org || { echo "deb [] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/extra.list && apt update && apt install -y mongodb-org; }
# #
# # and delete this line:
# #
# # RUN apt install -y mongodb-org || { echo "deb [trusted=yes] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/extra.list && apt update && apt install -y mongodb-org; }
# bin/builddockerimage mongod
# bin/rundocker mongod
# bin/exportmetrics Mongo
#
# # GRAFANA
# bin/artifactgenerator yamlfiles/vagrant/grafana.yaml
# # sed -i '' 's/apache2/apache2 -DFOREGROUND/g' artifacts/output/containers/apache/Dockerfile
# # add following lines to Dockerfile:
# #
# # RUN apt install -y wget gnupg
# # RUN wget -q -O - https://packages.grafana.com/gpg.key | apt-key add -
# #
# bin/builddockerimage grafana-server
# bin/rundocker grafana-server
# bin/exportmetrics Grafana
#
# bashscripts/rmfromknownhosts.sh localhost
#
# # PROMETHEUS
# bin/artifactgenerator yamlfiles/vagrant/prometheus.yaml
# awk '
# /^RUN apt install -y rsync$/ {
#   while ((getline line < "batchrun/install-dockerstrings/prometheus.Dockerfile") > 0) print line
#   next
# }
# { print }
# ' artifacts/output/containers/prometheus/Dockerfile > Dockerfile.new && mv Dockerfile.new artifacts/output/containers/prometheus/Dockerfile
# bin/builddockerimage prometheus
# bin/rundocker prometheus
# bin/exportmetrics Prometheus

# bashscripts/rmfromknownhosts.sh localhost

# # POSTGRESQL
# bin/artifactgenerator yamlfiles/vagrant/postgres.yaml
# sed -i '' 's/postgres/postgres -D/g' artifacts/output/containers/postgres/Dockerfile
# bin/builddockerimage postgres
# bin/rundocker postgres
# bin/exportmetrics Postgres
#
# # INFLUXDB
# bin/artifactgenerator yamlfiles/vagrant/influxdb.yaml
# # sed -i '' 's/apache2/apache2 -DFOREGROUND/g' artifacts/output/containers/apache/Dockerfile
# bin/builddockerimage influxdb
# bin/rundocker influxdb
# bin/exportmetrics Influxdb
#
# # VAULT
# bin/artifactgenerator yamlfiles/vagrant/vault.yaml
# sed -i '' 's/vault server/vault server -dev/g' artifacts/output/containers/apache/Dockerfile
# bin/builddockerimage vault
# bin/rundocker vault
# bin/exportmetrics Vault
#
# # RABBITMQ
# bin/artifactgenerator yamlfiles/vagrant/rabbitmq.yaml
# # sed -i '' 's/apache2/apache2 -DFOREGROUND/g' artifacts/output/containers/rabbitmq/Dockerfile
# bin/builddockerimage rabbitmq
# bin/rundocker rabbitmq
# bin/exportmetrics Rabbitmq
# # bin/removedocker rabbitmq  # remove container if necessary
#
# # CONSUL
# bin/artifactgenerator yamlfiles/vagrant/consul.yaml
# sed -i '' 's/consul agent/consul agent -dev/g' artifacts/output/containers/consul/Dockerfile
# bin/builddockerimage consul
# bin/rundocker consul
# bin/exportmetrics Consul
#
# # MINIO
# bin/artifactgenerator yamlfiles/vagrant/minio.yaml
# # sed -i '' 's/apache2/apache2 -DFOREGROUND/g' artifacts/output/containers/apache/Dockerfile
# bin/builddockerimage minio
# bin/rundocker minio
# bin/exportmetrics Minio


# nginx apache mariadb redis-server elk mongod postgres rabbitmq consul vault grafana-server influxdb prometheus minio


# bin/removedocker nginx
# bin/removedocker apache
# bin/removedocker mariadb
# bin/removedocker redis
# bin/removedocker elk
# bin/removedocker mongo
# bin/removedocker postgres
# bin/removedocker rabbitmq
# bin/removedocker consul
# bin/removedocker vault
# bin/removedocker grafana
# bin/removedocker influxdb
# bin/removedocker prometheus
# bin/removedocker minio
