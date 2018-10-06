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