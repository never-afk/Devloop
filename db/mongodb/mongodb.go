package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Client struct {
	*mongo.Client
	conf *Config
}

type CollectionHandel struct {
	*mongo.Collection
	timeOut time.Duration
}

func (c *CollectionHandel) ctx() context.Context {
	return ctx(c.timeOut)
}

func (c *CollectionHandel) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.Collection.InsertOne(c.ctx(), document, opts...)
}

func (c *CollectionHandel) InsertMany(documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.Collection.InsertMany(c.ctx(), documents, opts...)
}

func (c *CollectionHandel) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.Collection.FindOne(c.ctx(), filter, opts...)
}

func (c *CollectionHandel) Find(filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return c.Collection.Find(c.ctx(), filter, opts...)
}

func (c *CollectionHandel) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteOne(c.ctx(), filter, opts...)
}

func (c *CollectionHandel) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.Collection.DeleteMany(c.ctx(), filter, opts...)
}

func (c *CollectionHandel) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.Collection.UpdateOne(c.ctx(), filter, update, opts...)
}

func (c *CollectionHandel) UpdateMany(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.Collection.UpdateMany(c.ctx(), filter, update, opts...)
}

func (c *CollectionHandel) Upsert(data interface{}) (*mongo.UpdateResult, error) {
	q, u, err := parseData(data)

	if err != nil {
		return nil, err
	}
	opt := options.Update().SetUpsert(true)
	return c.Collection.UpdateOne(c.ctx(), q, u, opt)
}

func (c *CollectionHandel) UpsertMany(data []interface{}) (*mongo.BulkWriteResult, error) {

	update := make([]mongo.WriteModel, 0)
	for _, v := range data {
		q, u, err := parseData(v)
		if err != nil {
			return nil, err
		}
		update = append(update, mongo.NewUpdateOneModel().SetFilter(q).SetUpdate(u).SetUpsert(true))
	}

	opts := options.BulkWrite().SetOrdered(false)
	return c.BulkWrite(c.ctx(), update, opts)
}

// 解析数据
func parseData(data interface{}) (q bson.M, u bson.M, err error) {
	update := bson.M{}
	switch data.(type) {
	case bson.M:
		update = data.(bson.M)
	case bson.D:
		update = data.(bson.D).Map()
	default:
		var b []byte
		if b, err = bson.Marshal(data); err == nil {
			err = bson.Unmarshal(b, &update)
		}
	}

	// 无错误的情况下才处理
	if err == nil {
		q = bson.M{"_id": update["_id"]}
		// 主键不可更改
		delete(update, "_id")

		u = bson.M{}

		// 处理操作符号
		for k, v := range update {
			if strings.HasPrefix(k, "$") {
				u[k] = v
				delete(update, k)

				switch v.(type) {

				case bson.M:
					for field, _ := range v.(bson.M) {
						delete(update, field)
					}
				case bson.D:
					for _, d := range v.(bson.D) {
						delete(update, d.Key)
					}
				}

			}
		}

		// 数据不为空才放入set
		if len(update) > 0 {
			u["$set"] = update
		}
	}
	return
}

type Config struct {
	TimeOut      time.Duration
	DatabaseName string
	DSN          string
}

func Connect(config *Config) (c *Client, err error) {
	c = &Client{
		conf: config,
	}
	c.Client, err = mongo.Connect(ctx(c.conf.TimeOut), options.Client().ApplyURI(config.DSN))
	return
}

// 实例化操作对象
func (c *Client) Table(tableName string) *CollectionHandel {
	return &CollectionHandel{
		Collection: c.Database(c.conf.DatabaseName).Collection(tableName),
		timeOut:    c.conf.TimeOut,
	}
}

func ctx(duration time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), duration)
	return ctx
}
