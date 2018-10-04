package gatway

import (
	"testing"
)

func TestETH(t *testing.T) {
	eth := NewETHBrowser("", nil)

	trxs, err := eth.GetTrxs("0xddbd2b932c763ba5b1b7ae3b362eac3e8d40121a", 0, 99999999)
	if err != nil {
		t.Errorf("trx error : %s\n", err)
	}

	t.Logf("trx list: %v\n", trxs)

	num, err := eth.GetLatestBlockNum()
	if err != nil {
		t.Errorf("GetLatestBlockNum error : %s\n", err)
	}

	t.Logf("Get latest blocknum: %v\n", num)
}
