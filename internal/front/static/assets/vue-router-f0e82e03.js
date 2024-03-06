import"./vue-d82d38a0.js";import{B as rt,e as I,h as ze,c as N,K as qe,p as ae,a as ot,z as st,H as it,al as ct}from"./vue-runtime-core-ffab33f8.js";import{s as at,u as Q,h as lt,a as ut,r as ft}from"./vue-reactivity-557f7d4b.js";/*!
  * vue-router v4.2.4
  * (c) 2023 Eduardo San Martin Morote
  * @license MIT
  */const z=typeof window<"u";function ht(e){return e.__esModule||e[Symbol.toStringTag]==="Module"}const S=Object.assign;function le(e,t){const n={};for(const r in t){const o=t[r];n[r]=L(o)?o.map(e):e(o)}return n}const F=()=>{},L=Array.isArray,dt=/\/$/,pt=e=>e.replace(dt,"");function ue(e,t,n="/"){let r,o={},l="",d="";const m=t.indexOf("#");let c=t.indexOf("?");return m<c&&m>=0&&(c=-1),c>-1&&(r=t.slice(0,c),l=t.slice(c+1,m>-1?m:t.length),o=e(l)),m>-1&&(r=r||t.slice(0,m),d=t.slice(m,t.length)),r=yt(r??t,n),{fullPath:r+(l&&"?")+l+d,path:r,query:o,hash:d}}function mt(e,t){const n=t.query?e(t.query):"";return t.path+(n&&"?")+n+(t.hash||"")}function be(e,t){return!t||!e.toLowerCase().startsWith(t.toLowerCase())?e:e.slice(t.length)||"/"}function gt(e,t,n){const r=t.matched.length-1,o=n.matched.length-1;return r>-1&&r===o&&q(t.matched[r],n.matched[o])&&Ke(t.params,n.params)&&e(t.query)===e(n.query)&&t.hash===n.hash}function q(e,t){return(e.aliasOf||e)===(t.aliasOf||t)}function Ke(e,t){if(Object.keys(e).length!==Object.keys(t).length)return!1;for(const n in e)if(!vt(e[n],t[n]))return!1;return!0}function vt(e,t){return L(e)?ke(e,t):L(t)?ke(t,e):e===t}function ke(e,t){return L(t)?e.length===t.length&&e.every((n,r)=>n===t[r]):e.length===1&&e[0]===t}function yt(e,t){if(e.startsWith("/"))return e;if(!e)return t;const n=t.split("/"),r=e.split("/"),o=r[r.length-1];(o===".."||o===".")&&r.push("");let l=n.length-1,d,m;for(d=0;d<r.length;d++)if(m=r[d],m!==".")if(m==="..")l>1&&l--;else break;return n.slice(0,l).join("/")+"/"+r.slice(d-(d===r.length?1:0)).join("/")}var X;(function(e){e.pop="pop",e.push="push"})(X||(X={}));var Y;(function(e){e.back="back",e.forward="forward",e.unknown=""})(Y||(Y={}));function Rt(e){if(!e)if(z){const t=document.querySelector("base");e=t&&t.getAttribute("href")||"/",e=e.replace(/^\w+:\/\/[^\/]+/,"")}else e="/";return e[0]!=="/"&&e[0]!=="#"&&(e="/"+e),pt(e)}const Et=/^[^#]+#/;function Pt(e,t){return e.replace(Et,"#")+t}function wt(e,t){const n=document.documentElement.getBoundingClientRect(),r=e.getBoundingClientRect();return{behavior:t.behavior,left:r.left-n.left-(t.left||0),top:r.top-n.top-(t.top||0)}}const te=()=>({left:window.pageXOffset,top:window.pageYOffset});function St(e){let t;if("el"in e){const n=e.el,r=typeof n=="string"&&n.startsWith("#"),o=typeof n=="string"?r?document.getElementById(n.slice(1)):document.querySelector(n):n;if(!o)return;t=wt(o,e)}else t=e;"scrollBehavior"in document.documentElement.style?window.scrollTo(t):window.scrollTo(t.left!=null?t.left:window.pageXOffset,t.top!=null?t.top:window.pageYOffset)}function Ae(e,t){return(history.state?history.state.position-t:-1)+e}const he=new Map;function Ct(e,t){he.set(e,t)}function bt(e){const t=he.get(e);return he.delete(e),t}let kt=()=>location.protocol+"//"+location.host;function Ue(e,t){const{pathname:n,search:r,hash:o}=t,l=e.indexOf("#");if(l>-1){let m=o.includes(e.slice(l))?e.slice(l).length:1,c=o.slice(m);return c[0]!=="/"&&(c="/"+c),be(c,"")}return be(n,e)+r+o}function At(e,t,n,r){let o=[],l=[],d=null;const m=({state:u})=>{const g=Ue(e,location),R=n.value,k=t.value;let b=0;if(u){if(n.value=g,t.value=u,d&&d===R){d=null;return}b=k?u.position-k.position:0}else r(g);o.forEach(E=>{E(n.value,R,{delta:b,type:X.pop,direction:b?b>0?Y.forward:Y.back:Y.unknown})})};function c(){d=n.value}function f(u){o.push(u);const g=()=>{const R=o.indexOf(u);R>-1&&o.splice(R,1)};return l.push(g),g}function s(){const{history:u}=window;u.state&&u.replaceState(S({},u.state,{scroll:te()}),"")}function a(){for(const u of l)u();l=[],window.removeEventListener("popstate",m),window.removeEventListener("beforeunload",s)}return window.addEventListener("popstate",m),window.addEventListener("beforeunload",s,{passive:!0}),{pauseListeners:c,listen:f,destroy:a}}function Oe(e,t,n,r=!1,o=!1){return{back:e,current:t,forward:n,replaced:r,position:window.history.length,scroll:o?te():null}}function Ot(e){const{history:t,location:n}=window,r={value:Ue(e,n)},o={value:t.state};o.value||l(r.value,{back:null,current:r.value,forward:null,position:t.length-1,replaced:!0,scroll:null},!0);function l(c,f,s){const a=e.indexOf("#"),u=a>-1?(n.host&&document.querySelector("base")?e:e.slice(a))+c:kt()+e+c;try{t[s?"replaceState":"pushState"](f,"",u),o.value=f}catch(g){console.error(g),n[s?"replace":"assign"](u)}}function d(c,f){const s=S({},t.state,Oe(o.value.back,c,o.value.forward,!0),f,{position:o.value.position});l(c,s,!0),r.value=c}function m(c,f){const s=S({},o.value,t.state,{forward:c,scroll:te()});l(s.current,s,!0);const a=S({},Oe(r.value,c,null),{position:s.position+1},f);l(c,a,!1),r.value=c}return{location:r,state:o,push:m,replace:d}}function _t(e){e=Rt(e);const t=Ot(e),n=At(e,t.state,t.location,t.replace);function r(l,d=!0){d||n.pauseListeners(),history.go(l)}const o=S({location:"",base:e,go:r,createHref:Pt.bind(null,e)},t,n);return Object.defineProperty(o,"location",{enumerable:!0,get:()=>t.location.value}),Object.defineProperty(o,"state",{enumerable:!0,get:()=>t.state.value}),o}function vn(e){return e=location.host?e||location.pathname+location.search:"",e.includes("#")||(e+="#"),_t(e)}function xt(e){return typeof e=="string"||e&&typeof e=="object"}function Ve(e){return typeof e=="string"||typeof e=="symbol"}const T={path:"/",name:void 0,params:{},query:{},hash:"",fullPath:"/",matched:[],meta:{},redirectedFrom:void 0},De=Symbol("");var _e;(function(e){e[e.aborted=4]="aborted",e[e.cancelled=8]="cancelled",e[e.duplicated=16]="duplicated"})(_e||(_e={}));function K(e,t){return S(new Error,{type:e,[De]:!0},t)}function H(e,t){return e instanceof Error&&De in e&&(t==null||!!(e.type&t))}const xe="[^/]+?",Mt={sensitive:!1,strict:!1,start:!0,end:!0},Lt=/[.+*?^${}()[\]/\\]/g;function Nt(e,t){const n=S({},Mt,t),r=[];let o=n.start?"^":"";const l=[];for(const f of e){const s=f.length?[]:[90];n.strict&&!f.length&&(o+="/");for(let a=0;a<f.length;a++){const u=f[a];let g=40+(n.sensitive?.25:0);if(u.type===0)a||(o+="/"),o+=u.value.replace(Lt,"\\$&"),g+=40;else if(u.type===1){const{value:R,repeatable:k,optional:b,regexp:E}=u;l.push({name:R,repeatable:k,optional:b});const w=E||xe;if(w!==xe){g+=10;try{new RegExp(`(${w})`)}catch(M){throw new Error(`Invalid custom RegExp for param "${R}" (${w}): `+M.message)}}let _=k?`((?:${w})(?:/(?:${w}))*)`:`(${w})`;a||(_=b&&f.length<2?`(?:/${_})`:"/"+_),b&&(_+="?"),o+=_,g+=20,b&&(g+=-8),k&&(g+=-20),w===".*"&&(g+=-50)}s.push(g)}r.push(s)}if(n.strict&&n.end){const f=r.length-1;r[f][r[f].length-1]+=.7000000000000001}n.strict||(o+="/?"),n.end?o+="$":n.strict&&(o+="(?:/|$)");const d=new RegExp(o,n.sensitive?"":"i");function m(f){const s=f.match(d),a={};if(!s)return null;for(let u=1;u<s.length;u++){const g=s[u]||"",R=l[u-1];a[R.name]=g&&R.repeatable?g.split("/"):g}return a}function c(f){let s="",a=!1;for(const u of e){(!a||!s.endsWith("/"))&&(s+="/"),a=!1;for(const g of u)if(g.type===0)s+=g.value;else if(g.type===1){const{value:R,repeatable:k,optional:b}=g,E=R in f?f[R]:"";if(L(E)&&!k)throw new Error(`Provided param "${R}" is an array but it is not repeatable (* or + modifiers)`);const w=L(E)?E.join("/"):E;if(!w)if(b)u.length<2&&(s.endsWith("/")?s=s.slice(0,-1):a=!0);else throw new Error(`Missing required param "${R}"`);s+=w}}return s||"/"}return{re:d,score:r,keys:l,parse:m,stringify:c}}function Bt(e,t){let n=0;for(;n<e.length&&n<t.length;){const r=t[n]-e[n];if(r)return r;n++}return e.length<t.length?e.length===1&&e[0]===40+40?-1:1:e.length>t.length?t.length===1&&t[0]===40+40?1:-1:0}function Ht(e,t){let n=0;const r=e.score,o=t.score;for(;n<r.length&&n<o.length;){const l=Bt(r[n],o[n]);if(l)return l;n++}if(Math.abs(o.length-r.length)===1){if(Me(r))return 1;if(Me(o))return-1}return o.length-r.length}function Me(e){const t=e[e.length-1];return e.length>0&&t[t.length-1]<0}const It={type:0,value:""},Tt=/[a-zA-Z0-9_]/;function $t(e){if(!e)return[[]];if(e==="/")return[[It]];if(!e.startsWith("/"))throw new Error(`Invalid path "${e}"`);function t(g){throw new Error(`ERR (${n})/"${f}": ${g}`)}let n=0,r=n;const o=[];let l;function d(){l&&o.push(l),l=[]}let m=0,c,f="",s="";function a(){f&&(n===0?l.push({type:0,value:f}):n===1||n===2||n===3?(l.length>1&&(c==="*"||c==="+")&&t(`A repeatable param (${f}) must be alone in its segment. eg: '/:ids+.`),l.push({type:1,value:f,regexp:s,repeatable:c==="*"||c==="+",optional:c==="*"||c==="?"})):t("Invalid state to consume buffer"),f="")}function u(){f+=c}for(;m<e.length;){if(c=e[m++],c==="\\"&&n!==2){r=n,n=4;continue}switch(n){case 0:c==="/"?(f&&a(),d()):c===":"?(a(),n=1):u();break;case 4:u(),n=r;break;case 1:c==="("?n=2:Tt.test(c)?u():(a(),n=0,c!=="*"&&c!=="?"&&c!=="+"&&m--);break;case 2:c===")"?s[s.length-1]=="\\"?s=s.slice(0,-1)+c:n=3:s+=c;break;case 3:a(),n=0,c!=="*"&&c!=="?"&&c!=="+"&&m--,s="";break;default:t("Unknown state");break}}return n===2&&t(`Unfinished custom RegExp for param "${f}"`),a(),d(),o}function jt(e,t,n){const r=Nt($t(e.path),n),o=S(r,{record:e,parent:t,children:[],alias:[]});return t&&!o.record.aliasOf==!t.record.aliasOf&&t.children.push(o),o}function Gt(e,t){const n=[],r=new Map;t=Be({strict:!1,end:!0,sensitive:!1},t);function o(s){return r.get(s)}function l(s,a,u){const g=!u,R=zt(s);R.aliasOf=u&&u.record;const k=Be(t,s),b=[R];if("alias"in s){const _=typeof s.alias=="string"?[s.alias]:s.alias;for(const M of _)b.push(S({},R,{components:u?u.record.components:R.components,path:M,aliasOf:u?u.record:R}))}let E,w;for(const _ of b){const{path:M}=_;if(a&&M[0]!=="/"){const j=a.record.path,B=j[j.length-1]==="/"?"":"/";_.path=a.record.path+(M&&B+M)}if(E=jt(_,a,k),u?u.alias.push(E):(w=w||E,w!==E&&w.alias.push(E),g&&s.name&&!Ne(E)&&d(s.name)),R.children){const j=R.children;for(let B=0;B<j.length;B++)l(j[B],E,u&&u.children[B])}u=u||E,(E.record.components&&Object.keys(E.record.components).length||E.record.name||E.record.redirect)&&c(E)}return w?()=>{d(w)}:F}function d(s){if(Ve(s)){const a=r.get(s);a&&(r.delete(s),n.splice(n.indexOf(a),1),a.children.forEach(d),a.alias.forEach(d))}else{const a=n.indexOf(s);a>-1&&(n.splice(a,1),s.record.name&&r.delete(s.record.name),s.children.forEach(d),s.alias.forEach(d))}}function m(){return n}function c(s){let a=0;for(;a<n.length&&Ht(s,n[a])>=0&&(s.record.path!==n[a].record.path||!We(s,n[a]));)a++;n.splice(a,0,s),s.record.name&&!Ne(s)&&r.set(s.record.name,s)}function f(s,a){let u,g={},R,k;if("name"in s&&s.name){if(u=r.get(s.name),!u)throw K(1,{location:s});k=u.record.name,g=S(Le(a.params,u.keys.filter(w=>!w.optional).map(w=>w.name)),s.params&&Le(s.params,u.keys.map(w=>w.name))),R=u.stringify(g)}else if("path"in s)R=s.path,u=n.find(w=>w.re.test(R)),u&&(g=u.parse(R),k=u.record.name);else{if(u=a.name?r.get(a.name):n.find(w=>w.re.test(a.path)),!u)throw K(1,{location:s,currentLocation:a});k=u.record.name,g=S({},a.params,s.params),R=u.stringify(g)}const b=[];let E=u;for(;E;)b.unshift(E.record),E=E.parent;return{name:k,path:R,params:g,matched:b,meta:Kt(b)}}return e.forEach(s=>l(s)),{addRoute:l,resolve:f,removeRoute:d,getRoutes:m,getRecordMatcher:o}}function Le(e,t){const n={};for(const r of t)r in e&&(n[r]=e[r]);return n}function zt(e){return{path:e.path,redirect:e.redirect,name:e.name,meta:e.meta||{},aliasOf:void 0,beforeEnter:e.beforeEnter,props:qt(e),children:e.children||[],instances:{},leaveGuards:new Set,updateGuards:new Set,enterCallbacks:{},components:"components"in e?e.components||null:e.component&&{default:e.component}}}function qt(e){const t={},n=e.props||!1;if("component"in e)t.default=n;else for(const r in e.components)t[r]=typeof n=="object"?n[r]:n;return t}function Ne(e){for(;e;){if(e.record.aliasOf)return!0;e=e.parent}return!1}function Kt(e){return e.reduce((t,n)=>S(t,n.meta),{})}function Be(e,t){const n={};for(const r in e)n[r]=r in t?t[r]:e[r];return n}function We(e,t){return t.children.some(n=>n===e||We(e,n))}const Qe=/#/g,Ut=/&/g,Vt=/\//g,Dt=/=/g,Wt=/\?/g,Fe=/\+/g,Qt=/%5B/g,Ft=/%5D/g,Ye=/%5E/g,Yt=/%60/g,Xe=/%7B/g,Xt=/%7C/g,Ze=/%7D/g,Zt=/%20/g;function me(e){return encodeURI(""+e).replace(Xt,"|").replace(Qt,"[").replace(Ft,"]")}function Jt(e){return me(e).replace(Xe,"{").replace(Ze,"}").replace(Ye,"^")}function de(e){return me(e).replace(Fe,"%2B").replace(Zt,"+").replace(Qe,"%23").replace(Ut,"%26").replace(Yt,"`").replace(Xe,"{").replace(Ze,"}").replace(Ye,"^")}function en(e){return de(e).replace(Dt,"%3D")}function tn(e){return me(e).replace(Qe,"%23").replace(Wt,"%3F")}function nn(e){return e==null?"":tn(e).replace(Vt,"%2F")}function ee(e){try{return decodeURIComponent(""+e)}catch{}return""+e}function rn(e){const t={};if(e===""||e==="?")return t;const r=(e[0]==="?"?e.slice(1):e).split("&");for(let o=0;o<r.length;++o){const l=r[o].replace(Fe," "),d=l.indexOf("="),m=ee(d<0?l:l.slice(0,d)),c=d<0?null:ee(l.slice(d+1));if(m in t){let f=t[m];L(f)||(f=t[m]=[f]),f.push(c)}else t[m]=c}return t}function He(e){let t="";for(let n in e){const r=e[n];if(n=en(n),r==null){r!==void 0&&(t+=(t.length?"&":"")+n);continue}(L(r)?r.map(l=>l&&de(l)):[r&&de(r)]).forEach(l=>{l!==void 0&&(t+=(t.length?"&":"")+n,l!=null&&(t+="="+l))})}return t}function on(e){const t={};for(const n in e){const r=e[n];r!==void 0&&(t[n]=L(r)?r.map(o=>o==null?null:""+o):r==null?r:""+r)}return t}const ge=Symbol(""),Ie=Symbol(""),ne=Symbol(""),ve=Symbol(""),pe=Symbol("");function W(){let e=[];function t(r){return e.push(r),()=>{const o=e.indexOf(r);o>-1&&e.splice(o,1)}}function n(){e=[]}return{add:t,list:()=>e.slice(),reset:n}}function Je(e,t,n){const r=()=>{e[t].delete(n)};st(r),it(r),ct(()=>{e[t].add(n)}),e[t].add(n)}function yn(e){const t=I(ge,{}).value;t&&Je(t,"leaveGuards",e)}function Rn(e){const t=I(ge,{}).value;t&&Je(t,"updateGuards",e)}function $(e,t,n,r,o){const l=r&&(r.enterCallbacks[o]=r.enterCallbacks[o]||[]);return()=>new Promise((d,m)=>{const c=a=>{a===!1?m(K(4,{from:n,to:t})):a instanceof Error?m(a):xt(a)?m(K(2,{from:t,to:a})):(l&&r.enterCallbacks[o]===l&&typeof a=="function"&&l.push(a),d())},f=e.call(r&&r.instances[o],t,n,c);let s=Promise.resolve(f);e.length<3&&(s=s.then(c)),s.catch(a=>m(a))})}function fe(e,t,n,r){const o=[];for(const l of e)for(const d in l.components){let m=l.components[d];if(!(t!=="beforeRouteEnter"&&!l.instances[d]))if(sn(m)){const f=(m.__vccOpts||m)[t];f&&o.push($(f,n,r,l,d))}else{let c=m();o.push(()=>c.then(f=>{if(!f)return Promise.reject(new Error(`Couldn't resolve component "${d}" at "${l.path}"`));const s=ht(f)?f.default:f;l.components[d]=s;const u=(s.__vccOpts||s)[t];return u&&$(u,n,r,l,d)()}))}}return o}function sn(e){return typeof e=="object"||"displayName"in e||"props"in e||"__vccOpts"in e}function Te(e){const t=I(ne),n=I(ve),r=N(()=>t.resolve(Q(e.to))),o=N(()=>{const{matched:c}=r.value,{length:f}=c,s=c[f-1],a=n.matched;if(!s||!a.length)return-1;const u=a.findIndex(q.bind(null,s));if(u>-1)return u;const g=$e(c[f-2]);return f>1&&$e(s)===g&&a[a.length-1].path!==g?a.findIndex(q.bind(null,c[f-2])):u}),l=N(()=>o.value>-1&&un(n.params,r.value.params)),d=N(()=>o.value>-1&&o.value===n.matched.length-1&&Ke(n.params,r.value.params));function m(c={}){return ln(c)?t[Q(e.replace)?"replace":"push"](Q(e.to)).catch(F):Promise.resolve()}return{route:r,href:N(()=>r.value.href),isActive:l,isExactActive:d,navigate:m}}const cn=ze({name:"RouterLink",compatConfig:{MODE:3},props:{to:{type:[String,Object],required:!0},replace:Boolean,activeClass:String,exactActiveClass:String,custom:Boolean,ariaCurrentValue:{type:String,default:"page"}},useLink:Te,setup(e,{slots:t}){const n=ut(Te(e)),{options:r}=I(ne),o=N(()=>({[je(e.activeClass,r.linkActiveClass,"router-link-active")]:n.isActive,[je(e.exactActiveClass,r.linkExactActiveClass,"router-link-exact-active")]:n.isExactActive}));return()=>{const l=t.default&&t.default(n);return e.custom?l:qe("a",{"aria-current":n.isExactActive?e.ariaCurrentValue:null,href:n.href,onClick:n.navigate,class:o.value},l)}}}),an=cn;function ln(e){if(!(e.metaKey||e.altKey||e.ctrlKey||e.shiftKey)&&!e.defaultPrevented&&!(e.button!==void 0&&e.button!==0)){if(e.currentTarget&&e.currentTarget.getAttribute){const t=e.currentTarget.getAttribute("target");if(/\b_blank\b/i.test(t))return}return e.preventDefault&&e.preventDefault(),!0}}function un(e,t){for(const n in t){const r=t[n],o=e[n];if(typeof r=="string"){if(r!==o)return!1}else if(!L(o)||o.length!==r.length||r.some((l,d)=>l!==o[d]))return!1}return!0}function $e(e){return e?e.aliasOf?e.aliasOf.path:e.path:""}const je=(e,t,n)=>e??t??n,fn=ze({name:"RouterView",inheritAttrs:!1,props:{name:{type:String,default:"default"},route:Object},compatConfig:{MODE:3},setup(e,{attrs:t,slots:n}){const r=I(pe),o=N(()=>e.route||r.value),l=I(Ie,0),d=N(()=>{let f=Q(l);const{matched:s}=o.value;let a;for(;(a=s[f])&&!a.components;)f++;return f}),m=N(()=>o.value.matched[d.value]);ae(Ie,N(()=>d.value+1)),ae(ge,m),ae(pe,o);const c=ft();return ot(()=>[c.value,m.value,e.name],([f,s,a],[u,g,R])=>{s&&(s.instances[a]=f,g&&g!==s&&f&&f===u&&(s.leaveGuards.size||(s.leaveGuards=g.leaveGuards),s.updateGuards.size||(s.updateGuards=g.updateGuards))),f&&s&&(!g||!q(s,g)||!u)&&(s.enterCallbacks[a]||[]).forEach(k=>k(f))},{flush:"post"}),()=>{const f=o.value,s=e.name,a=m.value,u=a&&a.components[s];if(!u)return Ge(n.default,{Component:u,route:f});const g=a.props[s],R=g?g===!0?f.params:typeof g=="function"?g(f):g:null,b=qe(u,S({},R,t,{onVnodeUnmounted:E=>{E.component.isUnmounted&&(a.instances[s]=null)},ref:c}));return Ge(n.default,{Component:b,route:f})||b}}});function Ge(e,t){if(!e)return null;const n=e(t);return n.length===1?n[0]:n}const hn=fn;function En(e){const t=Gt(e.routes,e),n=e.parseQuery||rn,r=e.stringifyQuery||He,o=e.history,l=W(),d=W(),m=W(),c=at(T);let f=T;z&&e.scrollBehavior&&"scrollRestoration"in history&&(history.scrollRestoration="manual");const s=le.bind(null,i=>""+i),a=le.bind(null,nn),u=le.bind(null,ee);function g(i,p){let h,v;return Ve(i)?(h=t.getRecordMatcher(i),v=p):v=i,t.addRoute(v,h)}function R(i){const p=t.getRecordMatcher(i);p&&t.removeRoute(p)}function k(){return t.getRoutes().map(i=>i.record)}function b(i){return!!t.getRecordMatcher(i)}function E(i,p){if(p=S({},p||c.value),typeof i=="string"){const y=ue(n,i,p.path),O=t.resolve({path:y.path},p),D=o.createHref(y.fullPath);return S(y,O,{params:u(O.params),hash:ee(y.hash),redirectedFrom:void 0,href:D})}let h;if("path"in i)h=S({},i,{path:ue(n,i.path,p.path).path});else{const y=S({},i.params);for(const O in y)y[O]==null&&delete y[O];h=S({},i,{params:a(y)}),p.params=a(p.params)}const v=t.resolve(h,p),C=i.hash||"";v.params=s(u(v.params));const A=mt(r,S({},i,{hash:Jt(C),path:v.path})),P=o.createHref(A);return S({fullPath:A,hash:C,query:r===He?on(i.query):i.query||{}},v,{redirectedFrom:void 0,href:P})}function w(i){return typeof i=="string"?ue(n,i,c.value.path):S({},i)}function _(i,p){if(f!==i)return K(8,{from:p,to:i})}function M(i){return U(i)}function j(i){return M(S(w(i),{replace:!0}))}function B(i){const p=i.matched[i.matched.length-1];if(p&&p.redirect){const{redirect:h}=p;let v=typeof h=="function"?h(i):h;return typeof v=="string"&&(v=v.includes("?")||v.includes("#")?v=w(v):{path:v},v.params={}),S({query:i.query,hash:i.hash,params:"path"in v?{}:i.params},v)}}function U(i,p){const h=f=E(i),v=c.value,C=i.state,A=i.force,P=i.replace===!0,y=B(h);if(y)return U(S(w(y),{state:typeof y=="object"?S({},C,y.state):C,force:A,replace:P}),p||h);const O=h;O.redirectedFrom=p;let D;return!A&&gt(r,v,h)&&(D=K(16,{to:O,from:v}),Se(v,v,!0,!1)),(D?Promise.resolve(D):Re(O,v)).catch(x=>H(x)?H(x,2)?x:se(x):oe(x,O,v)).then(x=>{if(x){if(H(x,2))return U(S({replace:P},w(x.to),{state:typeof x.to=="object"?S({},C,x.to.state):C,force:A}),p||O)}else x=Pe(O,v,!0,P,C);return Ee(O,v,x),x})}function et(i,p){const h=_(i,p);return h?Promise.reject(h):Promise.resolve()}function ye(i){const p=J.values().next().value;return p&&typeof p.runWithContext=="function"?p.runWithContext(i):i()}function Re(i,p){let h;const[v,C,A]=dn(i,p);h=fe(v.reverse(),"beforeRouteLeave",i,p);for(const y of v)y.leaveGuards.forEach(O=>{h.push($(O,i,p))});const P=et.bind(null,i,p);return h.push(P),G(h).then(()=>{h=[];for(const y of l.list())h.push($(y,i,p));return h.push(P),G(h)}).then(()=>{h=fe(C,"beforeRouteUpdate",i,p);for(const y of C)y.updateGuards.forEach(O=>{h.push($(O,i,p))});return h.push(P),G(h)}).then(()=>{h=[];for(const y of A)if(y.beforeEnter)if(L(y.beforeEnter))for(const O of y.beforeEnter)h.push($(O,i,p));else h.push($(y.beforeEnter,i,p));return h.push(P),G(h)}).then(()=>(i.matched.forEach(y=>y.enterCallbacks={}),h=fe(A,"beforeRouteEnter",i,p),h.push(P),G(h))).then(()=>{h=[];for(const y of d.list())h.push($(y,i,p));return h.push(P),G(h)}).catch(y=>H(y,8)?y:Promise.reject(y))}function Ee(i,p,h){m.list().forEach(v=>ye(()=>v(i,p,h)))}function Pe(i,p,h,v,C){const A=_(i,p);if(A)return A;const P=p===T,y=z?history.state:{};h&&(v||P?o.replace(i.fullPath,S({scroll:P&&y&&y.scroll},C)):o.push(i.fullPath,C)),c.value=i,Se(i,p,h,P),se()}let V;function tt(){V||(V=o.listen((i,p,h)=>{if(!Ce.listening)return;const v=E(i),C=B(v);if(C){U(S(C,{replace:!0}),v).catch(F);return}f=v;const A=c.value;z&&Ct(Ae(A.fullPath,h.delta),te()),Re(v,A).catch(P=>H(P,12)?P:H(P,2)?(U(P.to,v).then(y=>{H(y,20)&&!h.delta&&h.type===X.pop&&o.go(-1,!1)}).catch(F),Promise.reject()):(h.delta&&o.go(-h.delta,!1),oe(P,v,A))).then(P=>{P=P||Pe(v,A,!1),P&&(h.delta&&!H(P,8)?o.go(-h.delta,!1):h.type===X.pop&&H(P,20)&&o.go(-1,!1)),Ee(v,A,P)}).catch(F)}))}let re=W(),we=W(),Z;function oe(i,p,h){se(i);const v=we.list();return v.length?v.forEach(C=>C(i,p,h)):console.error(i),Promise.reject(i)}function nt(){return Z&&c.value!==T?Promise.resolve():new Promise((i,p)=>{re.add([i,p])})}function se(i){return Z||(Z=!i,tt(),re.list().forEach(([p,h])=>i?h(i):p()),re.reset()),i}function Se(i,p,h,v){const{scrollBehavior:C}=e;if(!z||!C)return Promise.resolve();const A=!h&&bt(Ae(i.fullPath,0))||(v||!h)&&history.state&&history.state.scroll||null;return rt().then(()=>C(i,p,A)).then(P=>P&&St(P)).catch(P=>oe(P,i,p))}const ie=i=>o.go(i);let ce;const J=new Set,Ce={currentRoute:c,listening:!0,addRoute:g,removeRoute:R,hasRoute:b,getRoutes:k,resolve:E,options:e,push:M,replace:j,go:ie,back:()=>ie(-1),forward:()=>ie(1),beforeEach:l.add,beforeResolve:d.add,afterEach:m.add,onError:we.add,isReady:nt,install(i){const p=this;i.component("RouterLink",an),i.component("RouterView",hn),i.config.globalProperties.$router=p,Object.defineProperty(i.config.globalProperties,"$route",{enumerable:!0,get:()=>Q(c)}),z&&!ce&&c.value===T&&(ce=!0,M(o.location).catch(C=>{}));const h={};for(const C in T)Object.defineProperty(h,C,{get:()=>c.value[C],enumerable:!0});i.provide(ne,p),i.provide(ve,lt(h)),i.provide(pe,c);const v=i.unmount;J.add(i),i.unmount=function(){J.delete(i),J.size<1&&(f=T,V&&V(),V=null,c.value=T,ce=!1,Z=!1),v()}}};function G(i){return i.reduce((p,h)=>p.then(()=>ye(h)),Promise.resolve())}return Ce}function dn(e,t){const n=[],r=[],o=[],l=Math.max(t.matched.length,e.matched.length);for(let d=0;d<l;d++){const m=t.matched[d];m&&(e.matched.find(f=>q(f,m))?r.push(m):n.push(m));const c=e.matched[d];c&&(t.matched.find(f=>q(f,c))||o.push(c))}return[n,r,o]}function Pn(){return I(ne)}function wn(){return I(ve)}export{vn as a,wn as b,En as c,Rn as d,yn as o,Pn as u};