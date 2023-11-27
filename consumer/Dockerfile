FROM golang:latest 

# Set the working directory to /app 
WORKDIR /app 

# Copy the current directory contents into the container at /app 
COPY . .  

# Build the Go application 
RUN go build -o main . 

# Expose the port on which the app will run 
EXPOSE 8000 

# Command to run the executable 
CMD ["./main"] 