package ports

type MsgSender interface {
	Send(msg string) error
}
