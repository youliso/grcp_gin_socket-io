package mongo

import (
	"context"
	"github.com/qiniu/qmgo"
)

type mgoS struct {
	db  *qmgo.Client
	ctx *context.Context
}

var mMd = make(map[string]mgoS, 0)

func InitMongo(DbName, uri string, maxPoolSize *uint64) {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{
		Uri:         uri,
		MaxPoolSize: maxPoolSize,
	})
	if err != nil {
		println(err.Error())
	}
	mMd[DbName] = mgoS{
		db:  client,
		ctx: &ctx,
	}
}

func connect(DbName, db string) (*qmgo.Database, *context.Context) {
	ms := mMd[DbName].db.Database(db)
	return ms, mMd[DbName].ctx
}

func Count(DbName, db, collection string, query interface{}) (int64, error) {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	//defer ms.Close()
	return coll.Find(*ctx, query).Count()
}

func InsertOne(DbName, db, collection string, doc interface{}) (*qmgo.InsertOneResult, error) {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.InsertOne(*ctx, doc)
}

func InsertMany(DbName, db, collection string, docs ...interface{}) (*qmgo.InsertManyResult, error) {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.InsertMany(*ctx, docs)
}

func FindOne(DbName, db, collection string, query, selector, result interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.Find(*ctx, query).Select(selector).One(result)
}

func FindAll(DbName, db, collection string, query, selector, result interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.Find(*ctx, query).Select(selector).All(result)
}

func FindPage(DbName, db, collection string, page, limit int64, query, selector, result interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.Find(*ctx, query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}

func Update(DbName, db, collection string, selector, update interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.Update(*ctx, selector, update)
}

func Upsert(DbName, db, collection string, selector, update interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	_, err := coll.Upsert(*ctx, selector, update)
	return err
}

func UpdateAll(DbName, db, collection string, selector, update interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	_, err := coll.UpdateAll(*ctx, selector, update)
	return err
}

func Remove(DbName, db, collection string, selector interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	return coll.Remove(*ctx, selector)
}

func RemoveAll(DbName, db, collection string, selector interface{}) error {
	ms, ctx := connect(DbName, db)
	coll := ms.Collection(collection)
	_, err := coll.DeleteAll(*ctx, selector)
	return err
}
