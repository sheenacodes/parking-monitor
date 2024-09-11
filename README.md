todo

api docs
prometheus meyrics

retention rabbit mq

docker exec -it rabbitmq rabbitmqctl set_policy TTL ".*" '{"message-ttl":60000}' --apply-to queues
