package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/periBot/f1-race-leaderboard/internal/models"
)

// DynamoClient defines the subset of the DynamoDB API we use (for testing).
type DynamoClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// Cache provides TTL-based caching backed by DynamoDB.
type Cache struct {
	client    DynamoClient
	tableName string
	ttl       time.Duration
}

// New creates a new DynamoDB-backed cache.
func New(client DynamoClient, tableName string, ttl time.Duration) *Cache {
	return &Cache{
		client:    client,
		tableName: tableName,
		ttl:       ttl,
	}
}

// Get retrieves a cached item by key. Returns the raw JSON data and true if
// found and not expired, or nil and false on miss/expiry.
func (c *Cache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	out, err := c.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.tableName),
		Key: map[string]types.AttributeValue{
			"cache_key": &types.AttributeValueMemberS{Value: key},
		},
	})
	if err != nil {
		return nil, false, fmt.Errorf("dynamodb GetItem: %w", err)
	}

	if out.Item == nil {
		return nil, false, nil
	}

	var item models.CacheItem
	if err := attributevalue.UnmarshalMap(out.Item, &item); err != nil {
		return nil, false, fmt.Errorf("unmarshalling cache item: %w", err)
	}

	// DynamoDB TTL deletion is eventually consistent (up to 48h delay),
	// so we also check expiry in application code.
	if time.Now().Unix() > item.TTL {
		return nil, false, nil
	}

	return []byte(item.Data), true, nil
}

// Put stores data in the cache with the configured TTL.
func (c *Cache) Put(ctx context.Context, key string, data []byte) error {
	now := time.Now()
	item := models.CacheItem{
		CacheKey:  key,
		Data:      string(data),
		TTL:       now.Add(c.ttl).Unix(),
		FetchedAt: now.Unix(),
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("marshalling cache item: %w", err)
	}

	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("dynamodb PutItem: %w", err)
	}

	return nil
}
