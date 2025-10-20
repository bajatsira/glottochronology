package repository

import (
	"fmt"

	"LAB1/internal/app/ds"

	"github.com/sirupsen/logrus"
)

func (r *Repository) GetLangs() ([]ds.Lang, error) {
	var langs []ds.Lang
	err := r.db.Find(&langs).Error
	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	if err != nil {
		return nil, err
	}
	if len(langs) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return langs, nil
}

func (r *Repository) GetLang(id int) (ds.Lang, error) {
	lang := ds.Lang{}
	err := r.db.Where("id = ?", id).First(&lang).Error
	if err != nil {
		return ds.Lang{}, err
	}
	return lang, nil
}

func (r *Repository) GetLangsByName(name string) ([]ds.Lang, error) {
	var langs []ds.Lang
	err := r.db.Where("name ILIKE ?", "%"+name+"%").Find(&langs).Error
	if err != nil {
		return nil, err
	}
	return langs, nil
}

func (r *Repository) GetLangCount() int64 {
	var glottoID uint
	var count int64
	ResearcherID := 1 // пока захардкодили, позже будет из JWT

	// Находим текущий черновик заявки исследователя
	err := r.db.Model(&ds.Glotto{}).
		Where("researcher_id = ? AND status = ?", ResearcherID, "черновик").
		Select("id").
		First(&glottoID).Error
	if err != nil {
		return 0
	}

	// Считаем количество языков в этой заявке
	err = r.db.Model(&ds.GlottoLanguage{}).
		Where("glotto_id = ?", glottoID).
		Count(&count).Error
	if err != nil {
		logrus.Println("Error counting languages in glotto request:", err)
	}

	return count
}
