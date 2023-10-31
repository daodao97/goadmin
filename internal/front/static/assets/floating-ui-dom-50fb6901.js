import{r as Q,s as A}from"./floating-ui-core-a8b4680c.js";function p(t){var e;return(t==null||(e=t.ownerDocument)==null?void 0:e.defaultView)||window}function w(t){return p(t).getComputedStyle(t)}function X(t){return t instanceof p(t).Node}function b(t){return X(t)?(t.nodeName||"").toLowerCase():"#document"}function m(t){return t instanceof HTMLElement||t instanceof p(t).HTMLElement}function P(t){return typeof ShadowRoot<"u"&&(t instanceof p(t).ShadowRoot||t instanceof ShadowRoot)}function E(t){const{overflow:e,overflowX:n,overflowY:o,display:i}=w(t);return/auto|scroll|overlay|hidden|clip/.test(e+o+n)&&!["inline","contents"].includes(i)}function U(t){return["table","td","th"].includes(b(t))}function M(t){const e=O(),n=w(t);return n.transform!=="none"||n.perspective!=="none"||!!n.containerType&&n.containerType!=="normal"||!e&&!!n.backdropFilter&&n.backdropFilter!=="none"||!e&&!!n.filter&&n.filter!=="none"||["transform","perspective","filter"].some(o=>(n.willChange||"").includes(o))||["paint","layout","strict","content"].some(o=>(n.contain||"").includes(o))}function O(){return!(typeof CSS>"u"||!CSS.supports)&&CSS.supports("-webkit-backdrop-filter","none")}function F(t){return["html","body","#document"].includes(b(t))}const B=Math.min,S=Math.max,W=Math.round,T=t=>({x:t,y:t});function Y(t){const e=w(t);let n=parseFloat(e.width)||0,o=parseFloat(e.height)||0;const i=m(t),r=i?t.offsetWidth:n,s=i?t.offsetHeight:o,l=W(n)!==r||W(o)!==s;return l&&(n=r,o=s),{width:n,height:o,$:l}}function x(t){return t instanceof Element||t instanceof p(t).Element}function _(t){return x(t)?t:t.contextElement}function R(t){const e=_(t);if(!m(e))return T(1);const n=e.getBoundingClientRect(),{width:o,height:i,$:r}=Y(e);let s=(r?W(n.width):n.width)/o,l=(r?W(n.height):n.height)/i;return s&&Number.isFinite(s)||(s=1),l&&Number.isFinite(l)||(l=1),{x:s,y:l}}const Z=T(0);function z(t){const e=p(t);return O()&&e.visualViewport?{x:e.visualViewport.offsetLeft,y:e.visualViewport.offsetTop}:Z}function D(t,e,n,o){e===void 0&&(e=!1),n===void 0&&(n=!1);const i=t.getBoundingClientRect(),r=_(t);let s=T(1);e&&(o?x(o)&&(s=R(o)):s=R(t));const l=function(a,y,u){return y===void 0&&(y=!1),!(!u||y&&u!==p(a))&&y}(r,n,o)?z(r):T(0);let c=(i.left+l.x)/s.x,f=(i.top+l.y)/s.y,d=i.width/s.x,g=i.height/s.y;if(r){const a=p(r),y=o&&x(o)?p(o):o;let u=a.frameElement;for(;u&&o&&y!==a;){const h=R(u),L=u.getBoundingClientRect(),H=getComputedStyle(u),J=L.left+(u.clientLeft+parseFloat(H.paddingLeft))*h.x,K=L.top+(u.clientTop+parseFloat(H.paddingTop))*h.y;c*=h.x,f*=h.y,d*=h.x,g*=h.y,c+=J,f+=K,u=p(u).frameElement}}return A({width:d,height:g,x:c,y:f})}function V(t){return x(t)?{scrollLeft:t.scrollLeft,scrollTop:t.scrollTop}:{scrollLeft:t.pageXOffset,scrollTop:t.pageYOffset}}function v(t){var e;return(e=(X(t)?t.ownerDocument:t.document)||window.document)==null?void 0:e.documentElement}function j(t){return D(v(t)).left+V(t).scrollLeft}function C(t){if(b(t)==="html")return t;const e=t.assignedSlot||t.parentNode||P(t)&&t.host||v(t);return P(e)?e.host:e}function q(t){const e=C(t);return F(e)?t.ownerDocument?t.ownerDocument.body:t.body:m(e)&&E(e)?e:q(e)}function G(t,e){var n;e===void 0&&(e=[]);const o=q(t),i=o===((n=t.ownerDocument)==null?void 0:n.body),r=p(o);return i?e.concat(r,r.visualViewport||[],E(o)?o:[]):e.concat(o,G(o))}function N(t,e,n){let o;if(e==="viewport")o=function(i,r){const s=p(i),l=v(i),c=s.visualViewport;let f=l.clientWidth,d=l.clientHeight,g=0,a=0;if(c){f=c.width,d=c.height;const y=O();(!y||y&&r==="fixed")&&(g=c.offsetLeft,a=c.offsetTop)}return{width:f,height:d,x:g,y:a}}(t,n);else if(e==="document")o=function(i){const r=v(i),s=V(i),l=i.ownerDocument.body,c=S(r.scrollWidth,r.clientWidth,l.scrollWidth,l.clientWidth),f=S(r.scrollHeight,r.clientHeight,l.scrollHeight,l.clientHeight);let d=-s.scrollLeft+j(i);const g=-s.scrollTop;return w(l).direction==="rtl"&&(d+=S(r.clientWidth,l.clientWidth)-c),{width:c,height:f,x:d,y:g}}(v(t));else if(x(e))o=function(i,r){const s=D(i,!0,r==="fixed"),l=s.top+i.clientTop,c=s.left+i.clientLeft,f=m(i)?R(i):T(1);return{width:i.clientWidth*f.x,height:i.clientHeight*f.y,x:c*f.x,y:l*f.y}}(e,n);else{const i=z(t);o={...e,x:e.x-i.x,y:e.y-i.y}}return A(o)}function I(t,e){const n=C(t);return!(n===e||!x(n)||F(n))&&(w(n).position==="fixed"||I(n,e))}function tt(t,e,n){const o=m(e),i=v(e),r=n==="fixed",s=D(t,!0,r,e);let l={scrollLeft:0,scrollTop:0};const c=T(0);if(o||!o&&!r)if((b(e)!=="body"||E(i))&&(l=V(e)),m(e)){const f=D(e,!0,r,e);c.x=f.x+e.clientLeft,c.y=f.y+e.clientTop}else i&&(c.x=j(i));return{x:s.left+l.scrollLeft-c.x,y:s.top+l.scrollTop-c.y,width:s.width,height:s.height}}function k(t,e){return m(t)&&w(t).position!=="fixed"?e?e(t):t.offsetParent:null}function $(t,e){const n=p(t);if(!m(t))return n;let o=k(t,e);for(;o&&U(o)&&w(o).position==="static";)o=k(o,e);return o&&(b(o)==="html"||b(o)==="body"&&w(o).position==="static"&&!M(o))?n:o||function(i){let r=C(i);for(;m(r)&&!F(r);){if(M(r))return r;r=C(r)}return null}(t)||n}const et={convertOffsetParentRelativeRectToViewportRelativeRect:function(t){let{rect:e,offsetParent:n,strategy:o}=t;const i=m(n),r=v(n);if(n===r)return e;let s={scrollLeft:0,scrollTop:0},l=T(1);const c=T(0);if((i||!i&&o!=="fixed")&&((b(n)!=="body"||E(r))&&(s=V(n)),m(n))){const f=D(n);l=R(n),c.x=f.x+n.clientLeft,c.y=f.y+n.clientTop}return{width:e.width*l.x,height:e.height*l.y,x:e.x*l.x-s.scrollLeft*l.x+c.x,y:e.y*l.y-s.scrollTop*l.y+c.y}},getDocumentElement:v,getClippingRect:function(t){let{element:e,boundary:n,rootBoundary:o,strategy:i}=t;const r=[...n==="clippingAncestors"?function(c,f){const d=f.get(c);if(d)return d;let g=G(c).filter(h=>x(h)&&b(h)!=="body"),a=null;const y=w(c).position==="fixed";let u=y?C(c):c;for(;x(u)&&!F(u);){const h=w(u),L=M(u);L||h.position!=="fixed"||(a=null),(y?!L&&!a:!L&&h.position==="static"&&a&&["absolute","fixed"].includes(a.position)||E(u)&&!L&&I(c,u))?g=g.filter(H=>H!==u):a=h,u=C(u)}return f.set(c,g),g}(e,this._c):[].concat(n),o],s=r[0],l=r.reduce((c,f)=>{const d=N(e,f,i);return c.top=S(d.top,c.top),c.right=B(d.right,c.right),c.bottom=B(d.bottom,c.bottom),c.left=S(d.left,c.left),c},N(e,s,i));return{width:l.right-l.left,height:l.bottom-l.top,x:l.left,y:l.top}},getOffsetParent:$,getElementRects:async function(t){let{reference:e,floating:n,strategy:o}=t;const i=this.getOffsetParent||$,r=this.getDimensions;return{reference:tt(e,await i(n),o),floating:{x:0,y:0,...await r(n)}}},getClientRects:function(t){return Array.from(t.getClientRects())},getDimensions:function(t){return Y(t)},getScale:R,isElement:x,isRTL:function(t){return getComputedStyle(t).direction==="rtl"}},ot=(t,e,n)=>{const o=new Map,i={platform:et,...n},r={...i.platform,_c:o};return Q(t,e,{...i,platform:r})};export{ot as B};
