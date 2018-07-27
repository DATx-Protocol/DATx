package types

//const value
const (
	BlockIntervalMs uint64 = 500 //Millisecond
	BlockIntervalUs        = BlockIntervalMs * 1000
	BlockTimeEpoch  uint64 = 946684800000 // epoch is year 2000.

	MaxTrackedDposConfirmations int = 1024

	ProducerRepetitions int = 6
	MaxProducers        int = 125

	SystemAccountName string = "datx_chain"

	RootKey string = "Datx chain is a great project!"
)

//block status
const (
	Irreversible uint16 = 1 << iota //this block has already been applied before by this node and is considered irreversible
	Validate                        //this is a complete block signed by a valid producer and has been previously applied by this node and therefore validated but it is not yet irreversible
	Complete                        //this is a vomplete block signed by a valid producer but is not yet irreversible nor has it yet been applied by this node
	Incomplete                      //this is an incomplete block(either being produced by a producer or speculatively produced by a node)
)
