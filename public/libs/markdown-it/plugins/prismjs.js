// import {components} from "/public/libs/prismjs/components.js"
const NEW_LINE_EXP = /\n(?!$)/g;
const DEFAULTS = {
  plugins: [],
  init: () => {},
  defaultLanguageForUnknown: undefined,
  defaultLanguageForUnspecified: undefined,
  defaultLanguage: undefined,
};

function loadJS(file) {
  let js = $.ajax({ type: "GET", url: file, async: false }).responseText; //No need to append
}

function langAliasPreProcess(lang) {
  if (lang === "py") { return "python"; }
  if (lang === "js") { return "javascript"; }
  if (lang === "sh") { return "bash"; }
  if (lang === "html") { return "markup"; }
  if ((lang === "mma") || (lang === "mathematica")) { return "wolfram"; }
  return lang;
}

function loadPrismPlugin(name) {
  if (
    name === "toolbar" ||
    name === "line-numbers" ||
    name === "autoloader" ||
    name === "copy-to-clipboard" ||
    name === "code-header"
  ) {
    return;
  }
  try {
    require(`prismjs/plugins/${name}/prism-${name}-min`);
    // require(`prismjs/plugins/${name}`);
  } catch (e) {
    throw new Error(
      `Cannot load Prism plugin "${name}". Please check the spelling. ${e}`
    );
  }
}

// Checks whether an option represents a valid Prism language
function checkLanguageOption(options, optionName) {
  const language = options[optionName];
  if (language !== undefined && loadPrismLang(language) === undefined) {
    throw new Error(
      `Bad option ${optionName}: There is no Prism language '${language}'.`
    );
  }
}

