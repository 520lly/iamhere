package db

import(
   "fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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


