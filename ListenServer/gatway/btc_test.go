package gatway

import "testing"

func TestBTC(t *testing.T) {
	btc := NewBTCBrowser("", nil)

	trxlist, err := btc.GetTrxs("1AJbsFZ64EpEfS5UAjAfcUG8pH8Jn3rn1F")
	if err != nil {
		t.Errorf("BTC gettrx err : %v\n", err)
	}

	t.Logf("BTC trxlist:%v\n", trxlist)
}
