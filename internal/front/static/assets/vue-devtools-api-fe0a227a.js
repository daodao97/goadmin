function d(){return g().__VUE_DEVTOOLS_GLOBAL_HOOK__}function g(){return typeof navigator<"u"&&typeof window<"u"?window:typeof global<"u"?global:{}}const _=typeof Proxy=="function",h="devtools-plugin:setup",O="plugin:settings:set";let u,f;function S(){var o;return u!==void 0||(typeof window<"u"&&window.performance?(u=!0,f=window.performance):typeof global<"u"&&(!((o=global.perf_hooks)===null||o===void 0)&&o.performance)?(u=!0,f=global.perf_hooks.performance):u=!1),u}function w(){return S()?f.now():Date.now()}class y{constructor(n,s){this.target=null,this.targetQueue=[],this.onQueue=[],this.plugin=n,this.hook=s;const i={};if(n.settings)for(const e in n.settings){const t=n.settings[e];i[e]=t.defaultValue}const r=`__vue-devtools-plugin-settings__${n.id}`;let a=Object.assign({},i);try{const e=localStorage.getItem(r),t=JSON.parse(e);Object.assign(a,t)}catch{}this.fallbacks={getSettings(){return a},setSettings(e){try{localStorage.setItem(r,JSON.stringify(e))}catch{}a=e},now(){return w()}},s&&s.on(O,(e,t)=>{e===this.plugin.id&&this.fallbacks.setSettings(t)}),this.proxiedOn=new Proxy({},{get:(e,t)=>this.target?this.target.on[t]:(...l)=>{this.onQueue.push({method:t,args:l})}}),this.proxiedTarget=new Proxy({},{get:(e,t)=>this.target?this.target[t]:t==="on"?this.proxiedOn:Object.keys(this.fallbacks).includes(t)?(...l)=>(this.targetQueue.push({method:t,args:l,resolve:()=>{}}),this.fallbacks[t](...l)):(...l)=>new Promise(c=>{this.targetQueue.push({method:t,args:l,resolve:c})})})}async setRealTarget(n){this.target=n;for(const s of this.onQueue)this.target.on[s.method](...s.args);for(const s of this.targetQueue)s.resolve(await this.target[s.method](...s.args))}}function p(o,n){const s=o,i=g(),r=d(),a=_&&s.enableEarlyProxy;if(r&&(i.__VUE_DEVTOOLS_PLUGIN_API_AVAILABLE__||!a))r.emit(h,o,n);else{const e=a?new y(s,r):null;(i.__VUE_DEVTOOLS_PLUGINS__=i.__VUE_DEVTOOLS_PLUGINS__||[]).push({pluginDescriptor:s,setupFn:n,proxy:e}),e&&n(e.proxiedTarget)}}export{p as s};
