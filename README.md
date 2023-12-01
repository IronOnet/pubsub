# pubsub
Viswals publisher subscriber microservice challenge.

## APPROX ETA  
~14 hours 

## Difficulty 

- Medium 


## How To Run this Project 

- Go version: 1.20.6

### Without Docker (for the publisher and consumer services)


Open a new terminal and run

```bash
    docker-compose up 
```
This will make sure that the mysql, redis and rabbitmq containers are initialized

#### Publisher service 

To run the publisher service open another terminal and run 

```bash

cd publisher/ 

# download dependencies 
go mod download

# run the main file 
go run cmd/main.go 

```

#### Consumer service  

To run the consumer service,  at the root of the project run 

```bash 

cd consumer/ 

# download dependencies 
go mod download 

# run the main file 
go run cmd/main.go 

```

##### REST API

Once the server is started and the rabbitmq go routine has finished writing to the database, go to address (localhost:8000/user and localhost:8000/users) to 
interact with the API endpoints 

- user endpoint : Returns a single user 

This endpoint acccepts the following url format : /user/:id (example http://localhost:8000/user/10) 

- users endpoint: Returns an array of users 

This endpoint accepts the following url format : /users?limit=&skip= (example http://localhost:8000/users?limit=100&skip=1) 


## RUN with docker-compose 

This method does not work at the moment due to a latent bug that needs to be found. 