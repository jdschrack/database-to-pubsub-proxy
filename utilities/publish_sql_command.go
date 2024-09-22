package utilities

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/jdschrack/database-to-pubsub-proxy/models"
)

// Avro schema for SQLCommandEvent

//func main() {
//	// Simulate an intercepted SQL command
//	sqlCommand := "SELECT * FROM users;"
//	database := "example_db"
//	user := "db_user"
//
//	err := utilities.PublishSQLCommand(sqlCommand, database, user)
//	if err != nil {
//		log.Fatalf("Failed to publish SQL command: %v", err)
//	}
//}

func PublishSQLCommand(sqlCommand, database, user string) error {
	// Create a Pub/Sub client
	ctx := context.Background()

	record := models.SqlCommandEvent{
		SqlCommand: sqlCommand,
		Database:   database,
		Timestamp:  time.Now().Format(time.RFC3339),
		User:       user,
	}

	pubsubClient, err := pubsub.NewClient(ctx, "jh-sdb-dig-banno-capybara")
	if err != nil {
		return err
	}
	defer pubsubClient.Close() // Prepare the Avro record

	//Publish the Avro-encoded message to Pub/Sub
	topic := pubsubClient.Topic("symitar-sql-events")

	cfg, err := topic.Config(ctx)
	if err != nil {
		log.Fatalf("Failed to get topic config: %v", err)
		return err
	}

	log.Printf("Topic config: %+v", cfg)

	encoding := cfg.SchemaSettings.Encoding

	msg, err := CreateAvroMessage(record, encoding)

	if err != nil {
		log.Fatalf("Failed to create Avro message: %v", err)
		return err
	}

	result := topic.Publish(
		ctx, &pubsub.Message{
			Data: msg,
		},
	)

	// Wait for the result and log the message ID
	id, err := result.Get(ctx)
	if err != nil {
		return err
	}

	log.Printf("Published message with ID: %s", id)
	return nil
}
