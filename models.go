package main

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
)

type Event struct {
	Id         bson.ObjectId `json:"id"              bson:"_id"`
	AppId      string        `json:"appid"           bson:"appid"`
	AsteriskId string        `json:"asteriskid"	     bson:"asteriskid"`
	Name       string        `json:"name"			 bson:"name"`
	Data       string        `json:"data"			 bson:"data"`
	Type       int           `json:"type"			 bson:"type"`
}

func (event *Event) String() string {
	datas, _ := json.MarshalIndent(event, "", "")
	return string(datas)
}

func (event *Event) Json() []byte {
	datas, _ := json.MarshalIndent(event, "", "")
	return datas
}
