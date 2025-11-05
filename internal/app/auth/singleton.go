package auth

// GetCreatorID возвращает фиксированный ID создающего пользователя (singleton)
func GetCreatorID() uint {
	// Можно заменить на чтение из env или config
	return 1
}
