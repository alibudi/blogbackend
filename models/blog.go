package models

type Blog struct {
    Id        uint   `json:"id"`
    Title string `json:"title"`
    Desc  string `json:"desc"`
    Image  string `json:"image"`
	UserID string `Json:"userid"`
	User User `Json:"user";gorm:"foreignkey:UserID"`
}
