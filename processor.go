package finalto

import "github.com/quickfixgo/quickfix"

type ResponseProcessorInterface interface {
	RequestForPositionAck(msg *quickfix.Message)
}
