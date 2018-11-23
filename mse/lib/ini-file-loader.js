var eol = process.platform === "win32" ? "\r\n" : "\n"
  
function INI() {
    this.sections = {};
}
  
/**
 * 删除Section
 * @param sectionName
 */
INI.prototype.removeSection = function (sectionName) {
  
    sectionName =  sectionName.replace(/\[/g,'(');
    sectionName = sectionName.replace(/]/g,')');
  
    if (this.sections[sectionName]) {
        delete this.sections[sectionName];
    }
}
/**
 * 创建或者得到某个Section
 * @type {Function}
 */
INI.prototype.getOrCreateSection = INI.prototype.section = function (sectionName) {
  
    sectionName =  sectionName.replace(/\[/g,'(');
    sectionName = sectionName.replace(/]/g,')');
  
    if (!this.sections[sectionName]) {
        this.sections[sectionName] = {};
    }
    return this.sections[sectionName]
}
  
/**
 * 将INI转换成文本
 *
 * @returns {string}
 */
INI.prototype.encodeToIni = INI.prototype.toString = function encodeIni() {
    var _INI = this;
    var sectionOut = _INI.encodeSection(null, _INI);
    Object.keys(_INI.sections).forEach(function (k, _, __) {
        if (_INI.sections) {
            sectionOut += _INI.encodeSection(k, _INI.sections[k])
        }
    });
    return sectionOut;
}
  
/**
 *
 * @param section
 * @param obj
 * @returns {string}
 */
INI.prototype.encodeSection = function (section, obj) {
    var out = "";
    Object.keys(obj).forEach(function (k, _, __) {
        var val = obj[k]
        if (val && Array.isArray(val)) {
            val.forEach(function (item) {
                out += safe(k + "[]") + " = " + safe(item) + "\n"
            })
        } else if (val && typeof val === "object") {
        } else {
            out += safe(k) + " = " + safe(val) + eol
        }
    })
    if (section && out.length) {
        out = "[" + safe(section) + "]" + eol + out
    }
    return out+"\n";
}
function safe(val) {
    return (typeof val !== "string" || val.match(/[\r\n]/) || val.match(/^\[/) || (val.length > 1 && val.charAt(0) === "\"" && val.slice(-1) === "\"") || val !== val.trim()) ? JSON.stringify(val) : val.replace(/;/g, '\\;')
}
  
var regex = {
    section: /^\s*\[\s*([^\]]*)\s*\]\s*$/,
    param: /^\s*([\w\.\-\_]+)\s*=\s*(.*?)\s*$/,
    comment: /^\s*;.*$/
};
  
/**
 *
 * @param data
 * @returns {INI}
 */
exports.parse = function (data) {
    var value = new INI();
    var lines = data.split(/\r\n|\r|\n/);
    var section = null;
    lines.forEach(function (line) {
        if (regex.comment.test(line)) {
            return;
        } else if (regex.param.test(line)) {
            var match = line.match(regex.param);
            if (section) {
                section[match[1]] = match[2];
            } else {
                value[match[1]] = match[2];
            }
        } else if (regex.section.test(line)) {
            var match = line.match(regex.section);
            section = value.getOrCreateSection(match[1])
        } else if (line.length == 0 && section) {
            section = null;
        }
        ;
    });
    return value;
}
  
/**
 * 创建INI
 * @type {Function}
 */
exports.createINI = exports.create = function () {
    return new INI();
};
  
var fs = require('fs');
  
exports.loadFileSync =function(fileName/*,charset*/){
    return exports.parse(fs.readFileSync(fileName, "utf-8")) ;
}

const os = require('os');
exports.getConfigFile = function(){
    let platform = os.platform(); 
    if(platform == 'linux'){
         cfgFile = os.homedir + '/.local/share/datxos/noddatx/config/config.ini';
    }
    else if (platform == 'darwin'){
         cfgFile = os.homedir + '/Library/Application Support/datxos/noddatx/config/config.ini';
    }
    if(cfgFile == '' || cfgFile == null || cfgFile == undefined){
        return new INI();
    }

    __ini = exports.loadFileSync(cfgFile)
    let flag = __ini["across-chain-cfg-flag"];
    if(flag != '1'){
        fs.appendFileSync(cfgFile,'\n' + 
        'btc-muladdress = 2MwtNFrT9P1wDa3Hid6kVu9mh84cd59UHKN \n' + 
        'btc-mulscript = 522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102040e0d9141b06ad92f38d7a3d76cfcb6ada4c9e4c5b18d18f5539564a382640853ae \n' + 
        'btc-wif = cQX7uYfhPJJGPdoqgHQXLC4WQWcDXGH7DgLs4m89xWykYkg1V4iA \n' + 
        'btc-maxfee = 200000 \n' + 
        //'eth-endpoint = https://ropsten.infura.io/v3/c9bebaf8a3ca41d8a5818b9f5740c69f \n' + 
        //'eth-endpoint-ws = wss://ropsten.infura.io/ws \n' + 
        //'eth-muladdress = 0xF6Bb0E08E268Eb2826C076dEFbFf24283694a63c \n' + 
        'eth-endpoint = https://mainnet.infura.io/v3/c9bebaf8a3ca41d8a5818b9f5740c69f \n' + 
        'eth-endpoint-ws = wss://mainnet.infura.io/ws \n' + 
        'eth-myaddress = 0x054892113EA10CAa44187E9f4D180Fc5dbE4e15A \n' + 
        'eth-privatekey = 416F392F52B831E319AE44532FAB6123F43C777AF7E8FA26619FE8BD70EF1E81 \n' + 
        'eth-muladdress = 0xDaBBbacFF575a85E47EA1D1ff97f55B22Ab3184A \n' + 
        'eth-gasprice = 10000000000 \n' + 
        'eth-gaslimit = 300000 \n' + 
        'eos-chainid = aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906 \n' + 
        'eos-endpoint = http://213.239.208.37:8888 \n' + 
        'eos-account = alice \n' + 
        'eos-privateKey = 5JvJVQwFSYRbVncoKi3HwbN85vW3x3dmm9TkVpXALJgJCUXLFia \n' + 
        'eos-mulAccount = datxtest1112 \n' + 
        'listen-server-endpoint = http://127.0.0.1:8880 \n' + 
        'across-chain-endpoint = https://127.0.0.1:8080 \n' +
        'across-chain-cfg-flag = 1');

        __ini = exports.loadFileSync(cfgFile)
    }
    return __ini;
}

exports.getWallet = function(){
    let walletCfgFile = os.homedir + "/datxos-wallet/wallet_password.ini";
    __ini_wallet = exports.loadFileSync(walletCfgFile);

    return {wname : __ini_wallet["wallet-namer"], wpassword : __ini_wallet["wallet-password"]};
}