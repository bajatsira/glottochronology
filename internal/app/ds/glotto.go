package ds

type Order struct {
	ID    int
	Title string
}

type Lang struct {
	ID            int    `gorm:"primaryKey"`
	Name          string `gorm:"type:varchar(50);not null"`
	Family        string `gorm:"type:varchar(50);not null"`
	Subgroup      string `gorm:"type:varchar(100)"` // Ветвь языка
	WritingFamily string `gorm:"type:varchar(50)"`
	Description   string `gorm:"type:text"`
	PhotoURL      string `gorm:"type:varchar(255)"`
}
