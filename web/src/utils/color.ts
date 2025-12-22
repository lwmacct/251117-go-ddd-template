/**
 * 颜色工具函数
 * 提供颜色转换、处理和生成
 */

// ============================================================================
// 类型定义
// ============================================================================

export interface RGB {
  r: number;
  g: number;
  b: number;
}

export interface RGBA extends RGB {
  a: number;
}

export interface HSL {
  h: number;
  s: number;
  l: number;
}

export interface HSLA extends HSL {
  a: number;
}

export interface HSV {
  h: number;
  s: number;
  v: number;
}

// ============================================================================
// 解析函数
// ============================================================================

/**
 * 解析十六进制颜色
 * @example
 * parseHex('#ff0000') // { r: 255, g: 0, b: 0 }
 * parseHex('#f00') // { r: 255, g: 0, b: 0 }
 * parseHex('#ff0000ff') // { r: 255, g: 0, b: 0, a: 1 }
 */
export function parseHex(hex: string): RGBA {
  let h = hex.replace("#", "");

  // 处理短格式
  if (h.length === 3) {
    h = h[0] + h[0] + h[1] + h[1] + h[2] + h[2];
  } else if (h.length === 4) {
    h = h[0] + h[0] + h[1] + h[1] + h[2] + h[2] + h[3] + h[3];
  }

  const r = parseInt(h.slice(0, 2), 16);
  const g = parseInt(h.slice(2, 4), 16);
  const b = parseInt(h.slice(4, 6), 16);
  const a = h.length === 8 ? parseInt(h.slice(6, 8), 16) / 255 : 1;

  return { r, g, b, a };
}

/**
 * 解析 RGB/RGBA 字符串
 * @example
 * parseRgb('rgb(255, 0, 0)') // { r: 255, g: 0, b: 0, a: 1 }
 * parseRgb('rgba(255, 0, 0, 0.5)') // { r: 255, g: 0, b: 0, a: 0.5 }
 */
export function parseRgb(rgb: string): RGBA {
  const match = rgb.match(/rgba?\s*\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(?:,\s*([\d.]+))?\s*\)/i);

  if (!match) {
    throw new Error(`Invalid RGB string: ${rgb}`);
  }

  return {
    r: parseInt(match[1], 10),
    g: parseInt(match[2], 10),
    b: parseInt(match[3], 10),
    a: match[4] ? parseFloat(match[4]) : 1,
  };
}

/**
 * 解析 HSL/HSLA 字符串
 * @example
 * parseHsl('hsl(0, 100%, 50%)') // { h: 0, s: 100, l: 50, a: 1 }
 */
export function parseHsl(hsl: string): HSLA {
  const match = hsl.match(/hsla?\s*\(\s*([\d.]+)\s*,\s*([\d.]+)%\s*,\s*([\d.]+)%\s*(?:,\s*([\d.]+))?\s*\)/i);

  if (!match) {
    throw new Error(`Invalid HSL string: ${hsl}`);
  }

  return {
    h: parseFloat(match[1]),
    s: parseFloat(match[2]),
    l: parseFloat(match[3]),
    a: match[4] ? parseFloat(match[4]) : 1,
  };
}

/**
 * 解析任意颜色格式
 * @example
 * parseColor('#ff0000') // { r: 255, g: 0, b: 0, a: 1 }
 * parseColor('rgb(255, 0, 0)') // { r: 255, g: 0, b: 0, a: 1 }
 */
export function parseColor(color: string): RGBA {
  const trimmed = color.trim().toLowerCase();

  if (trimmed.startsWith("#")) {
    return parseHex(trimmed);
  }

  if (trimmed.startsWith("rgb")) {
    return parseRgb(trimmed);
  }

  if (trimmed.startsWith("hsl")) {
    return hslToRgb(parseHsl(trimmed));
  }

  // 尝试作为命名颜色处理
  const namedColor = NAMED_COLORS[trimmed as keyof typeof NAMED_COLORS];
  if (namedColor) {
    return parseHex(namedColor);
  }

  throw new Error(`Invalid color: ${color}`);
}

// ============================================================================
// 转换函数
// ============================================================================

/**
 * RGB 转十六进制
 * @example
 * rgbToHex({ r: 255, g: 0, b: 0 }) // '#ff0000'
 * rgbToHex({ r: 255, g: 0, b: 0, a: 0.5 }) // '#ff000080'
 */
