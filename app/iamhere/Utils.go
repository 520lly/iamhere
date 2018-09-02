package iamhere

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	. "github.com/520lly/iamhere/app/modules"
	"gopkg.in/mgo.v2/bson"
)

func JsonToString(i interface{}) string {
	ret, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return string(ret)
}

func NewError(context string) error {
	return errors.New(context)
}

func CheckInRangeInt(num int, bottom int, top int) (ret bool) {
	return num > bottom && num < top
}

func CheckInRangefloat64(num float64, bottom float64, top float64) (ret bool) {
	return num > bottom && num < top
}

func CreateTimeStampUnix() int64 {
	return time.Now().Unix()
}

func ConvertString2Float64(s string) (float64, error) {
	var ret float64 = 0
	if len(s) != 0 {
		ret, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return ret, err
		}
	} else {
		return ret, NewError("empty string")
	}
	return ret, nil
}
func BsonToString(b bson.ObjectId) string {
	return string(b.Hex())
}

func StringToBson(s string) bson.ObjectId {
	return bson.ObjectIdHex(s)
}

func ValidateAreaCategory(category int) bool {
	return (category > CategoryMinimum && category < CategoryMaximum)
}

func ValidateAreaType(t int) bool {
	return (t > TypeMinimum && t < TypeMaximum)
}
