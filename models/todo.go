package models

type Todo struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `json:"title"`
	IsDone    bool   `json:"is_done"`
	ImagePath string `json:"image_path"`
}
