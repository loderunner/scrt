import{_ as s,r as i,o as a,c,a as e,b as n,F as r,e as t,d as l}from"./app.3fd7ff3b.js";const d={},h=e("h1",{id:"git",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#git","aria-hidden":"true"},"#"),t(" Git")],-1),p=e("p",null,[t("Use the "),e("code",null,"git"),t(" storage type to create and access a store in a git repository. "),e("code",null,"scrt"),t(" will clone the repository in memory, checkout the given branch (or the default branch if no branch is given), read the store in the file at the given path, and will commit and push any modifications to the remote.")],-1),_=e("h3",{id:"options",tabindex:"-1"},[e("a",{class:"header-anchor",href:"#options","aria-hidden":"true"},"#"),t(" Options")],-1),u=e("strong",null,[e("code",null,"--git-url")],-1),g=t(" (required): a git-compatible repository URL. Most git-compatible URLs and protocols can be used. See "),m={href:"https://git-scm.com/docs/git-clone#_git_urls",target:"_blank",rel:"noopener noreferrer"},f=e("code",null,"git clone",-1),b=t(" documentation"),k=t(" to learn more."),v=e("p",null,[e("strong",null,[e("code",null,"--git-path")]),t(" (required): the path to the store file inside the the git repository, relative to the repository root. A repository can contain multiple scrt stores, at different paths.")],-1),w=e("p",null,[e("strong",null,[e("code",null,"--git-branch"),t(":")]),t(" the name of the branch to checkout after cloning (or initializing). If no branch is given, the default branch from the remote will be used, or "),e("code",null,"main"),t(" if a new repository is initialized.")],-1),x=e("strong",null,[e("code",null,"--git-checkout"),t(":")],-1),y=t(" a git revision to checkout. If specified, the revision will be checked out in a "),E={href:"https://git-scm.com/docs/git-checkout#_detached_head",target:"_blank",rel:"noopener noreferrer"},I=t('"detached HEAD"'),L=t(" and pushing will not work; making updates ("),N=e("code",null,"init",-1),V=t(", "),z=e("code",null,"set",-1),A=t(" or "),B=e("code",null,"unset",-1),U=t(") will be impossible."),q=l(`<p><strong><code>--git-message</code>:</strong> the message of the git commit. A default message will be used if this is not set.</p><h3 id="example" tabindex="-1"><a class="header-anchor" href="#example" aria-hidden="true">#</a> Example</h3><div class="language-bash ext-sh"><pre class="language-bash"><code>scrt init --storage<span class="token operator">=</span>git <span class="token punctuation">\\</span>
          --password<span class="token operator">=</span>p4ssw0rd <span class="token punctuation">\\</span>
          --git-url<span class="token operator">=</span>git@github.com:githubuser/secrets.git <span class="token punctuation">\\</span>
          --git-path<span class="token operator">=</span>store.scrt
</code></pre></div><div class="custom-container tip"><p class="custom-container-title">TIP</p><p><code>scrt</code> will initialize a new repo if none can be cloned.</p></div>`,4);function F(R,S){const o=i("ExternalLinkIcon");return a(),c(r,null,[h,p,_,e("p",null,[u,g,e("a",m,[f,b,n(o)]),k]),v,w,e("p",null,[x,y,e("a",E,[I,n(o)]),L,N,V,z,A,B,U]),q],64)}var C=s(d,[["render",F],["__file","git.html.vue"]]);export{C as default};
