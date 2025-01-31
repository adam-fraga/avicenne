# Use the latest Go image as the base
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the source files to the working directory
COPY . .

# Build the application
RUN go build -o /app/avicenne main.go

# Run the application when the container starts
CMD ["./avicenne"]

