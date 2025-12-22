/**
 * 字符串工具函数
 * 提供常用的字符串处理功能
 */

// ============================================================================
// 截断与填充
// ============================================================================

export interface TruncateOptions {
  /** 最大长度 */
  length?: number;
  /** 省略符号，默认 "..." */
  omission?: string;
  /** 截断位置: start | middle | end */
  position?: "start" | "middle" | "end";
}

/**
 * 截断字符串
 * @example
 * truncate("Hello World", { length: 8 }) // "Hello..."
 * truncate("Hello World", { length: 8, position: "start" }) // "...World"
 * truncate("Hello World", { length: 8, position: "middle" }) // "Hel...ld"
 */
export function truncate(str: string, options: TruncateOptions = {}): string {
  const { length = 30, omission = "...", position = "end" } = options;

  if (str.length <= length) {
    return str;
  }

  const omissionLength = omission.length;
  const charsToShow = length - omissionLength;

  if (charsToShow <= 0) {
    return omission.slice(0, length);
  }

  switch (position) {
    case "start":
      return omission + str.slice(-charsToShow);
    case "middle": {
      const frontChars = Math.ceil(charsToShow / 2);
      const backChars = Math.floor(charsToShow / 2);
      return str.slice(0, frontChars) + omission + str.slice(-backChars);
    }
    case "end":
    default:
      return str.slice(0, charsToShow) + omission;
  }
}

/**
 * 左填充字符串
 * @example
 * padStart("5", 3, "0") // "005"
 */
export function padStart(str: string | number, length: number, char: string = " "): string {
  return String(str).padStart(length, char);
}

/**
 * 右填充字符串
 * @example
 * padEnd("5", 3, "0") // "500"
 */
export function padEnd(str: string | number, length: number, char: string = " "): string {
  return String(str).padEnd(length, char);
}

// ============================================================================
// 大小写转换
// ============================================================================

/**
 * 首字母大写
 * @example
 * capitalize("hello") // "Hello"
 */
export function capitalize(str: string): string {
  if (!str) return str;
  return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
}

/**
 * 每个单词首字母大写
 * @example
 * titleCase("hello world") // "Hello World"
 */
export function titleCase(str: string): string {
  if (!str) return str;
  return str
    .toLowerCase()
    .split(/\s+/)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
}

/**
 * 转换为驼峰命名
 * @example
 * camelCase("hello-world") // "helloWorld"
 * camelCase("hello_world") // "helloWorld"
 */
export function camelCase(str: string): string {
  if (!str) return str;
  return str.toLowerCase().replace(/[-_\s]+(.)?/g, (_, char) => (char ? char.toUpperCase() : ""));
}

/**
 * 转换为帕斯卡命名（大驼峰）
 * @example
 * pascalCase("hello-world") // "HelloWorld"
 */
export function pascalCase(str: string): string {
  const camel = camelCase(str);
  return camel.charAt(0).toUpperCase() + camel.slice(1);
}

/**
 * 转换为蛇形命名
 * @example
 * snakeCase("helloWorld") // "hello_world"
 * snakeCase("HelloWorld") // "hello_world"
 */
export function snakeCase(str: string): string {
  if (!str) return str;
  return str
    .replace(/([a-z])([A-Z])/g, "$1_$2")
    .replace(/[-\s]+/g, "_")
    .toLowerCase();
}

/**
 * 转换为短横线命名
 * @example
 * kebabCase("helloWorld") // "hello-world"
 * kebabCase("HelloWorld") // "hello-world"
 */
export function kebabCase(str: string): string {
  if (!str) return str;
  return str
    .replace(/([a-z])([A-Z])/g, "$1-$2")
    .replace(/[_\s]+/g, "-")
    .toLowerCase();
}

/**
 * 转换为常量命名
 * @example
 * constantCase("helloWorld") // "HELLO_WORLD"
 */
