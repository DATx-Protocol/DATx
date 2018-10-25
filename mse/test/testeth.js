ethapi = require('../lib/ethapi');
const Web3 = require('web3');
const Tx = require('ethereumjs-tx');

const eth_endpoint = 'https://ropsten.infura.io/v3/c9bebaf8a3ca41d8a5818b9f5740c69f';
const web3 = new Web3(new Web3.providers.HttpProvider(eth_endpoint));

//托管账户
//address1:0x054892113EA10CAa44187E9f4D180Fc5dbE4e15A
//prikey1:416F392F52B831E319AE44532FAB6123F43C777AF7E8FA26619FE8BD70EF1E81

//address2:0x66Af4e3d52CdFB3b629b7C8d4bdd221052a58ab5
//prikey2:A94C93F07EA4DA1989C73058B369D240C7D613BF944D12404BD17126C42F7A0D

//address3:0x418c15E6777eF44ffF52b83adADAb5b2A64C9e68
//prikey3:75394070B62191CA0B30DD0A3B1A3AC4975E1515AAC189D5B18C48EBE63D162F

//multiaddress:0x03f2216C5fBeE2F881333F0eD5e0C7247eddF9C4


//用户账号
//address:0xDEdBe1ACe5f723Fe50cf0015B5EFb7392Efc118c
//prikey:4BE956D26884B3F6FC758E43B4EB60582D883DE271E442F872136362041A8565



(async function() {
    var from = '0x66Af4e3d52CdFB3b629b7C8d4bdd221052a58ab5';
    var to = '0xDEdBe1ACe5f723Fe50cf0015B5EFb7392Efc118c'
    var nonce = await web3.eth.getTransactionCount(from);
    var privateKey = Buffer.from('A94C93F07EA4DA1989C73058B369D240C7D613BF944D12404BD17126C42F7A0D', 'hex');

    var rawTx = {
    nonce: web3.utils.toHex(nonce),
    to: to,
    gasPrice: web3.utils.toHex(2000000000),
    gasLimit: web3.utils.toHex(300000),
    value: web3.utils.toHex(160e15),
    data: '',
    chainId: 3
    };
    var tx = new Tx(rawTx);
    tx.sign(privateKey);
    var serializedTx = tx.serialize();
    var rawparam = serializedTx.toString('hex');

    let result = await web3.eth.sendSignedTransaction('0x' + rawparam);
    console.log('trxid:' + result.transactionHash);
})();




//withdraw to 0xdedbe1ace5f723fe50cf0015b5efb7392efc118c
// ethapi.withdraw('0x054892113ea10caa44187e9f4d180fc5dbe4e15a','0x03f2216C5fBeE2F881333F0eD5e0C7247eddF9C4'
// ,'0x054892113ea10caa44187e9f4d180fc5dbe4e15a',300000,'416F392F52B831E319AE44532FAB6123F43C777AF7E8FA26619FE8BD70EF1E81'
// ,ethapi.fromAscii('wdvff'))

// ethapi.withdraw('0x66af4e3d52cdfb3b629b7c8d4bdd221052a58ab5','0x03f2216C5fBeE2F881333F0eD5e0C7247eddF9C4'
// ,'0x054892113ea10caa44187e9f4d180fc5dbe4e15a',300000,'A94C93F07EA4DA1989C73058B369D240C7D613BF944D12404BD17126C42F7A0D'
// ,ethapi.fromAscii('wdvff'))