export function rgbToHex(rgb: RGB | RGBA): string {
  const toHex = (n: number) => Math.round(n).toString(16).padStart(2, "0");

  const hex = `#${toHex(rgb.r)}${toHex(rgb.g)}${toHex(rgb.b)}`;

  if ("a" in rgb && rgb.a !== 1) {
    return hex + toHex(rgb.a * 255);
  }

  return hex;
}

/**
 * RGB 转 HSL
 * @example
 * rgbToHsl({ r: 255, g: 0, b: 0 }) // { h: 0, s: 100, l: 50 }
 */
export function rgbToHsl(rgb: RGB): HSL {
  const r = rgb.r / 255;
  const g = rgb.g / 255;
  const b = rgb.b / 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  const l = (max + min) / 2;

  let h = 0;
  let s = 0;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

    switch (max) {
      case r:
        h = ((g - b) / d + (g < b ? 6 : 0)) / 6;
        break;
      case g:
        h = ((b - r) / d + 2) / 6;
        break;
      case b:
        h = ((r - g) / d + 4) / 6;
        break;
    }
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    l: Math.round(l * 100),
  };
}

/**
 * HSL 转 RGB
 * @example
 * hslToRgb({ h: 0, s: 100, l: 50 }) // { r: 255, g: 0, b: 0, a: 1 }
 */
export function hslToRgb(hsl: HSL | HSLA): RGBA {
  const h = hsl.h / 360;
  const s = hsl.s / 100;
  const l = hsl.l / 100;

  let r: number, g: number, b: number;

  if (s === 0) {
    r = g = b = l;
  } else {
    const hue2rgb = (p: number, q: number, t: number) => {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1 / 6) return p + (q - p) * 6 * t;
      if (t < 1 / 2) return q;
      if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
      return p;
    };

    const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
    const p = 2 * l - q;

    r = hue2rgb(p, q, h + 1 / 3);
    g = hue2rgb(p, q, h);
    b = hue2rgb(p, q, h - 1 / 3);
  }

  return {
    r: Math.round(r * 255),
    g: Math.round(g * 255),
    b: Math.round(b * 255),
    a: "a" in hsl ? hsl.a : 1,
  };
}

/**
 * RGB 转 HSV
 * @example
 * rgbToHsv({ r: 255, g: 0, b: 0 }) // { h: 0, s: 100, v: 100 }
 */
export function rgbToHsv(rgb: RGB): HSV {
  const r = rgb.r / 255;
  const g = rgb.g / 255;
  const b = rgb.b / 255;

  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  const d = max - min;

  let h = 0;
  const s = max === 0 ? 0 : d / max;
  const v = max;

  if (max !== min) {
    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }
    h /= 6;
  }

  return {
    h: Math.round(h * 360),
    s: Math.round(s * 100),
    v: Math.round(v * 100),
  };
}

/**
 * HSV 转 RGB
 * @example
 * hsvToRgb({ h: 0, s: 100, v: 100 }) // { r: 255, g: 0, b: 0 }
 */
export function hsvToRgb(hsv: HSV): RGB {
  const h = hsv.h / 360;
  const s = hsv.s / 100;
  const v = hsv.v / 100;

  let r: number, g: number, b: number;

  const i = Math.floor(h * 6);
  const f = h * 6 - i;
  const p = v * (1 - s);
  const q = v * (1 - f * s);
  const t = v * (1 - (1 - f) * s);

  switch (i % 6) {
    case 0:
      r = v;
      g = t;
      b = p;
      break;
    case 1:
      r = q;
      g = v;
      b = p;
      break;
    case 2:
      r = p;
      g = v;
      b = t;
      break;
    case 3:
      r = p;
      g = q;
      b = v;
      break;
    case 4:
      r = t;
      g = p;
      b = v;
      break;
    default:
      r = v;
      g = p;
      b = q;
      break;
  }

  return {
    r: Math.round(r * 255),
    g: Math.round(g * 255),
    b: Math.round(b * 255),
  };
}

// ============================================================================
// 格式化函数
// ============================================================================

/**
 * 格式化为 RGB 字符串
 * @example
 * formatRgb({ r: 255, g: 0, b: 0 }) // 'rgb(255, 0, 0)'
 * formatRgb({ r: 255, g: 0, b: 0, a: 0.5 }) // 'rgba(255, 0, 0, 0.5)'
 */
export function formatRgb(rgb: RGB | RGBA): string {
  if ("a" in rgb && rgb.a !== 1) {
    return `rgba(${rgb.r}, ${rgb.g}, ${rgb.b}, ${rgb.a})`;
  }
  return `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})`;
}

