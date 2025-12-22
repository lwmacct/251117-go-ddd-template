/**
 * DOM 工具函数
 * 提供 DOM 操作、事件处理和样式管理
 */

// ============================================================================
// 类型定义
// ============================================================================

export type ElementTarget = Element | string | null | undefined;

export interface ScrollOptions {
  behavior?: ScrollBehavior;
  block?: ScrollLogicalPosition;
  inline?: ScrollLogicalPosition;
}

export interface Position {
  x: number;
  y: number;
}

export interface Size {
  width: number;
  height: number;
}

export interface Rect extends Position, Size {
  top: number;
  right: number;
  bottom: number;
  left: number;
}

// ============================================================================
// 元素查询
// ============================================================================

/**
 * 获取元素
 * @example
 * getElement('#my-id') // Element
 * getElement(document.body) // Element
 */
export function getElement(target: ElementTarget): Element | null {
  if (!target) {
    return null;
  }

  if (typeof target === "string") {
    return document.querySelector(target);
  }

  return target;
}

/**
 * 获取所有匹配元素
 * @example
 * getElements('.my-class') // Element[]
 */
export function getElements(selector: string): Element[] {
  return Array.from(document.querySelectorAll(selector));
}

/**
 * 检查元素是否存在
 * @example
 * elementExists('#my-id') // true/false
 */
export function elementExists(target: ElementTarget): boolean {
  return getElement(target) !== null;
}

/**
 * 等待元素出现
 * @example
 * await waitForElement('#dynamic-element')
 */
export function waitForElement(selector: string, timeout: number = 5000): Promise<Element> {
  return new Promise((resolve, reject) => {
    const element = document.querySelector(selector);
    if (element) {
      resolve(element);
      return;
    }

    const observer = new MutationObserver(() => {
      const element = document.querySelector(selector);
      if (element) {
        observer.disconnect();
        resolve(element);
      }
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true,
    });

    setTimeout(() => {
      observer.disconnect();
      reject(new Error(`Element "${selector}" not found within ${timeout}ms`));
    }, timeout);
  });
}

// ============================================================================
// 类名操作
// ============================================================================

/**
 * 添加类名
 * @example
 * addClass('#my-id', 'active')
 * addClass(element, ['a', 'b'])
 */
export function addClass(target: ElementTarget, ...classNames: (string | string[])[]): void {
  const element = getElement(target);
  if (!element) return;

  const classes = classNames.flat().filter(Boolean);
  element.classList.add(...classes);
}

/**
 * 移除类名
 * @example
 * removeClass('#my-id', 'active')
 */
export function removeClass(target: ElementTarget, ...classNames: (string | string[])[]): void {
  const element = getElement(target);
  if (!element) return;

  const classes = classNames.flat().filter(Boolean);
  element.classList.remove(...classes);
}

/**
 * 切换类名
 * @example
 * toggleClass('#my-id', 'active')
 * toggleClass('#my-id', 'active', true) // 强制添加
 */
export function toggleClass(target: ElementTarget, className: string, force?: boolean): boolean {
  const element = getElement(target);
  if (!element) return false;

  return element.classList.toggle(className, force);
}

/**
 * 检查是否包含类名
 * @example
 * hasClass('#my-id', 'active') // true/false
 */
export function hasClass(target: ElementTarget, className: string): boolean {
  const element = getElement(target);
  if (!element) return false;

  return element.classList.contains(className);
}

/**
 * 替换类名
 * @example
 * replaceClass('#my-id', 'old-class', 'new-class')
 */
export function replaceClass(target: ElementTarget, oldClass: string, newClass: string): boolean {
  const element = getElement(target);
  if (!element) return false;

  return element.classList.replace(oldClass, newClass);
}

// ============================================================================
// 样式操作
// ============================================================================

/**
 * 获取样式
 * @example
 * getStyle('#my-id', 'color') // 'rgb(0, 0, 0)'
 */
export function getStyle(target: ElementTarget, property: string): string {
  const element = getElement(target);
  if (!element) return "";

  return getComputedStyle(element).getPropertyValue(property);
}

/**
 * 设置样式
 * @example
 * setStyle('#my-id', 'color', 'red')
 * setStyle('#my-id', { color: 'red', fontSize: '16px' })
 */
export function setStyle(
  target: ElementTarget,
  property: string | Record<string, string | number>,
  value?: string | number
): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  if (typeof property === "string") {
    element.style.setProperty(property, typeof value === "number" ? `${value}px` : (value as string));
  } else {
    for (const [key, val] of Object.entries(property)) {
      const cssKey = key.replace(/([A-Z])/g, "-$1").toLowerCase();
      element.style.setProperty(cssKey, typeof val === "number" ? `${val}px` : val);
    }
  }
}

