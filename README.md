## Message broker (MB) 

### Requirements

#### Phase #0

[x] Operator should be able to define queues in message broker.
[x] Operator should be able to define topics in message broker:
    [X] A topic supports regex definition for queues where to send the messages.

- Provider should be able to send messages to a queue or a topic
[x] Consumer should be able to receive messages from a subscribed queues.

#### Phase #1

- Operator should be able to define fanout exchanges that would send out message to all known queues