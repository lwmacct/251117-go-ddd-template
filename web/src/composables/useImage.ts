/**
 * Image Composable
 * 提供图片加载和处理的响应式管理
 */

import { ref, computed, watch, onMounted, onUnmounted, type Ref, type ComputedRef, type MaybeRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseImageOptions {
  /** 是否立即加载 */
  immediate?: boolean;
  /** 加载延迟（毫秒） */
  delay?: number;
  /** 加载失败时的回退图片 */
  fallback?: string;
  /** 加载中的占位图片 */
  placeholder?: string;
  /** 跨域设置 */
  crossOrigin?: "anonymous" | "use-credentials" | "";
  /** 引用策略 */
  referrerPolicy?: string;
  /** 加载超时（毫秒） */
  timeout?: number;
  /** 加载成功回调 */
  onLoad?: (image: HTMLImageElement) => void;
  /** 加载失败回调 */
  onError?: (error: Error) => void;
}

export interface UseImageReturn {
  /** 图片元素 */
  image: Ref<HTMLImageElement | null>;
  /** 是否正在加载 */
  isLoading: Ref<boolean>;
  /** 是否加载完成 */
  isReady: Ref<boolean>;
  /** 加载错误 */
  error: Ref<Error | null>;
  /** 图片宽度 */
  width: Ref<number>;
  /** 图片高度 */
  height: Ref<number>;
  /** 宽高比 */
  aspectRatio: ComputedRef<number>;
  /** 手动加载 */
  load: () => Promise<HTMLImageElement>;
  /** 取消加载 */
  abort: () => void;
}

// ============================================================================
// useImage - 图片加载
// ============================================================================

/**
 * 图片加载
 * @example
 * const { isLoading, isReady, error, width, height } = useImage('/image.jpg')
 *
 * // 懒加载
 * const { load, isReady } = useImage('/large-image.jpg', { immediate: false })
 * await load()
 */
export function useImage(src: MaybeRef<string>, options: UseImageOptions = {}): UseImageReturn {
  const { immediate = true, delay = 0, fallback, crossOrigin, referrerPolicy, timeout, onLoad, onError } = options;

  const image = ref<HTMLImageElement | null>(null);
  const isLoading = ref(false);
  const isReady = ref(false);
  const error = ref<Error | null>(null);
  const width = ref(0);
  const height = ref(0);

  const aspectRatio = computed(() => {
    if (height.value === 0) return 0;
    return width.value / height.value;
  });

  let abortController: AbortController | null = null;
  let timeoutTimer: ReturnType<typeof setTimeout> | null = null;

  const getSrc = () => {
    return typeof src === "object" && "value" in src ? src.value : src;
  };

  const load = (): Promise<HTMLImageElement> => {
    return new Promise((resolve, reject) => {
      abort();

      const currentSrc = getSrc();
      if (!currentSrc) {
        reject(new Error("No image source provided"));
        return;
      }

      isLoading.value = true;
      isReady.value = false;
      error.value = null;

      abortController = new AbortController();

      const loadImage = () => {
        const img = new Image();

        if (crossOrigin) img.crossOrigin = crossOrigin;
        if (referrerPolicy) img.referrerPolicy = referrerPolicy;

        // 设置超时
        if (timeout) {
          timeoutTimer = setTimeout(() => {
            abort();
            const err = new Error(`Image load timeout: ${currentSrc}`);
            error.value = err;
            isLoading.value = false;
            onError?.(err);

            if (fallback) {
              img.src = fallback;
            } else {
              reject(err);
            }
          }, timeout);
        }

        img.onload = () => {
          if (timeoutTimer) clearTimeout(timeoutTimer);

          image.value = img;
          width.value = img.naturalWidth;
          height.value = img.naturalHeight;
          isLoading.value = false;
          isReady.value = true;
          error.value = null;

          onLoad?.(img);
          resolve(img);
        };

        img.onerror = () => {
          if (timeoutTimer) clearTimeout(timeoutTimer);

          const err = new Error(`Failed to load image: ${currentSrc}`);
          error.value = err;
          isLoading.value = false;
          onError?.(err);

          if (fallback) {
            img.src = fallback;
          } else {
            reject(err);
          }
        };

        img.src = currentSrc;
      };

      if (delay > 0) {
        setTimeout(loadImage, delay);
      } else {
        loadImage();
      }
    });
  };

  const abort = () => {
    if (abortController) {
      abortController.abort();
      abortController = null;
    }
    if (timeoutTimer) {
      clearTimeout(timeoutTimer);
      timeoutTimer = null;
    }
    isLoading.value = false;
  };

  // 监听 src 变化
  if (typeof src === "object" && "value" in src) {
    watch(src, () => {
      if (immediate) {
        load().catch(() => {});
      }
    });
  }

  // 自动加载
  if (immediate) {
    onMounted(() => {
      load().catch(() => {});
    });
  }

  onUnmounted(() => {
    abort();
  });

  return {
    image,
    isLoading,
    isReady,
    error,
    width,
    height,
    aspectRatio,
    load,
    abort,
  };
}

// ============================================================================
// useImagePreload - 图片预加载
// ============================================================================

export interface UseImagePreloadReturn {
  /** 是否正在加载 */
  isLoading: Ref<boolean>;
  /** 加载进度 (0-1) */
  progress: Ref<number>;
  /** 已加载数量 */
  loaded: Ref<number>;
  /** 总数量 */
  total: Ref<number>;
  /** 加载失败的图片 */
  failed: Ref<string[]>;
  /** 开始预加载 */
  preload: () => Promise<void>;
  /** 取消预加载 */
  abort: () => void;
}

/**
 * 图片预加载
 * @example
 * const images = ['/img1.jpg', '/img2.jpg', '/img3.jpg']
 * const { preload, progress, isLoading } = useImagePreload(images)
 *
 * await preload()
 * console.log('所有图片已预加载')
 */
export function useImagePreload(sources: MaybeRef<string[]>): UseImagePreloadReturn {
  const isLoading = ref(false);
  const progress = ref(0);
  const loaded = ref(0);
  const total = ref(0);
  const failed = ref<string[]>([]);

  let aborted = false;

  const getSources = () => {
    return typeof sources === "object" && "value" in sources ? sources.value : sources;
  };

  const preload = async (): Promise<void> => {
    const urls = getSources();
    if (urls.length === 0) return;

    aborted = false;
    isLoading.value = true;
    progress.value = 0;
    loaded.value = 0;
    total.value = urls.length;
    failed.value = [];

    const loadPromises = urls.map((url) => {
      return new Promise<void>((resolve) => {
        if (aborted) {
          resolve();
          return;
        }

        const img = new Image();

        img.onload = () => {
          if (!aborted) {
            loaded.value++;
            progress.value = loaded.value / total.value;
          }
          resolve();
        };

        img.onerror = () => {
          if (!aborted) {
            loaded.value++;
            progress.value = loaded.value / total.value;
            failed.value.push(url);
          }
          resolve();
        };

        img.src = url;
      });
    });

    await Promise.all(loadPromises);
    isLoading.value = false;
  };

  const abort = () => {
    aborted = true;
    isLoading.value = false;
  };

  return {
    isLoading,
    progress,
    loaded,
    total,
    failed,
    preload,
    abort,
  };
}

// ============================================================================
// useLazyImage - 懒加载图片
// ============================================================================

export interface UseLazyImageOptions extends UseImageOptions {
  /** 触发加载的阈值 */
  threshold?: number;
  /** 根元素边距 */
  rootMargin?: string;
  /** 根元素 */
  root?: HTMLElement | null;
}

export interface UseLazyImageReturn extends UseImageReturn {
  /** 目标元素引用 */
  targetRef: Ref<HTMLElement | null>;
  /** 是否可见 */
  isVisible: Ref<boolean>;
}

/**
 * 懒加载图片
 * @example
 * const { targetRef, isReady, isVisible } = useLazyImage('/large-image.jpg')
 *
 * // 在模板中使用
 * // <div ref="targetRef">
 * //   <img v-if="isReady" :src="src" />
 * // </div>
 */
export function useLazyImage(src: MaybeRef<string>, options: UseLazyImageOptions = {}): UseLazyImageReturn {
  const { threshold = 0.1, rootMargin = "50px", root = null, ...imageOptions } = options;

  const targetRef = ref<HTMLElement | null>(null);
  const isVisible = ref(false);

  const imageState = useImage(src, {
    ...imageOptions,
    immediate: false,
  });

  let observer: IntersectionObserver | null = null;

  onMounted(() => {
    if (!targetRef.value) return;

    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            isVisible.value = true;
            imageState.load().catch(() => {});

            // 加载后停止观察
            if (observer && targetRef.value) {
              observer.unobserve(targetRef.value);
            }
          }
        });
      },
      {
        threshold,
        rootMargin,
        root,
      }
    );

    observer.observe(targetRef.value);
  });

  onUnmounted(() => {
    if (observer) {
      observer.disconnect();
      observer = null;
    }
  });

  return {
    ...imageState,
    targetRef,
    isVisible,
  };
}

