package finalto

import (
	"bytes"
	"fmt"
	"github.com/quickfixgo/quickfix"
	"io"
	"os"
)

type Client struct {
	requestChan    chan FixApiRequest
	closeCheckChan chan bool

	cfgFileName string
	app         *TradeApplication

	//过程变量
	appSettings    *quickfix.Settings
	fileLogFactory quickfix.LogFactory
	//结果变量
	initiator *quickfix.Initiator
}

func NewClient(cfgFileName string) (*Client, error) {

	result := Client{
		requestChan:    make(chan FixApiRequest, 500),
		closeCheckChan: make(chan bool),
		cfgFileName:    cfgFileName,
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
	result.app = NewTradeClient(&result)

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
		case <-cli.requestChan:
			//收到请求了就发送出去
			fmt.Printf("Dispatch pushPing Stop! ")
			return
		}
	}
}

// 接收返回并解析到不同的处理器去
func (cli *Client) RunResponse() {
	for {
		select {
		case isClose := <-cli.closeCheckChan:
			if isClose {
				return
			}
		case <-cli.requestChan:
			//收到请求了就发送出去
			fmt.Printf("Dispatch pushPing Stop! ")
			return
		}
	}
}
