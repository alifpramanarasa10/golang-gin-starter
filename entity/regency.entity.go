package entity

const (
	// regencyTableName is a variable for user table name
	regencyTableName = "main.regencies"
)

type ParamsLocation struct {
	Limit   int64  `json:"limit"`
	Offset  int64  `json:"offset"`
	Keyword string `json:"keyword"`
}

// Regency entity
type Regency struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	AltName    string    `json:"alt_name"`
	Latitude   string    `json:"latitude"`
	Longitude  string    `json:"longitude"`
	ProvinceID int64     `json:"province_id"`
	Province   *Province `gorm:"foreignKey:ProvinceID"`
}

// TableName represents table name on db, need to define it because the db has multi schema
func (u *Regency) TableName() string {
	return regencyTableName
}
