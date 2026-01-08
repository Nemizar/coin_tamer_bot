package telegram

import (
	"fmt"

	"github.com/patrickmn/go-cache"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
)

type UserState string

const (
	UserStateWaitingForCategory UserState = "waiting_for_category"
)

type PendingTransaction struct {
	Amount transaction.Amount
}

func (b *Bot) savePendingTransaction(
	chatID int64,
	amount transaction.Amount,
) {
	b.cache.Set(
		pendingTransactionKey(chatID),
		PendingTransaction{
			Amount: amount,
		},
		cache.DefaultExpiration,
	)

	b.cache.Set(
		userStateKey(chatID),
		UserStateWaitingForCategory,
		cache.DefaultExpiration,
	)
}

func (b *Bot) getUserState(chatID int64) (UserState, error) {
	res, ok := b.cache.Get(userStateKey(chatID))
	if !ok {
		return "", fmt.Errorf("state not found")
	}

	if us, ok := res.(UserState); ok {
		return us, nil
	}

	return "", fmt.Errorf("invalid state")
}

func (b *Bot) getPendingTransaction(chatID int64) (PendingTransaction, error) {
	res, ok := b.cache.Get(pendingTransactionKey(chatID))
	if !ok {
		return PendingTransaction{}, fmt.Errorf("not found pending transaction")
	}

	if pt, ok := res.(PendingTransaction); ok {
		return pt, nil
	}

	return PendingTransaction{}, fmt.Errorf("pending transaction is incorrect")
}

func userStateKey(chatID int64) string {
	return fmt.Sprintf("user-state-%d", chatID)
}

func pendingTransactionKey(chatID int64) string {
	return fmt.Sprintf("pending-transaction-%d", chatID)
}
