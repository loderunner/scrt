import{_ as e,d as t}from"./app.ce2c9135.js";const a={},s=t(`<h1 id="get" tabindex="-1"><a class="header-anchor" href="#get" aria-hidden="true">#</a> get</h1><div class="language-text ext-text"><pre class="language-text"><code>scrt get [flags] key
</code></pre></div><p>Retrieve the value associated to the key in the store, if it exists. Returns an error if no value is associated to the key.</p><h3 id="example" tabindex="-1"><a class="header-anchor" href="#example" aria-hidden="true">#</a> Example</h3><p>Retrieve the value associated to the key <code>greeting</code> in the store, using implicit store configuration (configuration file or environment variables).</p><div class="language-bash ext-sh"><pre class="language-bash"><code>scrt get greeting

<span class="token comment"># Output: Hello World</span>
</code></pre></div>`,6);function r(i,n){return s}var c=e(a,[["render",r],["__file","get.html.vue"]]);export{c as default};
