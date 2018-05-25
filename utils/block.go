package utils

type BlockHeader struct {
	//producer account
	Producer string

	//Previous block id
	Previous string

	//producer signature
	ProSignature string
}

type Block struct {
	//block header
	BlockHeader

	//transcation pool
	Transcations []*Transcation
}
