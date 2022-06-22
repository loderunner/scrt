import{_ as n,r as s,o as a,c as l,b as e,d as r,w as i,a as c,e as t}from"./app.58b453f6.js";const d={},h=c(`<h1 id="global" tabindex="-1"><a class="header-anchor" href="#global" aria-hidden="true">#</a> Global</h1><p>Use <code>scrt --help</code> to output a full help message.</p><div class="language-text ext-text"><pre class="language-text"><code>A secret manager for the command-line

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
</code></pre></div><h3 id="global-options" tabindex="-1"><a class="header-anchor" href="#global-options" aria-hidden="true">#</a> Global options</h3>`,4),u=e("strong",null,[e("code",null,"-c")],-1),p=t(", "),g=e("strong",null,[e("code",null,"--config"),t(":")],-1),_=t(" Path to a YAML "),f=t("Configuration file"),m=e("p",null,[e("strong",null,[e("code",null,"--storage"),t(":")]),t(" storage type, see Reference for details.")],-1),v=e("p",null,[e("strong",null,[e("code",null,"-p")]),t(", "),e("strong",null,[e("code",null,"--password"),t(":")]),t(" password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.")],-1);function b(x,y){const o=s("RouterLink");return a(),l("div",null,[h,e("p",null,[u,p,g,_,r(o,{to:"/guide/configuration.html"},{default:i(()=>[f]),_:1})]),m,v])}var w=n(d,[["render",b],["__file","global.html.vue"]]);export{w as default};
