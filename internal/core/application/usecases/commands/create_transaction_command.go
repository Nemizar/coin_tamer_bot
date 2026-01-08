package commands

import (
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/shared"
	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/transaction"
)

type CreateTransactionCommand interface {
	UserID() shared.ID
	Amount() transaction.Amount
	CategoryID() shared.ID
}

type createTransactionCommand struct {
	userID     shared.ID
	amount     transaction.Amount
	categoryID shared.ID
}

func NewCreateTransactionCommand(userID shared.ID, amount transaction.Amount, categoryID shared.ID) (CreateTransactionCommand, error) {
	return &createTransactionCommand{userID: userID, amount: amount, categoryID: categoryID}, nil
}

func (c createTransactionCommand) UserID() shared.ID {
	return c.userID
}

func (c createTransactionCommand) Amount() transaction.Amount {
	return c.amount
}

func (c createTransactionCommand) CategoryID() shared.ID {
	return c.categoryID
}
