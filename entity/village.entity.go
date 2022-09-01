package entity

const (
	// villageTableName is a variable for user table name
	villageTableName = "main.villages"
)

// Village entity
type Village struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	AltName    string    `json:"alt_name"`
	Latitude   string    `json:"latitude"`
	Longitude  string    `json:"longitude"`
	DistrictID int64     `json:"district_id"`
	District   *District `gorm:"foreignKey:DistrictID"`
}

// TableName represents table name on db, need to define it because the db has multi schema
func (u *Village) TableName() string {
	return villageTableName
}
