FROM golang:1.20.6-alpine

# Set the working directory to /app 
WORKDIR /app 

# Copy 
COPY go.mod go.sum ./

RUN go mod download 

# Copy the rest of the application cod to the working directory 
COPY . ./

# Build the Go application 
RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go 

# Expose the port on which the app will run 
EXPOSE 8000 

# Command to run the executable 
CMD ["./main"] 