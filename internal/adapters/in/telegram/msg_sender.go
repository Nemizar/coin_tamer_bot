package telegram

import (
	"encoding/json"
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) sendMsg(chatID int64, text string) error {
	_, err := b.bot.Send(tgbotapi.NewMessage(chatID, text))
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) deleteMessage(chatID int64, messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)

	var unmarshalTypeError *json.UnmarshalTypeError

	_, err := b.bot.Send(deleteMsg)
	if err != nil && !errors.As(err, &unmarshalTypeError) {
		return err
	}

	return nil
}

func (b *Bot) sendReplyMarkup(chatID int64, text string, replyMarkup *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = replyMarkup

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) sendMessageAndDeleteInlineKeyboard(chatID int64, prevMsgID int, text string) error {
	err := b.sendMsg(chatID, text)
	if err != nil {
		return err
	}

	err = b.deleteMessage(chatID, prevMsgID)
	if err != nil {
		return err
	}

	return nil
}
