# rabbitmq type

This type manages rabbitmq connections and manage sending

## Required variables

The following env variables should be definned when using this type:

* **RABBITMQ_HOST** defines rabbitmq host, it must be a string, it's default value is "localhost" 
* **RABBITMQ_PORT** defines rabbitmq port, it must be a valid port number, it's default value is 5672 
* **RABBITMQ_USER** defines rabbitmq user, it's default value is "guest" 
* **RABBITMQ_PASSWORD** defines rabbitmq user's password, it's default value is "guest" 


## Available functions

### NewConfig -> Define a new config

Defines a new rabbitmq config, it does not test rabbitmq instance connection

### SendMessage -> Sends a message ([] byte) through required queue

Example:
```
config, _ := NewConfig()
queueName := "test"
testString := []byte("This is a Test")

sendError := config.SendMessage(queueName, testString)
```

### ReceiveMessages -> Receives messages ([] byte) through required queue until context is closed

This function should be called as goroutine

Example:
```
config, _ := NewConfig()
queueName := "test"
messagesReceived := make(chan []byte)

ctx, cancel := context.WithCancel(context.Background())

receiveErrors := make(chan error)

go config.ReceiveMessages(ctx, queueName, messagesReceived, receiveErrors)

cancel()

messageReceived := <-messagesReceived
```
