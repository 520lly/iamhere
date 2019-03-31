package db

import (
	"fmt"

	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Init mgo and the common DAO
const (
	QueryOperationAND string = "$and" //and Operation
	QueryOperationOR  string = "$or"  //or Operation
)

type (
	DBCP *mgo.Collection
)

// database connection
var DBSession *mgo.Session
var DBQuery *mgo.Query

//database sheets definition
var DBCAreaMessages *mgo.Collection
var DBCOceanMessages *mgo.Collection
var DBCAccounts *mgo.Collection
var DBCAreas *mgo.Collection

// sessions
var DBSessions *mgo.Collection
var logger echo.Logger

//Initialization of DAO
func Init(url, dbname string, log echo.Logger) {
	logger = log
	logger.Debug("Initialization of DB. url, ", url, ", dbname, ", dbname)
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
	DBSession, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	DBSession.SetMode(mgo.Monotonic, true)

	//AreaMessages
	DBCAreaMessages = DBSession.DB(dbname).C("areas_messages")
	//DBCOceanMessages
	DBCOceanMessages = DBSession.DB(dbname).C("ocean_messages")
	//DBCAccounts
	DBCAccounts = DBSession.DB(dbname).C("accounts")
	//All areas collection
	DBCAreas = DBSession.DB(dbname).C("areas")

}

func CreateNewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}

func CreateGeoIndex(collection *mgo.Collection) error {
	fmt.Println("CreateGeoIndex:")
	// ensure
	// Creating the indexes
	index := mgo.Index{
		Key: []string{"$2dsphere:location"},
	}
	if err := collection.EnsureIndex(index); err != nil {
		fmt.Println("err:", err)
		return err
	}
	return nil
}

func DeleteMessageWithID(id bson.ObjectId) error {
	if c := findCollectionWithID(id); c != nil {
		return deleteItemWithID(c, id)
	}
	return NewError("wrong Message ID")
}

func DeleteAreaWithID(id bson.ObjectId) error {
	return deleteItemWithID(DBCAreas, id)
}

func DeleteAccountWithID(id bson.ObjectId) error {
	return deleteItemWithID(DBCAccounts, id)
}

func deleteItemWithID(collection *mgo.Collection, id bson.ObjectId) error {
	if err := collection.RemoveId(id); err != nil {
		return err
	}
	return nil
}

func findCollectionWithID(id bson.ObjectId) *mgo.Collection {
	var msgs []*Message
	//get all area messages.
	err := DBCAreaMessages.Find(bson.M{"_id": id}).All(&msgs)
	if err == nil {
		if len(msgs) != 0 {
			return DBCAreaMessages
		}
	}
	err = DBCOceanMessages.Find(bson.M{"_id": id}).All(&msgs)
	if err == nil {
		if len(msgs) != 0 {
			return DBCOceanMessages
		}
	}
	return nil

}

func FindAreas(lon float64, lat float64, rad float64) (ret []*Area) {
	var areas []*Area
	if CheckInRangefloat64(lon, LongitudeMinimum, LongitudeMaximum) && CheckInRangefloat64(lat, LatitudeMinimum, LatitudeMaximum) && CheckInRangefloat64(rad, RadiusMinimum, RadiusMaximum) {
		if err := DBCAreas.Find(bson.M{
			"location": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"Type":        "Point",
						"coordinates": []float64{lon, lat},
					},
					"$maxDistance": rad,
				},
			},
		}).All(&areas); err != nil {
			return nil
		}
	} else {
		return nil
	}
	return areas
}

func FindAllUsers() ([]*User, error) {
	return findAllUsers()
}

