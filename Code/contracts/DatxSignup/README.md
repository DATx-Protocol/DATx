## Compile smart contract：
	cd (your path)/contracts/DatxSignup/
	datxiocpp -o DatxSignup.wast DatxSignup.cpp
	datxiocpp -g DatxSignup.abi  DatxSignup.cpp

## Deploy smart contracts:
###	
    (1)create system accounts
        cldatx create account datxio datxio.bpay  your_public_key
        cldatx create account datxio datxio.msig  your_public_key
        cldatx create account datxio datxio.names your_public_key
        cldatx create account datxio datxio.ram   your_public_key
        cldatx create account datxio datxio.save  your_public_key
        cldatx create account datxio datxio.stake your_public_key
        cldatx create account datxio datxio.token your_public_key
        cldatx create account datxio datxio.vpay  your_public_key
        cldatx create account datxio datxio.veri  your_public_key

        cldatx create account datxio datxio.dbtc  your_public_key
        cldatx create account datxio datxio.deth  your_public_key
        cldatx create account datxio datxio.deos  your_public_key

        cldatx create account datxio datxio.charg your_public_key
        cldatx create account datxio datxio.sign  your_public_key

    (2)create & issue DATX token
        cldatx set contract datxio.token (your path)/build/contracts/DatxToken -p datxio.token
        cldatx push action datxio.token create '["datxio", "1000000000.0000 DATX", 0, 0, 0]' -p datxio.token
        cldatx push action datxio.token issue  '["datxio" ,"1000000000.0000 DATX", "memo"]'  -p datxio

    (3)deploy system contract
        cldatx set contract datxio       (your path)/build/contracts/DatxSystem
        
    (4)test system newaccount
        cldatx system newaccount datxio heyaoahsdfge your_public_key --stake-net '0.001 DATX' --stake-cpu '0.02 DATX' --buy-ram-kbytes 3

    (5)add permission
        cldatx set account permission datxio.sign active '{"threshold": 1,"keys": [{"key": "your_public_key","weight": 1}],"accounts": [{"permission":{"actor":"datxio.sign","permission":"datxio.code"},"weight":1}]}' owner -p datxio.sign
	(6)deploy smart contracts
		cldatx set contract datxio.sign (your path)/contracts/DatxSignup  -p datxio.sign
	(7)transfer
		cldatx push action datxio.sign transfer '{"from":"account","to":"datxio.sign","quantity":"1.0000 DATX","memo":"your_name your_public_key"}' -p account
		your_name should be less than 13 characters and only contains the following symbol .12345abcdefghijklmnopqrstuvwxyz

	