// leanote 通用方法

//统计字数（中文按单字算，英文按单词算)
function calcWords(str) {
  sLen = 0;
  try {
    str = str.replace(/<[^>]+>/g, ""); //过滤所有的HTML标签
    //先将回车换行符做特殊处理
    str = str.replace(/(\r\n+|\s+|　+)/g, "龘");
    //处理英文字符数字，连续字母、数字、英文符号视为一个单词
    str = str.replace(/[\x00-\xff]/g, "m");
    //合并字符m，连续字母、数字、英文符号视为一个单词
    str = str.replace(/m+/g, "*");
    //去掉回车换行符
    str = str.replace(/龘+/g, "");
    //返回字数
    sLen = str.length;
  } catch (e) {}
  return sLen;
}

const ATTR_TIMEOUT = 'timeout';
const ATTR_TIMEOUT_VALUE = 2000;
// 判断元素是否被锁上
function nodeIsLocked(node) {
  if ($(node)[0].hasAttribute(ATTR_TIMEOUT)) {
    let timeout = $(node).attr(ATTR_TIMEOUT);
    if (Number(timeout) > Date.now()) {
      return true;
    }
  }
  return false;
}

function lockNode(node) {
  $(node).attr(ATTR_TIMEOUT, Date.now() + ATTR_TIMEOUT_VALUE);
}

function unlockNode(node) {
  $(node).removeAttr(ATTR_TIMEOUT);
}
