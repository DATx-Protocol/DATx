package server

import (
	"testing"
)

func TestGetExpiredTrxs(t *testing.T) {
	ex := NewExtract("datxos")
	ex.Startup()
	defer ex.Close()

	expire, err := ex.getExpiredTrxs()
	if err != nil {
		t.Errorf("getExpiredTrxs err: %v\n", err)
	}

	t.Logf("getExpiredTrxs data:%v\n", expire)

	for _, v := range expire {
		erro := ex.pushExtractAction(v)
		if erro != nil {
			t.Errorf("pushExtractAction err: %v\n", erro)
		}
	}
}
