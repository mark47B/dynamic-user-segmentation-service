package core

type Slug struct {
	ID   uint32 `json:"id" gorm:"primary_key"`
	Name string `json:"name"`
}
