package gatway

import (
	"testing"
)

func TestEOS(t *testing.T) {
	eos := NewEOSBrowser("https://api.eosmonitor.io/v1/", nil)

	res, err := eos.GetBlocks(12159298)
	if err != nil {
		t.Errorf("blocks error : %s\n", err)
	}
	t.Logf("\nblocks :%v\n", res)

	res, errs := eos.GetBlocks(12758916)
	if errs != nil {
		t.Errorf("blocks error : %s\n", errs)
	}
	t.Logf("\nblocks :%v\n", res)

	trx, err := eos.GetTransaction("4510d3738f472bd3cb518eedbdbdf4c0ed0a8f3948768fd482b997f4993ee441")
	if err != nil {
		t.Errorf("trx error : %s\n", err)
	}
	t.Logf("\ntrx :%v\n", trx)

	account, err := eos.GetAccounts("valuenetwork")
	if err != nil {
		t.Errorf("account error : %s\n", err)
	}
	t.Logf("\naccount :%v\n", account)
}
