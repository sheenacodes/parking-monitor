#!/bin/bash

check_rabbitmq_up() {
    #curl -s -f -o /dev/null http://localhost:15672/api/healthchecks/node

    rabbitmq-diagnostics -q status | grep -q "running"
}

#
echo "Waiting for RabbitMQ to start..."
until check_rabbitmq_up; do
    echo "RabbitMQ is not available yet. Waiting..."
    sleep 5
done

echo "RabbitMQ is up and running."

# Create queues for vehicle entry log and exit log
rabbitmqadmin declare queue name=entryqueue durable=true
rabbitmqadmin declare queue name=exitqueue durable=true

echo "Queues created."
