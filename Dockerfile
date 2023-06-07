FROM golang:1.20.3-alpine

# set the working directory
WORKDIR /app

# copy the go.mod and sum files
COPY go.mod go.sum ./

# Download any dependencies
RUN go mod download

# Copy the source tree
COPY . .

# Build
RUN go build -o hotel-api .

# Expose the port
EXPOSE 3000

# Set the entry point
CMD ["./hotel-api"]

