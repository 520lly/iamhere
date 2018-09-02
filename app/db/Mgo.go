package db

import (
	"fmt"
	//. "github.com/520lly/iamhere/config"
	//. "github.com/spf13/viper"
	. "github.com/520lly/iamhere/app/iamhere"
	. "github.com/520lly/iamhere/app/modules"
	"github.com/labstack/gommon/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"strings"
)

// Init mgo and the common DAO

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
var logger log.Logger

//Initialization of DAO
func Init(url, dbname string) {
	logger.SetLevel(log.DEBUG)
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

//func FindAreaWithLocation(lon float64, lat float64) (ret *Area) {
//    if CheckInRangefloat64(lon, LongitudeMinimum, LongitudeMaximum) && CheckInRangefloat64(lat, LatitudeMinimum, LatitudeMaximum) {
//        areas := findAllArea()
//        for _, area := range areas {
//            var areaMatchs []*Area
//            if err := DBCAreas.Find(bson.M{
//                "location": bson.M{
//                    "$nearSphere": bson.M{
//                        "$geometry": bson.M{
//                            "Type":        "Point",
//                            "coordinates": []float64{lon, lat},
//                        },
//                        "$maxDistance": area.Radius,
//                    },
//                },
//            }).All(&areaMatchs); err != nil {
//                return nil
//            }
//            if len(areaMatchs) == 1 && areaMatchs[0].ID == area.ID {
//                return areaMatchs[0]
//            }
//        }
//    } else {
//        return nil
//    }
//    return nil
//}

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
	if err := collection.Find(nil).Limit(num).All(&msgs); err != nil {
		return nil
	}
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
	err := collection.Insert(i)
	return Err(err)
}

//////////////////////////////////////////////////////////////////////////////
//Update functions
//////////////////////////////////////////////////////////////////////////////
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
	err := collection.Update(query, i)
	return Err(err)
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
	err := collection.Update(GetIdBsonQ(id), i)
	return Err(err)
}

//func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
//    err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
//    return Err(err)
//}

func UpdateByIdField(collection *mgo.Collection, id bson.ObjectId, field string, value interface{}) bool {
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
	return collection.FindId(id).One(i)
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