func findAllUsers() ([]*User, error) {
	var users []*User
	//get all areas
	if err := DBCAccounts.Find(nil).All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

//@name FindMsgWithID
//@brief Get a specific Message with its ID
func FindMsgWithID(id bson.ObjectId) (*Message, error) {
	var msg Message
	if err := DBCAreaMessages.Find(id).All(&msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

//@name FindMsgsWithGeoLocation
//@brief Find messages with one specific geo location
func FindMsgsWithGeoLocation(dbcName string, geo GeoJson, page int, size int) ([]*Message, error) {
	if dbcName == "Ocean" {
		return findMsgsWithGeoLocation(DBCOceanMessages, geo, page, size, Config.ApiConfig.MixAcessDistanceLimit, false)
	} else {
		return findMsgsWithGeoLocation(DBCAreaMessages, geo, page, size, Config.ApiConfig.MixAcessDistanceLimit, false)
	}
}

func findMsgsWithGeoLocation(collection *mgo.Collection, geo GeoJson, page int, size int, distance int, limit bool) ([]*Message, error) {
	var msgs []*Message
	if err := collection.Pipe([]bson.M{
		bson.M{
			"location": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"Type":        "Point",
						"coordinates": []float64{geo.Coordinates[0], geo.Coordinates[1]},
					},
					"$maxDistance": distance,
				},
			},
		},
		bson.M{
			"$match": bson.M{
				"limitaccess": limit,
			},
		},
		bson.M{
			"$sort": bson.M{"timestamp": 1},
		},
		bson.M{
			"$skip": Config.ApiConfig.RandomItemLimit * page,
		},
		bson.M{
			"$sample": bson.M{"size": size},
		},
	}).All(&msgs); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("msgs:", len(msgs))
	return msgs, nil
}

func GetUserMessages(uid string, page int, size int) ([]*Message, error) {
	logger.Debug("userId:[", uid, "], page:[", page, "], size [", size, "]")
	var msgsAreas []*Message
	if err := DBCAreaMessages.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"userid": uid,
			},
		},
		bson.M{
			"$sort": bson.M{"timestamp": 1},
		},
	}).All(&msgsAreas); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("msgsAreas :", len(msgsAreas))
	var msgsOcean []*Message
	if err := DBCOceanMessages.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"userid": uid,
			},
		},
		bson.M{
			"$sort": bson.M{"timestamp": 1},
		},
	}).All(&msgsOcean); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("msgsOcean :", len(msgsOcean))
	msgsTotal := append(msgsAreas, msgsOcean...)
	logger.Debug("msgsTotal :", len(msgsTotal), "Config.ApiConfig.RandomItemLimit * page: ", Config.ApiConfig.RandomItemLimit*page)
	if len(msgsTotal) > Config.ApiConfig.RandomItemLimit*page+size {
		return msgsTotal[Config.ApiConfig.RandomItemLimit*page : Config.ApiConfig.RandomItemLimit*page+size], nil
	} else if len(msgsTotal) > Config.ApiConfig.RandomItemLimit*page {
		return msgsTotal[Config.ApiConfig.RandomItemLimit*page : len(msgsTotal)], nil
	} else {
		return nil, nil
	}
}

//@name FindMsgsWith1Feild
//@brief Find messages with one specific field
func FindMsgsWith1Feild(dbcName string, key string, value string, page int, size int) ([]*Message, error) {
	if dbcName == "Ocean" {
		return findMsgsWith1Feild(DBCOceanMessages, key, value, page, size)
	} else {
		return findMsgsWith1Feild(DBCAreaMessages, key, value, page, size)
	}
}

func findMsgsWith1Feild(collection *mgo.Collection, key string, value string, page int, size int) ([]*Message, error) {
	logger.Debug("key:[", key, "]", "value:[", value, "]")
	var msgs []*Message
	if err := collection.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				key:           value,
				"limitaccess": false,
				//"available": false,
			},
		},
		bson.M{
			"$sort": bson.M{"timestamp": 1},
		},
		bson.M{
			"$skip": Config.ApiConfig.RandomItemLimit * page,
		},
		bson.M{
			"$sample": bson.M{"size": size},
		},
	}).All(&msgs); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("msgs:", len(msgs))
	return msgs, nil
}

//@name FindMsgsWith2Feild
//@brief Find messages with 2 specific field
func FindMsgsWith2Feild(dbcName string, conditions map[string]string, page int, size int) ([]*Message, error) {
	if dbcName == "Ocean" {
		return findMsgsWith2Feild(DBCOceanMessages, conditions, page, size)
	} else {
		return findMsgsWith2Feild(DBCAreaMessages, conditions, page, size)
	}
}

func findMsgsWith2Feild(collection *mgo.Collection, m map[string]string, page int, size int) ([]*Message, error) {
	logger.Debug("m:", m, "  size:", size)
	var msgs []*Message
	if err := collection.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				m["key1"]:     m["value1"],
				m["key2"]:     m["value2"],
				"limitaccess": false,
			},
		},
		bson.M{
			"$sort": bson.M{"timestamp": 1},
		},
		bson.M{
			"$skip": Config.ApiConfig.RandomItemLimit * page,
		},
		bson.M{
			"$sample": bson.M{"size": size},
		},
	}).All(&msgs); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("msgs:", len(msgs))
	return msgs, nil
}

