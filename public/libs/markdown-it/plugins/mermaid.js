function randomString(e) {
  // e 长度
  e = e || 32;
  var t = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz",
    a = t.length,
    n = "";
  for (i = 0; i < e; i++) n += t.charAt(Math.floor(Math.random() * a));
  return n;
}

// 判断 fence 是否是 mermaid 或者 chart
function mermaid_rule(state, startLine, endLine, silent) {
  var marker,
    len,
    params,
    nextLine,
    mem,
    token,
    markup,
    haveEndMarker = false,
    pos = state.bMarks[startLine] + state.tShift[startLine],
    max = state.eMarks[startLine];

  // if it's indented more than 3 spaces, it should be a code block
  if (state.sCount[startLine] - state.blkIndent >= 4) {
    return false;
  }
  if (pos + 3 > max) {
    return false;
  }

  marker = state.src.charCodeAt(pos);

  if (marker !== 0x7e /* ~ */ && marker !== 0x60 /* ` */) {
    return false;
  }

  // scan marker length
  mem = pos;
  pos = state.skipChars(pos, marker);

  len = pos - mem;

  if (len < 3) {
    return false;
  }

  markup = state.src.slice(mem, pos);
  params = state.src.slice(pos, max);
  if (params != "mermaid" && params != "chart") {
    return false;
  }

  if (marker === 0x60 /* ` */) {
    if (params.indexOf(String.fromCharCode(marker)) >= 0) {
      return false;
    }
  }

  // Since start is found, we can report success here in validation mode
  if (silent && (params === "mermaid" || params === "chart")) {
    return true;
  }

  // search end of block
  nextLine = startLine;

  for (;;) {
    nextLine++;
    if (nextLine >= endLine) {
      // unclosed block should be autoclosed by end of document.
      // also block seems to be autoclosed by end of parent
      break;
    }

    pos = mem = state.bMarks[nextLine] + state.tShift[nextLine];
    max = state.eMarks[nextLine];

    if (pos < max && state.sCount[nextLine] < state.blkIndent) {
      // non-empty line with negative indent should stop the list:
      // - ```
      //  test
      break;
    }

    if (state.src.charCodeAt(pos) !== marker) {
      continue;
    }

    if (state.sCount[nextLine] - state.blkIndent >= 4) {
      // closing fence should be indented less than 4 spaces
      continue;
    }

    pos = state.skipChars(pos, marker);

    // closing code fence must be at least as long as the opening one
    if (pos - mem < len) {
      continue;
    }

    // make sure tail has spaces only
    pos = state.skipSpaces(pos);

    if (pos < max) {
      continue;
    }

    haveEndMarker = true;
    // found!
    break;
  }

  // If a fence has heading spaces, they should be removed from its inner block
  len = state.sCount[startLine];

  state.line = nextLine + (haveEndMarker ? 1 : 0);

  token = state.push("fence", "code", 0);
  token.info = params;
  token.content = state.getLines(startLine + 1, nextLine, len, true);
  token.markup = markup;
  token.map = [startLine, state.line];
  if (params === "mermaid" || params === "chart") {
    console.log("规则匹配啊，草泥马---" + token.info);
  }
  return true;
}

const htmlEntities = (str) =>
  String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");

(function (global, factory) {
  typeof exports === "object" && typeof module !== "undefined"
    ? (module.exports = factory())
    : typeof define === "function" && define.amd
    ? define(["mermaid"], factory)
    : ((global =
        typeof globalThis !== "undefined" ? globalThis : global || self),
      (global.markdownitMermaid = factory()));
})(this, function () {
  "use strict";

  const MermaidChart = function (code) {
    try {
      mermaid.mermaidAPI.render(randomString(16), code, (sc) => {
        code = sc;
      });
      return `<div class="mermaid">${code}</div>`;
    } catch (err) {
      return `<pre>${htmlEntities(err.name)}: ${htmlEntities(
        err.message
      )}</pre>`;
    }
  };

  const ChartjsChart = function (code) {
    try {
      const json = JSON.parse(code);
      return `<canvas class="chartjs">${JSON.stringify(json)}</canvas>`;
    } catch (err) {
      return `<pre>${htmlEntities(err.name)}: ${htmlEntities(
        err.message
      )}</pre>`;
    }
  };

  function mermaid_plugin(md, opts) {
    Object.assign(MermaidPlugIn.default, opts);
    const {
      token: _token = "mermaid",
      ...dictionary
    } =  MermaidPlugIn.default.dictionary;
    mermaid.initialize(MermaidPlugIn.default);
    const defaultRenderer = md.renderer.rules.fence.bind(md.renderer.rules);

    var mermaidRender = function (tokens, idx, options, env, slf) {
      console.log("cnm 为什么不执行");
      let token = tokens[idx];
      let code = token.content.trim();
      if (token.info.trim() === "mermaid") {
        return MermaidChart(code.replace(/(.*?)[ \n](.*?)([ \n])/, replacer));
      } else if (token.info.trim() === "chart") {
        return ChartjsChart(code);
      }
      return defaultRenderer(tokens, idx, options, env, self);
    };

    md.block.ruler.before("fence", "mermaid", mermaid_rule);
    md.renderer.rules.mermaid = mermaidRender;
  }

  const MermaidPlugIn = (md, opts) => {
    Object.assign(MermaidPlugIn.default, opts);
    const {
      token: _token = "mermaid",
      ...dictionary
    } = MermaidPlugIn.default.dictionary;
    // const dictionary = swapObj(_dictionary);
    mermaid.initialize(MermaidPlugIn.default);

    const defaultRenderer = md.renderer.rules.fence.bind(md.renderer.rules);

    function replacer(_, p1, p2, p3) {
      p1 = dictionary[p1] ?? p1;
      p2 = dictionary[p2] ?? p2;
      return p2 === "" ? `${p1}\n` : `${p1} ${p2}${p3}`;
    }

    md.renderer.rules.fence = (tokens, idx, opts, env, self) => {
      const token = tokens[idx];
      const code = token.content.trim();
      if (token.info.trim() === _token) {
        return MermaidChart(code.replace(/(.*?)[ \n](.*?)([ \n])/, replacer));
      } else if (token.info.trim() === "chart") {
        return ChartjsChart(code);
      }
      return defaultRenderer(tokens, idx, opts, env, self);
    };
  };

  MermaidPlugIn.default = {
    startOnLoad: false,
    securityLevel: "true",
    theme: "default",
    flowchart: {
      htmlLabels: false,
      useMaxWidth: true,
    },
    dictionary: {
      token: "mermaid",
    },
  };

  return MermaidPlugIn;
});
