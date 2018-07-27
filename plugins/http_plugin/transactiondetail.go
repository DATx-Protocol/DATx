package httpplugin

//TransactionDetail struct apply to block explore
type TransactionDetail struct {
	TrxHash     string  //transaction Id  common.hash to string
	BlockHeight int64   // block height
	TimeStamp   string  //timestamp
	Pending     string  //mode of transaction
	Amount      float64 //amount of transaction
	AccountFrom string  //sender account
	TrxType     string  //type of transaction
}
