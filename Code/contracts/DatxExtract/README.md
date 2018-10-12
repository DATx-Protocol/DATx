## Compile smart contract：
	cd ~/datx/Code/contracts/（path）
	datxiocpp -o datxextract.wast DatxExtract.cpp
	datxiocpp -g datxextract.abi DatxExtract.hpp

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
		for example: cldatx create account datxio datx.extract DATX76kRKSJJVUb2bYLQUwjfSvoQsqU7mwzGCwbD9TtGfukhPwN43E
		
	(5)deploy smart contracts
		cldatx set contract datx.extract../../../contracts/DatxExtract  -p datx.extract


## push action
	cldatx push action datx.extract recordtrx '{"trxid":70b4643bf0648e47784bb115255ee96de9bade0b1479a7abae68b1e627e9a611,"producer":"bp1","category":"BTC"}' -p bp1

	cldatx push action datx.extract setverifiers '{"verifiers":["v1","v2","v3","v4","v5"]}' -p datx.extract

    cldatx push action datx.extract setdoing '{"trxid":70b4643bf0648e47784bb115255ee96de9bade0b1479a7abae68b1e627e9a611,"producer":"bp1","verifier":"verifier1"}' -p verifier1

    cldatx push action datx.extract setsuccess '{"trxid":70b4643bf0648e47784bb115255ee96de9bade0b1479a7abae68b1e627e9a611,"producer":"bp1"}' -p bp1

	cldatx push action datx.extract updateexpire '' -p datx.extract



## get table
	cldatx get table code scope  table_name
	for example:cldatx get table datx.extract datx.extract record
				cldatx get table datx.extract datx.extract success
                cldatx get table datx.extract datx.extract expiration

## make sure every step you wallet is unlock
	cldatx wallet unlock --name (name) --password (password)

## start PrivateChain：
###	(1)kdatxd：
		kdatxd --http-server-address=127.0.0.1:8900
		
###	(2)noddatx:
		noddatx -e -p datxio --plugin  datxio::core_api_plugin --plugin datxio::history_api_plugin --replay-blockchain --verbose-http-errors
		
## when you restart your noddatx you should remove data first ：
		cd  ~/.local/share/datxio/noddatx/
		sudo rm -rf data
## then open your wallet and unlock it	
		cldatx wallet open -n (name)
		cldatx wallet unlock --name (name) --password (password)

	