/**
 * 移除样式
 * @example
 * removeStyle('#my-id', 'color')
 */
export function removeStyle(target: ElementTarget, property: string): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  element.style.removeProperty(property);
}

/**
 * 获取 CSS 变量值
 * @example
 * getCSSVariable('--primary-color') // '#3498db'
 */
export function getCSSVariable(name: string, element?: Element): string {
  const target = element || document.documentElement;
  return getComputedStyle(target).getPropertyValue(name).trim();
}

/**
 * 设置 CSS 变量
 * @example
 * setCSSVariable('--primary-color', '#3498db')
 */
export function setCSSVariable(name: string, value: string, element?: Element): void {
  const target = (element || document.documentElement) as HTMLElement;
  target.style.setProperty(name, value);
}

// ============================================================================
// 属性操作
// ============================================================================

/**
 * 获取属性
 * @example
 * getAttribute('#my-id', 'data-value') // '123'
 */
export function getAttribute(target: ElementTarget, name: string): string | null {
  const element = getElement(target);
  if (!element) return null;

  return element.getAttribute(name);
}

/**
 * 设置属性
 * @example
 * setAttribute('#my-id', 'data-value', '123')
 * setAttribute('#my-id', { 'data-a': '1', 'data-b': '2' })
 */
export function setAttribute(target: ElementTarget, name: string | Record<string, string>, value?: string): void {
  const element = getElement(target);
  if (!element) return;

  if (typeof name === "string") {
    element.setAttribute(name, value || "");
  } else {
    for (const [key, val] of Object.entries(name)) {
      element.setAttribute(key, val);
    }
  }
}

/**
 * 移除属性
 * @example
 * removeAttribute('#my-id', 'data-value')
 */
export function removeAttribute(target: ElementTarget, name: string): void {
  const element = getElement(target);
  if (!element) return;

  element.removeAttribute(name);
}

/**
 * 检查属性是否存在
 * @example
 * hasAttribute('#my-id', 'disabled') // true/false
 */
export function hasAttribute(target: ElementTarget, name: string): boolean {
  const element = getElement(target);
  if (!element) return false;

  return element.hasAttribute(name);
}

/**
 * 获取 data 属性
 * @example
 * getDataAttribute('#my-id', 'value') // 等同于 data-value
 */
export function getDataAttribute(target: ElementTarget, name: string): string | undefined {
  const element = getElement(target) as HTMLElement;
  if (!element) return undefined;

  return element.dataset[name];
}

/**
 * 设置 data 属性
 * @example
 * setDataAttribute('#my-id', 'value', '123')
 */
export function setDataAttribute(target: ElementTarget, name: string, value: string): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  element.dataset[name] = value;
}

// ============================================================================
// 尺寸和位置
// ============================================================================

/**
 * 获取元素边界矩形
 * @example
 * getRect('#my-id')
 * // { x: 100, y: 50, width: 200, height: 100, top: 50, ... }
 */
export function getRect(target: ElementTarget): Rect | null {
  const element = getElement(target);
  if (!element) return null;

  const rect = element.getBoundingClientRect();
  return {
    x: rect.x,
    y: rect.y,
    width: rect.width,
    height: rect.height,
    top: rect.top,
    right: rect.right,
    bottom: rect.bottom,
    left: rect.left,
  };
}

/**
 * 获取元素尺寸
 * @example
 * getSize('#my-id') // { width: 200, height: 100 }
 */
export function getSize(target: ElementTarget): Size | null {
  const rect = getRect(target);
  if (!rect) return null;

  return {
    width: rect.width,
    height: rect.height,
  };
}

/**
 * 获取元素位置（相对于文档）
 * @example
 * getOffset('#my-id') // { x: 100, y: 200 }
 */
export function getOffset(target: ElementTarget): Position | null {
  const element = getElement(target) as HTMLElement;
  if (!element) return null;

  const rect = element.getBoundingClientRect();
  return {
    x: rect.left + window.scrollX,
    y: rect.top + window.scrollY,
  };
}

/**
 * 获取窗口尺寸
 * @example
 * getWindowSize() // { width: 1920, height: 1080 }
 */
export function getWindowSize(): Size {
  return {
    width: window.innerWidth,
    height: window.innerHeight,
  };
}

/**
 * 获取文档尺寸
 * @example
 * getDocumentSize() // { width: 1920, height: 5000 }
 */
export function getDocumentSize(): Size {
  return {
    width: Math.max(document.body.scrollWidth, document.documentElement.scrollWidth),
    height: Math.max(document.body.scrollHeight, document.documentElement.scrollHeight),
  };
}

/**
 * 获取滚动位置
 * @example
 * getScrollPosition() // { x: 0, y: 100 }
 */
