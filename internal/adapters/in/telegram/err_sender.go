package telegram

func (b *Bot) sendValidationError(chatID int64) {
	if err := b.sendMsg(chatID, "Неверная сумма транзакции"); err != nil {
		b.logger.Error(
			"Ошибка отправки сообщения о валидации суммы",
			"err", err.Error(),
		)
	}
}

func (b *Bot) sendCategoriesError(chatID int64) {
	if err := b.sendMsg(
		chatID,
		"Не удалось получить категории. Повторите попытку добавления транзакции",
	); err != nil {
		b.logger.Error(
			"Ошибка отправки сообщения о получении категорий",
			"err", err.Error(),
		)
	}
}
