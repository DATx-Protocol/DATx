CLMSE
===========================================================

# btc
----------------------------------------------------------
## 比特币生成密钥对 clmse btc genKeyPair -t
#### 参数：
    -t 是否为测试网络
### 返回：
    
    { wif: 'cSWdyDcJscQGztzK9RLa8zxPpvAYBTVoG7cgrvW1DVM2VxxunJC2',
      pubkey:'0386fd567e5e6ae2758bbc795ce9ec47c99a8a9358844e7268191c3ebc62338097',
      prikey:'9318bb127939b343fb4f46517b858efd7731a05becbd4e18a68cf8e9d3410d77' 
    }


## 比特币从WIF中提取密钥对 clmse btc getKeysFromWIF cSWdyDcJscQGztzK9RLa8zxPpvAYBTVoG7cgrvW1DVM2VxxunJC2 -t
#### 参数：
    <wif>
    -t 是否为测试网络
### 返回：
    
    { pubkey:'0386fd567e5e6ae2758bbc795ce9ec47c99a8a9358844e7268191c3ebc62338097',
      prikey:'9318bb127939b343fb4f46517b858efd7731a05becbd4e18a68cf8e9d3410d77' 
    }


## 比特币生成P2PKH型地址 clmse btc genP2PKHAddr 0386fd567e5e6ae2758bbc795ce9ec47c99a8a9358844e7268191c3ebc62338097 -t
#### 参数：
    <public key>
    -t 是否为测试网络
### 返回：
    
    n3zRpYAQRbTHaPBKezvswacFTCSsRb28Hx



## 比特币生成多重签名地址 clmse btc genMultiSigAddr '["03a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb","038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac9","02040e0d9141b06ad92f38d7a3d76cfcb6ada4c9e4c5b18d18f5539564a3826408"]' 2 -t
#### 参数：
    <public keys>
    <num>
    -t 是否为测试网络
### 返回：
    
    { address: '2MwtNFrT9P1wDa3Hid6kVu9mh84cd59UHKN',
      script:'522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102040e0d9141b06ad92f38d7a3d76cfcb6ada4c9e4c5b18d18f5539564a382640853ae' 
    }



