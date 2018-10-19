## Compile smart contract：
	cd ~/datx/Code/contracts/（path）
	datxoscpp -o DatxRecharge.wast DatxRecharge.cpp
	datxoscpp -g DatxRecharge.abi DatxRecharge.cpp

## Deploy smart contracts:
###	(1)creat wallet
		cldatx wallet create -n (name) --to-console
		
	(2)create key
		cldatx create key --to-console
		
	(3)import your key to wallet
		<1>cldatx wallet import -n (name) --private-key (PrivateKey)
		<2>cldatx wallet import -n (name) --private-key 5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3 (system account datxos)
		
	(4)create an account
		cldatx create account datxos (name) (PublicKey)
		for example: cldatx create account datxos datxos.charg DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E
		
	(5)deploy smart contracts
		cldatx set contract datxos.charg ../../../contracts/DatxRecharge  -p testdatxos.charg

	(6)add permission
		cldatx set account permission datxos.dtoke active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}' owner -p datxos.dtoke

		cldatx set account permission datxos.dbtc active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}' owner -p datxos.dbtc

		cldatx set account permission datxos.deos active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}' owner -p datxos.deos

		cldatx set account permission datxos.deth active '{"threshold": 1,"keys": [{"key": "DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E","weight": 1}],"accounts": [{"permission":{"actor":"datxos.charg","permission":"datxos.code"},"weight":1}]}' owner -p datxos.deth

## push action
	cldatx push action datxos.charg recorduser '{"datxaddress":"datxuser","address":"38ZnTpSdCKq3BqexJFCvDKCq6AZvdUwKtQ","bpname":"bp1"}' -p datxos.charg

	cldatx push action datxos.charg charge '{"bpname":"bp1","hash":"70b4643bf0648e47784bb115255ee96de9bade0b1479a7abae68b1e627e9a611","from":"38ZnTpSdCKq3BqexJFCvDKCq6AZvdUwKtQ","to":"01XXXXXXXXXXXXX","blocknum":120,"quantity":"10","category":"BTC","memo":"this is first test"}' -p datxos.charg

	cldatx push action datxos.charg charge '{"bpname":"bp1","hash":"30b4643bf0648e47784bb115255ee96de9bade0b1479a7abae68b1e627e9a611","from":"11a","to":"01aa","blocknum":120,"quantity":"10","category":"EOS","memo":"datx1"}' -p datxos.charg

## get table
	cldatx get table code scope  table_name
	for example:cldatx get table datxos.charg datxos.charg record
				cldatx get table datxos.charg datxos.charg transaction

## make sure every step you wallet is unlock
	cldatx wallet unlock --name (name) --password (password)

## start PrivateChain：
###	(1)kdatxd：
		kdatxd --http-server-address=127.0.0.1:8900
		
###	(2)noddatx:
		noddatx -e -p datxos --accessory  datxos::core_api_accessory --accessory datxos::history_api_accessory --replay-blockchain --verbose-http-errors
		
## when you restart your noddatx you should remove data first ：
		cd  ~/.local/share/datxos/noddatx/
		sudo rm -rf data
## then open your wallet and unlock it	
		cldatx wallet open -n (name)
		cldatx wallet unlock --name (name) --password (password)

	
