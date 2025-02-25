#!/bin/bash

set -e # Exit on error

echo "🚀 Loading data into MySQL and MongoDB..."

# MySQL Import
MYSQL_CONTAINER="mysql"
MYSQL_USER="root"
MYSQL_PASSWORD="root"
MYSQL_DATABASE="payment_registration_system"
SQL_FILE="src/internal/tests/component/data/relational/database.sql"

echo "📂 Copying SQL file to MySQL container..."
docker cp "$SQL_FILE" "$MYSQL_CONTAINER:/database.sql"

echo "💾 Importing database.sql into MySQL..."
docker exec "$MYSQL_CONTAINER" sh -c "mysql -u $MYSQL_USER -p$MYSQL_PASSWORD $MYSQL_DATABASE < /database.sql" &&
    echo "✅ MySQL database imported successfully!" ||
    {
        echo "❌ MySQL import failed!"
        exit 1
    }

# MongoDB Import
MONGO_CONTAINER="mongodb"
MONGO_URI="mongodb://mongo_user:mongo_password@mongodb:27017/payment_registration_system?authSource=admin"
DATA_DIR="src/internal/tests/component/data/non_relational"

echo "📂 Processing JSON files for MongoDB..."
for file in "$DATA_DIR"/*.json; do
    if [ -f "$file" ]; then
        COLLECTION_NAME=$(basename "$file" .json) # Extract filename without extension
        echo "📦 Copying $file to MongoDB container..."
        docker cp "$file" "$MONGO_CONTAINER:/data/$(basename "$file")"

        echo "📥 Importing $COLLECTION_NAME into MongoDB..."
        docker exec "$MONGO_CONTAINER" mongoimport \
            --uri "$MONGO_URI" \
            --collection "$COLLECTION_NAME" \
            --file "/data/$(basename "$file")" \
            --jsonArray &&
            echo "✅ Successfully imported $COLLECTION_NAME into MongoDB!" ||
            {
                echo "❌ Failed to import $COLLECTION_NAME into MongoDB!"
                exit 1
            }
    fi
done

echo "🎉 Data import completed successfully!"