func FindUsersWithFeild(collection *mgo.Collection, m map[string]string) ([]*User, error) {
	logger.Debug("m:", m)
	var users []*User
	if err := collection.Find(
		bson.M{"$or": []bson.M{
			bson.M{m["key1"]: m["value1"]},
			bson.M{m["key2"]: m["value2"]},
			bson.M{m["key3"]: m["value3"]},
		},
		}).All(&users); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("users:", len(users))
	return users, nil
}

func FindUsersWithPW(collection *mgo.Collection, m map[string]string) ([]*User, error) {
	logger.Debug("m:", m)
	var users []*User
	if err := collection.Find(bson.M{
		"$and": []bson.M{
			bson.M{"$or": []bson.M{
				bson.M{m["key1"]: m["value1"]},
				bson.M{m["key2"]: m["value2"]},
				bson.M{m["key3"]: m["value3"]},
			},
			},
			bson.M{m["key4"]: m["value4"]},
		},
	}).All(&users); err != nil {
		logger.Error("err:", err.Error())
		return nil, err
	}
	logger.Debug("users:", len(users))
	return users, nil
}

func FindUserWithID(id bson.ObjectId) *User {
	var user *User
	if err := findItemWithID(DBCAccounts, id, &user); err != nil {
		return nil
	}
	return user
}

func findItemWithID(collection *mgo.Collection, id bson.ObjectId, i interface{}) error {
	if err := collection.FindId(id).One(i); err != nil {
		return err
	}
	return nil
}

func FindAllArea() ([]*Area, error) {
	return findAllArea()
}

func findAllArea() ([]*Area, error) {
	var areas []*Area
	//get all areas
	if err := DBCAreas.Find(nil).All(&areas); err != nil {
		return nil, err
	}
	return areas, nil
}

func GetRandomMessages(collection *mgo.Collection, num int) []*Message {
	var msgs []*Message
	if num <= 0 {
		//don't use limit and return all Message
		if err := collection.Find(nil).All(&msgs); err != nil {
			return nil
		}
	} else {
      //Randomly selects the specified number of documents from the its input
		if err := collection.Pipe([]bson.M{{"$sample": bson.M{"size": num}}}).All(&msgs); err != nil {
			return nil
		}
	}
	logger.Debug("Found msgs:", len(msgs))
	return msgs
}

func GetSpecifiedLocationMessages(collection *mgo.Collection, lon float64, lat float64, rad float64, num int) []*Message {
	var msgs []*Message
	if CheckInRangefloat64(lon, LongitudeMinimum, LongitudeMaximum) && CheckInRangefloat64(lat, LatitudeMinimum, LatitudeMaximum) && CheckInRangefloat64(rad, RadiusMinimum, RadiusMaximum) {
		if err := DBCAreaMessages.Find(bson.M{
			"location": bson.M{
				"$nearSphere": bson.M{
					"$geometry": bson.M{
						"Type":        "Point",
						"coordinates": []float64{lon, lat},
					},
					"$maxDistance": rad,
				},
			},
		}).Limit(num).All(&msgs); err != nil {
			return nil
		}
	} else {
		return nil
	}
	return msgs
}

func Insert(collection *mgo.Collection, i interface{}) bool {
	if err := collection.Insert(i); err != nil {
		logger.Error("UpdateById Err:", err.Error())
		return false
	}
	return true
}

//////////////////////////////////////////////////////////////////////////////
//Update functions
//////////////////////////////////////////////////////////////////////////////
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
	if err := collection.Update(query, i); err != nil {
		logger.Error("UpdateById Err:", err.Error())
		return false
	}
	return true
}

func PushNewCarePoint(collection *mgo.Collection, id bson.ObjectId, value interface{}) bool {
	logger.Debug("value: ", value)
	if err := collection.UpdateId(id, bson.M{"$push": bson.M{"carepoints": bson.M{"$each": value}}}); err != nil {
		logger.Error("UpdateById Err:", err.Error())
		return false
	}
	return true
}

