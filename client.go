package finalto

import (
	"bytes"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"io"
	"os"
)

type Client struct {
	requestChan    chan quickfix.Messagable
	responseChan   chan *quickfix.Message
	closeCheckChan chan bool

	cfgFileName string
	app         *TradeApplication
	processor   ResponseProcessorInterface

	//过程变量
	appSettings    *quickfix.Settings
	fileLogFactory quickfix.LogFactory
	//结果变量
	initiator *quickfix.Initiator
}

func NewClient(cfgFileName string, processor ResponseProcessorInterface) (*Client, error) {

	result := Client{
		requestChan:    make(chan quickfix.Messagable, 500),
		responseChan:   make(chan *quickfix.Message, 500),
		closeCheckChan: make(chan bool),
		cfgFileName:    cfgFileName,
		processor:      processor,
	}

	//----------------------------------------

	cfg, err := os.Open(result.cfgFileName)
	if err != nil {
		return nil, err
	}
	defer cfg.Close()

	stringData, readErr := io.ReadAll(cfg)
	if readErr != nil {
		return nil, readErr
	}

	//1. 解析cfg配置文件到一个struct里
	appSettings, err := quickfix.ParseSettings(bytes.NewReader(stringData))
	if err != nil {
		return nil, err
	}
	result.appSettings = appSettings

	//
	result.app = NewTradeApplication(&result)

	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		return nil, err
	}
	result.fileLogFactory = fileLogFactory

	//2. 构造了一个跟server的链接，其中配置是在appSettings中
	initiator, err := quickfix.NewInitiator(result.app, quickfix.NewMemoryStoreFactory(), result.appSettings, result.fileLogFactory)
	if err != nil {
		return nil, err
	}
	result.initiator = initiator

	return &result, nil
}

func (cli *Client) Start() {
	err := cli.initiator.Start()
	if err != nil {
		return
	}
	go cli.RunRequest()
	go cli.RunResponse()
}

func (cli *Client) Stop() {
	cli.initiator.Stop()
	cli.closeCheckChan <- true
}

// 发送请求
func (cli *Client) RunRequest() {
	for {
		select {
		case isClose := <-cli.closeCheckChan:
			if isClose {
				return
			}
		case request := <-cli.requestChan:
			quickfix.Send(request)
			return
		}
	}
}

func (cli *Client) RunResponse() {
	for {
		select {
		case isClose := <-cli.closeCheckChan:
			if isClose {
				return
			}
		case msg := <-cli.responseChan:
			msgType, _ := msg.MsgType()
			if msgType == string(enum.MsgType_REQUEST_FOR_POSITIONS_ACK) {
				//TOD 8=FIX.4.49=15135=AO34=249=DemoOrders52=20241121-03:31:08.86356=TH_AL_25249_Demo1=25249581=1710=position4721=f3d3275e-89ad-4926-bd01-eeebb9f3ff42728=1729=210=049
				PosMaintRptID, _ := msg.Body.GetString(721)
				PosReqID, _ := msg.Body.GetString(710)
				TotalNumPosReports, _ := msg.Body.GetInt(727)
				PosReqResult, _ := msg.Body.GetInt(728)
				PosReqStatus, _ := msg.Body.GetInt(729)
				Account, _ := msg.Body.GetString(1)
				AccountType, _ := msg.Body.GetInt(581)
				Symbol, _ := msg.Body.GetString(55)
				Text, _ := msg.Body.GetString(58)

				result := RequestForPositionAck{
					PosMaintRptID:      PosMaintRptID,
					PosReqID:           PosReqID,
					TotalNumPosReports: TotalNumPosReports,
					PosReqResult:       PosReqResult,
					PosReqStatus:       PosReqStatus,
					Account:            Account,
					AccountType:        AccountType,
					Symbol:             Symbol,
					Text:               Text,
				}
				cli.processor.RequestForPositionAck(result)
			} else if msgType == string(enum.MsgType_POSITION_REPORT) {
				PosMaintRptID, _ := msg.Body.GetString(721)
				PosReqID, _ := msg.Body.GetString(710)

				PosReqResult, _ := msg.Body.GetInt(728)
				ClearingBusinessDate, _ := msg.Body.GetString(715)
				Account, _ := msg.Body.GetString(1)
				AccountType, _ := msg.Body.GetInt(581)
				Symbol, _ := msg.Body.GetString(55)
				ContractMultiplier, _ := msg.Body.GetInt(231)
				SettlPrice, _ := msg.Body.GetString(730)
				SettlPriceType, _ := msg.Body.GetInt(731)
				PriorSetlPrice, _ := msg.Body.GetString(734)
				LongPos, _ := msg.Body.GetInt(704)
				ShortPos, _ := msg.Body.GetInt(705)

				result := PositionReport{
					PosMaintRptID:        PosMaintRptID,
					PosReqID:             PosReqID,
					PosReqResult:         PosReqResult,
					ClearingBusinessDate: ClearingBusinessDate,
					Account:              Account,
					AccountType:          AccountType,
					Symbol:               Symbol,
					ContractMultiplier:   ContractMultiplier,
					SettlPrice:           SettlPrice,
					SettlPriceType:       SettlPriceType,
					PriorSetlPrice:       PriorSetlPrice,
					LongPos:              LongPos,
					ShortPos:             ShortPos,
				}
				cli.processor.PositionReport(result)
			}
			return
		}
	}
}
