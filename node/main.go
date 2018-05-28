package main

import (
	"datx_chain/node/msg"
	"datx_chain/plugins/chain_plugin"
	"datx_chain/plugins/http_plugin"
	"datx_chain/plugins/p2p_plugin"
	"fmt"
	"io/ioutil"
	"log"
)

const messageId = 0

func MyProtocol() p2p.Protocol {
	return p2p.Protocol{
		Name:    "MyProtocol",
		Version: 1,
		Length:  1,
		Run:     msgHandler,
	}
}

func main() {
	// nodekey, _ := crypto.GenerateKey()
	// srv := p2p.Server{
	// 	Config: p2p.Config{
	// 		MaxPeers:   10,
	// 		PrivateKey: nodekey,
	// 		Name:       "datx",
	// 		ListenAddr: ":30300",
	// 		Protocols:  []p2p.Protocol{MyProtocol()},
	// 	},
	// }

	// if err := srv.Start(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	//http server
	httpServer := http_plugin.NewHttpPlugin()
	httpServer.Init()
	httpServer.Open()
	defer httpServer.Close()

	//test chain_plugin
	sd := chain_plugin.GetInstance()
	sd.Init()
	log.Printf("chain %v", sd.Config)

	select {}
}

func msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
	for {
		recmsg, err := ws.ReadMsg()
		if err != nil {
			return err
		}

		var myMessage msg.BaseMsg
		err = recmsg.Decode(&myMessage)
		if err != nil {
			// handle decode error
			continue
		}

		switch myMessage.GetMsgType() {
		case msg.StatusMsg:
			err := p2p.SendItems(ws, msg.StatusMsg, "My status OK")
			if err != nil {
				return err
			}
		case msg.BigFileMsg:
			Body := myMessage.GetMsgBody()
			werr := ioutil.WriteFile("test.png", Body, 0666)
			if werr != nil {
				fmt.Print("create file failed!")
			}

			fmt.Println("recv:", Body)
			err := p2p.SendItems(ws, msg.ResponseCode, "File Recv Success")
			if err != nil {
				return err
			}
		default:
			fmt.Println("recv:", myMessage)
		}
	}

	return nil
}
