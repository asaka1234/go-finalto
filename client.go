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
				cli.processor.RequestForPositionAck(msg)
			}
			return
		}
	}
}
