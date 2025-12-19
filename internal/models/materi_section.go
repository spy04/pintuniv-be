package models

type MateriSection struct {
	ID       uint
	MateriID uint
	Title    string
	Content  string `gorm:"type:longtext"` // markdown
	Order    int
	IsFree   bool `gorm:"default:false"` // 🔑 SATU-SATUNYA LOCK
}
