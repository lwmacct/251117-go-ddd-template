/**
 * 拖拽排序 Composable
 * 提供列表拖拽排序功能
 */

import { ref, computed, type Ref } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseSortableOptions<T> {
  /** 拖拽开始回调 */
  onDragStart?: (item: T, index: number) => void;
  /** 拖拽结束回调 */
  onDragEnd?: (item: T, oldIndex: number, newIndex: number) => void;
  /** 排序变化回调 */
  onSort?: (items: T[]) => void;
  /** 是否禁用拖拽 */
  disabled?: boolean;
  /** 拖拽手柄选择器 */
  handle?: string;
}

export interface UseDraggableReturn<T> {
  /** 是否正在拖拽 */
  isDragging: Ref<boolean>;
  /** 当前拖拽的项索引 */
  dragIndex: Ref<number | null>;
  /** 当前悬停的项索引 */
  overIndex: Ref<number | null>;
  /** 获取拖拽项的属性 */
  getDragItemProps: (index: number) => DragItemProps;
  /** 获取拖拽区域的属性 */
  getDragContainerProps: () => DragContainerProps;
  /** 移动项 */
  moveItem: (fromIndex: number, toIndex: number) => void;
  /** 重置拖拽状态 */
  reset: () => void;
}

export interface DragItemProps {
  draggable: boolean;
  "data-index": number;
  onDragstart: (e: DragEvent) => void;
  onDragend: (e: DragEvent) => void;
  onDragover: (e: DragEvent) => void;
  onDragenter: (e: DragEvent) => void;
  onDragleave: (e: DragEvent) => void;
  onDrop: (e: DragEvent) => void;
}

