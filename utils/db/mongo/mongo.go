package mongo

import (
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

var mMd = make(map[string]*mgo.Session, 0)

func InitMongo(DbName, database, uri, name, pwd string, maxPoolSize int, timeout time.Duration) {
	dialInfo := &mgo.DialInfo{
		Addrs:     []string{uri},
		Timeout:   timeout,
		Source:    database,
		Username:  name,
		Password:  pwd,
		PoolLimit: maxPoolSize,
	}
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		println(err.Error())
	}
	mMd[DbName] = s
}

func connect(DbName, db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := mMd[DbName].Copy()
	c := ms.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

func getDb(DbName, db string) (*mgo.Session, *mgo.Database) {
	ms := mMd[DbName].Copy()
	return ms, ms.DB(db)
}

func IsEmpty(DbName, db, collection string) bool {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	count, err := c.Count()
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

func Count(DbName, db, collection string, query interface{}) (int, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

func Insert(DbName, db, collection string, docs ...interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Insert(docs...)
}

func FindOne(DbName, db, collection string, query, selector, result interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).One(result)
}

func FindAll(DbName, db, collection string, query, selector, result interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).All(result)
}

func FindPage(DbName, db, collection string, page, limit int, query, selector, result interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}

func FindIter(DbName, db, collection string, query interface{}) *mgo.Iter {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Find(query).Iter()
}

func Update(DbName, db, collection string, selector, update interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Update(selector, update)
}

func Upsert(DbName, db, collection string, selector, update interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

func UpdateAll(DbName, db, collection string, selector, update interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	_, err := c.UpdateAll(selector, update)
	return err
}

func Remove(DbName, db, collection string, selector interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	return c.Remove(selector)
}

func RemoveAll(DbName, db, collection string, selector interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	_, err := c.RemoveAll(selector)
	return err
}

//insert one or multi documents
func BulkInsert(DbName, db, collection string, docs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Insert(docs...)
	return bulk.Run()
}

func BulkRemove(DbName, db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()

	bulk := c.Bulk()
	bulk.Remove(selector...)
	return bulk.Run()
}

func BulkRemoveAll(DbName, db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.RemoveAll(selector...)
	return bulk.Run()
}

func BulkUpdate(DbName, db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Update(pairs...)
	return bulk.Run()
}

func BulkUpdateAll(DbName, db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.UpdateAll(pairs...)
	return bulk.Run()
}

func BulkUpsert(DbName, db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Upsert(pairs...)
	return bulk.Run()
}

func PipeAll(DbName, db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.All(result)
}

func PipeOne(DbName, db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.One(result)
}

func PipeIter(DbName, db, collection string, pipeline interface{}, allowDiskUse bool) *mgo.Iter {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}

	return pipe.Iter()

}

func Explain(DbName, db, collection string, pipeline, result interface{}) error {
	ms, c := connect(DbName, db, collection)
	defer ms.Close()
	pipe := c.Pipe(pipeline)
	return pipe.Explain(result)
}
func GridFSCreate(DbName, db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(DbName, db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Create(name)
}

func GridFSFindOne(DbName, db, prefix string, query, result interface{}) error {
	ms, d := getDb(DbName, db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).One(result)
}

func GridFSFindAll(DbName, db, prefix string, query, result interface{}) error {
	ms, d := getDb(DbName, db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).All(result)
}

func GridFSOpen(DbName, db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(DbName, db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Open(name)
}

func GridFSRemove(DbName, db, prefix, name string) error {
	ms, d := getDb(DbName, db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Remove(name)
}