export function getScrollPosition(target?: ElementTarget): Position {
  if (target) {
    const element = getElement(target);
    if (element) {
      return {
        x: element.scrollLeft,
        y: element.scrollTop,
      };
    }
  }

  return {
    x: window.scrollX || document.documentElement.scrollLeft,
    y: window.scrollY || document.documentElement.scrollTop,
  };
}

// ============================================================================
// 滚动操作
// ============================================================================

/**
 * 滚动到指定位置
 * @example
 * scrollTo(0, 500) // 滚动到 y=500
 * scrollTo('#section', { behavior: 'smooth' })
 */
export function scrollTo(target: ElementTarget | number, options?: ScrollOptions | number): void {
  if (typeof target === "number") {
    const y = typeof options === "number" ? options : target;
    const x = typeof options === "number" ? target : 0;

    window.scrollTo({
      left: x,
      top: y,
      behavior: "smooth",
    });
  } else {
    const element = getElement(target);
    if (!element) return;

    element.scrollIntoView(options as ScrollIntoViewOptions);
  }
}

/**
 * 滚动到顶部
 * @example
 * scrollToTop()
 * scrollToTop({ behavior: 'smooth' })
 */
export function scrollToTop(options?: ScrollOptions): void {
  window.scrollTo({
    top: 0,
    behavior: options?.behavior || "smooth",
  });
}

/**
 * 滚动到底部
 * @example
 * scrollToBottom()
 */
export function scrollToBottom(options?: ScrollOptions): void {
  const { height } = getDocumentSize();
  window.scrollTo({
    top: height,
    behavior: options?.behavior || "smooth",
  });
}

/**
 * 检查元素是否在视口内
 * @example
 * isInViewport('#my-id') // true/false
 */
export function isInViewport(target: ElementTarget, threshold: number = 0): boolean {
  const element = getElement(target);
  if (!element) return false;

  const rect = element.getBoundingClientRect();
  const { width, height } = getWindowSize();

  return (
    rect.top >= -threshold &&
    rect.left >= -threshold &&
    rect.bottom <= height + threshold &&
    rect.right <= width + threshold
  );
}

/**
 * 检查元素是否部分在视口内
 * @example
 * isPartiallyInViewport('#my-id') // true/false
 */
export function isPartiallyInViewport(target: ElementTarget): boolean {
  const element = getElement(target);
  if (!element) return false;

  const rect = element.getBoundingClientRect();
  const { width, height } = getWindowSize();

  return rect.top < height && rect.bottom > 0 && rect.left < width && rect.right > 0;
}

// ============================================================================
// 焦点操作
// ============================================================================

/**
 * 设置焦点
 * @example
 * focus('#my-input')
 */
export function focus(target: ElementTarget): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  element.focus();
}

/**
 * 移除焦点
 * @example
 * blur('#my-input')
 */
export function blur(target: ElementTarget): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  element.blur();
}

/**
 * 获取当前焦点元素
 * @example
 * getActiveElement() // HTMLInputElement
 */
export function getActiveElement(): Element | null {
  return document.activeElement;
}

/**
 * 检查元素是否有焦点
 * @example
 * hasFocus('#my-input') // true/false
 */
export function hasFocus(target: ElementTarget): boolean {
  const element = getElement(target);
  if (!element) return false;

  return document.activeElement === element;
}

// ============================================================================
// 元素创建和操作
// ============================================================================

/**
 * 创建元素
 * @example
 * createElement('div', { className: 'box', id: 'my-box' }, 'Hello')
 */
export function createElement<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  attributes?: Record<string, string>,
  content?: string | Element | Element[]
): HTMLElementTagNameMap[K] {
  const element = document.createElement(tag);

  if (attributes) {
    for (const [key, value] of Object.entries(attributes)) {
      if (key === "className") {
        element.className = value;
      } else if (key === "style" && typeof value === "object") {
        Object.assign(element.style, value);
      } else {
        element.setAttribute(key, value);
      }
    }
  }

  if (content) {
    if (typeof content === "string") {
      element.textContent = content;
    } else if (Array.isArray(content)) {
      element.append(...content);
    } else {
      element.appendChild(content);
    }
  }

  return element;
}

/**
 * 移除元素
 * @example
 * removeElement('#my-id')
 */
export function removeElement(target: ElementTarget): void {
  const element = getElement(target);
  if (!element) return;

  element.remove();
}

/**
 * 克隆元素
 * @example
 * cloneElement('#my-id', true) // 深克隆
 */
export function cloneElement<T extends Element>(target: ElementTarget, deep: boolean = true): T | null {
  const element = getElement(target);
  if (!element) return null;

  return element.cloneNode(deep) as T;
}

/**
 * 在元素前插入
 * @example
 * insertBefore('#target', newElement)
 */
