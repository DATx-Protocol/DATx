	
## create account
	cldatx create account datxos datxos.dbtc (DATX**********)
	cldatx create account datxos datxos.deth (DATX**********)
	cldatx create account datxos datxos.deos (DATX**********)
## deploy DatxToken
	cldatx set contract datxos.dbtc DatxToken -p datxos.dbtc
	cldatx set contract datxos.deth DatxToken -p datxos.deth
	cldatx set contract datxos.deos DatxToken -p datxos.deos

## create token
	cldatx push action datxos.dbtc create '[ "datxos", "21000000.0000 DBTC", 0, 0, 0]' -p datxos.dbtc
	cldatx push action datxos.deth create '[ "datxos", "102000000.0000 DETH", 0, 0, 0]' -p datxos.deth
    cldatx push action datxos.deos create '[ "datxos", "10000000000.0000 DEOS", 0, 0, 0]' -p datxos.deos
## transfer token to account
	cldatx  push action datxos.dbtc issue '["datxos.dbtc", "21000000.0000 DBTC", "memo "]' -p datxos
	cldatx  push action datxos.deth issue '["datxos.deth", "102000000.0000 DETH", "memo "]' -p datxos
	cldatx  push action datxos.deos issue '["datxos.deos", "10000000000.0000 DEOS", "memo "]' -p datxos

## transfer to someone
    cldatx push action datxos.dbtc transfer '{"from":"datxos.dbtc","to":"datxos","quantity":"1.0000 DBTC","memo":"test"}' -p datxos.dbtc



	cldatx push action datxos.deth transfer '{"from":"datxos.deth","to":"test","quantity":"30.0000 DBTC","memo":"test"}' -p datxos.deth

	cldatx push action test withdraw '{"from":"test","to":"datxos.deth","quantity":"6.0000 DETH","memo":"test"}' -p test

	cldatx push action datxos.deos transfer '{"from":"datxos.deos","to":"test","quantity":"30.0000 DEOS","memo":"test"}' -p datxos.deos


	cldatx push action test withdraw '{"from":"test","to":"datxos.deos","quantity":"6.0000 DEOS","memo":"test"}' -p test