// ============================================================================
// useProgressiveImage - 渐进式图片加载
// ============================================================================

export interface UseProgressiveImageOptions {
  /** 加载成功回调 */
  onLoad?: (image: HTMLImageElement) => void;
  /** 加载失败回调 */
  onError?: (error: Error) => void;
}

export interface UseProgressiveImageReturn {
  /** 当前显示的图片 URL */
  currentSrc: Ref<string>;
  /** 是否正在加载缩略图 */
  isLoadingThumb: Ref<boolean>;
  /** 是否正在加载原图 */
  isLoadingFull: Ref<boolean>;
  /** 是否完全加载 */
  isReady: Ref<boolean>;
  /** 加载阶段 */
  stage: Ref<"placeholder" | "thumbnail" | "full">;
  /** 错误信息 */
  error: Ref<Error | null>;
}

/**
 * 渐进式图片加载
 * @example
 * const { currentSrc, stage, isReady } = useProgressiveImage({
 *   placeholder: '/placeholder.svg',
 *   thumbnail: '/thumb.jpg',
 *   full: '/full.jpg'
 * })
 *
 * // stage: 'placeholder' -> 'thumbnail' -> 'full'
 */
export function useProgressiveImage(
  images: {
    placeholder?: string;
    thumbnail?: string;
    full: string;
  },
  options: UseProgressiveImageOptions = {}
): UseProgressiveImageReturn {
  const { onLoad, onError } = options;

  const currentSrc = ref(images.placeholder || "");
  const isLoadingThumb = ref(false);
  const isLoadingFull = ref(false);
  const isReady = ref(false);
  const stage = ref<"placeholder" | "thumbnail" | "full">("placeholder");
  const error = ref<Error | null>(null);

  const loadImage = (src: string): Promise<HTMLImageElement> => {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.onload = () => resolve(img);
      img.onerror = () => reject(new Error(`Failed to load: ${src}`));
      img.src = src;
    });
  };

  onMounted(async () => {
    try {
      // 加载缩略图
      if (images.thumbnail) {
        isLoadingThumb.value = true;
        await loadImage(images.thumbnail);
        currentSrc.value = images.thumbnail;
        stage.value = "thumbnail";
        isLoadingThumb.value = false;
      }

      // 加载原图
      isLoadingFull.value = true;
      const fullImg = await loadImage(images.full);
      currentSrc.value = images.full;
      stage.value = "full";
      isLoadingFull.value = false;
      isReady.value = true;
      onLoad?.(fullImg);
    } catch (err) {
      const loadError = err instanceof Error ? err : new Error(String(err));
      error.value = loadError;
      isLoadingThumb.value = false;
      isLoadingFull.value = false;
      onError?.(loadError);
    }
  });

  return {
    currentSrc,
    isLoadingThumb,
    isLoadingFull,
    isReady,
    stage,
    error,
  };
}

