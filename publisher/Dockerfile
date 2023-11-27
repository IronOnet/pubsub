FROM golang:latest  

# Set the working directory to /app 
WORKDIR /app 

# Copy the content of the the current folder to /app 
COPY . .  

# Build the Go application 
RUN go build cmd/main.go

EXPOSE 5672 

# Command to run the executable 
CMD ["./main"] 