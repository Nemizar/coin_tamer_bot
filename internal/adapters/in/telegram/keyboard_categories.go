package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/category"
)

func newCategoriesInlineKeyboard(categories []*category.Category) tgbotapi.InlineKeyboardMarkup {
	var keyboardRows [][]tgbotapi.InlineKeyboardButton
	row := make([]tgbotapi.InlineKeyboardButton, 0, 3)

	for _, c := range categories {
		button := tgbotapi.NewInlineKeyboardButtonData(c.Name(), c.ID().String())
		row = append(row, button)

		if len(row) == 3 {
			keyboardRows = append(keyboardRows, row)
			row = make([]tgbotapi.InlineKeyboardButton, 0, 3)
		}
	}

	if len(row) > 0 {
		keyboardRows = append(keyboardRows, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)
}