/**
 * 格式化为 HSL 字符串
 * @example
 * formatHsl({ h: 0, s: 100, l: 50 }) // 'hsl(0, 100%, 50%)'
 */
export function formatHsl(hsl: HSL | HSLA): string {
  if ("a" in hsl && hsl.a !== 1) {
    return `hsla(${hsl.h}, ${hsl.s}%, ${hsl.l}%, ${hsl.a})`;
  }
  return `hsl(${hsl.h}, ${hsl.s}%, ${hsl.l}%)`;
}

// ============================================================================
// 颜色操作
// ============================================================================

/**
 * 调亮颜色
 * @example
 * lighten('#ff0000', 20) // '#ff6666'
 */
export function lighten(color: string, amount: number): string {
  const rgba = parseColor(color);
  const hsl = rgbToHsl(rgba);

  hsl.l = Math.min(100, hsl.l + amount);

  const result = hslToRgb({ ...hsl, a: rgba.a });
  return rgbToHex(result);
}

/**
 * 调暗颜色
 * @example
 * darken('#ff0000', 20) // '#990000'
 */
export function darken(color: string, amount: number): string {
  const rgba = parseColor(color);
  const hsl = rgbToHsl(rgba);

  hsl.l = Math.max(0, hsl.l - amount);

  const result = hslToRgb({ ...hsl, a: rgba.a });
  return rgbToHex(result);
}

/**
 * 调整饱和度
 * @example
 * saturate('#ff0000', 20) // 更鲜艳
 * desaturate('#ff0000', 20) // 更灰暗
 */
export function saturate(color: string, amount: number): string {
  const rgba = parseColor(color);
  const hsl = rgbToHsl(rgba);

  hsl.s = Math.min(100, hsl.s + amount);

  const result = hslToRgb({ ...hsl, a: rgba.a });
  return rgbToHex(result);
}

export function desaturate(color: string, amount: number): string {
  return saturate(color, -amount);
}

/**
 * 设置透明度
 * @example
 * setAlpha('#ff0000', 0.5) // '#ff000080'
 */
export function setAlpha(color: string, alpha: number): string {
  const rgba = parseColor(color);
  rgba.a = Math.max(0, Math.min(1, alpha));
  return rgbToHex(rgba);
}

/**
 * 反转颜色
 * @example
 * invert('#ff0000') // '#00ffff'
 */
export function invert(color: string): string {
  const rgba = parseColor(color);
  return rgbToHex({
    r: 255 - rgba.r,
    g: 255 - rgba.g,
    b: 255 - rgba.b,
    a: rgba.a,
  });
}

/**
 * 转为灰度
 * @example
 * grayscale('#ff0000') // '#4c4c4c'
 */
export function grayscale(color: string): string {
  const rgba = parseColor(color);
  const gray = Math.round(0.299 * rgba.r + 0.587 * rgba.g + 0.114 * rgba.b);
  return rgbToHex({ r: gray, g: gray, b: gray, a: rgba.a });
}

/**
 * 混合两种颜色
 * @example
 * mix('#ff0000', '#0000ff', 0.5) // '#800080'
 */
export function mix(color1: string, color2: string, ratio: number = 0.5): string {
  const c1 = parseColor(color1);
  const c2 = parseColor(color2);

  const r = Math.round(c1.r * (1 - ratio) + c2.r * ratio);
  const g = Math.round(c1.g * (1 - ratio) + c2.g * ratio);
  const b = Math.round(c1.b * (1 - ratio) + c2.b * ratio);
  const a = c1.a * (1 - ratio) + c2.a * ratio;

  return rgbToHex({ r, g, b, a });
}

/**
 * 获取补色
 * @example
 * complement('#ff0000') // '#00ffff'
 */
export function complement(color: string): string {
  const rgba = parseColor(color);
  const hsl = rgbToHsl(rgba);

  hsl.h = (hsl.h + 180) % 360;

  const result = hslToRgb({ ...hsl, a: rgba.a });
  return rgbToHex(result);
}

// ============================================================================
// 颜色分析
// ============================================================================

/**
 * 计算亮度 (0-1)
 * @example
 * getLuminance('#ffffff') // 1
 * getLuminance('#000000') // 0
 */
