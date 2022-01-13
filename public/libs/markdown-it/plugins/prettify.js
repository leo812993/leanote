(function (global, factory) {
    typeof exports === "object" && typeof module !== "undefined"
      ? (module.exports = factory())
      : typeof define === "function" && define.amd
        ? define(['google-code-prettify'], factory)
        : ((global =
          typeof globalThis !== "undefined" ? globalThis : global || self),
          (global.markdownitPrettify = factory()));
  })(this, function () {
    "use strict";

    function markdownItPrettify(md, useroptions) {
        md.options.highlight = (text, lang) => prettyPrint();
    }

    return markdownItPrettify;
});