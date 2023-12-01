FROM golang:1.20.6-alpine

# Set the working directory to /app 
WORKDIR /app 

# Copy go.mod and go.sum 
COPY go.mod go.sum ./ 

RUN go mod download 

# Copy the content of the the current folder to /app 
COPY . ./

# Build the Go application 
RUN CGO_ENABLED=0 GOOS=linux go build cmd/main.go



# Command to run the executable 
CMD ["./main"] 