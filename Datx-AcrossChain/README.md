DATX——ACROSSCHAIN
===========================================================

# 提供接口
----------------------------------------------------------
## 比特币生成密钥对 https://localhost:8080/btc/genKeyPairs?IsTestnet=1
#### 参数：
    IsTestnet 是否为测试网络
### 返回：
    {"wif":"KyDjtpbGEGPrT8FRXq71agsJuFcTdew3hHCkZEdNDXPCp1j5rvbZ","pubkey":"038ef11818e18eeca101cf417898a3b425c4c71f8406a0290c017a22f1159a9ccd","prikey":"3bb0a110ec975fce5d71159d5a555700bf2fa15c054123b03f030e24a7d1e9e5"}
    wif：钱包形式的密钥对
    pubkey：公钥hex形式
    prikey：私钥hex形式


## 比特币根据公钥生成地址 https://localhost:8080/btc/genP2PKHAddr?isTestnet=1&pubkey=038ef11818e18eeca101cf417898a3b425c4c71f8406a0290c017a22f1159a9ccd
### 参数：
    IsTestnet 是否为测试网络
    pubkey    公钥的hex形式
### 返回：
    mtwvNmzYWqcx1JUPonX6sMxjSyeu9nBEck
    直接返回地址的hex形式


## 比特币生成多重签名地址 https://localhost:8080/btc/genMulSigAddr?isTestnet=1&pubkeys=03a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb，038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac9,02c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e&num=2
### 参数：
    IsTestnet 是否为测试网络
    pubkeys   多重签名公钥组
    num       阈值
### 返回：
    {"address":"2N8J8pLpVjkX1KHKHDWdadEYzD9ttLcueBX","script":"522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb2102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e52ae"}
    address     多重签名地址hex形式
    script      锁定脚本的hex形式


## 比特币发起提现请求 https://localhost:8080/btc/withdraw?isTestnet=1&to=n4fc3bKTVrVRveBrHZf5Zv4wGGBTf5sdHg&value=10000000&fee=100000&trxid=asdada
### 参数：
    IsTestnet 是否为测试网络
    to        目的地址
    value     提现金额，以聪为单位
    fee       手续费，以聪为单位，用户承担
    trxid     datx链上交易ID
### 返回：
    transaction id


## 比特币对交易签名 https://localhost:8080/btc/signTrx?IsTestnet=1&trxSerialize=aaaaaaaaaaa&to=bbbbbbbb&trxid=asdafff
### 参数：
    IsTestnet       是否为测试网络
    trxSerialize    交易的hex形式
    trxid           datx链上的交易ID
### 返回：
    020000000144bd36311739e090ad98616e93882c321676c76eb48d9a9583a6cdc19dfa7fce01000000b60047304402201d4212a73b92bdef0b61199eae990b0439d303edfe01e1ec9fa37d15a04c3037022065ddaf3d42369d07f07c7083d52d9f87b0905459e910ff0745e91dff9cdf255e0100004c69522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53aeffffffff0280969800000000001976a914fdedbbbf546a9344fffa80f433b2ebd9b60b245988ac809698000000000017a9140535e10b0244b05ae907edeabb43f1cd80e0f1d68700000000


## 以太坊提现请求 https://localhost:8080/eth/withdraw?to=ccc&value=100000&data=data&trxid=asdasd
### 参数：
    to              提现目标地址
    value           提现金额（wei）
    data            附加信息
    trxid           datx链上的交易ID

### 返回：
    success


## EOS创建多重签名账户 https://localhost:8080/eos/createAccount?creator=alice&accountName=jack&auth={threshold: 2,ccounts: [ { permission: {actor:'bob',permission:'active'}, weight: 1 },{ permission: {actor:'alice',permission:'active'}, weight: 1 } ],waits: []}
### 参数：
    creator     创建人，本节点EOS账户
    accountName 创建账户名称
    auth        账户权限设置，json格式

### 返回：
    包含交易id的json格式结果

## EOS向账户直接发币（eosio.token直接发，非转账）https://localhost:8080/eos/issue?to=aaaa&value=1000.00 EOS&memo=aaaa
### 参数：
    to      收款账户
    value   金额与类型
    memo    备注
### 返回参数：
    包含交易id的json格式结果

## EOS发起提现请求 https://localhost:8080/eos/withdraw?to=aaaaa&value=1000.0000 EOS&trxid=aaaa
### 参数：
    to      收款账户
    value   金额与类型
    trxid   datx链上的交易ID
### 返回参数：
    包含交易id的json格式结果

## EOS确认提案   https://localhost:8080/eos/confirm?propser=aaa&proposeName=bbbb
### 参数：
    proposer    提案发起账户
    proposeName 提案名称
    trxid       datx链上的交易ID
### 返回参数：
    包含交易id的json格式结果