export function getLuminance(color: string): number {
  const rgba = parseColor(color);

  const toLinear = (c: number) => {
    const sRGB = c / 255;
    return sRGB <= 0.03928 ? sRGB / 12.92 : Math.pow((sRGB + 0.055) / 1.055, 2.4);
  };

  return 0.2126 * toLinear(rgba.r) + 0.7152 * toLinear(rgba.g) + 0.0722 * toLinear(rgba.b);
}

/**
 * 计算对比度
 * @example
 * getContrast('#ffffff', '#000000') // 21
 */
export function getContrast(color1: string, color2: string): number {
  const l1 = getLuminance(color1);
  const l2 = getLuminance(color2);

  const lighter = Math.max(l1, l2);
  const darker = Math.min(l1, l2);

  return (lighter + 0.05) / (darker + 0.05);
}

/**
 * 判断是否为深色
 * @example
 * isDark('#000000') // true
 * isDark('#ffffff') // false
 */
export function isDark(color: string): boolean {
  return getLuminance(color) < 0.5;
}

/**
 * 判断是否为浅色
 * @example
 * isLight('#ffffff') // true
 */
export function isLight(color: string): boolean {
  return !isDark(color);
}

/**
 * 获取适合的文本颜色
 * @example
 * getTextColor('#000000') // '#ffffff'
 * getTextColor('#ffffff') // '#000000'
 */
export function getTextColor(
  backgroundColor: string,
  lightText: string = "#ffffff",
  darkText: string = "#000000"
): string {
  return isDark(backgroundColor) ? lightText : darkText;
}

// ============================================================================
// 颜色生成
// ============================================================================

/**
 * 生成随机颜色
 * @example
 * randomColor() // '#a1b2c3'
 */
export function randomColor(): string {
  const r = Math.floor(Math.random() * 256);
  const g = Math.floor(Math.random() * 256);
  const b = Math.floor(Math.random() * 256);
  return rgbToHex({ r, g, b });
}

/**
 * 生成渐变色数组
 * @example
 * generateGradient('#ff0000', '#0000ff', 5)
 * // ['#ff0000', '#bf003f', '#7f007f', '#3f00bf', '#0000ff']
 */
export function generateGradient(startColor: string, endColor: string, steps: number): string[] {
  const colors: string[] = [];

  for (let i = 0; i < steps; i++) {
    const ratio = i / (steps - 1);
    colors.push(mix(startColor, endColor, ratio));
  }

  return colors;
}

/**
 * 生成调色板（基于主色）
 * @example
 * generatePalette('#3498db')
 * // { 50: '...', 100: '...', ..., 900: '...' }
 */
export function generatePalette(
  baseColor: string
): Record<50 | 100 | 200 | 300 | 400 | 500 | 600 | 700 | 800 | 900, string> {
  return {
    50: lighten(baseColor, 45),
    100: lighten(baseColor, 40),
    200: lighten(baseColor, 30),
    300: lighten(baseColor, 20),
    400: lighten(baseColor, 10),
    500: baseColor,
    600: darken(baseColor, 10),
    700: darken(baseColor, 20),
    800: darken(baseColor, 30),
    900: darken(baseColor, 40),
  };
}

// ============================================================================
// 命名颜色
// ============================================================================

