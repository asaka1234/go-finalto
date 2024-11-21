package finalto

type ResponseProcessorInterface interface {
	RequestForPositionAck(msg RequestForPositionAck)
	PositionReport(msg PositionReport)
}
