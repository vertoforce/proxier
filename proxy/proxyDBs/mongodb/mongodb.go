// Package mongodb is a proxydb using mongodb
package mongodb

import (
	"context"

	"github.com/vertoforce/proxier/proxy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ProxyDB Is a ProxyDB using mongodb
type ProxyDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

// New Create new MongoDBProxyDB from connectString
func New(ctx context.Context, connectString, database, collection string) (*ProxyDB, error) {
	clientOptions := options.Client().ApplyURI(connectString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return NewFromDB(ctx, client, database, collection)
}

// NewFromDB Create new MongoDBProxyDB from mongo client
func NewFromDB(ctx context.Context, client *mongo.Client, database, collection string) (*ProxyDB, error) {
	err := client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := &ProxyDB{}
	db.client = client
	db.database = client.Database(database)
	db.collection = db.database.Collection(collection)

	return db, nil
}

// GetProxies Get all proxies in the collection
func (db *ProxyDB) GetProxies(ctx context.Context) ([]*proxy.Proxy, error) {
	proxies := []*proxy.Proxy{}

	cursor, err := db.collection.Find(ctx, bson.D{})
	if err == mongo.ErrNoDocuments {
		return proxies, nil
	}
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		p := proxy.Proxy{}
		err = cursor.Decode(&p)
		if err != nil {
			continue
		}

		proxies = append(proxies, &p)
	}

	return proxies, err
}

// StoreProxy Store a proxy in the collection
func (db *ProxyDB) StoreProxy(ctx context.Context, proxy *proxy.Proxy) error {
	_, err := db.collection.InsertOne(ctx, proxy)
	return err
}

// DelProxy Delete a proxy in the collection
func (db *ProxyDB) DelProxy(ctx context.Context, proxy *proxy.Proxy) error {
	_, err := db.collection.DeleteOne(ctx, proxy)
	return err
}

// Clear Delete ALL proxies from collection
func (db *ProxyDB) Clear(ctx context.Context) error {
	_, err := db.collection.DeleteMany(ctx, bson.D{})
	return err
}
