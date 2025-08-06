#/bin/bash

/opt/bitnami/kafka/bin/kafka-topics.sh \
  --create \
  --if-not-exists \
  --topic $KAFKA_TOPIC \
  --partitions 1 \
  --replication-factor 1 \
  --bootstrap-server $KAFKA_HOST
