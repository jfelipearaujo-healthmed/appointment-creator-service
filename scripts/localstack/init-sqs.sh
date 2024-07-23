#!/bin/sh

echo "Initializing SQS queues..."

awslocal sqs create-queue \
    --queue-name AppointmentQueue.fifo \
    --attributes FifoQueue="true"

echo "SQS queues initialized!"