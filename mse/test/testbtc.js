
 bitcoinapi = require('../lib/bitcoinapi');
 bitcoin = require('bitcoinjs-lib');

// 托管账户的3组公私钥
// describe('getKeysFromWIF',function(){
//     keyPair1 = bitcoinapi.getKeysFromWIF('cV6dJVkPZfmmDg5XHNpExd8fgMTaoYgv7NusH26ZEFj2vcA9xvhz',{IsTestnet:1})
//     keyPair2 = bitcoinapi.getKeysFromWIF('cQX7uYfhPJJGPdoqgHQXLC4WQWcDXGH7DgLs4m89xWykYkg1V4iA',{IsTestnet:1})
//     keyPair3 = bitcoinapi.getKeys5FromWIF('cMpRvr5XW4Uw4hxGb7TbcG1NqUqAw77CUqAywim6htmpFREXbfbu',{IsTestnet:1})

//     expect(keyPair1.pubkey).to.equal('03a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb')
//     expect(keyPair1.prikey).to.equal('e041a10f944dcf1978d9e53c32e120ff656b13ed9b22c5b2ae583903f77660d3')
//     expect(keyPair2.pubkey).to.equal('038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac9')
//     expect(keyPair2.prikey).to.equal('57ab6ad9fb00db996dce408bd36a15b65488b57da904e9a6cf9eed44a4c3d5f0')
//     expect(keyPair3.pubkey).to.equal('02c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e')
//     expect(keyPair3.prikey).to.equal('07100d72612bed3f3cf9283d26047e5cc611e8e074f2e9c0dd8ec02f78d2a85c')

// })


//用户账号
//wif:cQg6EvRRka4BtLpzFuuVKUptToqiJj3SzhjDZEcak5eybApzAEFc
//prikey:5c48e64645f87ae2f9fb0a48977345b973f333b121a0d561563938b3c8aacc51
//pubkey:0298225bd8d722cf0af4ee6a943e566951f4fac3f578f231a5cbb8a0b93f4a0e49
//address:myDAgFpwT3sTkppYKeS6LwKMCipM218EKE


//托管账户地址
//address:2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD
//script:522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53ae


//往托管账户转账
(async function(){
    network = bitcoin.networks.testnet;
    to = '2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD';
    from = 'myDAgFpwT3sTkppYKeS6LwKMCipM218EKE';
    fee = 1e3;
    value = 100e3;
    memo = (new Date()).toLocaleString();

    chunk = await bitcoinapi.getUTXOS(from, {IsTestnet:1});
    utxos = JSON.parse(chunk.data).unspent_outputs;

    txb = new bitcoin.TransactionBuilder(network);
    txb.addInput(utxos[0].tx_hash_big_endian, utxos[0].tx_output_n);

    txb.addOutput(to, value - fee);
    txb.addOutput(from, utxos[0].value - value);

    data = Buffer.from(memo, 'utf8');
    embed = bitcoin.payments.embed({ data : [data],network : network });
    txb.addOutput(embed.output, 0);

    //sign
    userwif = 'cQg6EvRRka4BtLpzFuuVKUptToqiJj3SzhjDZEcak5eybApzAEFc';
    keyPair = bitcoin.ECPair.fromWIF(userwif, network);
    txb.sign(0,keyPair);

    result = await bitcoinapi.broadcastTrx(txb.buildIncomplete(),{IsTestnet:1});
    console.log('trxid:' + result);
})();






//多重签名从托管账户提现
// (async function(){
//     script = '522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53ae'
//     //buildTrx
//     trx = await bitcoinapi.buildTrx('2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD','myDAgFpwT3sTkppYKeS6LwKMCipM218EKE',2e3,1e3,{IsTestnet:1})

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




