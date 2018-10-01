	
## create account
	cldatx create account datxio datxio.dbtc (DATX**********)
	cldatx create account datxio datxio.deth (DATX**********)
	cldatx create account datxio datxio.deos (DATX**********)
## deploy DatxToken
	cldatx set contract datxio.dbtc DatxToken -p datxio.dbtc
	cldatx set contract datxio.deth DatxToken -p datxio.deth
	cldatx set contract datxio.dbtc DatxToken -p datxio.dbtc

## create token
	cldatx push action datxio.dbtc create '[ "datxio", "21000000.0000 DBTC", 0, 0, 0]' -p datxio.dbtc
	cldatx push action datxio.deth create '[ "datxio", "102000000.0000 DETH", 0, 0, 0]' -p datxio.deth
	cldatx push action datxio.dbtc create '[ "datxio", "10000000000.0000 DEOS", 0, 0, 0]' -p datxio.dbtc

## transfer token to account
	cldatx  push action datxio.dbtc issue '["datxio.dbtc", "21000000.0000 DBTC", "memo "]' -p datxio
	cldatx  push action datxio.deth issue '["datxio.dbtc", "102000000.0000 DETH", "memo "]' -p datxio
	cldatx  push action datxio.dbtc issue '["datxio.dbtc", "10000000000.0000 DEOS", "memo "]' -p datxio

## transfer to someone
    cldatx push action datxio.dbtc transfer '{"from":"datxio.dbtc","to":"datxio","quantity":"1.0000 DBTC","memo":"test"}' -p datxio.dbtc
