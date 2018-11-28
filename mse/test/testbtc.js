
 bitcoinapi = require('../lib/bitcoinapi');
 bitcoin = require('bitcoinjs-lib');

// 托管账户的3组公私钥
// describe('getKeysFromWIF',function(){
//     keyPair1 = bitcoinapi.getKeysFromWIF('cV6dJVkPZfmmDg5XHNpExd8fgMTaoYgv7NusH26ZEFj2vcA9xvhz',{IsTestnet:1})
//     keyPair2 = bitcoinapi.getKeysFromWIF('cQX7uYfhPJJGPdoqgHQXLC4WQWcDXGH7DgLs4m89xWykYkg1V4iA',{IsTestnet:1})
//     keyPair3 = bitcoinapi.getKeysFromWIF('cV7ReEJu9bSNpTc1DprnqT3S53k1LzYVaJTxJWUsmQMNZsaFSnc8',{IsTestnet:1})

//     expect(keyPair1.pubkey).to.equal('03a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb')
//     expect(keyPair2.pubkey).to.equal('038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac9')
//     expect(keyPair3.pubkey).to.equal('02040e0d9141b06ad92f38d7a3d76cfcb6ada4c9e4c5b18d18f5539564a3826408')

// })


//用户账号
//wif:cQg6EvRRka4BtLpzFuuVKUptToqiJj3SzhjDZEcak5eybApzAEFc
//pubkey:0298225bd8d722cf0af4ee6a943e566951f4fac3f578f231a5cbb8a0b93f4a0e49
//address:myDAgFpwT3sTkppYKeS6LwKMCipM218EKE


//托管账户地址
//address:2MwtNFrT9P1wDa3Hid6kVu9mh84cd59UHKN
//script:522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102040e0d9141b06ad92f38d7a3d76cfcb6ada4c9e4c5b18d18f5539564a382640853ae


//往托管账户转账
(async function(){
    network = bitcoin.networks.testnet;
    to = '2MwtNFrT9P1wDa3Hid6kVu9mh84cd59UHKN';
    from = 'myDAgFpwT3sTkppYKeS6LwKMCipM218EKE';
    fee = 10e3;
    value = 1000e3;
    memo = (new Date()).toLocaleString();

    chunk = await bitcoinapi.getUTXOS(from, {IsTestnet:1});
    utxos = JSON.parse(chunk.data).unspent_outputs;

    trx = await bitcoinapi.buildTrx(from,to,value,fee,{IsTestnet:1},'Memo');

    //sign
    userwif = 'cQg6EvRRka4BtLpzFuuVKUptToqiJj3SzhjDZEcak5eybApzAEFc';
    insDecode = bitcoinapi.decodeInput(trx);
    keyPair = bitcoin.ECPair.fromWIF(userwif, network);
    inputCount = trx.ins.length;
    txb = bitcoin.TransactionBuilder.fromTransaction(trx, network);
    for (var i = 0; i < inputCount; i++) {
        try{
        txb.sign(i, keyPair);
        }
        catch(e){
        // do nothing
        }
    }

    trx = txb.buildIncomplete();

    result = await bitcoinapi.broadcastTrx(trx,{IsTestnet:1});
    console.log('trxid:' + result);
})();






//多重签名从托管账户提现
// (async function(){
//     script = '522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102040e0d9141b06ad92f38d7a3d76cfcb6ada4c9e4c5b18d18f5539564a382640853ae'
//     //buildTrx
//     trx = await bitcoinapi.buildTrx('2MwtNFrT9P1wDa3Hid6kVu9mh84cd59UHKN','myDAgFpwT3sTkppYKeS6LwKMCipM218EKE',9e3,1e3,{IsTestnet:1},'dea9be5288291d3f3ab3d35354f5540fbcfd76e08c10e49b13b5d65237a347eb')

//     //sign locally
//     trx = await bitcoinapi.signTrx(trx,'cV6dJVkPZfmmDg5XHNpExd8fgMTaoYgv7NusH26ZEFj2vcA9xvhz',script,{IsTestnet:1})

//     //sign remotly
//     trxSerialize = trx.toHex()
//     trx = bitcoinapi.getTrxFromHex(trxSerialize)
//     trx = await bitcoinapi.signTrx(trx,'cQX7uYfhPJJGPdoqgHQXLC4WQWcDXGH7DgLs4m89xWykYkg1V4iA',script,{IsTestnet:1})
//     trxSerialize = trx.toHex()

//     //broadcast trx
//     trx = bitcoinapi.getTrxFromHex(trxSerialize)
//     bitcoinapi.broadcastTrx(trx) 
// })()