export function constantCase(str: string): string {
  return snakeCase(str).toUpperCase();
}

// ============================================================================
// URL 与 Slug
// ============================================================================

export interface SlugifyOptions {
  /** 分隔符，默认 "-" */
  separator?: string;
  /** 是否转小写，默认 true */
  lowercase?: boolean;
  /** 是否移除特殊字符，默认 true */
  strict?: boolean;
}

/**
 * 转换为 URL 友好的 slug
 * @example
 * slugify("Hello World!") // "hello-world"
 * slugify("你好世界") // "ni-hao-shi-jie" (如果安装了拼音库)
 */
export function slugify(str: string, options: SlugifyOptions = {}): string {
  const { separator = "-", lowercase = true, strict = true } = options;

  let result = str.trim();

  // 转小写
  if (lowercase) {
    result = result.toLowerCase();
  }

  // 替换空格和常见分隔符
  result = result.replace(/[\s_]+/g, separator);

  // 移除特殊字符
  if (strict) {
    result = result.replace(/[^\w\u4e00-\u9fa5-]/g, "");
  }

  // 移除连续分隔符
  result = result.replace(new RegExp(`${separator}+`, "g"), separator);

  // 移除首尾分隔符
  result = result.replace(new RegExp(`^${separator}|${separator}$`, "g"), "");

  return result;
}

// ============================================================================
// 搜索与匹配
// ============================================================================

/**
 * 检查字符串是否包含子串（不区分大小写）
 * @example
 * containsIgnoreCase("Hello World", "world") // true
 */
export function containsIgnoreCase(str: string, search: string): boolean {
  return str.toLowerCase().includes(search.toLowerCase());
}

/**
 * 高亮搜索关键词
 * @example
 * highlight("Hello World", "world") // "Hello <mark>World</mark>"
 */
export function highlight(str: string, search: string, tag: string = "mark"): string {
  if (!search) return str;
  const regex = new RegExp(`(${escapeRegExp(search)})`, "gi");
  return str.replace(regex, `<${tag}>$1</${tag}>`);
}

/**
 * 转义正则表达式特殊字符
 * @example
 * escapeRegExp("hello.world") // "hello\\.world"
 */
