package types

import (
	"math/big"
	time "time"
)

//BlockTime struct
type BlockTime struct {
	//block time
	Time *big.Int

	//block interval that is ms
	Interval uint64

	//Epoch is year 2000
	Epoch uint64

	Slot uint64
}

//ToTime convert string to time
func ToTime(v string) time.Time {
	t, _ := time.Parse(time.RFC3339, v)
	return t
}

//NewBlockTime new
func NewBlockTime(t time.Time) *BlockTime {
	var res BlockTime
	res.Interval = BlockIntervalMs
	res.Epoch = BlockTimeEpoch

	res.Time = big.NewInt(t.Unix())

	unix := uint64(t.Unix())
	res.Slot = (unix*1000 - res.Epoch) / res.Interval

	return &res
}

//SetTime set time
func (bt *BlockTime) SetTime(t time.Time) {
	bt.Time = big.NewInt(t.Unix())

	unix := uint64(t.Unix())
	bt.Slot = (unix*1000 - bt.Epoch) / bt.Interval
}

//Less compare
func (bt *BlockTime) Less(in *BlockTime) bool {
	return bt.Time.Cmp(in.Time) == -1 //bt.Time.Sub(in.Time) < time.Duration(0)
}

//Plus plus slot
func (bt *BlockTime) Plus() {
	bt.Slot++
}

//AddPlusTime plus time
func (bt *BlockTime) AddPlusTime(d int64) uint64 {
	if bt == nil {
		return 0
	}

	newTime := big.NewInt(0)
	newTime.Add(bt.Time, big.NewInt(d))

	return newTime.Uint64()
}

//String print
func (bt *BlockTime) String() string {
	var res string
	bd, err := bt.Time.MarshalText()
	if err != nil {
		return res
	}

	res = string(bd)

	return res
}

//MaxTime max time
func MaxTime() time.Time {
	return time.Unix(1<<63-62135596801, 999999999)
}
