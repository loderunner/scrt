import{_ as s,r,o as l,c as i,a as e,b as n,w as a,F as c,d,e as t}from"./app.3fd7ff3b.js";const h={},u=d(`<h1 id="global" tabindex="-1"><a class="header-anchor" href="#global" aria-hidden="true">#</a> Global</h1><p>Use <code>scrt --help</code> to output a full help message.</p><div class="language-text ext-text"><pre class="language-text"><code>A secret manager for the command-line

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

Flags:
  -c, --config string     configuration file
  -h, --help              help for scrt
  -p, --password string   master password to unlock the store
      --storage string    storage type
  -v, --version           version for scrt

Use &quot;scrt [command] --help&quot; for more information about a command.
</code></pre></div><h3 id="global-options" tabindex="-1"><a class="header-anchor" href="#global-options" aria-hidden="true">#</a> Global options</h3>`,4),g=e("strong",null,[e("code",null,"-c")],-1),_=t(", "),p=e("strong",null,[e("code",null,"--config"),t(":")],-1),m=t(" Path to a YAML "),f=t("Configuration file"),b=e("strong",null,[e("code",null,"--storage"),t(":")],-1),v=t(" storage type, see "),y=t("Storage types"),x=t(" for details."),k=e("p",null,[e("strong",null,[e("code",null,"-p")]),t(", "),e("strong",null,[e("code",null,"--password"),t(":")]),t(" password to the store. The argument will be used to derive a key, to decrypt and encrypt the data in the store.")],-1);function w(L,A){const o=r("RouterLink");return l(),i(c,null,[u,e("p",null,[g,_,p,m,n(o,{to:"/guide/configuration.html"},{default:a(()=>[f]),_:1})]),e("p",null,[b,v,n(o,{to:"/reference/storage.html"},{default:a(()=>[y]),_:1}),x]),k],64)}var N=s(h,[["render",w],["__file","global.html.vue"]]);export{N as default};
