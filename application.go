package finalto

import (
	"fmt"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
)

// TradeClient implements the quickfix.Application interface
type TradeApplication struct {
	isLogon bool
	cli     *Client
}

func NewTradeApplication(cli *Client) *TradeApplication {
	return &TradeApplication{
		cli: cli,
	}
}

// OnCreate implemented as part of Application interface
func (e TradeApplication) OnCreate(sessionID quickfix.SessionID) {
	fmt.Printf("----create----\n")
}

// OnLogon implemented as part of Application interface
func (e TradeApplication) OnLogon(sessionID quickfix.SessionID) {
	fmt.Printf("----logon----%s\n", sessionID.String())
	e.isLogon = true
}

// OnLogout implemented as part of Application interface
func (e TradeApplication) OnLogout(sessionID quickfix.SessionID) {
	fmt.Printf("----logout----%s\n", sessionID.Qualifier)
}

// FromAdmin implemented as part of Application interface
func (e TradeApplication) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return nil
}

// ToAdmin implemented as part of Application interface
func (e TradeApplication) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	msgType, _ := msg.MsgType()
	uName, _ := e.cli.appSettings.GlobalSettings().Setting("Username")
	uPwd, _ := e.cli.appSettings.GlobalSettings().Setting("Password")
	if msgType == string(enum.MsgType_LOGON) {
		msg.Body.Set(field.NewUsername(uName))
		msg.Body.Set(field.NewPassword(uPwd))
	}
}

// ToApp implemented as part of Application interface
func (e TradeApplication) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf(fmt.Sprintf("Sending: %s", msg.String()))
	return
}

// FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e TradeApplication) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	e.cli.responseChan <- msg
	return
}
