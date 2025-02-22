# This script performs the setup of the necessary containers of the databases, to later run the tests.
# Usage: ./run_test.sh
#
# Arguments:
#   None
#
# Note:
#   Ensure you have the necessary permissions to run docker commands.
#!/bin/bash

DB_NAME="payment_registration_system"
MONGO_URI="mongodb://root:root@localhost:27017/$DB_NAME?authSource=admin"
MONGO_CONTAINER="test_mongodb"

# Check Docker is installed
check_docker() {
    echo "Checking Docker status..."
    if ! docker info >/dev/null 2>&1; then
        echo "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    if docker-compose ps >/dev/null 2>&1; then
        echo "Docker is running. Proceeding with the tests..."
    else
        echo "Docker is not running. Please start Docker and try again."
        docker-compose down -v
        docker-compose up -d
    fi
}

# Run docker-compose
run_docker_compose() {
    echo "Running docker-compose..."
    # First, in case of conflict, remove and stop the containers
    docker-compose down -v
    # Then, run the containers
    docker-compose up -d
}

# Setup MongoDB
# Create the database and the collection
# Create the user and permissions
setup_mongodb() {
    echo "Setting up MongoDB..."
    echo "Drop the database if it exists..."
    docker exec test_mongodb mongosh admin --username root --password root --eval "db.getSiblingDB('payment_registration_system').dropDatabase();"

    echo "Create the user of the database..."
    docker exec test_mongodb mongosh admin --username root --password root --eval "
        db = db.getSiblingDB('admin');
        db.createUser({
            user: 'testuser',
            pwd: 'testpassword',
            roles: [
                { role: 'readWrite', db: 'payment_registration_system' },
                { role: 'dbAdmin', db: 'payment_registration_system' }
            ]
        });
    "
    echo "Loading data..."
    for file in data/non_relational/*.json; do
        if [ -f "$file" ]; then
            collection_name=$(basename "$file" .json) # Extract filename without .json extension
            echo "Importing $file into $collection_name collection in MongoDB..."

            docker cp "$file" test_mongodb:/data/"$(basename "$file")"

            docker exec test_mongodb mongoimport --uri "mongodb://testuser:testpassword@localhost:27017/payment_registration_system?authSource=admin" \
                --collection "$collection_name" --file "/data/$collection_name.json" --jsonArray
            echo "Imported $collection_name into MongoDB. Using the following command:"
            echo "docker exec test_mongodb mongoimport --uri 'mongodb://testuser:testpassword@localhost:27017/payment_registration_system?authSource=admin' --collection '$collection_name' --file '/data/$collection_name.json' --jsonArray"
        fi
    done

    echo "MongoDB setup complete!"
}

# Setup MySQL
# Create the database and the tables
# Create the user and permissions
setup_mysql() {
    echo "Setting up MySQL..."
    echo "Drop the database if it exists..."
    docker exec test_mysql mysql -u root -proot -e "DROP DATABASE IF EXISTS payment_registration_system;"

    echo "Create the database's user and database..."
    docker exec test_mysql mysql -u root -proot -e "
        CREATE USER IF NOT EXISTS 'testuser'@'%' IDENTIFIED BY 'testpassword';
        CREATE DATABASE IF NOT EXISTS payment_registration_system;
        GRANT ALL PRIVILEGES ON payment_registration_system.* TO 'testuser'@'%';
        FLUSH PRIVILEGES;
    "

    echo "Loading data..."
    docker exec test_mysql mysql -u root -proot -h 127.0.0.1 -P 3306 payment_registration_system <./data/relational/database.sql

    echo "MySQL setup complete!"
}

# Function to wait for a specified duration
wait_for() {
    local duration=$1
    echo "Waiting for $duration seconds..."
    sleep $duration
}

# Run go test
run_go_test() {
    echo "Running go test..."
    go test -v ./...

    if [ $? -eq 0 ]; then
        echo "All tests passed!"
    else
        echo "Some tests failed. Please check the logs."
        exit 1
    fi
}

# Main function
main() {
    check_docker
    run_docker_compose
    wait_for 10
    setup_mongodb
    # setup_mysql For now, we handle the MySQL setup in Go tests
    wait_for 10
    run_go_test
}

main
