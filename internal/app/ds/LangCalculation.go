/*package ds

import (
	//"database/sql"
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
*/

package ds

import (
	"time"
)

type Lang struct {
	ID            uint   `gorm:"primaryKey;column:id"`
	Name          string `gorm:"column:name"`
	Family        string `gorm:"column:family"`
	Subgroup      string `gorm:"column:subgroup"`
	WritingFamily string `gorm:"type:varchar(50)"`
	Status        string `gorm:"column:status;default:'активен'"` // для удаления языка

	Description string `gorm:"column:description"`
	PhotoKey    string `gorm:"column:photo_key"`
}

func (Lang) TableName() string {
	return "langs"
}

type LangCalculationLanguage struct {
	ID                uint `gorm:"primaryKey;column:id"`
	LangCalculationID uint `gorm:"column:lang_calculation_id"`
	LanguageID        uint `gorm:"column:language_id"`
	IsBase            bool `gorm:"column:is_base"`

	Language Lang `gorm:"foreignKey:LanguageID;references:ID"`
}

func (LangCalculationLanguage) TableName() string {
	return "lang_calculation_languages"
}

type LangCalculation struct {
	ID             uint      `gorm:"primaryKey;column:id"`
	Status         string    `gorm:"column:status"`
	BaseLanguageID uint      `gorm:"column:base_language_id"`
	DateCreate     time.Time `gorm:"column:date_create"`
	DateUpdate     time.Time `gorm:"column:date_update"`
	DateFinish     time.Time `gorm:"column:date_finish"`
	ResearcherID   uint      `gorm:"column:researcher_id"`
	LinguistID     uint      `gorm:"column:linguist_id"`
	ResultYearsAgo int       `gorm:"column:result_years_ago"`
	SimilarityRate float64   `gorm:"column:similarity_rate"`

	Languages []LangCalculationLanguage `gorm:"foreignKey:LangCalculationID;references:ID"`
}

func (LangCalculation) TableName() string {
	return "lang_calculation"
}