(function (global, factory) {
  typeof exports === "object" && typeof module !== "undefined"
    ? (module.exports = factory())
    : typeof define === "function" && define.amd
    ? define(["prismjs"], factory)
    : ((global =
        typeof globalThis !== "undefined" ? globalThis : global || self),
      (global.markdownitPrismjs = factory()));
})(this, function () {
  "use strict";

  function loadPrismLang(lang) {
    if (!lang) return undefined;
    if (typeof components === "undefined")
      throw new Error("Prism.components has not loaded yet.");
    if (!(lang in components.languages)) return undefined;
    let langObject = Prism.languages[lang];
    if (langObject === undefined) {
      // require(["prismjs/components/prism-" + lang + '.min'], () => {
      //     requirejs_done = true;
      //     langObject = Prism.languages[lang];
      // });
      // loadJS("//cdn.jsdelivr.net/npm/prismjs@1/components/prism-"+lang+'.min.js');
      loadJS("/libs/prismjs/components/prism-" + lang + ".min.js");
      langObject = Prism.languages[lang];
    }
    return langObject;
  }

  // Select the language to use for highlighting, based on the provided options and the specified language.
  function selectLanguage(options, lang) {
    let langToUse = lang;
    if (
      langToUse === "" &&
      options.defaultLanguageForUnspecified !== undefined
    ) {
      langToUse = options.defaultLanguageForUnspecified;
      return [langToUse, undefined];
    }
    if (langToUse === 'cpp') {loadPrismLang('c');} // cpp dependency c
    if (langToUse === 'php') {loadPrismLang('markup-templating');} // php dependency markup-templating
    let prismLang = loadPrismLang(langToUse);
    // if (prismLang === undefined && options.defaultLanguageForUnknown !== undefined ) {
    //   langToUse = options.defaultLanguageForUnknown;
    //   prismLang = loadPrismLang(langToUse);
    // }
    return [langToUse, prismLang];
  }

  // https://stackoverflow.com/questions/59508413/static-html-generation-with-prismjs-how-to-enable-line-numbers
  function highlight(md, opts, text, lang) {
    let langToUse = langAliasPreProcess(lang.toLowerCase());
    let prismLang;
    [langToUse, prismLang] = selectLanguage(opts, langToUse);
    let code = prismLang
      ? Prism.highlight(text, prismLang, langToUse)
      : md.utils.escapeHtml(text);
    let code_list = code.split("\n");
    code_list.pop(); // 多了一个空行，移除

    if (opts.plugins.indexOf("line-numbers") >= 0 && code_list.length > 1) {
      code = code_list
        .map((line, num) => {
          // return `<li class="L${num%10} prettyprint linenums"><code class="language-${langToUse}">${line}</code></li>` // 和行号的颜色进行区分，只能这么包一下
          return `<li class="L${
            num % 10
          } prettyprint linenums"><span>${line}</span></li>`;
        })
        .join("");
      code = `<ol class="linenums">${code}</ol>`;
    } else if (code_list.length == 1) {
      code = `<code class="language-${langToUse}">${code}</code>`;
    }

    if (opts.plugins.indexOf("toolbar") >= 0) {
      code = `<div class="code-toolbar">${code}<div class="toolbar"></div></div>`;
    }

    if (code_list.length == 1) {
      code = `<pre>${code}</pre>`;
    } else if (opts.withClassLangeuePrefix) {
      code = `<pre class="prettyprint linenums language-${langToUse}">${code}</pre>`;
    } else {
      code = `<pre class="prettyprint linenums">${code}</pre>`;
    }

    if (opts.plugins.indexOf("code-header") >= 0) {
      code = `<div class="highlight"><div class="code-header"><span text-data=${lang}><i class="fa-fw fas fa-code small"></i></span><button aria-label="copy" data-original-title="" title=""><i class="far fa-clipboard"></i></button></div>${code}</div>`;
    }

    return code;
  }
  // Patch the `<pre>` and `<code>` tags produced by the `existingRule` for fenced code blocks.
  // function applyCodeAttributes(markdownit, options, existingRule) {
  //   function languageClass(md, lang) {
  //     return md.options.langPrefix + md.utils.escapeHtml(lang);
  //   }

  //     return (tokens, idx, renderOptions, env, self) => {
  //       const fenceToken = tokens[idx];
  //       const info = fenceToken.info
  //         ? markdownit.utils.unescapeAll(fenceToken.info).trim()
  //         : "";
  //       const lang = info.split(/(\s+)/g)[0];
  //       const [langToUse] = selectLanguage(options, lang);
  //       if (!langToUse) {
  //         return existingRule(tokens, idx, renderOptions, env, self);
  //       } else {
  //         fenceToken.info = langToUse;
  //         const existingResult = existingRule(
  //           tokens,
  //           idx,
  //           renderOptions,
  //           env,
  //           self
  //         );
  //         const langClass = languageClass(markdownit, langToUse);
  //         console.log(existingResult);
  //         return existingResult.replace(
  //           /<((?:pre|code)[^>]*?)(?:\s+class="([^"]*)"([^>]*))?>/g,
  //           (match, tagStart, existingClasses, tagEnd) =>
  //             existingClasses?.includes(langClass)
  //               ? match
  //               : `<${tagStart} class="${existingClasses ? `${existingClasses} ` : ""}${langClass}"${tagEnd || ""}>`
  //         );
  //       }
  //     };
  //   }

  function markdownItPrism(md, useroptions) {
    const opts = Object.assign({}, DEFAULTS, useroptions);

    checkLanguageOption(opts, "defaultLanguage");
    checkLanguageOption(opts, "defaultLanguageForUnknown");
    checkLanguageOption(opts, "defaultLanguageForUnspecified");
    opts.defaultLanguageForUnknown =
      opts.defaultLanguageForUnknown || opts.defaultLanguage;
    opts.defaultLanguageForUnspecified =
      opts.defaultLanguageForUnspecified || opts.defaultLanguage;

    opts.plugins.forEach(loadPrismPlugin);

    opts.init(Prism);

    // register ourselves as highlighter
    // md.options.highlight = (text, lang) => highlight(md, opts, text, lang);
    md.renderer.rules["fence"] = function (tokens, idx, options, env, slf) {
      opts.withClassLangeuePrefix = true;
      options = Object.assign(options, DEFAULTS, opts);

      let info = tokens[idx].info
          ? md.utils.unescapeAll(tokens[idx].info).trim()
          : "",
        langName = "",
        langAttrs = "";
      if (info) {
        let arr = info.split(/(\s+)/g);
        langName = arr[0];
        langAttrs = arr.slice(2).join("");
      }

      return highlight(md, options, tokens[idx].content, langName);
    };

    // [Highlight function not called for all code blocks](https://github.com/markdown-it/markdown-it/issues/130)
    md.renderer.rules["code_block"] = function (
      tokens,
      idx,
      options,
      env,
      slf
    ) {
      opts.withClassLangeuePrefix = false;
      options = Object.assign(options, DEFAULTS, opts); // 函数传递进来的 options = md.options
      let lang = "";
      let arr = tokens[idx].content.match(/^\`\`\`(\w+)\n/);
      if (arr) {
        if (arr.length >= 2) {
          lang = arr[1];
        }
      }
      return highlight(md, options, tokens[idx].content, lang);
    };
    // md.renderer.rules.fence = applyCodeAttributes(md, options, md.renderer.rules.fence || (() => ""));
  }

  return markdownItPrism;
});
