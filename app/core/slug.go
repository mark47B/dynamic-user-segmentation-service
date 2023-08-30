package core

type Slug struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type SlugRequestAdd struct {
	Name string `json:"name"`
}
