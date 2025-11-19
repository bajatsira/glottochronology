package repository

import (
	"errors"
	"fmt"
	"time"

	"LAB1/internal/app/ds"

	"gorm.io/gorm"
)

// GetGlottosFiltered — список заявок (исключая 'черновик' и 'удалён'), фильтр по статусу и диапазону дат формирования (date_create)
func (r *Repository) GetGlottosFiltered(status string, dateFrom, dateTo *time.Time) ([]ds.LangCalculation, error) {
	var glottos []ds.LangCalculation
	q := r.db.Preload("Languages.Language").Model(&ds.LangCalculation{}).
		Where("status <> ?", "удалён").Where("status <> ?", "черновик")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if dateFrom != nil {
		q = q.Where("date_create >= ?", *dateFrom)
	}
	if dateTo != nil {
		q = q.Where("date_create <= ?", *dateTo)
	}
	if err := q.Find(&glottos).Error; err != nil {
		return nil, err
	}
	return glottos, nil
}

// FormGlotto — формирование заявки создателем (проверки минимальные)
func (r *Repository) FormGlotto(id uint, researcherID uint) error {
	var g ds.LangCalculation
	if err := r.db.Preload("Languages").First(&g, id).Error; err != nil {
		return err
	}
	if g.Status != "черновик" {
		return fmt.Errorf("заявка не в статусе черновик")
	}
	// проверка: должно быть хотя бы 1 язык
	var cnt int64
	if err := r.db.Model(&ds.LangCalculationLanguage{}).Where("lang_calculation_id = ?", id).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return fmt.Errorf("недостаточно языков для формирования")
	}
	return r.db.Model(&ds.LangCalculation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      "сформирован",
		"date_update": time.Now(),
		"date_create": g.DateCreate,
	}).Error
}

// CompleteGlotto — завершение/отклонение модератором; action = "завершить" или "отклонить"
func (r *Repository) CompleteGlotto(id uint, moderatorID uint, action string) error {
	var g ds.LangCalculation
	if err := r.db.Preload("Languages").First(&g, id).Error; err != nil {
		return err
	}
	if g.Status != "сформирован" {
		return fmt.Errorf("только сформированную заявку можно завершить/отклонить")
	}
	now := time.Now()
	updates := map[string]interface{}{
		"linguist_id": moderatorID,
		"date_finish": now,
		"date_update": now,
	}
	if action == "отклонить" {
		updates["status"] = "отклонён"
	} else {
		updates["status"] = "завершён"
		// Пример вычисления: similarity и result_years_ago (заглушки — замените формулой из лабы)
		updates["similarity_rate"] = 0.0
		updates["result_years_ago"] = 0
	}
	return r.db.Model(&ds.LangCalculation{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteGlottoLanguage — удалить связь glotto_languages по glotto_id + language_id
func (r *Repository) DeleteGlottoLanguage(glottoID, languageID uint) error {
	return r.db.Where("lang_calculation_id = ? AND language_id = ?", glottoID, languageID).Delete(&ds.LangCalculationLanguage{}).Error
}

// UpdateGlotto — обновить поля glotto (без системных)
func (r *Repository) UpdateGlotto(id uint, updates map[string]interface{}) error {
	updates["date_update"] = time.Now()
	return r.db.Model(&ds.LangCalculation{}).Where("id = ?", id).Updates(updates).Error
}

// GetLangCalculationsByResearcher возвращает все расчёты (заявки) конкретного исследователя
// с подгрузкой связанных языков (через LangCalculationLanguage → Lang)
// и с фильтром по статусу (если status != "" и не пустой)
func (r *Repository) GetLangCalculationsByResearcher(researcherID uint, status string) ([]ds.LangCalculation, error) {
	var calculations []ds.LangCalculation

	// Базовый запрос с Preload всех связанных сущностей
	query := r.db.
		Preload("Languages.Language").
		Where("researcher_id = ?", researcherID)

	// Если передан конкретный статус — фильтруем по нему
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Исключаем удалённые и (по желанию) черновики, если они не нужны в этом списке
	query = query.Where("status != ?", "удалён")

	// Выполняем запрос
	if err := query.Find(&calculations).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return calculations, nil
}

func (r *Repository) CompleteLangCalculation(id uint, moderatorID uint, action string) error {
	var calc ds.LangCalculation

	// Ищем расчёт и сразу проверяем статус
	result := r.db.First(&calc, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("calculation not found")
		}
		return result.Error
	}

	// Разрешаем только из статуса "на модерации"
	if calc.Status != "на модерации" {
		return fmt.Errorf("invalid status transition")
	}

	// Обновляем поля
	now := time.Now()
	if action == "завершить" {
		calc.Status = "завершён"
		calc.DateFinish = now
	} else if action == "отклонить" {
		calc.Status = "отклонён"
	}

	calc.LinguistID = moderatorID

	calc.DateUpdate = now

	return r.db.Save(&calc).Error
}
