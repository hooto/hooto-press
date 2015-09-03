// CodeMirror, copyright (c) by Marijn Haverbeke and others
// Distributed under an MIT license: http://codemirror.net/LICENSE
(function(e){typeof exports=="object"&&typeof module=="object"?e(require("../../lib/codemirror")):typeof define=="function"&&define.amd?define(["../../lib/codemirror"],e):e(CodeMirror)})(function(e){function n(e,t){for(var n=0,r=e.length;n<r;++n)t(e[n])}function r(e,t){if(!Array.prototype.indexOf){var n=e.length;while(n--)if(e[n]===t)return!0;return!1}return e.indexOf(t)!=-1}function i(n,r,i,s){var o=n.getCursor(),u=i(n,o);if(/\b(?:string|comment)\b/.test(u.type))return;u.state=e.innerMode(n.getMode(),u.state).state,/^[\w$_]*$/.test(u.string)?u.end>o.ch&&(u.end=o.ch,u.string=u.string.slice(0,o.ch-u.start)):u={start:o.ch,end:o.ch,string:"",state:u.state,type:u.string=="."?"property":null};var a=u;while(a.type=="property"){a=i(n,t(o.line,a.start));if(a.string!=".")return;a=i(n,t(o.line,a.start));if(!f)var f=[];f.push(a)}return{list:p(u,f,r,s),from:t(o.line,u.start),to:t(o.line,u.end)}}function s(e,t){return i(e,c,function(e,t){return e.getTokenAt(t)},t)}function o(e,t){var n=e.getTokenAt(t);return t.ch==n.start+1&&n.string.charAt(0)=="."?(n.end=n.start,n.string=".",n.type="property"):/^\.[\w$_]*$/.test(n.string)&&(n.type="property",n.start++,n.string=n.string.replace(/\./,"")),n}function u(e,t){return i(e,h,o,t)}function p(e,t,i,s){function h(e){e.lastIndexOf(u,0)==0&&!r(o,e)&&o.push(e)}function p(e){typeof e=="string"?n(a,h):e instanceof Array?n(f,h):e instanceof Function&&n(l,h);for(var t in e)h(t)}var o=[],u=e.string,c=s&&s.globalScope||window;if(t&&t.length){var d=t.pop(),v;if(d.type&&d.type.indexOf("variable")===0){s&&s.additionalContext&&(v=s.additionalContext[d.string]);if(!s||s.useGlobalScope!==!1)v=v||c[d.string]}else d.type=="string"?v="":d.type=="atom"?v=1:d.type=="function"&&(c.jQuery==null||d.string!="$"&&d.string!="jQuery"||typeof c.jQuery!="function"?c._!=null&&d.string=="_"&&typeof c._=="function"&&(v=c._()):v=c.jQuery());while(v!=null&&t.length)v=v[t.pop().string];v!=null&&p(v)}else{for(var m=e.state.localVars;m;m=m.next)h(m.name);for(var m=e.state.globalVars;m;m=m.next)h(m.name);(!s||s.useGlobalScope!==!1)&&p(c),n(i,h)}return o}var t=e.Pos;e.registerHelper("hint","javascript",s),e.registerHelper("hint","coffeescript",u);var a="charAt charCodeAt indexOf lastIndexOf substring substr slice trim trimLeft trimRight toUpperCase toLowerCase split concat match replace search".split(" "),f="length concat join splice push pop shift unshift slice reverse sort indexOf lastIndexOf every some filter forEach map reduce reduceRight ".split(" "),l="prototype apply call bind".split(" "),c="break case catch continue debugger default delete do else false finally for function if in instanceof new null return switch throw true try typeof var void while with".split(" "),h="and break catch class continue delete do else extends false finally for if in instanceof isnt new no not null of off on or return switch then throw true try typeof until void while with yes".split(" ")});