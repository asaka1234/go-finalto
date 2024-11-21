package finalto

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type ResponseProcessor struct {
}

func (m ResponseProcessor) RequestForPositionAck(msg RequestForPositionAck) {
	fmt.Printf("-RequestForPositionAck---%+v\n", msg)
}

func (m ResponseProcessor) PositionReport(msg PositionReport) {

}

//----------------------------------------------------------------------------------

func New() *Client {
	processor := ResponseProcessor{}
	client, err := NewClient("./finalto.cfg", processor)
	if err != nil {
		return nil
	}
	client.Start() //启动
	return client
}

func TestRequestForPosition(t *testing.T) {
	resp := New()
	resp.RequestForPosition("position4", "25249")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGHUP)
	<-quit
}
