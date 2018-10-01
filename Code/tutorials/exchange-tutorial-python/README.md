The following steps must be taken for the example script to work.

0. Create wallet
0. Create account for datxio.token
0. Create account for scott
0. Create account for exchange
0. Set token contract on datxio.token
0. Create DATX token
0. Issue initial tokens to scott

**Note**:
Deleting the `transactions.txt` file will prevent replay from working.


### Create wallet
`cldatx wallet create`

### Create account steps
`cldatx create key`

`cldatx create key`

`cldatx wallet import  --private-key <private key from step 1>`

`cldatx wallet import  --private-key <private key from step 2>`

`cldatx create account datxio <account_name> <public key from step 1> <public key from step 2>`

### Set contract steps
`cldatx set contract datxio.token /contracts/datxio.token -p datxio.token@active`

### Create DATX token steps
`cldatx push action datxio.token create '{"issuer": "datxio.token", "maximum_supply": "100000.0000 DATX", "can_freeze": 1, "can_recall": 1, "can_whitelist": 1}' -p datxio.token@active`

### Issue token steps
`cldatx push action datxio.token issue '{"to": "scott", "quantity": "900.0000 DATX", "memo": "testing"}' -p datxio.token@active`
