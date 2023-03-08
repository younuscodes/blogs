#!/bin/bash

# Build the Docker image for the app.
#docker build -t myapp .

# Start the Docker services using docker-compose.
docker-compose up -d

# Wait for the app to start up.
echo "Waiting for app to start up..."
#until $(curl --output /dev/null --silent --head --fail http://localhost:8123); do
    #printf '.'
    sleep 30
#done
echo "Insert Data..."
curl -X POST -H "Content-Type: application/json" -d '{"title":"My first article", "content":"This is some content", "author":"John Doe"}' http://localhost:8123/articles

echo "Verify Inserted Data..."
curl http://localhost:8123/articles

# Run the tests.
#echo "Running tests..."
#docker run --network container:myapp myapp go test ./...

# Stop the Docker services.
#echo "Stopping services..."
#docker-compose down
