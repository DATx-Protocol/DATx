ethapi = require('../lib/ethapi')
expect = require('chai').expect;

//withdraw to 0xdedbe1ace5f723fe50cf0015b5efb7392efc118c
ethapi.withdraw('0x054892113ea10caa44187e9f4d180fc5dbe4e15a','0xf6bb0e08e268eb2826c076defbff24283694a63c'
,'0x054892113ea10caa44187e9f4d180fc5dbe4e15a',300000,'416F392F52B831E319AE44532FAB6123F43C777AF7E8FA26619FE8BD70EF1E81'
,ethapi.fromAscii('asdadadaqqqqqq'));

ethapi.withdraw('0x66af4e3d52cdfb3b629b7c8d4bdd221052a58ab5','0xf6bb0e08e268eb2826c076defbff24283694a63c'
,'0x054892113ea10caa44187e9f4d180fc5dbe4e15a',300000,'A94C93F07EA4DA1989C73058B369D240C7D613BF944D12404BD17126C42F7A0D'
,ethapi.fromAscii('asdadadaqqqqqq'))