export function insertBefore(target: ElementTarget, newElement: Element): void {
  const element = getElement(target);
  if (!element || !element.parentNode) return;

  element.parentNode.insertBefore(newElement, element);
}

/**
 * 在元素后插入
 * @example
 * insertAfter('#target', newElement)
 */
export function insertAfter(target: ElementTarget, newElement: Element): void {
  const element = getElement(target);
  if (!element || !element.parentNode) return;

  element.parentNode.insertBefore(newElement, element.nextSibling);
}

/**
 * 包裹元素
 * @example
 * wrap('#my-id', 'div', { className: 'wrapper' })
 */
export function wrap<K extends keyof HTMLElementTagNameMap>(
  target: ElementTarget,
  wrapperTag: K,
  attributes?: Record<string, string>
): HTMLElementTagNameMap[K] | null {
  const element = getElement(target);
  if (!element || !element.parentNode) return null;

  const wrapper = createElement(wrapperTag, attributes);
  element.parentNode.insertBefore(wrapper, element);
  wrapper.appendChild(element);

  return wrapper;
}

/**
 * 解除包裹
 * @example
 * unwrap('#wrapped-element')
 */
export function unwrap(target: ElementTarget): void {
  const element = getElement(target);
  if (!element || !element.parentNode) return;

  const parent = element.parentNode;
  const grandparent = parent.parentNode;

  if (!grandparent) return;

  while (parent.firstChild) {
    grandparent.insertBefore(parent.firstChild, parent);
  }

  grandparent.removeChild(parent);
}

// ============================================================================
// 事件工具
// ============================================================================

/**
 * 添加事件监听
 * @example
 * on('#my-id', 'click', handleClick)
 * on(window, 'resize', handleResize, { passive: true })
 */
export function on<K extends keyof HTMLElementEventMap>(
  target: ElementTarget | Window | Document,
  event: K,
  handler: (e: HTMLElementEventMap[K]) => void,
  options?: AddEventListenerOptions
): () => void {
  const element = typeof target === "string" ? getElement(target) : target;
  if (!element) return () => {};

  element.addEventListener(event, handler as EventListener, options);

  return () => {
    element.removeEventListener(event, handler as EventListener, options);
  };
}

/**
 * 移除事件监听
 * @example
 * off('#my-id', 'click', handleClick)
 */
export function off<K extends keyof HTMLElementEventMap>(
  target: ElementTarget | Window | Document,
  event: K,
  handler: (e: HTMLElementEventMap[K]) => void,
  options?: EventListenerOptions
): void {
  const element = typeof target === "string" ? getElement(target) : target;
  if (!element) return;

  element.removeEventListener(event, handler as EventListener, options);
}

/**
 * 一次性事件监听
 * @example
 * once('#my-id', 'click', handleClick)
 */
export function once<K extends keyof HTMLElementEventMap>(
  target: ElementTarget | Window | Document,
  event: K,
  handler: (e: HTMLElementEventMap[K]) => void
): () => void {
  return on(target, event, handler, { once: true });
}

/**
 * 触发事件
 * @example
 * trigger('#my-id', 'click')
 * trigger('#my-id', 'custom-event', { detail: { foo: 'bar' } })
 */
export function trigger(target: ElementTarget, eventName: string, options?: CustomEventInit): void {
  const element = getElement(target);
  if (!element) return;

  const event = new CustomEvent(eventName, options);
  element.dispatchEvent(event);
}

// ============================================================================
// 可见性
// ============================================================================

/**
 * 显示元素
 * @example
 * show('#my-id')
 */
export function show(target: ElementTarget): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  element.style.display = "";
}

/**
 * 隐藏元素
 * @example
 * hide('#my-id')
 */
export function hide(target: ElementTarget): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  element.style.display = "none";
}

/**
 * 切换显示/隐藏
 * @example
 * toggle('#my-id')
 */
export function toggle(target: ElementTarget, force?: boolean): void {
  const element = getElement(target) as HTMLElement;
  if (!element) return;

  const isHidden = element.style.display === "none" || getComputedStyle(element).display === "none";

  const shouldShow = force !== undefined ? force : isHidden;

  element.style.display = shouldShow ? "" : "none";
}

/**
 * 检查元素是否可见
 * @example
 * isVisible('#my-id') // true/false
 */
export function isVisible(target: ElementTarget): boolean {
  const element = getElement(target) as HTMLElement;
  if (!element) return false;

  return !!(element.offsetWidth || element.offsetHeight || element.getClientRects().length);
}

/**
 * 检查元素是否隐藏
 * @example
 * isHidden('#my-id') // true/false
 */
export function isHidden(target: ElementTarget): boolean {
  return !isVisible(target);
}
