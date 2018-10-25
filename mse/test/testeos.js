
const Eos = require('eosjs');
const chainid = '038f4b0fc8ff18a4f0842a8f0564611f6e96e8535901dd45e43ac8691a1c4dca';
const httpEndpoint = 'http://jungle.cryptolions.io:18888';
const eos = Eos({httpEndpoint, chainid, keyProvider: '5JytBc6QzyivnJzNpMXu3xv23mZCu8NJLawAkoxyuMctb7HXNwT',expireInSeconds:60,
broadcast: true,
debug: true, 
sign: true});

(async function(){
    auth = { threshold: 2, accounts: [{permission: {actor: "dotcdatx1234", permission: "active"}, weight: 1},
    {permission: {actor: "avazudatx123", permission: "active"}, weight: 1},
    {permission: {actor: "teebikdatx12", permission: "active"}, weight: 1}]};

    op_data = {
         account: 'datxmultisig',
         permission: 'active',
         parent: 'owner',
         auth: auth
         };

    const result = await eos.transaction(tr => {
        tr.updateauth(op_data, {authorization: 'datxmultisig@owner'});
    })
    //result = await eos.updateauth(op_data);
    console.log(result);
})();


// var keyProvider = ['5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3', //eosio的私钥
// '5Jf86MboxXs7x2Kfrm4UNqeoQsuN7PyEw7U8GPbBAjbfL3E1ckQ', //eosio.token的私钥
// '5JvJVQwFSYRbVncoKi3HwbN85vW3x3dmm9TkVpXALJgJCUXLFia',//alice的私钥
// '5JFsGvwD63dc8G2bzV7xmspHeRH1VbJpn2mmzUhd3UUi6RHNYVM']; //bob的私钥

// (async function() {
//     //create account
//     var auth = { threshold: 2,
//         accounts: [ { permission: {actor:'bob',permission:'active'}, weight: 1 },{ permission: {actor:'alice',permission:'active'}, weight: 1 } ],
//         waits: []}
//     obj = await eosapi.createAccount('eosio',keyProvider[0],'colin'
//     ,auth)
//     //obj && obj.transaction_id

//     //issue
//     issue = await eosapi.issue('colin','100.0000 EOS','')
//     //issue && issue.transaction_id

//     //propose
//     auth = [{actor: 'alice', permission: 'active'},{actor: 'bob', permission: 'active'}]
//     proposeName = await eosapi.propose('alice','colin','bob','50.0000 EOS',auth,keyProvider[2])

//     //alice confirm
//     confirm1 = await eosapi.confirm('alice',proposeName,'alice',keyProvider[2])

//     //bob comfirm
//     confirm2 = await eosapi.confirm('alice',proposeName,'bob',keyProvider[3])
    
//     //bob exec
//     result = await eosapi.exec('alice',proposeName,'bob',keyProvider[3])
// })()