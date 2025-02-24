package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MySQL settings
const (
	mysqlDSN       = "user:password@tcp(127.0.0.1:3306)/benchmark"
	mysqlTableName = "test"
)

// MongoDB settings
const (
	mongoURI        = "mongodb://localhost:27017"
	mongoDatabase   = "benchmark"
	mongoCollection = "test"
)

func main() {
	fmt.Println("Starting benchmark...")

	// Get MAX_INSERTS from environment variable or set default to 1000
	maxInserts := getMaxInserts()

	// WaitGroup to synchronize goroutines
	var wg sync.WaitGroup

	// Add two goroutines to the WaitGroup
	wg.Add(2)

	// Run MySQL benchmark in a separate goroutine
	go func() {
		defer wg.Done() // Signal that this goroutine is done
		benchmarkMySQL(maxInserts)
	}()

	// Run MongoDB benchmark in a separate goroutine
	go func() {
		defer wg.Done() // Signal that this goroutine is done
		benchmarkMongoDB(maxInserts)
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	fmt.Println("Benchmarking completed.")
}

// Get MAX_INSERTS from environment variable with a default value
func getMaxInserts() int {
	maxInsertsStr := os.Getenv("MAX_INSERTS")
	if maxInsertsStr == "" {
		fmt.Println("MAX_INSERTS environment variable not set, using default: 1000")
		return 1000
	}
	maxInserts, err := strconv.Atoi(maxInsertsStr)
	if err != nil {
		log.Fatalf("Error converting MAX_INSERTS to integer: %v", err)
	}
	fmt.Printf("MAX_INSERTS set to: %d\n", maxInserts)
	return maxInserts
}

func benchmarkMySQL(maxInserts int) {
	fmt.Println("\nBenchmarking MySQL...")

	// Open connection to MySQL
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatal("Error connecting to MySQL: ", err)
	}
	defer db.Close()

	// Prepare the table
	_, err = db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INT AUTO_INCREMENT PRIMARY KEY, value VARCHAR(255))", mysqlTableName))
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	// Perform write benchmark
	start := time.Now()
	for i := 0; i < maxInserts; i++ {
		_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (value) VALUES (?)", mysqlTableName), fmt.Sprintf("Value %d", i))
		if err != nil {
			log.Fatal("Error inserting data into MySQL: ", err)
		}
	}
	writeDuration := time.Since(start)

	// Perform read benchmark
	start = time.Now()
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", mysqlTableName))
	if err != nil {
		log.Fatal("Error reading data from MySQL: ", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var value string
		err := rows.Scan(&id, &value)
		if err != nil {
			log.Fatal("Error scanning row: ",
				err)
		}
	}
	readDuration := time.Since(start)

	fmt.Printf("MySQL Write Duration: %v\n", writeDuration)
	fmt.Printf("MySQL Read Duration: %v\n", readDuration)
}

func benchmarkMongoDB(maxInserts int) {
	fmt.Println("\nBenchmarking MongoDB...")

	// Create a MongoDB client
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
		}
	}()

	// Ping the MongoDB server to ensure the connection is established
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	// Select the database and collection
	collection := client.Database(mongoDatabase).Collection(mongoCollection)

	// Clean up the collection before the test
	err = collection.Drop(context.TODO())

	if err != nil {
		log.Fatalf("Error dropping collection: %v", err)
	}

	// Write benchmark: Insert documents
	start := time.Now()
	for i := 0; i < maxInserts; i++ {
		_, err = collection.InsertOne(context.TODO(), bson.D{
			{Key: "name", Value: fmt.Sprintf("item-%d", i)},
			{Key: "value", Value: i},
		})
		if err != nil {
			log.Fatalf("Error inserting document into MongoDB: %v", err)
		}
	}
	writeDuration := time.Since(start)

	// Read benchmark: Retrieve all documents
	start = time.Now()
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Error reading documents from MongoDB: %v", err)
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor
	for cursor.Next(ctx) {
		var result bson.D
		if err := cursor.Decode(&result); err != nil {
			log.Fatalf("Error decoding document: %v", err)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("Cursor error: %v", err)
	}

	readDuration := time.Since(start)

	fmt.Printf("MongoDB Write Duration: %v\n", writeDuration)
	fmt.Printf("MongoDB Read Duration: %v\n", readDuration)
}
