package ds

import (
	"database/sql"
	"time"
)

type Glotto struct {
	ID             uint      `gorm:"primaryKey"`
	Status         string    `gorm:"type:varchar(20);not null"`
	BaseLanguageID uint      `gorm:"not null"`
	DateCreate     time.Time `gorm:"not null"`
	DateUpdate     time.Time
	DateFinish     sql.NullTime    `gorm:"default:null"`
	ResultYearsAgo sql.NullInt64   `gorm:"default:null"`
	SimilarityRate sql.NullFloat64 `gorm:"default:null"`

	ResearcherID uint `gorm:"not null"`
	LinguistID   uint

	Researcher Users `gorm:"foreignKey:ResearcherID"`
	Linguist   Users `gorm:"foreignKey:LinguistID"`
}
