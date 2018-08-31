package db

import (
	//"fmt"
	//. "github.com/520lly/iamhere/config"
	//. "github.com/spf13/viper"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	//"strings"
)

// Init mgo and the common DAO

// database connection
var Session *mgo.Session

//database sheets definition
var AreaMessages *mgo.Collection
var OceanMessages *mgo.Collection
var Accounts *mgo.Collection
var Areas *mgo.Collection

// sessions
var Sessions *mgo.Collection

//Initialization of DAO
func Init(url, dbname string) {
	if url == "" {
		//TODO load config file
		url = "localhost"
	}
	if dbname == "" {
		dbname = "iamhere"
	}

	// get db config from host, port, username, password
	//url = "mongodb://" + usernameAndPassword + host + ":" + port + "/" + dbname

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	Session.SetMode(mgo.Monotonic, true)

	//AreaMessages
	AreaMessages = Session.DB(dbname).C("areas_messages")

}

func close() {
	Session.Close()
}

func Insert(collection *mgo.Collection, i interface{}) bool {
	err := collection.Insert(i)
	return Err(err)
}
