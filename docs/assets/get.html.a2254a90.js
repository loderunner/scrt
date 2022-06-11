import{_ as e,o as t,c as a,d as s}from"./app.daf49d13.js";const r={},i=s(`<h1 id="get" tabindex="-1"><a class="header-anchor" href="#get" aria-hidden="true">#</a> get</h1><div class="language-text ext-text"><pre class="language-text"><code>scrt get [flags] key
</code></pre></div><p>Retrieve the value associated to the key in the store, if it exists. Returns an error if no value is associated to the key.</p><h3 id="example" tabindex="-1"><a class="header-anchor" href="#example" aria-hidden="true">#</a> Example</h3><p>Retrieve the value associated to the key <code>greeting</code> in the store, using implicit store configuration (configuration file or environment variables).</p><div class="language-bash ext-sh"><pre class="language-bash"><code>scrt get greeting

<span class="token comment"># Output: Hello World</span>
</code></pre></div>`,6),n=[i];function o(c,l){return t(),a("div",null,n)}var h=e(r,[["render",o],["__file","get.html.vue"]]);export{h as default};
