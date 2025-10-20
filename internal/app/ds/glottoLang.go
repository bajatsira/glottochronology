package ds

type GlottoLanguage struct {
	ID         uint `gorm:"primaryKey"`
	GlottoID   uint `gorm:"not null"`               // Заявка
	LanguageID uint `gorm:"not null"`               // Язык
	IsBase     bool `gorm:"not null;default:false"` // Флаг базового языка

	Language Lang   `gorm:"foreignKey:LanguageID"`
	Glottoh  Glotto `gorm:"foreignKey:GlottoID"`
}
