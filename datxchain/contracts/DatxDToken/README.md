## Compile smart contract：
	cd ~/datx/Code/contracts/（path）
	datxoscpp -o DatxDToken.wast DatxDToken.cpp
	datxoscpp -g DatxDToken.abi DatxDToken.cpp

## create account
	cldatx create account datxos datxos.dtoke (DATX**********)
	
## deploy DatxToken
	cldatx set contract datxos.dtoke DatxDToken -p datxos.dtoke

## create token
	cldatx push action datxos.dtoke create '[ "datxos.dbtc", "21000000.0000 DBTC", 0, 0, 0]' -p datxos.dtoke
	cldatx push action datxos.dtoke create '[ "datxos.deth", "102000000.0000 DETH", 0, 0, 0]' -p datxos.dtoke
    cldatx push action datxos.dtoke create '[ "datxos.deos", "10000000000.0000 DEOS", 0, 0, 0]' -p datxos.dtoke
## transfer token to account
	cldatx  push action datxos.dtoke issue '["datxos.dbtc", "21000000.0000 DBTC", "memo "]' -p datxos.dbtc
	cldatx  push action datxos.dtoke issue '["datxos.deth", "102000000.0000 DETH", "memo "]' -p datxos.deth
	cldatx  push action datxos.dtoke issue '["datxos.deos", "10000000000.0000 DEOS", "memo "]' -p datxos.deos

## transfer to someone
    cldatx push action datxos.dtoke transfer '{"from":"datxos.dbtc","to":"test","quantity":"1.0000 DBTC","memo":"test"}' -p datxos.dbtc

	cldatx push action test extract '{"from":"test","to":"datxos.dbtc","quantity":"6.0000 DBTC","memo":"test"}' -p test

	cldatx push action datxos.dtoke transfer '{"from":"datxos.deth","to":"test","quantity":"30.0000 DETH","memo":"test"}' -p datxos.deth

	cldatx push action test extract '{"from":"test","to":"datxos.deth","quantity":"6.0000 DETH","memo":"test"}' -p test

	cldatx push action datxos.dtoke transfer '{"from":"datxos.deos","to":"datxos.deos","quantity":"30.0000 DEOS","memo":"test"}' -p datxos.deos

	cldatx push action test extract '{"from":"test","to":"datxos.deos","quantity":"6.0000 DEOS","memo":"test"}' -p test