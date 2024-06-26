#!/bin/bash

print_message() {
  blue='\033[0;34m'
  nc='\033[0m' # No Color
  echo -e "${blue}$1${nc}"
}

print_message "Exercising message processor"
print_message "----------------------------"

print_message "Starting rabbitmq"
docker-compose up -d

rabbitmq_is_ready() {
  docker-compose exec rabbitmq rabbitmqadmin list queues > /dev/null
}

until rabbitmq_is_ready; do
  sleep 1
  print_message "Waiting rabbitmq to start..."
done

print_message "Starting message processor"
go build -o message_processor example/process_payments_events.go
./message_processor &> message_processor_output.txt &
MESSAGE_PROCESSOR_PID=$!

print_message "Sending messages to rabbitmq"
for i in {1..10} ; do
    echo "Publishing WithdrawalCreated event with amount $i"
    docker-compose exec rabbitmq rabbitmqadmin publish routing_key="message-processor-example-queue" payload="
      {
    		\"type\": \"WithdrawalCreated\",
    		\"data\": {
    			\"id\": \"e728f3a7-b92f-46fe-b080-524442065cb3\",
    			\"user_id\": \"a4fe943b-9e38-4741-84f5-5d3931db2d0e\",
    			\"psp_id\": \"9f4a9742-8e59-4bb7-975f-c504f8ce4efd\",
    			\"amount\": $i,
    			\"currency\": \"usd\",
    			\"beneficiary\": {
    			  \"id\": \"02ec607d-5016-456d-ae2d-a2f4d0354421\"
    			}
    		}
    	}
    "
done

print_message "Terminating message processor"
kill -n 15 $MESSAGE_PROCESSOR_PID

print_message "Shutting down rabbitmq"
docker-compose down

print_message "Printing message processor output"
cat message_processor_output.txt
