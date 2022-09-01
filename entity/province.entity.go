package entity

const (
	// provinceTableName is a variable for user table name
	provinceTableName = "main.provinces"
)

// Province entity
type Province struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	AltName   string `json:"alt_name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// TableName represents table name on db, need to define it because the db has multi schema
func (u *Province) TableName() string {
	return provinceTableName
}
