
 bitcoinapi = require('../lib/bitcoinapi');
// expect = require('chai').expect;

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

// describe('genP2PKHAddr',function(){
//     addr = bitcoinapi.genP2PKHAddr('03a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb',{IsTestnet:1})
//     expect(addr).to.equal('n4fc3bKTVrVRveBrHZf5Zv4wGGBTf5sdHg')
//  })

// describe('genMulSigAddr',function(){
//     pubkeys = ['03a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb',
//                '038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac9',
//                '02c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e']
//     addr = bitcoinapi.genMulSigAddr(pubkeys,2,{IsTestnet:1})
//     expect(addr.address).to.equal('2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD')
//     expect(addr.script).to.equal('522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53ae')
// })


// describe('simulate bitcoin withdraw',function(){
    // (async function(){
    //     script = '522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53ae'
    //     //buildTrx
    //     trx = await bitcoinapi.buildTrx('2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD','n4fc3bKTVrVRveBrHZf5Zv4wGGBTf5sdHg',2e3,1e3,{IsTestnet:1})

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
// })
const bitcoin = require('bitcoinjs-lib');
//result = bitcoin.classify.output(Buffer.from('6a0d626974636f696e6a732d6c6962','hex'));
//console.log(result);
//console.log(bitcoin.script.toASM(Buffer.from('6a0d626974636f696e6a732d6c6962','hex')));
result = bitcoin.script.decompile(Buffer.from('6a0d626974636f696e6a732d6c6962','hex')).slice(1).toString('utf8');
//console.log(result);

//result = Buffer.from('626974636f696e6a732d6c6962','hex').toString('utf8');
console.log(result);