export function escapeRegExp(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

/**
 * 模糊匹配（简单实现）
 * @example
 * fuzzyMatch("HelloWorld", "hw") // true
 * fuzzyMatch("HelloWorld", "wh") // false
 */
export function fuzzyMatch(str: string, pattern: string): boolean {
  const strLower = str.toLowerCase();
  const patternLower = pattern.toLowerCase();

  let patternIndex = 0;
  for (let i = 0; i < strLower.length && patternIndex < patternLower.length; i++) {
    if (strLower[i] === patternLower[patternIndex]) {
      patternIndex++;
    }
  }

  return patternIndex === patternLower.length;
}

// ============================================================================
// 格式化
// ============================================================================

/**
 * 手机号码格式化（脱敏）
 * @example
 * maskPhone("13812345678") // "138****5678"
 */
export function maskPhone(phone: string): string {
  if (!phone || phone.length < 7) return phone;
  return phone.slice(0, 3) + "****" + phone.slice(-4);
}

/**
 * 邮箱格式化（脱敏）
 * @example
 * maskEmail("hello@example.com") // "h****@example.com"
 */
export function maskEmail(email: string): string {
  if (!email) return email;
  const [local, domain] = email.split("@");
  if (!domain) return email;
  const maskedLocal = local.charAt(0) + "****";
  return maskedLocal + "@" + domain;
}

/**
 * 身份证号格式化（脱敏）
 * @example
 * maskIdCard("110101199001011234") // "1101****1234"
 */
export function maskIdCard(idCard: string): string {
  if (!idCard || idCard.length < 8) return idCard;
  return idCard.slice(0, 4) + "****" + idCard.slice(-4);
}

/**
 * 银行卡号格式化
 * @example
 * formatBankCard("6222021234567890123") // "6222 0212 3456 7890 123"
 */
export function formatBankCard(cardNo: string): string {
  if (!cardNo) return cardNo;
  return cardNo.replace(/\s/g, "").replace(/(\d{4})(?=\d)/g, "$1 ");
}

/**
 * 银行卡号脱敏
 * @example
 * maskBankCard("6222021234567890123") // "**** **** **** 0123"
 */
export function maskBankCard(cardNo: string): string {
  if (!cardNo || cardNo.length < 4) return cardNo;
  return "**** **** **** " + cardNo.slice(-4);
}

// ============================================================================
// 其他工具
// ============================================================================

/**
 * 移除 HTML 标签
 * @example
 * stripHtml("<p>Hello <b>World</b></p>") // "Hello World"
 */
export function stripHtml(str: string): string {
  return str.replace(/<[^>]*>/g, "");
}

/**
 * 转义 HTML 特殊字符
 * @example
 * escapeHtml("<script>alert('xss')</script>") // "&lt;script&gt;..."
 */
export function escapeHtml(str: string): string {
  const htmlEntities: Record<string, string> = {
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#39;",
  };
  return str.replace(/[&<>"']/g, (char) => htmlEntities[char]);
}

/**
 * 反转义 HTML 特殊字符
 */
export function unescapeHtml(str: string): string {
  const htmlEntities: Record<string, string> = {
    "&amp;": "&",
    "&lt;": "<",
    "&gt;": ">",
    "&quot;": '"',
    "&#39;": "'",
  };
  return str.replace(/&(amp|lt|gt|quot|#39);/g, (match) => htmlEntities[match]);
}

/**
 * 生成随机字符串
 * @example
 * randomString(8) // "aB3xK9mN"
 */
export function randomString(
  length: number = 8,
  chars: string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
): string {
  let result = "";
  for (let i = 0; i < length; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return result;
}

/**
 * 计算字符串的字节长度（UTF-8）
 * @example
 * byteLength("hello") // 5
 * byteLength("你好") // 6
 */
export function byteLength(str: string): number {
  let length = 0;
  for (let i = 0; i < str.length; i++) {
    const code = str.charCodeAt(i);
    if (code <= 0x7f) {
      length += 1;
    } else if (code <= 0x7ff) {
      length += 2;
    } else if (code <= 0xffff) {
      length += 3;
    } else {
      length += 4;
    }
  }
  return length;
}

/**
 * 检查字符串是否为空或仅包含空白字符
 * @example
 * isBlank("") // true
 * isBlank("  ") // true
 * isBlank("hello") // false
 */
export function isBlank(str: string | null | undefined): boolean {
  return str == null || str.trim().length === 0;
}

/**
 * 检查字符串是否非空且包含非空白字符
 */
export function isNotBlank(str: string | null | undefined): boolean {
  return !isBlank(str);
}

/**
 * 字符串模板替换
 * @example
 * template("Hello, {name}!", { name: "World" }) // "Hello, World!"
 */
export function template(str: string, data: Record<string, string | number>): string {
  return str.replace(/\{(\w+)\}/g, (_, key) => String(data[key] ?? ""));
}

/**
 * 统计子串出现次数
 * @example
 * countOccurrences("hello hello world", "hello") // 2
 */
export function countOccurrences(str: string, search: string): number {
  if (!search) return 0;
  let count = 0;
  let pos = 0;
  while ((pos = str.indexOf(search, pos)) !== -1) {
    count++;
    pos += search.length;
  }
  return count;
}

/**
 * 反转字符串
 * @example
 * reverse("hello") // "olleh"
 */
export function reverse(str: string): string {
  return [...str].reverse().join("");
}

/**
 * 移除字符串中的重复空格
 * @example
 * normalizeSpaces("hello   world") // "hello world"
 */
export function normalizeSpaces(str: string): string {
  return str.replace(/\s+/g, " ").trim();
}
