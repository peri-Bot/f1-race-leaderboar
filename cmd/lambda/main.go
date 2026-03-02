package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/periBot/f1-race-leaderboard/internal/cache"
	"github.com/periBot/f1-race-leaderboard/internal/f1client"
	"github.com/periBot/f1-race-leaderboard/internal/handler"
)

func main() {
	// Read configuration from environment variables set by Terraform.
	tableName := os.Getenv("DYNAMODB_TABLE")
	if tableName == "" {
		tableName = "f1-race-cache"
	}

	ttlSeconds := 300 // default 5 minutes
	if v := os.Getenv("CACHE_TTL_SECONDS"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			ttlSeconds = parsed
		}
	}

	// Initialise AWS SDK.
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// Wire up dependencies.
	dynamoClient := dynamodb.NewFromConfig(cfg)
	cacheLayer := cache.New(dynamoClient, tableName, time.Duration(ttlSeconds)*time.Second)
	f1Client := f1client.New()
	h := handler.New(cacheLayer, f1Client)

	log.Printf("Starting F1 Race Leaderboard Lambda (table=%s, ttl=%ds)", tableName, ttlSeconds)

	// Start Lambda runtime.
	lambda.Start(h.HandleRequestWithHealth)
}
