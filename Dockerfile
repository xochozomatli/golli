# Start from golang base image
FROM golang:alpine as builder

# Add Maintainer info
LABEL maintainer="Armand Villaverde <aavillaverde11@gmail.com>"

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Setup folders
WORKDIR /app

COPY go.mod go.sum .env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "./main" ]
