package msg

const (
	StatusMsg uint64 = 1 << iota
	TxMsg
	NewBlockMsg
	BlockHeaderMsg
	GetBlockHeaderMsg
	BlockBodiesMsg
	GetBlockBodiesMsg
	ReceiptsMsg
	GetReceiptsMsg
	BigFileMsg
	GetBigFileMs
	ResponseCode
)

type BaseMsg interface {
	GetMsgType() uint64
	GetMsgBody() []byte
}

type TestMsg struct {
	msgType uint64
	msgBody []byte
}
