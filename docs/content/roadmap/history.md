# 历史记录工具 Composable

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

<!--TOC-->

## Table of Contents

- [需求背景](#需求背景) `:29+4`
- [已实现功能](#已实现功能) `:33+24`
  - [useHistory](#usehistory) `:35+7`
  - [useManualHistory](#usemanualhistory) `:42+5`
  - [useTimestampedHistory](#usetimestampedhistory) `:47+5`
  - [useSnapshot](#usesnapshot) `:52+5`
- [使用方式](#使用方式) `:57+78`
  - [自动历史记录](#自动历史记录) `:59+24`
  - [配置选项](#配置选项) `:83+10`
  - [手动历史记录](#手动历史记录) `:93+18`
  - [状态快照](#状态快照) `:111+24`
- [API](#api) `:135+28`
  - [useHistory 返回值](#usehistory-返回值) `:137+15`
  - [useSnapshot 返回值](#usesnapshot-返回值) `:152+11`
- [代码位置](#代码位置) `:163+7`

<!--TOC-->

## 需求背景

需要为编辑器、表单等场景提供撤销/重做功能。

## 已实现功能

### useHistory

- 自动历史记录
- 撤销/重做
- 容量限制
- 深拷贝支持

### useManualHistory

- 手动提交历史
- 精确控制记录点

### useTimestampedHistory

- 带时间戳的历史
- 记录变更时间

### useSnapshot

- 命名快照
- 保存/恢复状态

## 使用方式

### 自动历史记录

```typescript
import { ref } from "vue";
import { useHistory } from "@/composables/useHistory";

const state = ref({ text: "", count: 0 });
const { current, undo, redo, canUndo, canRedo, historyCount } = useHistory(state);

// 修改值会自动记录历史
current.value.text = "Hello";
current.value.count++;

// 撤销
if (canUndo.value) {
  undo();
}

// 重做
if (canRedo.value) {
  redo();
}
```

### 配置选项

```typescript
const { undo, redo } = useHistory(state, {
  capacity: 50, // 最多保留 50 条历史
  deep: true, // 深度监听
  immediate: true, // 立即记录初始状态
});
```

### 手动历史记录

```typescript
import { useManualHistory } from "@/composables/useHistory";

const { current, commit, undo, redo } = useManualHistory({
  name: "",
  items: [],
});

// 修改多个字段
current.value.name = "New Name";
current.value.items.push({ id: 1 });

// 手动提交（作为一个操作记录）
commit();
```

### 状态快照

```typescript
import { useSnapshot } from "@/composables/useHistory";

const state = ref({
  /* ... */
});
const { save, restore, list, has, remove } = useSnapshot(state);

// 保存当前状态
save("before-edit");

// 做一些修改...

// 恢复到保存点
if (has("before-edit")) {
  restore("before-edit");
}

// 列出所有快照
console.log(list.value); // ["before-edit"]
```

## API

### useHistory 返回值

| 属性         | 类型                   | 说明           |
| ------------ | ---------------------- | -------------- |
| current      | `Ref<T>`               | 当前值         |
| canUndo      | `ComputedRef<boolean>` | 是否可撤销     |
| canRedo      | `ComputedRef<boolean>` | 是否可重做     |
| historyCount | `ComputedRef<number>`  | 历史数量       |
| undo         | `() => void`           | 撤销           |
| redo         | `() => void`           | 重做           |
| clear        | `() => void`           | 清除历史       |
| commit       | `() => void`           | 手动提交       |
| reset        | `() => void`           | 重置到初始状态 |
| go           | `(delta) => void`      | 跳转           |

### useSnapshot 返回值

| 方法    | 类型                    | 说明     |
| ------- | ----------------------- | -------- |
| save    | `(name) => void`        | 保存快照 |
| restore | `(name) => boolean`     | 恢复快照 |
| remove  | `(name) => boolean`     | 删除快照 |
| has     | `(name) => boolean`     | 检查快照 |
| list    | `ComputedRef<string[]>` | 快照列表 |
| clear   | `() => void`            | 清除所有 |

## 代码位置

```
web/src/
└── composables/
    └── useHistory.ts    # 历史记录工具
```
