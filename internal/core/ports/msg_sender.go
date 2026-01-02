package ports

import (
	"context"

	"github.com/Nemizar/coin_tamer_bot/internal/core/domain/models/user"
)

type MsgSender interface {
	Send(
		ctx context.Context,
		identity user.User,
		message string,
	) error
}
