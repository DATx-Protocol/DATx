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
        'btc-muladdress = 2MsimupueVskjJMy79kKGP5uzfCWfZuK8TD \n' + 
        'btc-mulscript = 522103a4ac53ded034de0ce8e1a5aa8cae967a7c33f8ef807ee31d0a972fbcd912c8cb21038e6c355aa3a7b0a3338215e1fb952c1c255eab07012c800a151f8fd7bb9feac92102c8a936b526d91e6047569ec8fd53779a2368a150d63cea655fc9c7ba66d2199e53ae \n' + 
        'btc-wif = cQX7uYfhPJJGPdoqgHQXLC4WQWcDXGH7DgLs4m89xWykYkg1V4iA \n' + 
        'btc-maxfee = 200000 \n' + 
        'eth-endpoint = https://ropsten.infura.io/v3/c9bebaf8a3ca41d8a5818b9f5740c69f \n' + 
        'eth-endpoint-ws = wss://ropsten.infura.io/ws \n' + 
        'eth-myaddress = 0x054892113ea10caa44187e9f4d180fc5dbe4e15a \n' + 
        'eth-privatekey = 416F392F52B831E319AE44532FAB6123F43C777AF7E8FA26619FE8BD70EF1E81 \n' + 
        'eth-muladdress = 0xF6Bb0E08E268Eb2826C076dEFbFf24283694a63c \n' + 
        'eth-gasprice = 2000000000 \n' + 
        'eth-gaslimit = 300000 \n' + 
        'eos-chainid = cf057bbfb72640471fd910bcb67639c22df9f92470936cddc1ade0e2f2e7dc4f \n' + 
        'eos-endpoint = http://127.0.0.1:8888 \n' + 
        'eos-account = alice \n' + 
        'eos-privateKey = 5JvJVQwFSYRbVncoKi3HwbN85vW3x3dmm9TkVpXALJgJCUXLFia \n' + 
        'eos-mulAccount = jacky \n' + 
        'listen-server-endpoint = http://127.0.0.1:8880 \n' + 
        'across-chain-endpoint = https://127.0.0.1:8080 \n' +
        'across-chain-cfg-flag = 1');

        __ini = exports.loadFileSync(cfgFile)
    }
    return __ini;
}