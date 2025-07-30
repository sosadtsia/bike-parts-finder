#!/bin/bash

# Generate a random UUID for the cluster ID
CLUSTER_ID=$(kafka-storage random-uuid)

# Format the storage directories
kafka-storage format --cluster-id=$CLUSTER_ID --config /etc/kafka/kraft/server.properties

# Start Kafka server
exec /etc/confluent/docker/run