// ============================================================================
// useImageValidation - 图片验证
// ============================================================================

export interface ImageValidationResult {
  /** 是否有效 */
  valid: boolean;
  /** 错误消息 */
  errors: string[];
  /** 图片信息 */
  info?: {
    width: number;
    height: number;
    type: string;
    size: number;
    aspectRatio: number;
  };
}

export interface UseImageValidationOptions {
  /** 最大宽度 */
  maxWidth?: number;
  /** 最大高度 */
  maxHeight?: number;
  /** 最小宽度 */
  minWidth?: number;
  /** 最小高度 */
  minHeight?: number;
  /** 最大文件大小（字节） */
  maxSize?: number;
  /** 允许的类型 */
  allowedTypes?: string[];
  /** 固定宽高比 */
  aspectRatio?: number;
  /** 宽高比容差 */
  aspectRatioTolerance?: number;
}

/**
 * 验证图片文件
 * @example
 * const result = await validateImage(file, {
 *   maxWidth: 1920,
 *   maxHeight: 1080,
 *   maxSize: 5 * 1024 * 1024, // 5MB
 *   allowedTypes: ['image/jpeg', 'image/png']
 * })
 *
 * if (!result.valid) {
 *   console.log('验证失败:', result.errors)
 * }
 */
export async function validateImage(
  file: File,
  options: UseImageValidationOptions = {}
): Promise<ImageValidationResult> {
  const {
    maxWidth,
    maxHeight,
    minWidth,
    minHeight,
    maxSize,
    allowedTypes,
    aspectRatio,
    aspectRatioTolerance = 0.1,
  } = options;

  const errors: string[] = [];

  // 验证类型
  if (allowedTypes && !allowedTypes.includes(file.type)) {
    errors.push(`不支持的文件类型: ${file.type}`);
  }

  // 验证大小
  if (maxSize && file.size > maxSize) {
    errors.push(`文件大小超过限制: ${(file.size / 1024 / 1024).toFixed(2)}MB`);
  }

  // 加载图片获取尺寸
  const img = await new Promise<HTMLImageElement>((resolve, reject) => {
    const image = new Image();
    image.onload = () => resolve(image);
    image.onerror = () => reject(new Error("无法加载图片"));
    image.src = URL.createObjectURL(file);
  });

  const width = img.naturalWidth;
  const height = img.naturalHeight;
  const currentAspectRatio = width / height;

  // 清理 URL
  URL.revokeObjectURL(img.src);

  // 验证尺寸
  if (maxWidth && width > maxWidth) {
    errors.push(`图片宽度超过限制: ${width}px > ${maxWidth}px`);
  }

  if (maxHeight && height > maxHeight) {
    errors.push(`图片高度超过限制: ${height}px > ${maxHeight}px`);
  }

  if (minWidth && width < minWidth) {
    errors.push(`图片宽度不足: ${width}px < ${minWidth}px`);
  }

  if (minHeight && height < minHeight) {
    errors.push(`图片高度不足: ${height}px < ${minHeight}px`);
  }

  // 验证宽高比
  if (aspectRatio) {
    const diff = Math.abs(currentAspectRatio - aspectRatio);
    if (diff > aspectRatioTolerance) {
      errors.push(`宽高比不符合要求: ${currentAspectRatio.toFixed(2)} != ${aspectRatio.toFixed(2)}`);
    }
  }

  return {
    valid: errors.length === 0,
    errors,
    info: {
      width,
      height,
      type: file.type,
      size: file.size,
      aspectRatio: currentAspectRatio,
    },
  };
}

