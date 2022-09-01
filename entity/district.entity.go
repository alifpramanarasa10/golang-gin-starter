package entity

const (
	// districtTableName is a variable for user table name
	districtTableName = "main.districts"
)

// District entity
type District struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	AltName   string   `json:"alt_name"`
	Latitude  string   `json:"latitude"`
	Longitude string   `json:"longitude"`
	RegencyID int64    `json:"regency_id"`
	Regency   *Regency `gorm:"foreignKey:RegencyID"`
}

// TableName represents table name on db, need to define it because the db has multi schema
func (u *District) TableName() string {
	return districtTableName
}
