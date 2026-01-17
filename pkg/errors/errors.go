package errors

type Kind int

const (
	UserError Kind = iota + 1
	SystemError
	InternalError
)

type Error struct {
	Kind    Kind
	Message string
}

func (e Error) Error() string {
	return e.Message
}
