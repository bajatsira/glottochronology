package repository

import (
	//"database/sql"
	"errors"
	"fmt"
	"time"

	"LAB1/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *Repository) GetLanguageByID(id int) (ds.GlottoLanguage, error) {
	var language ds.GlottoLanguage
	err := r.db.Preload("Language").First(&language, id).Error
	if err != nil {
		return ds.GlottoLanguage{}, err
	}
	return language, nil
}

/*func (r *Repository) GetLanguageByID(id int) (*ds.GlottoLanguage, error) {
    // Сырой SQL-запрос вручную. Делаем JOIN между таблицей связи и таблицей языков.
    query := `
        SELECT
            gl.id AS glotto_lang_id,
            gl.glotto_id,
            gl.language_id,
            gl.is_base,
            l.id,
            l.name,
            l.family,
            l.subgroup,
            l.writing_family,
            l.description,
            l.photo_url
        FROM glotto_languages gl
        JOIN langs l ON gl.language_id = l.id
        WHERE gl.id = $1
    `

    // Создаём курсор (указатель на строку результата)
    row := r.db.Raw(query, id).Row()

    // Объявляем структуру результата
    result := &ds.GlottoLanguage{}
    var lang ds.Lang

    // Ручное считывание полей
    err := row.Scan(
        &result.ID,
        &result.GlottoID,
        &result.LanguageID,
        &result.IsBase,
        &lang.ID,
        &lang.Name,
        &lang.Family,
        &lang.Subgroup,
        &lang.WritingFamily,
        &lang.Description,
        &lang.PhotoURL,
    )
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil // Если записи нет
        }
        return nil, err
    }

    // Присваиваем язык в связанную структуру
    result.Language = lang

    return result, nil
}
*/

/*
AddServiceToDraft:
- ищет черновик (researcher_id + status='черновик')
- если нет — создаёт новый Glotto (status='черновик')
- затем добавляет запись в glotto_languages (m-m)
- защита от дублей через OnConflict DoNothing (UNIQUE (glotto_id, language_id))
*/
func (r *Repository) AddServiceToDraft(researcherID uint, languageID uint) error {
	var g ds.Glotto

	// 1) попытаться найти существующий черновик
	err := r.db.Where("researcher_id = ? AND status = ?", researcherID, "черновик").First(&g).Error
	if err != nil {
		/*if errors.Is(err, gorm.ErrRecordNotFound) {
			// создаём новый черновик
			g = ds.Glotto{
				Status:         "черновик",
				BaseLanguageID: languageID, // при первом добавлении можно считать этот язык базовым
				DateCreate:     time.Now(),
				ResearcherID:   researcherID,
			}
			logrus.Infof("Черновик не найден, создаю новый для researcher_id=%d", researcherID)

			if err = r.db.Create(&g).Error; err != nil {
				logrus.Errorf("Ошибка при создании черновика: %v", err)
				return err
			}

			logrus.Infof("Черновик успешно создан, ID=%d", g.ID)

		} else {
			return err
		}*/
		if errors.Is(err, gorm.ErrRecordNotFound) {
			g = ds.Glotto{
				Status:         "черновик",
				BaseLanguageID: languageID,
				DateCreate:     time.Now(),
				ResearcherID:   researcherID,
				LinguistID:     1, // <-- здесь указать существующего пользователя из таблицы users
			}
			if err = r.db.Create(&g).Error; err != nil {
				return err
			}
		}

	}

	// 2) добавляем в m-m (glotto_languages), избегаем дублей
	glLang := ds.GlottoLanguage{
		GlottoID:   g.ID,
		LanguageID: languageID,
		IsBase:     false,
	}

	// Используем ON CONFLICT DO NOTHING (если уникальный индекс настроен на (glotto_id, language_id))
	if err := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&glLang).Error; err != nil {
		return err
	}

	return nil
}

// GetDraftByResearcher — получает черновик с подгрузкой языков (и самих Lang)
func (r *Repository) GetDraftByResearcher(researcherID uint) (ds.Glotto, error) {
	var g ds.Glotto
	err := r.db.Preload("Languages.Language").
		Where("researcher_id = ? AND status = ?", researcherID, "черновик").
		First(&g).Error
	if err != nil {
		return ds.Glotto{}, err
	}
	return g, nil
}

// DeleteDraftSQL — логическое удаление заявки через raw SQL UPDATE (требование в методичке)
func (r *Repository) DeleteDraftSQL(glottoID uint) error {
	// Устанавливаем статус и время finished_at
	if err := r.db.Exec("UPDATE glottos SET status = $1, finished_at = now() WHERE id = $2", "удалён", glottoID).Error; err != nil {
		return err
	}
	return nil
}

// Пример альтернативного метода: получение glotto по ID с проверкой статуса (если нужно)
func (r *Repository) GetDraftByID(id uint) (ds.Glotto, error) {
	var g ds.Glotto
	err := r.db.Preload("Languages.Language").First(&g, id).Error
	if err != nil {
		return ds.Glotto{}, err
	}
	return g, nil
}