const NAMED_COLORS = {
  aliceblue: "#f0f8ff",
  antiquewhite: "#faebd7",
  aqua: "#00ffff",
  aquamarine: "#7fffd4",
  azure: "#f0ffff",
  beige: "#f5f5dc",
  bisque: "#ffe4c4",
  black: "#000000",
  blanchedalmond: "#ffebcd",
  blue: "#0000ff",
  blueviolet: "#8a2be2",
  brown: "#a52a2a",
  burlywood: "#deb887",
  cadetblue: "#5f9ea0",
  chartreuse: "#7fff00",
  chocolate: "#d2691e",
  coral: "#ff7f50",
  cornflowerblue: "#6495ed",
  cornsilk: "#fff8dc",
  crimson: "#dc143c",
  cyan: "#00ffff",
  darkblue: "#00008b",
  darkcyan: "#008b8b",
  darkgoldenrod: "#b8860b",
  darkgray: "#a9a9a9",
  darkgreen: "#006400",
  darkkhaki: "#bdb76b",
  darkmagenta: "#8b008b",
  darkolivegreen: "#556b2f",
  darkorange: "#ff8c00",
  darkorchid: "#9932cc",
  darkred: "#8b0000",
  darksalmon: "#e9967a",
  darkseagreen: "#8fbc8f",
  darkslateblue: "#483d8b",
  darkslategray: "#2f4f4f",
  darkturquoise: "#00ced1",
  darkviolet: "#9400d3",
  deeppink: "#ff1493",
  deepskyblue: "#00bfff",
  dimgray: "#696969",
  dodgerblue: "#1e90ff",
  firebrick: "#b22222",
  floralwhite: "#fffaf0",
  forestgreen: "#228b22",
  fuchsia: "#ff00ff",
  gainsboro: "#dcdcdc",
  ghostwhite: "#f8f8ff",
  gold: "#ffd700",
  goldenrod: "#daa520",
  gray: "#808080",
  green: "#008000",
  greenyellow: "#adff2f",
  honeydew: "#f0fff0",
  hotpink: "#ff69b4",
  indianred: "#cd5c5c",
  indigo: "#4b0082",
  ivory: "#fffff0",
  khaki: "#f0e68c",
  lavender: "#e6e6fa",
  lavenderblush: "#fff0f5",
  lawngreen: "#7cfc00",
  lemonchiffon: "#fffacd",
  lightblue: "#add8e6",
  lightcoral: "#f08080",
  lightcyan: "#e0ffff",
  lightgoldenrodyellow: "#fafad2",
  lightgray: "#d3d3d3",
  lightgreen: "#90ee90",
  lightpink: "#ffb6c1",
  lightsalmon: "#ffa07a",
  lightseagreen: "#20b2aa",
  lightskyblue: "#87cefa",
  lightslategray: "#778899",
  lightsteelblue: "#b0c4de",
  lightyellow: "#ffffe0",
  lime: "#00ff00",
  limegreen: "#32cd32",
  linen: "#faf0e6",
  magenta: "#ff00ff",
  maroon: "#800000",
  mediumaquamarine: "#66cdaa",
  mediumblue: "#0000cd",
  mediumorchid: "#ba55d3",
  mediumpurple: "#9370db",
  mediumseagreen: "#3cb371",
  mediumslateblue: "#7b68ee",
  mediumspringgreen: "#00fa9a",
  mediumturquoise: "#48d1cc",
  mediumvioletred: "#c71585",
  midnightblue: "#191970",
  mintcream: "#f5fffa",
  mistyrose: "#ffe4e1",
  moccasin: "#ffe4b5",
  navajowhite: "#ffdead",
  navy: "#000080",
  oldlace: "#fdf5e6",
  olive: "#808000",
  olivedrab: "#6b8e23",
  orange: "#ffa500",
  orangered: "#ff4500",
  orchid: "#da70d6",
  palegoldenrod: "#eee8aa",
  palegreen: "#98fb98",
  paleturquoise: "#afeeee",
  palevioletred: "#db7093",
  papayawhip: "#ffefd5",
  peachpuff: "#ffdab9",
  peru: "#cd853f",
  pink: "#ffc0cb",
  plum: "#dda0dd",
  powderblue: "#b0e0e6",
  purple: "#800080",
  rebeccapurple: "#663399",
  red: "#ff0000",
  rosybrown: "#bc8f8f",
  royalblue: "#4169e1",
  saddlebrown: "#8b4513",
  salmon: "#fa8072",
  sandybrown: "#f4a460",
  seagreen: "#2e8b57",
  seashell: "#fff5ee",
  sienna: "#a0522d",
  silver: "#c0c0c0",
  skyblue: "#87ceeb",
  slateblue: "#6a5acd",
  slategray: "#708090",
  snow: "#fffafa",
  springgreen: "#00ff7f",
  steelblue: "#4682b4",
  tan: "#d2b48c",
  teal: "#008080",
  thistle: "#d8bfd8",
  tomato: "#ff6347",
  turquoise: "#40e0d0",
  violet: "#ee82ee",
  wheat: "#f5deb3",
  white: "#ffffff",
  whitesmoke: "#f5f5f5",
  yellow: "#ffff00",
  yellowgreen: "#9acd32",
} as const;

/**
 * 获取命名颜色
 * @example
 * getNamedColor('red') // '#ff0000'
 */
export function getNamedColor(name: string): string | undefined {
  return NAMED_COLORS[name.toLowerCase() as keyof typeof NAMED_COLORS];
}

/**
 * 检查是否为有效颜色
 * @example
 * isValidColor('#ff0000') // true
 * isValidColor('invalid') // false
 */
export function isValidColor(color: string): boolean {
  try {
    parseColor(color);
    return true;
  } catch {
    return false;
  }
}
