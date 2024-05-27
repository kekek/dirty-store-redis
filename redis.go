package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// RedisDirtyStore is a Redis-based implementation of the DirtyStore interface.
type RedisDirtyStore struct {
	client *redis.Client
	ctx    context.Context
	prefix string
}

// NewRedisDirtyStore creates a new RedisDirtyStore.
func NewRedisDirtyStore(rdb *redis.Client, keyPrefix string, source []string) *RedisDirtyStore {
	ctx := context.Background()
	r := &RedisDirtyStore{
		client: rdb,
		ctx:    ctx,
		prefix: keyPrefix,
	}

	r.Write(source...)

	return r
}

// Write writes sensitive words to the storage.
func (r *RedisDirtyStore) Write(words ...string) error {
	for _, word := range words {
		err := r.client.SAdd(r.ctx, r.getStoreKey(), word).Err()
		if err != nil {
			return err
		}
	}
	r.upVersion()
	return nil
}

// Read reads sensitive words in an iterative manner.
func (r *RedisDirtyStore) Read() <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		iter := r.client.SScan(r.ctx, r.getStoreKey(), 0, "", 0).Iterator()
		for iter.Next(r.ctx) {
			out <- iter.Val()
		}
		if err := iter.Err(); err != nil {
			fmt.Println("Error reading from Redis:", err)
		}
	}()
	return out
}

// ReadAll retrieves all sensitive word data.
func (r *RedisDirtyStore) ReadAll() ([]string, error) {
	words, err := r.client.SMembers(r.ctx, r.getStoreKey()).Result()
	if err != nil {
		return nil, err
	}
	return words, nil
}

// Remove removes sensitive words.
func (r *RedisDirtyStore) Remove(words ...string) error {
	for _, word := range words {
		err := r.client.SRem(r.ctx, r.getStoreKey(), word).Err()
		if err != nil {
			return err
		}
	}
	r.upVersion()
	return nil
}

// Version returns the data storage version number.
func (r *RedisDirtyStore) Version() uint64 {
	ver, _ := r.client.Get(r.ctx, r.getVersionKey()).Uint64()
	return ver
}

func (r *RedisDirtyStore) upVersion() uint64 {
	ver, err := r.client.Incr(r.ctx, r.getVersionKey()).Uint64()
	if err != nil {
		log.Println(err)
	}
	return ver
}

func (r *RedisDirtyStore) getVersionKey() string {
	return r.prefix + "version"
}

func (r *RedisDirtyStore) getStoreKey() string {
	return r.prefix + "store"
}