func PullNewCarePoint(collection *mgo.Collection, id bson.ObjectId, value interface{}) bool {
	logger.Debug("value: ", value)
	if err := collection.UpdateId(id, bson.M{"$pull": bson.M{"carepoints": value.([]GeoJson)[0]}}); err != nil {
		logger.Error("UpdateById Err:", err.Error())
		return false
	}
	return true
}

//func Upsert(collection *mgo.Collection, query interface{}, i interface{}) bool {
//    _, err := collection.Upsert(query, i)
//    return Err(err)
//}
//func UpdateAll(collection *mgo.Collection, query interface{}, i interface{}) bool {
//    _, err := collection.UpdateAll(query, i)
//    return Err(err)
//}

func UpdateById(collection *mgo.Collection, id bson.ObjectId, i interface{}) bool {
	//if err := collection.Update(GetIdBsonQ(id), i); err != nil {
	if err := collection.UpdateId(id, i); err != nil {
		logger.Error("UpdateById Err:", err.Error())
		return false
	}
	return true
}

//func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
//    err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
//    return Err(err)
//}

func UpdateByIdField(collection *mgo.Collection, id bson.ObjectId, field string, value interface{}) bool {
	logger.Debug("UpdateByIdField, field: ", field, " <----> value: ", value)
	return UpdateById(collection, id, bson.M{"$set": bson.M{field: value}})
}

//func UpdateByIdAndUserIdMap(collection *mgo.Collection, id, userId string, v bson.M) bool {
//    return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
//}

//func UpdateByIdAndUserIdField2(collection *mgo.Collection, id, userId bson.ObjectId, field string, value interface{}) bool {
//    return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": bson.M{field: value}})
//}
//func UpdateByIdAndUserIdMap2(collection *mgo.Collection, id, userId bson.ObjectId, v bson.M) bool {
//    return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": v})
//}

////
//func UpdateByQField(collection *mgo.Collection, q interface{}, field string, value interface{}) bool {
//    _, err := collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
//    return Err(err)
//}
//func UpdateByQI(collection *mgo.Collection, q interface{}, v interface{}) bool {
//    _, err := collection.UpdateAll(q, bson.M{"$set": v})
//    return Err(err)
//}

func GetIdAndUserIdQ(id bson.ObjectId) bson.M {
	return bson.M{"_id": id}
}
func GetIdBsonQ(id bson.ObjectId) bson.M {
	return bson.M{"_id": id}
}

//DB error handle
func Err(err error) bool {
	if err != nil {
		//fmt.Println(err)
		if err.Error() == "not found" {
			return true
		}
		return false
	}
	return true
}

func Get(collection *mgo.Collection, id bson.ObjectId, i interface{}) error {
	if c, ok := i.(Message); ok {
		logger.Debug("c = ", c)
		var msg Message
		collection.FindId(id).One(&msg)
		i = &msg
		return nil
	}
	return NewError("nil")
}

func GetOneItemWithID(collection *mgo.Collection, id bson.ObjectId, i interface{}) interface{} {
	var ret interface{}
	switch i.(type) {
	case Message:
		var msg Message
		collection.FindId(id).One(msg)
		logger.Debug("msg:", msg)
		ret = &msg
		logger.Debug("msg:", ret)
	case User:
		var user User
		collection.FindId(id).One(&user)
		ret = &user
		logger.Debug("user:", ret)
	}
	return ret
}

// remove duplicated data via field
func Distinct(collection *mgo.Collection, q bson.M, field string, i interface{}) {
	collection.Find(q).Distinct(field, i)
}

//----------------------
func Count(collection *mgo.Collection, q interface{}) int {
	cnt, err := collection.Find(q).Count()
	if err != nil {
		Err(err)
	}
	return cnt
}

func Has(collection *mgo.Collection, q interface{}) bool {
	if Count(collection, q) > 0 {
		return true
	}
	return false
}

func close() {
	DBSession.Close()
}

// checking mognodb lost connection
// check every Operation
func CheckMongoSessionLost() {
	err := DBSession.Ping()
	if err != nil {
		logger.Error("Lost connection to db!")
		DBSession.Refresh()
		err = DBSession.Ping()
		if err == nil {
			logger.Warn("Reconnect to db successful.")
		} else {
			logger.Error("Reconnect to db falied!!!!")
		}
	}
}
