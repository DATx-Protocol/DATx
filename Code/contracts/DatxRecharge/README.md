## Compile smart contract：
	cd ~/datx/Code/contracts/（path）
	datxiocpp -o DatxRecharge.wast DatxRecharge.cpp
	datxiocpp -g DatxRecharge.abi DatxRecharge.cpp

## Deploy smart contracts:
###	(1)creat wallet
		cldatx wallet create -n (name) --to-console
		
	(2)create key
		cldatx create key --to-console
		
	(3)import your key to wallet
		<1>cldatx wallet import -n (name) --private-key (PrivateKey)
		<2>cldatx wallet import -n (name) --private-key 5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3 (system account datxio)
		
	(4)create an account
		cldatx create account datxio (name) (PublicKey)
		for example: cldatx create account datxio datxio.charg DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E
		
	(5)deploy smart contracts
		cldatx set contract datxio.charg ../../../contracts/DatxRecharge  -p testdatxio.charg

	(6)add permission
		cldatx set account permission dtoken active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxio.charg","permission":"datxio.code"},"weight":1}]}' owner -p dtoken

		cldatx set account permission dtoken active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxio.charg","permission":"datxio.code"},"weight":1}]}' owner -p dtoken

		cldatx set account permission dtoken active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxio.charg","permission":"datxio.code"},"weight":1}]}' owner -p dtoken


## push action
	cldatx push action test recorduser '{"datxaddress":"datxuser","address":"38ZnTpSdCKq3BqexJFCvDKCq6AZvdUwKtQ","bpname":"bp1"}' -p test

	cldatx push action test charge '{"bpname":"bp1","hash":"70b4643bf0648e47784bb115255ee96de9bade0b1479a7abae68b1e627e9a611","from":"38ZnTpSdCKq3BqexJFCvDKCq6AZvdUwKtQ","to":"01XXXXXXXXXXXXX","blocknum":120,"quantity":"10","category":"BTC","memo":"this is first test"}' -p test

## get table
	cldatx get table code scope  table_name
	for example:cldatx get table datxio.charg datxio.charg record
				cldatx get table datxio.charg datxio.charg transaction

## make sure every step you wallet is unlock
	cldatx wallet unlock --name (name) --password (password)

## start PrivateChain：
###	(1)kdatxd：
		kdatxd --http-server-address=127.0.0.1:8900
		
###	(2)noddatx:
		noddatx -e -p datxio --accessory  datxio::core_api_accessory --accessory datxio::history_api_accessory --replay-blockchain --verbose-http-errors
		
## when you restart your noddatx you should remove data first ：
		cd  ~/.local/share/datxio/noddatx/
		sudo rm -rf data
## then open your wallet and unlock it	
		cldatx wallet open -n (name)
		cldatx wallet unlock --name (name) --password (password)

	
