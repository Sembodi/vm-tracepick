#!/bin/bash

# Go to GoProject root directory
cd ..
make clean
make prep
make artifactgenerator
make builddocker
make rundocker
make removedocker
make export

bin/removedocker nginx 2> /dev/null 1> /dev/null && echo 'nginx container stopped and removed (if there was one)'
bin/removedocker apache 2> /dev/null 1> /dev/null && echo 'apache container stopped and removed (if there was one)'
bin/removedocker mariadb 2> /dev/null 1> /dev/null && echo 'mariadb container stopped and removed (if there was one)'
bin/removedocker redis-server 2> /dev/null 1> /dev/null && echo 'redis container stopped and removed (if there was one)'
bin/removedocker mongod 2> /dev/null 1> /dev/null && echo 'mongo container stopped and removed (if there was one)'
# bin/removedocker grafana-server 2> /dev/null 1> /dev/null && echo 'grafana container stopped and removed (if there was one)'

bashscripts/rmfromknownhosts.sh localhost

# NGINX:
bin/artifactgenerator yamlfiles/aws/nginx.yaml
sed -i '' 's/daemon on/daemon off/g' artifacts/output/containers/nginx/entrypoint.sh
bin/builddockerimage nginx
bin/rundocker nginx
bin/exportmetrics Nginx

bashscripts/rmfromknownhosts.sh localhost

# APACHE:
bin/artifactgenerator yamlfiles/aws/apache.yaml
sed -i '' 's/apachectl start/apachectl -DFOREGROUND/g' artifacts/output/containers/apache/entrypoint.sh
bin/builddockerimage apache
bin/rundocker apache
bin/exportmetrics Apache

bashscripts/rmfromknownhosts.sh localhost

# MARIADB:
bin/artifactgenerator yamlfiles/aws/mariadb.yaml
sed -i '' 's/mariadbd/mariadbd --console/g' artifacts/output/containers/mariadb/entrypoint.sh
bin/builddockerimage mariadb
bin/rundocker mariadb
bin/exportmetrics Mariadb

bashscripts/rmfromknownhosts.sh localhost

# REDIS:
bin/artifactgenerator yamlfiles/aws/redis.yaml
sed -i '' 's/daemonize yes/daemonize no/g' artifacts/output/containers/redis-server/entrypoint.sh
bin/builddockerimage redis-server
bin/rundocker redis-server
bin/exportmetrics Redis

bashscripts/rmfromknownhosts.sh localhost

# # # ELASTICSEARCH&KIBANA:
# # bin/artifactgenerator yamlfiles/aws/elk.yaml
# # # sed -i '' 's/apache2/apache2 -DFOREGROUND/g' artifacts/output/containers/apache/Dockerfile
# # bin/builddockerimage elk
# # bin/rundocker elk
# # bin/exportmetrics Elk
#
# MONGO:
bin/artifactgenerator yamlfiles/aws/mongo.yaml
sed -i '' 's/mongod --fork/mongod/g' artifacts/output/containers/mongod/entrypoint.sh
echo 'change dockerfile of mongod!!!'
# add following lines to Dockerfile:
#
# # RUN apt install -y wget gnupg
# # RUN wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | apt-key add -
# # RUN apt install -y mongodb-org || { echo "deb [] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/extra.list && apt update && apt install -y mongodb-org; }
# #
# # and delete this line:
# #
# # RUN apt install -y mongodb-org || { echo "deb [trusted=yes] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/extra.list && apt update && apt install -y mongodb-org; }
#
# bin/builddockerimage mongod
# bin/rundocker mongod
# bin/exportmetrics Mongo
#
# bin/artifactgenerator yamlfiles/aws/grafana.yaml
# # sed -i '' 's/apache2/apache2 -DFOREGROUND/g' artifacts/output/containers/apache/Dockerfile
# # add following lines to Dockerfile:
# #
# # RUN apt install -y wget gnupg
# # RUN wget -q -O - https://packages.grafana.com/gpg.key | apt-key add -
# #
# bin/builddockerimage grafana-server
# bin/rundocker grafana-server
# bin/exportmetrics Grafana




# nginx apache mariadb redis-server elk mongod postgres rabbitmq consul vault grafana-server influxdb prometheus minio


# bin/removedocker nginx
# bin/removedocker apache
# bin/removedocker mariadb
# bin/removedocker redis
# bin/removedocker elk
# bin/removedocker mongo
# bin/removedocker grafana
