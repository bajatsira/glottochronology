package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct {
	ID    int
	Title string
}

type Lang struct {
	ID            int
	Name          string
	Family        string
	Title         string
	WritingFamily string
}

func (r *Repository) GetOrders() ([]Order, error) {
	orders := []Order{
		{
			ID:    1,
			Title: "first order",
		},
		{
			ID:    2,
			Title: "second order",
		},
		{
			ID:    3,
			Title: "third order",
		},
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("заказ не найден")
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}

func (r *Repository) GetLang(id int) (Lang, error) {
	langs, err := r.GetLangs()
	if err != nil {
		return Lang{}, err
	}

	for _, lang := range langs {
		if lang.ID == id {
			return lang, nil
		}
	}
	return Lang{}, fmt.Errorf("язык не найден")
}

func (r *Repository) GetLangs() ([]Lang, error) {
	langs := []Lang{
		{
			ID:            1,
			Name:          "Муири",
			Family:        "Нахско-дагестанская",
			Title:         "Даргинская",
			WritingFamily: "Аджам, латиница, кириллица",
		},
		{
			ID:            2,
			Name:          "Кубачинский",
			Family:        "Нахско-дагестанская",
			Title:         "Даргинская",
			WritingFamily: "Аджам, латиница, кириллица",
		},
		{
			ID:            3,
			Name:          "Арабский",
			Family:        "Семитская",
			Title:         "Западносемитская",
			WritingFamily: "Арабский алфавит",
		},
	}
	if len(langs) == 0 {
		return nil, fmt.Errorf("массив языков пустой")
	}

	return langs, nil
}

func (r *Repository) GetLangsByName(name string) ([]Lang, error) {
	langs, err := r.GetLangs()
	if err != nil {
		return []Lang{}, err
	}

	var result []Lang
	for _, lang := range langs {
		if strings.Contains(strings.ToLower(lang.Name), strings.ToLower(name)) {
			result = append(result, lang)
		}
	}

	return result, nil
}

type ChronosData struct {
	BaseLanguage      Lang   // Базовый язык для сравнения (отмечен галочкой)
	SelectedLanguages []Lang // Список выбранных языков для сравнения с базовым
	Comment           string
	Divergence        string // Заглушка для результата расчёта
}

func (r *Repository) GetChronosData() ChronosData {
	langs, err := r.GetLangs()
	if err != nil {
		return ChronosData{}
	}

	// Заглушка: базовый язык - Муири, сравниваемые - Кубачинский
	return ChronosData{
		BaseLanguage: langs[0], // Муири как базовый (отмечен галочкой)
		SelectedLanguages: []Lang{
			langs[1], // Кубачинский
		},
		Comment:    "Анализ расхождения даргинских языков",
		Divergence: "Заглушка: Расхождение составляет 5000 лет",
	}
}
