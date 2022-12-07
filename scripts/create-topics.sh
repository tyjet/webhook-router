#!/bin/sh

# Requires Zookeeper and the Kafka broker to be running
brokers=$(echo dump | nc localhost 2181 | grep brokers)
if [ -z $brokers ]; then 
    echo "Broker ain't not yet running"
    exit 1
fi

echo "Broker is doing it thang"

docker exec broker kafka-topics --bootstrap-server broker:9092 --create --topic pull_request
docker exec broker kafka-topics --bootstrap-server broker:9092 --create --topic pull_request_review
