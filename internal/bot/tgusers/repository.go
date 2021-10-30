package tgusers

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	tb "gopkg.in/tucnak/telebot.v2"
)

const collectionName = "tgusers"

type Repository struct {
	mongoURI   string
	dbName     string
	collection *mongo.Collection
}

type Connection struct {
	Ctx        context.Context
	Client     *mongo.Client
	Disconnect func()
	Collection *mongo.Collection
}

func NewRepository(mongoUrl, dbName string) (*Repository, error) {
	r := &Repository{
		mongoURI: mongoUrl,
		dbName:   dbName,
	}
	con, err := r.getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Disconnect()
	err = con.Client.Ping(con.Ctx, readpref.Primary())
	return r, err
}

func (r *Repository) getConnection() (*Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(r.mongoURI))
	disconnect := func() {
		err := client.Disconnect(ctx)
		cancel()
		if err != nil {
			log.Println(err.Error())
		}
	}
	return &Connection{
		Ctx:        ctx,
		Client:     client,
		Disconnect: disconnect,
		Collection: client.Database(r.dbName).Collection(collectionName),
	}, err
}

func (r Repository) Ping() error {
	con, err := r.getConnection()
	if err != nil {
		return err
	}
	defer con.Disconnect()
	return con.Client.Ping(con.Ctx, readpref.Primary())
}

func (r *Repository) Save(user *tb.User) error {
	con, err := r.getConnection()
	if err != nil {
		return err
	}
	defer con.Disconnect()
	_, err = con.Collection.InsertOne(con.Ctx, user)
	return err
}

func (r Repository) FindAll() ([]*tb.User, error) {
	con, err := r.getConnection()
	if err != nil {
		return nil, err
	}
	defer con.Disconnect()
	cur, err := con.Collection.Find(con.Ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var users []*tb.User
	err = cur.All(con.Ctx, &users)
	return users, err
}
