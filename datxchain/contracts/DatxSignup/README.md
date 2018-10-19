## Compile smart contract：
	cd (your path)/contracts/DatxSignup/
	datxoscpp -o DatxSignup.wast DatxSignup.cpp
	datxoscpp -g DatxSignup.abi  DatxSignup.cpp

## Deploy smart contracts:
###	
    (1)create system accounts
        cldatx create account datxos datxos.bpay  your_public_key
        cldatx create account datxos datxos.msig  your_public_key
        cldatx create account datxos datxos.names your_public_key
        cldatx create account datxos datxos.ram   your_public_key
        cldatx create account datxos datxos.save  your_public_key
        cldatx create account datxos datxos.stake your_public_key
        cldatx create account datxos datxos.token your_public_key
        cldatx create account datxos datxos.vpay  your_public_key
        cldatx create account datxos datxos.veri  your_public_key

        cldatx create account datxos datxos.dbtc  your_public_key
        cldatx create account datxos datxos.deth  your_public_key
        cldatx create account datxos datxos.deos  your_public_key

        cldatx create account datxos datxos.charg your_public_key
        cldatx create account datxos datxos.sign  your_public_key

    (2)create & issue DATX token
        cldatx set contract datxos.token (your path)/build/contracts/DatxToken -p datxos.token
        cldatx push action datxos.token create '["datxos", "1000000000.0000 DATX", 0, 0, 0]' -p datxos.token
        cldatx push action datxos.token issue  '["datxos" ,"1000000000.0000 DATX", "memo"]'  -p datxos

    (3)deploy system contract
        cldatx set contract datxos       (your path)/build/contracts/DatxSystem
        
    (4)test system newaccount
        cldatx system newaccount datxos heyaoahsdfge your_public_key --stake-net '0.001 DATX' --stake-cpu '0.02 DATX' --buy-ram-kbytes 3

    (5)add permission
        cldatx set account permission datxos.sign active '{"threshold": 1,"keys": [{"key": "your_public_key","weight": 1}],"accounts": [{"permission":{"actor":"datxos.sign","permission":"datxos.code"},"weight":1}]}' owner -p datxos.sign
	(6)deploy smart contracts
		cldatx set contract datxos.sign (your path)/contracts/DatxSignup  -p datxos.sign
	(7)transfer
		cldatx push action datxos.sign transfer '{"from":"account","to":"datxos.sign","quantity":"1.0000 DATX","memo":"your_name your_public_key"}' -p account
		your_name should be less than 13 characters and only contains the following symbol .12345abcdefghijklmnopqrstuvwxyz

	