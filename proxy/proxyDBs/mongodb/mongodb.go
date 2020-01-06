// Package mongodb is a proxydb using mongodb
package mongodb

import (
	"context"
	"proxy/proxy"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBProxyDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func New(ctx context.Context, connectString, database, collection string) (*MongoDBProxyDB, error) {
	clientOptions := options.Client().ApplyURI(connectString)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	db := &MongoDBProxyDB{}
	db.client = client
	db.database = client.Database(database)
	db.collection = db.database.Collection(collection)

	return db, nil
}

func (db *MongoDBProxyDB) GetProxies(ctx context.Context) ([]proxy.Proxy, error) {
	proxies := []proxy.Proxy{}

	cursor, err := db.collection.Find(ctx, bson.D{})
	if err == mongo.ErrNoDocuments {
		return proxies, nil
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		p := proxy.Proxy{}
		err = cursor.Decode(&p)
		if err != nil {
			continue
		}

		proxies = append(proxies, p)
	}

	return proxies, err
}

func (db *MongoDBProxyDB) StoreProxy(ctx context.Context, proxy proxy.Proxy) error {
	_, err := db.collection.InsertOne(ctx, proxy)
	return err
}

func (db *MongoDBProxyDB) DelProxy(ctx context.Context, proxy proxy.Proxy) error {
	_, err := db.collection.DeleteOne(ctx, proxy)
	return err
}
