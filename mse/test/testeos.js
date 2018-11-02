//多重签名账户
// 账户名：datxmultisig
//Active Key: avazudatx123@active, dotcdatx1234@active, teebikdatx12@active
//Owner Key: avazudatx123@owner, dotcdatx1234@owner, teebikdatx12@owner


// 三个托管账户
// 账户名：dotcdatx1234 
// 公钥：EOS7PpgrAJYDA7Gvm9ZEtWZinuHSMoAVpRfwDANXtwMj3aJgTKAY2
// 私钥：5J4TL6uqQSGssBErRjge4QdsuqhYnhd5BeN64c9YZa8vSwssELH
 
// 账户名：avazudatx123
// 公钥：EOS7maepAkgWFUnysQvQWPdH3y51mTaSe3HjiPwgd4Zoc9cJ41oFg
// 私钥: 5JwXDKoeZL4xhAy5B9pR2ErfUt3DD9c9AL2uwNj4TSj3okknKej

// 账户名：teebikdatx12 
// 公钥: EOS6KDEG3T4owpoEkrquSzCasSD6NexHQ8xjES1TFWUosognHgCbw
// 私钥: 5KTk5Y7sYR1gKEa6CskbbX1GZADYUinZ4c31aq8Zz2fqssSfD1F

//用户账户：datxuser1234
// 公钥: EOS6DpURszfu17r83dfn8Pb9xaAs2Vko1ZMhx75FZPsX1vorWYffP
// 私钥: 5JytBc6QzyivnJzNpMXu3xv23mZCu8NJLawAkoxyuMctb7HXNwT


const Eos = require('eosjs');
const chainId = '038f4b0fc8ff18a4f0842a8f0564611f6e96e8535901dd45e43ac8691a1c4dca';
const httpEndpoint = 'http://jungle.cryptolions.io:18888';
config = {httpEndpoint, chainId, verbose: true, forceActionDataHex: false}
eos = Eos(Object.assign(config, {keyProvider: () => userProvidedKey}));

let userProvidedKey = null;

(async function(){
    userProvidedKey = '5JytBc6QzyivnJzNpMXu3xv23mZCu8NJLawAkoxyuMctb7HXNwT';
    result = await eos.transfer('datxuser1234', 'datxmultisig', '1.0000 EOS', 'test');
    console.log('trxid: ' + result.transaction_id);
})();


//set auth
// (async function(){
//     auth = { threshold: 2, accounts: [{permission: {actor: "dotcdatx1234", permission: "owner"}, weight: 1},
//     {permission: {actor: "avazudatx123", permission: "owner"}, weight: 1},
//     {permission: {actor: "teebikdatx12", permission: "owner"}, weight: 1}]};

//     op_data = {
//          account: 'datxmultisig',
//          permission: 'owner',
//          parent: '',
//          auth: auth
//          };

//     const result = await eos.transaction(tr => {
//         tr.updateauth(op_data, {authorization: 'datxmultisig@owner'});
//     })
//     console.log(result);
// })();


