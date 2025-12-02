# Geolocation Composable

<!--TOC-->

- [需求背景](#需求背景) `:30:33`
- [已实现功能](#已实现功能) `:34:35`
  - [位置获取](#位置获取) `:36:41`
  - [距离计算](#距离计算) `:42:46`
  - [坐标工具](#坐标工具) `:47:52`
  - [地图链接](#地图链接) `:53:58`
- [使用方式](#使用方式) `:59:60`
  - [基础用法](#基础用法) `:61:79`
  - [持续监听](#持续监听) `:80:95`
  - [边界检测](#边界检测) `:96:112`
  - [距离计算](#距离计算-1) `:113:130`
  - [坐标转换](#坐标转换) `:131:147`
  - [生成地图链接](#生成地图链接) `:148:161`
- [API](#api) `:162:163`
  - [useGeolocation](#usegeolocation) `:164:183`
  - [GeolocationCoordinates](#geolocationcoordinates) `:184:195`
  - [calculateDistance](#calculatedistance) `:196:203`
- [代码位置](#代码位置) `:204:210`

<!--TOC-->

> **状态**: ✅ 已完成
> **优先级**: 中
> **完成日期**: 2024-11-30

## 需求背景

前端需要获取和跟踪用户地理位置，用于地图、导航、位置服务等场景。

## 已实现功能

### 位置获取

- `useGeolocation` - 获取和跟踪地理位置
- `useGeolocationWatch` - 持续监听位置变化
- `useGeolocationBounds` - 边界检测

### 距离计算

- `calculateDistance` - 计算两点距离
- `formatDistance` - 格式化距离

### 坐标工具

- `dmsToDecimal` - 度分秒转十进制度
- `decimalToDms` - 十进制度转度分秒
- `formatCoordinates` - 格式化坐标

### 地图链接

- `getGoogleMapsUrl` - Google Maps 链接
- `getAppleMapsUrl` - Apple Maps 链接
- `getNavigationUrl` - 跨平台导航链接

## 使用方式

### 基础用法

```typescript
import { useGeolocation } from "@/composables/useGeolocation";

const { coords, isLoading, error, isSupported, getCurrentPosition } = useGeolocation();

// 检查支持
if (!isSupported.value) {
  console.log("浏览器不支持地理位置");
}

// 使用响应式坐标
console.log(coords.value?.latitude, coords.value?.longitude);

// 手动获取位置
const position = await getCurrentPosition();
```

### 持续监听

```typescript
import { useGeolocationWatch } from "@/composables/useGeolocation";

useGeolocationWatch({
  onPositionChange: (coords) => {
    updateUserLocation(coords.latitude, coords.longitude);
  },
  onError: (error) => {
    console.error("定位失败:", error.message);
  },
  minDistance: 10, // 移动超过 10 米才触发
});
```

### 边界检测

```typescript
import { useGeolocationBounds } from "@/composables/useGeolocation";

const { isInBounds } = useGeolocationBounds({
  bounds: {
    north: 40.8,
    south: 40.6,
    east: -73.9,
    west: -74.1,
  },
  onEnter: () => showNotification("您已进入配送区域"),
  onLeave: () => showNotification("您已离开配送区域"),
});
```

### 距离计算

```typescript
import { calculateDistance, formatDistance } from "@/composables/useGeolocation";

// 计算两点距离
const distance = calculateDistance(
  { latitude: 40.7128, longitude: -74.006 }, // 纽约
  { latitude: 34.0522, longitude: -118.2437 }, // 洛杉矶
  "km",
);
// 约 3936 公里

// 格式化距离
formatDistance(1500); // '1.5 km'
formatDistance(500); // '500 m'
```

### 坐标转换

```typescript
import { dmsToDecimal, decimalToDms, formatCoordinates } from "@/composables/useGeolocation";

// 度分秒转十进制
const lat = dmsToDecimal(40, 42, 46, "N"); // 40.7128

// 十进制转度分秒
const dms = decimalToDms(40.7128, "lat");
// { degrees: 40, minutes: 42, seconds: 46, direction: 'N' }

// 格式化坐标
formatCoordinates(40.7128, -74.006);
// '40.7128°N, 74.0060°W'
```

### 生成地图链接

```typescript
import { getGoogleMapsUrl, getNavigationUrl } from "@/composables/useGeolocation";

// Google Maps 链接
const googleUrl = getGoogleMapsUrl(40.7128, -74.006);
// 'https://www.google.com/maps?q=40.7128,-74.006&z=15'

// 跨平台导航链接
const navUrl = getNavigationUrl(40.7128, -74.006);
// iOS 返回 Apple Maps，其他返回 Google Maps
```

## API

### useGeolocation

| 选项               | 类型    | 默认值 | 说明             |
| ------------------ | ------- | ------ | ---------------- |
| immediate          | boolean | true   | 是否立即获取位置 |
| enableHighAccuracy | boolean | true   | 是否启用高精度   |
| maximumAge         | number  | 30000  | 缓存最大时间(ms) |
| timeout            | number  | 27000  | 超时时间(ms)     |

| 返回值             | 类型                            | 说明         |
| ------------------ | ------------------------------- | ------------ |
| coords             | `Ref<GeolocationCoordinates>  ` | 当前坐标     |
| locatedAt          | `Ref<number>                  ` | 位置时间戳   |
| isLoading          | `Ref<boolean>                 ` | 是否正在定位 |
| error              | `Ref<GeolocationPositionError>` | 错误信息     |
| isSupported        | `Ref<boolean>                 ` | 是否支持     |
| getCurrentPosition | `() => Promise                ` | 获取当前位置 |
| pause              | `() => void                   ` | 暂停监听     |
| resume             | `() => void                   ` | 恢复监听     |

### GeolocationCoordinates

| 属性             | 类型         | 说明          |
| ---------------- | ------------ | ------------- |
| latitude         | number       | 纬度          |
| longitude        | number       | 经度          |
| altitude         | number\|null | 海拔(米)      |
| accuracy         | number       | 精度(米)      |
| altitudeAccuracy | number\|null | 海拔精度(米)  |
| heading          | number\|null | 方向(0-360度) |
| speed            | number\|null | 速度(米/秒)   |

### calculateDistance

| 参数 | 类型                  | 说明         |
| ---- | --------------------- | ------------ |
| from | {latitude, longitude} | 起点坐标     |
| to   | {latitude, longitude} | 终点坐标     |
| unit | 'm'\|'km'\|'mi'\|'ft' | 单位，默认米 |

## 代码位置

```
web/src/
└── composables/
    └── useGeolocation.ts    # Geolocation Composable
```
