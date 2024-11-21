package finalto

import (
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44rfp "github.com/quickfixgo/fix44/requestforpositions"
	"time"
)

func (cli *Client) RequestForPosition() {
	cusZone := time.FixedZone("UTC", 0*60*60)
	ctime := time.Now().In(cusZone)
	cdate := ctime.Format("20060102-15:04:05")

	request := fix44rfp.New(field.NewPosReqID("position4"), //710
		field.NewPosReqType(enum.PosReqType_POSITIONS), //724
		field.NewAccount("25249"),                      //1
		field.NewAccountType(enum.AccountType_ACCOUNT_IS_CARRIED_ON_CUSTOMER_SIDE_OF_THE_BOOKS), //581
		field.NewClearingBusinessDate(cdate),                                                    //715  20060826-15:31:22
		field.NewTransactTimeNoMillis(ctime),                                                    //60   20060826-15:31:22
	)

	relatedSym := fix44rfp.NewNoPartyIDsRepeatingGroup()
	request.SetNoPartyIDs(relatedSym) //453

	request.Header.SetMsgType(enum.MsgType_REQUEST_FOR_POSITIONS) //35
	//request.Header.SetMsgSeqNum(123)              //34

	senderCompID, _ := cli.appSettings.GlobalSettings().Setting("SenderCompID")
	request.Header.SetSenderCompID(senderCompID) //49

	targetCompID, _ := cli.appSettings.GlobalSettings().Setting("TargetCompID")
	request.Header.SetTargetCompID(targetCompID) //56

	request.Header.SetSendingTime(ctime) //52

	//---------------------------------------
	cli.requestChan <- request
}
