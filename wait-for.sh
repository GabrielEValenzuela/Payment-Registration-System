#!/bin/sh

# Wait for MySQL
while ! nc -z mysql 3306; do
    echo "Waiting for MySQL..."
    sleep 2
done
echo "MySQL is ready!"

# Wait for MongoDB
while ! nc -z mongodb 27017; do
    echo "Waiting for MongoDB..."
    sleep 2
done
echo "MongoDB is ready!"

exec "$@"
