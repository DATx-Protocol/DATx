## Compile smart contract：
	cd ~/datx/Code/contracts/（path）
	datxiocpp -o DatxDToken.wast DatxDToken.cpp
	datxiocpp -g DatxDToken.abi DatxDToken.cpp

## create account
	cldatx create account datxio dtoken (DATX**********)
	
## deploy DatxToken
	cldatx set contract dtoken DatxDToken -p dtoken

## create token
	cldatx push action dtoken create '[ "datxio.dbtc", "21000000.0000 DBTC", 0, 0, 0]' -p dtoken
	cldatx push action dtoken create '[ "datxio.deth", "102000000.0000 DETH", 0, 0, 0]' -p dtoken
    cldatx push action dtoken create '[ "datxio.deos", "10000000000.0000 DEOS", 0, 0, 0]' -p dtoken
## transfer token to account
	cldatx  push action dtoken issue '["datxio.dbtc", "21000000.0000 DBTC", "memo "]' -p datxio.dbtc
	cldatx  push action dtoken issue '["datxio.deth", "102000000.0000 DETH", "memo "]' -p datxio.deth
	cldatx  push action dtoken issue '["datxio.deos", "10000000000.0000 DEOS", "memo "]' -p datxio.deos

## transfer to someone
    cldatx push action dtoken transfer '{"from":"datxio.dbtc","to":"test","quantity":"1.0000 DBTC","memo":"test"}' -p datxio.dbtc

	cldatx push action test extract '{"from":"test","to":"datxio.dbtc","quantity":"6.0000 DBTC","memo":"test"}' -p test

	cldatx push action dtoken transfer '{"from":"datxio.deth","to":"test","quantity":"30.0000 DETH","memo":"test"}' -p datxio.deth

	cldatx push action test extract '{"from":"test","to":"datxio.deth","quantity":"6.0000 DETH","memo":"test"}' -p test

	cldatx push action dtoken transfer '{"from":"datxio.deos","to":"datxio.deos","quantity":"30.0000 DEOS","memo":"test"}' -p datxio.deos

	cldatx push action test extract '{"from":"test","to":"datxio.deos","quantity":"6.0000 DEOS","memo":"test"}' -p test