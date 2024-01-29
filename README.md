# pubsub
Viswals publisher subscriber microservice challenge.

## APPROX ETA  UNTIL Completion (initial) 
~14 hours 

## Difficulty 

- Medium 


## How To Run this Project 

- Go version: 1.20.6

### Docker (for the publisher and consumer services)


Open a new terminal and run the following command at the root directory

```bash
    docker-compose up 
```

You will have to wait at least 1 minute for all the services to be bootstraped and the data 
transfer from the publisher service to the consumer to start.


## REST API

Once the server is started and the rabbitmq go routine has finished writing to the database, go to address (localhost:8000/user and localhost:8000/users) to 
interact with the API endpoints 

### Endpoints 

- http://localhost/8000/user/:id : returns a single user with the specified Id
- http://localhost/8000/users/?limit=100&skip=1 returns an array of users if the "limit" and "skip" parameters are not provided it defaults to a limit of 100, and a skip of 1


## Frontend & Server Sent events

Not supported.