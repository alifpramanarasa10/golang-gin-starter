package resource

import "gin-starter/entity"

// Province is a struct for province
type Province struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ProvinceListResponse is a struct for province list response
type ProvinceListResponse struct {
	List  []*Province `json:"list"`
	Total int64       `json:"total"`
}

// NewProvinceResponse create new NewProvinceResponse
func NewProvinceResponse(province *entity.Province) *Province {
	return &Province{
		ID:   province.ID,
		Name: province.Name,
	}
}

// Regency is a struct for regency
type Regency struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetRegencyByProvinceIDRequest is a struct for get regency by province id request
type GetRegencyByProvinceIDRequest struct {
	ProvinceID int64 `uri:"province_id" json:"province_id" binding:"required"`
}

// RegencyListResponse is a struct for regency list response
type RegencyListResponse struct {
	List  []*Regency `json:"list"`
	Total int64      `json:"total"`
}

// NewRegencyResponse create new NewRegencyResponse
func NewRegencyResponse(regency *entity.Regency) *Regency {
	return &Regency{
		ID:   regency.ID,
		Name: regency.Name,
	}
}

// District is a struct for district
type District struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetDistrictByRegencyIDRequest is a struct for get district by regency id request
type GetDistrictByRegencyIDRequest struct {
	RegencyID int64 `uri:"regency_id" json:"regency_id" binding:"required"`
}

// DistrictListResponse is a struct for district list response
type DistrictListResponse struct {
	List  []*District `json:"list"`
	Total int64       `json:"total"`
}

// NewDistrictResponse create new NewDistrictResponse
func NewDistrictResponse(district *entity.District) *District {
	return &District{
		ID:   district.ID,
		Name: district.Name,
	}
}

// Village is a struct for village
type Village struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetVillageByDistrictIDRequest is a struct for get village by district id request
type GetVillageByDistrictIDRequest struct {
	DistrictID int64 `uri:"district_id" json:"district_id" binding:"required"`
}

// VillageListResponse is a struct for village list response
type VillageListResponse struct {
	List  []*Village `json:"list"`
	Total int64      `json:"total"`
}

// NewVillageResponse create new NewVillageResponse
func NewVillageResponse(village *entity.Village) *Village {
	return &Village{
		ID:   village.ID,
		Name: village.Name,
	}
}
