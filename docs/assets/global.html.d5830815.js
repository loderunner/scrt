import{_ as a,r as l,o as r,c as i,a as e,b as s,w as n,d as c,e as t}from"./app.daf49d13.js";const d={},h=c(`<h1 id="global" tabindex="-1"><a class="header-anchor" href="#global" aria-hidden="true">#</a> Global</h1><p>Use <code>scrt --help</code> to output a full help message.</p><div class="language-text ext-text"><pre class="language-text"><code>A secret manager for the command-line

Usage:
  scrt [command]

Available Commands:
  init        Initialize a new store
  set         Associate a key to a value in a store
  get         Retrieve the value associated to key from a store
  list        List all the keys in a store
  unset       Remove the value associated to key in a store
  storage     List storage types and options
  help        Help about any command
  completion  Generate the autocompletion script for the specified shell

Flags:
  -c, --config string     configuration file
  -h, --help              help for scrt
  -p, --password string   master password to unlock the store
      --storage string    storage type
  -v, --verbose           verbose output
      --version           version for scrt
</code></pre></div><h3 id="global-options" tabindex="-1"><a class="header-anchor" href="#global-options" aria-hidden="true">#</a> Global options</h3>`,4),u=e("strong",null,[e("code",null,"-c")],-1),p=t(", "),g=e("strong",null,[e("code",null,"--config"),t(":")],-1),_=t(" Path to a YAML "),f=t("Configuration file"),m=e("strong",null,[e("code",null,"--storage"),t(":")],-1),v=t(" storage type, see "),b=t("Storage types"),y=t(" for details."),x=e("p",null,[e("strong",null,[e("code",null,"-p")]),t(", "),e("strong",null,[e("code",null,"--password"),t(":")]),t(" password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.")],-1);function k(w,L){const o=l("RouterLink");return r(),i("div",null,[h,e("p",null,[u,p,g,_,s(o,{to:"/guide/configuration.html"},{default:n(()=>[f]),_:1})]),e("p",null,[m,v,s(o,{to:"/reference/storage.html"},{default:n(()=>[b]),_:1}),y]),x])}var C=a(d,[["render",k],["__file","global.html.vue"]]);export{C as default};