// ============================================================================
// useImageCompression - 图片压缩
// ============================================================================

export interface UseImageCompressionOptions {
  /** 最大宽度 */
  maxWidth?: number;
  /** 最大高度 */
  maxHeight?: number;
  /** 压缩质量 (0-1) */
  quality?: number;
  /** 输出类型 */
  type?: "image/jpeg" | "image/png" | "image/webp";
}

export interface UseImageCompressionReturn {
  /** 是否正在压缩 */
  isCompressing: Ref<boolean>;
  /** 压缩进度 */
  progress: Ref<number>;
  /** 错误信息 */
  error: Ref<Error | null>;
  /** 压缩图片 */
  compress: (file: File) => Promise<Blob>;
  /** 压缩多张图片 */
  compressMultiple: (files: File[]) => Promise<Blob[]>;
}

/**
 * 图片压缩
 * @example
 * const { compress, isCompressing } = useImageCompression({
 *   maxWidth: 1920,
 *   quality: 0.8,
 *   type: 'image/jpeg'
 * })
 *
 * const compressedBlob = await compress(file)
 */
export function useImageCompression(options: UseImageCompressionOptions = {}): UseImageCompressionReturn {
  const { maxWidth = 1920, maxHeight = 1080, quality = 0.8, type = "image/jpeg" } = options;

  const isCompressing = ref(false);
  const progress = ref(0);
  const error = ref<Error | null>(null);

  const compress = async (file: File): Promise<Blob> => {
    return new Promise((resolve, reject) => {
      const img = new Image();

      img.onload = () => {
        URL.revokeObjectURL(img.src);

        let { naturalWidth: width, naturalHeight: height } = img;

        // 计算缩放比例
        const ratio = Math.min(maxWidth / width, maxHeight / height, 1);
        width = Math.round(width * ratio);
        height = Math.round(height * ratio);

        // 创建 canvas
        const canvas = document.createElement("canvas");
        canvas.width = width;
        canvas.height = height;

        const ctx = canvas.getContext("2d");
        if (!ctx) {
          reject(new Error("Failed to get canvas context"));
          return;
        }

        // 绘制图片
        ctx.drawImage(img, 0, 0, width, height);

        // 导出
        canvas.toBlob(
          (blob) => {
            if (blob) {
              resolve(blob);
            } else {
              reject(new Error("Failed to compress image"));
            }
          },
          type,
          quality
        );
      };

      img.onerror = () => {
        URL.revokeObjectURL(img.src);
        reject(new Error("Failed to load image"));
      };

      img.src = URL.createObjectURL(file);
    });
  };

  const compressMultiple = async (files: File[]): Promise<Blob[]> => {
    isCompressing.value = true;
    progress.value = 0;
    error.value = null;

    const results: Blob[] = [];

    try {
      for (let i = 0; i < files.length; i++) {
        const blob = await compress(files[i]);
        results.push(blob);
        progress.value = (i + 1) / files.length;
      }
    } catch (err) {
      error.value = err instanceof Error ? err : new Error(String(err));
    }

    isCompressing.value = false;
    return results;
  };

  return {
    isCompressing,
    progress,
    error,
    compress,
    compressMultiple,
  };
}
