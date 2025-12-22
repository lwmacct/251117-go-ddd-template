# 拖拽排序 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:27+4`
- [已实现功能](#已实现功能) `:31+20`
  - [useSortable](#usesortable) `:33+7`
  - [useFileDrop](#usefiledrop) `:40+7`
  - [getDragItemClasses](#getdragitemclasses) `:47+4`
- [使用方式](#使用方式) `:51+64`
  - [列表拖拽排序](#列表拖拽排序) `:53+34`
  - [文件拖放](#文件拖放) `:87+28`
- [API](#api) `:115+31`
  - [useSortable 选项](#usesortable-选项) `:117+9`
  - [useSortable 返回值](#usesortable-返回值) `:126+10`
  - [useFileDrop 选项](#usefiledrop-选项) `:136+10`
- [代码位置](#代码位置) `:146+7`

<!--TOC-->

## 需求背景

需要为列表提供拖拽排序功能，以及文件拖放上传功能。

## 已实现功能

### useSortable

- 列表拖拽排序
- 拖拽开始/结束回调
- 排序变化回调
- 可禁用拖拽

### useFileDrop

- 文件拖放上传
- 文件类型验证
- 文件大小验证
- 单/多文件支持

### getDragItemClasses

- 拖拽状态样式类生成

## 使用方式

### 列表拖拽排序

```typescript
import { ref } from "vue";
import { useSortable, getDragItemClasses } from "@/composables/useDraggable";

const items = ref(["Item 1", "Item 2", "Item 3"]);

const { isDragging, dragIndex, overIndex, getDragItemProps } = useSortable(items, {
  onSort: (newItems) => {
    console.log("排序后:", newItems);
  },
});
```

```vue
<template>
  <div class="sortable-list">
    <div v-for="(item, index) in items" :key="item" v-bind="getDragItemProps(index)" :class="getDragItemClasses(index, dragIndex, overIndex)">
      {{ item }}
    </div>
  </div>
</template>

<style scoped>
.is-dragging {
  opacity: 0.5;
}
.is-over {
  border-color: #1976d2;
}
</style>
```

### 文件拖放

```typescript
import { useFileDrop } from "@/composables/useDraggable";

const { isDragOver, files, getDropZoneProps, clearFiles } = useFileDrop({
  accept: ["image/*", ".pdf"],
  maxSize: 10 * 1024 * 1024, // 10MB
  onDrop: (files) => {
    uploadFiles(files);
  },
  onError: (error) => {
    toast.error(error);
  },
});
```

```vue
<template>
  <div class="drop-zone" :class="{ 'is-drag-over': isDragOver }" v-bind="getDropZoneProps()">
    <p v-if="files.length === 0">拖放文件到此处，或点击选择</p>
    <ul v-else>
      <li v-for="file in files" :key="file.name">{{ file.name }} ({{ formatFileSize(file.size) }})</li>
    </ul>
  </div>
</template>
```

## API

### useSortable 选项

| 选项        | 类型     | 说明         |
| ----------- | -------- | ------------ |
| onDragStart | Function | 拖拽开始回调 |
| onDragEnd   | Function | 拖拽结束回调 |
| onSort      | Function | 排序变化回调 |
| disabled    | boolean  | 是否禁用     |

### useSortable 返回值

| 属性             | 类型                | 说明           |
| ---------------- | ------------------- | -------------- |
| isDragging       | `Ref<boolean>`      | 是否正在拖拽   |
| dragIndex        | Ref<number \| null> | 拖拽项索引     |
| overIndex        | Ref<number \| null> | 悬停项索引     |
| getDragItemProps | Function            | 获取拖拽项属性 |
| moveItem         | Function            | 移动项         |

### useFileDrop 选项

| 选项     | 类型     | 说明           |
| -------- | -------- | -------------- |
| accept   | string[] | 接受的文件类型 |
| multiple | boolean  | 是否多文件     |
| maxSize  | number   | 最大文件大小   |
| onDrop   | Function | 拖放回调       |
| onError  | Function | 错误回调       |

## 代码位置

```
web/src/
└── composables/
    └── useDraggable.ts    # 拖拽排序
```
