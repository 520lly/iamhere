package modules

const (
	RspOK     int = 0
	RspFailed int = -1
)
const (
	ReasonSuccess         string = "Success"
	ReasonFailureParam    string = "Wrong parameter"
	ReasonMissingParam    string = "Missing parameter"
	ReasonFailureAPIKey   string = "Wrong APIKey"
	ReasonFailueGeneral   string = "Failure in general"
	ReasonDuplicate       string = "Parameter duplicated"
	ReasonInsertFailure   string = "Insert failed"
	ReasonWrongPw         string = "Wrong Password "
	ReasonOperationFailed string = "Operation Failure "
)

type GeoJson struct {
	Type        string    `josn:"type"`
	Coordinates []float64 `json:"coordinates"`
}
