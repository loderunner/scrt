import{_ as e,o as n,c as s,a}from"./app.58b453f6.js";const t={},o=a(`<h3 id="secure-secret-management-from-the-command-line" tabindex="-1"><a class="header-anchor" href="#secure-secret-management-from-the-command-line" aria-hidden="true">#</a> Secure secret management from the command line</h3><div class="language-bash ext-sh"><pre class="language-bash"><code><span class="token comment"># Set store password in environment</span>
$ <span class="token builtin class-name">export</span> <span class="token assign-left variable">SCRT_PASSWORD</span><span class="token operator">=</span>*******

<span class="token comment"># Add a secret to the store</span>
$ scrt <span class="token builtin class-name">set</span> greeting <span class="token string">&#39;Good news, everyone!&#39;</span>

<span class="token comment"># Retrieve the secret from the store</span>
$ scrt get greeting
Good news, everyone<span class="token operator">!</span>
</code></pre></div>`,2),r=[o];function c(l,i){return n(),s("div",null,r)}var p=e(t,[["render",c],["__file","index.html.vue"]]);export{p as default};
