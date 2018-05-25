package utils

type TranscationHeader struct {
	//transcation expiration ,the seconds from 1970.
	Expiration int64

	//reference the latest block number
	RefBlockNum uint16

	//reference block prefix
	RefBlockPerfix uint32
}

type Transcation struct {
	//transcation header, inherited from TranscationHeader
	TranscationHeader

	//action list
	Actions []*Action
}

//transcation constructor
func NewTrx(time int64) *Transcation {
	return &Transcation{
		TranscationHeader: TranscationHeader{
			Expiration: time,
		},
	}
}
