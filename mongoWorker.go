package main

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type MongoWorker struct {
	MongoSession  *mgo.Session
	MongoDatabase *mgo.Database
}

func NewMongoWorker() *MongoWorker {

	mongoWorker := &MongoWorker{}

	return mongoWorker

}

func (mongoWorker *MongoWorker) Open(host string) (err error) {
	//
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	mongoWorker.MongoSession = session
	database := mongoWorker.MongoSession.DB("notifications")

	mongoWorker.MongoDatabase = database
	return nil
}

func (mongoWorker *MongoWorker) Close() (err error) {
	if mongoWorker.MongoSession != nil {
		mongoWorker.MongoSession.Close()
	}
	return nil
}

func (mongoWorker *MongoWorker) Fetch() (results []Event, err error) {
	collection := mongoWorker.MongoDatabase.C("events")
	err = collection.Find(bson.M{"transport": "smtp"}).All(&results)

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (mongoWorker *MongoWorker) Delete(id bson.ObjectId) (err error) {
	collection := mongoWorker.MongoDatabase.C("events")
	err = collection.Remove(bson.M{"_id": id})
	return nil
}
