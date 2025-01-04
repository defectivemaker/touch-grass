ARG CACHEBUST=1

# Build stage
FROM golang:1.21 AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN make docker-server

# Run stage
FROM alpine:3.14
WORKDIR /app
COPY --from=build /app/server-indicum .
RUN mkdir -p /app/certs

# Install PostgreSQL client and bash
RUN apk add --no-cache postgresql-client bash

# Copy wait-for-db script
COPY docker-wait-for-db.sh /app/wait-for-db.sh
RUN chmod +x /app/wait-for-db.sh

EXPOSE 8081 8888

# Use wait-for-db script as entrypoint
ENTRYPOINT ["/app/wait-for-db.sh"]
CMD ["db", "5432", "/app/server-indicum"]