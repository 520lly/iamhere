package iamhere

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"time"

	. "github.com/520lly/iamhere/app/modules"
	"github.com/dgrijalva/jwt-go"
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

func CreateRandomNickname() string {
	return string(krand(12, 3))
}

func GetJWTSecretCode() []byte {
	return []byte(Config.ApiConfig.Secret)
}

func CreateNewJWTToken() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

//convert a string to hash in MD5
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//Check string empty
func CheckStringNotEmpty(str string) bool {
	return (len(str) != 0)
}

//Check bson Object empty
func CheckBsonObjNotEmpty(b bson.ObjectId) bool {
	return (len(b.Hex()) != 0)
}

//convert string to bson.ObjectId
func ConvertString2BsonObjectId(str string) bson.ObjectId {
	return bson.ObjectIdHex(str)
}

/**
* size random size
* kind 0    // pure number
       1    // low class alphabet
       2    // upper class alphabet
       3    // number、lower, upper class
*/
func krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