export interface DragContainerProps {
  onDragover: (e: DragEvent) => void;
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 列表拖拽排序
 * @example
 * const items = ref(['Item 1', 'Item 2', 'Item 3'])
 * const { getDragItemProps, isDragging } = useSortable(items)
 *
 * // 模板中
 * // <div
 * //   v-for="(item, index) in items"
 * //   :key="item"
 * //   v-bind="getDragItemProps(index)"
 * // >
 * //   {{ item }}
 * // </div>
 */
export function useSortable<T>(
  items: Ref<T[]>,
  options: UseSortableOptions<T> = {}
): UseDraggableReturn<T> {
  const { onDragStart, onDragEnd, onSort, disabled = false } = options;

  const isDragging = ref(false);
  const dragIndex = ref<number | null>(null);
  const overIndex = ref<number | null>(null);

  // 移动项
  const moveItem = (fromIndex: number, toIndex: number) => {
    if (fromIndex === toIndex) return;
    if (fromIndex < 0 || toIndex < 0) return;
    if (fromIndex >= items.value.length || toIndex >= items.value.length) return;

    const newItems = [...items.value];
    const [movedItem] = newItems.splice(fromIndex, 1);
    newItems.splice(toIndex, 0, movedItem);
    items.value = newItems;

    onSort?.(newItems);
  };

  // 重置状态
  const reset = () => {
    isDragging.value = false;
    dragIndex.value = null;
    overIndex.value = null;
  };

  // 获取拖拽项属性
  const getDragItemProps = (index: number): DragItemProps => ({
    draggable: !disabled,
    "data-index": index,

    onDragstart: (e: DragEvent) => {
      if (disabled) return;

      isDragging.value = true;
      dragIndex.value = index;

      // 设置拖拽数据
      if (e.dataTransfer) {
        e.dataTransfer.effectAllowed = "move";
        e.dataTransfer.setData("text/plain", String(index));
      }

      onDragStart?.(items.value[index], index);
    },

    onDragend: (e: DragEvent) => {
      if (disabled) return;

      const oldIndex = dragIndex.value;
      const newIndex = overIndex.value;

      if (
        oldIndex !== null &&
        newIndex !== null &&
        oldIndex !== newIndex
      ) {
        onDragEnd?.(items.value[oldIndex], oldIndex, newIndex);
      }

      reset();
    },

    onDragover: (e: DragEvent) => {
      if (disabled) return;
      e.preventDefault();

      if (e.dataTransfer) {
        e.dataTransfer.dropEffect = "move";
      }

      overIndex.value = index;
    },

    onDragenter: (e: DragEvent) => {
      if (disabled) return;
      e.preventDefault();
      overIndex.value = index;
    },

    onDragleave: (e: DragEvent) => {
      if (disabled) return;
      // 检查是否真的离开了元素
      const target = e.currentTarget as HTMLElement;
      const related = e.relatedTarget as HTMLElement;

      if (!target.contains(related)) {
        if (overIndex.value === index) {
          overIndex.value = null;
        }
      }
    },

    onDrop: (e: DragEvent) => {
      if (disabled) return;
      e.preventDefault();

      const fromIndex = dragIndex.value;
      const toIndex = index;

      if (fromIndex !== null && fromIndex !== toIndex) {
        moveItem(fromIndex, toIndex);
      }
    },
  });

  // 获取容器属性
  const getDragContainerProps = (): DragContainerProps => ({
    onDragover: (e: DragEvent) => {
      if (disabled) return;
      e.preventDefault();
    },
  });

  return {
    isDragging,
    dragIndex,
    overIndex,
    getDragItemProps,
    getDragContainerProps,
    moveItem,
    reset,
  };
}

// ============================================================================
// 拖拽状态样式
// ============================================================================

/**
 * 获取拖拽项的样式类
 * @example
 * const classes = getDragItemClasses(index, dragIndex.value, overIndex.value)
 */
export function getDragItemClasses(
  index: number,
  dragIndex: number | null,
  overIndex: number | null
): Record<string, boolean> {
  return {
    "is-dragging": dragIndex === index,
    "is-over": overIndex === index && dragIndex !== index,
    "is-above": overIndex === index && dragIndex !== null && dragIndex > index,
    "is-below": overIndex === index && dragIndex !== null && dragIndex < index,
  };
}

// ============================================================================
// 文件拖放
// ============================================================================

export interface UseFileDropOptions {
  /** 接受的文件类型 */
  accept?: string[];
  /** 是否允许多文件 */
  multiple?: boolean;
  /** 最大文件大小（字节） */
  maxSize?: number;
  /** 文件拖入回调 */
  onDrop?: (files: File[]) => void;
  /** 错误回调 */
  onError?: (error: string) => void;
}

/**
 * 文件拖放
 * @example
 * const dropZone = ref<HTMLElement>()
 * const { isDragOver, files, getDropZoneProps } = useFileDrop({
 *   accept: ['image/*'],
 *   onDrop: (files) => handleFiles(files)
 * })
 */
export function useFileDrop(options: UseFileDropOptions = {}) {
  const {
    accept = [],
    multiple = true,
    maxSize,
    onDrop,
    onError,
  } = options;

  const isDragOver = ref(false);
  const files = ref<File[]>([]);

  // 验证文件类型
  const validateFileType = (file: File): boolean => {
    if (accept.length === 0) return true;

    return accept.some((type) => {
      if (type.endsWith("/*")) {
        // 通配符类型，如 "image/*"
        const baseType = type.slice(0, -2);
        return file.type.startsWith(baseType);
      }
      return file.type === type || file.name.endsWith(type);
    });
  };

  // 验证文件大小
  const validateFileSize = (file: File): boolean => {
    if (!maxSize) return true;
    return file.size <= maxSize;
  };

  // 处理拖放的文件
  const handleFiles = (droppedFiles: FileList | File[]) => {
    const validFiles: File[] = [];
    const fileArray = Array.from(droppedFiles);

    for (const file of fileArray) {
      if (!validateFileType(file)) {
        onError?.(`文件类型不支持: ${file.name}`);
        continue;
      }

      if (!validateFileSize(file)) {
        onError?.(`文件太大: ${file.name}`);
        continue;
      }

      validFiles.push(file);

      if (!multiple) break;
    }

    if (validFiles.length > 0) {
      files.value = multiple ? validFiles : [validFiles[0]];
      onDrop?.(files.value);
    }
  };

  // 获取拖放区域属性
  const getDropZoneProps = () => ({
    onDragenter: (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      isDragOver.value = true;
    },

    onDragover: (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      isDragOver.value = true;

      if (e.dataTransfer) {
        e.dataTransfer.dropEffect = "copy";
      }
    },

    onDragleave: (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();

      const target = e.currentTarget as HTMLElement;
      const related = e.relatedTarget as HTMLElement;

      if (!target.contains(related)) {
        isDragOver.value = false;
      }
    },

    onDrop: (e: DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      isDragOver.value = false;

      if (e.dataTransfer?.files) {
        handleFiles(e.dataTransfer.files);
      }
    },
  });

  // 清空文件
  const clearFiles = () => {
    files.value = [];
  };

  return {
    isDragOver,
    files,
    getDropZoneProps,
    handleFiles,
    clearFiles,
  };
}
