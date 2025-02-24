#!/bin/bash
set -e

echo "Creating MySQL database and user..."

# Execute commands as root user
mysql -v --user="root" --password="$MYSQL_ROOT_PASSWORD" <<-EOSQL
    CREATE DATABASE IF NOT EXISTS payment_registration_system;
    CREATE USER IF NOT EXISTS 'app-user'@'%' IDENTIFIED BY 'app-pwd';
    GRANT ALL PRIVILEGES ON payment_registration_system.* TO 'app-user'@'%';
    FLUSH PRIVILEGES;
EOSQL

echo "MySQL initialization complete."
