import{_ as a,r as l,o as r,c as i,a as e,b as s,w as n,F as c,d,e as t}from"./app.ce2c9135.js";const h={},u=d(`<h1 id="global" tabindex="-1"><a class="header-anchor" href="#global" aria-hidden="true">#</a> Global</h1><p>Use <code>scrt --help</code> to output a full help message.</p><div class="language-text ext-text"><pre class="language-text"><code>A secret manager for the command-line

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
</code></pre></div><h3 id="global-options" tabindex="-1"><a class="header-anchor" href="#global-options" aria-hidden="true">#</a> Global options</h3>`,4),p=e("strong",null,[e("code",null,"-c")],-1),g=t(", "),_=e("strong",null,[e("code",null,"--config"),t(":")],-1),f=t(" Path to a YAML "),m=t("Configuration file"),v=e("strong",null,[e("code",null,"--storage"),t(":")],-1),b=t(" storage type, see "),y=t("Storage types"),x=t(" for details."),k=e("p",null,[e("strong",null,[e("code",null,"-p")]),t(", "),e("strong",null,[e("code",null,"--password"),t(":")]),t(" password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.")],-1);function w(L,A){const o=l("RouterLink");return r(),i(c,null,[u,e("p",null,[p,g,_,f,s(o,{to:"/guide/configuration.html"},{default:n(()=>[m]),_:1})]),e("p",null,[v,b,s(o,{to:"/reference/storage.html"},{default:n(()=>[y]),_:1}),x]),k],64)}var N=a(h,[["render",w],["__file","global.html.vue"]]);export{N as default};
