package db

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://admin:<password>@cluster0.nqn8h.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

func GetClient(ctx context.Context) (*mongo.Client, error) {
  client, _ := mongo.NewClient(options.Client().ApplyURI(connectionString))
  err := client.Connect(ctx)
  if err == nil {
      err = client.Ping(ctx, nil)
  }
  if err == nil {
      return client, nil
  }
  return nil, err
}
