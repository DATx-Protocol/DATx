package main

import (
	"io/ioutil"
	"io"
	"fmt"
	"os"

	"datx_chain/node/msg"
	"datx_chain/p2p"

	"github.com/ethereum/go-ethereum/crypto"
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
	nodekey, _ := crypto.GenerateKey()
	srv := p2p.Server{
		Config: p2p.Config{
			MaxPeers:   10,
			PrivateKey: nodekey,
			Name:       "datx",
			ListenAddr: ":30300",
			Protocols:  []p2p.Protocol{MyProtocol()},
		},
	}

	if err := srv.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
			if ioutil.WriteFile("test.png", Body, os.O_APPEND), err != nil{
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
