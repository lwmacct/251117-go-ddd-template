/**
 * Composables 统一导出
 * 提供所有可组合函数的统一入口
 */

// ============================================================================
// 异步与状态
// ============================================================================

export {
  useAsyncState,
  useAsyncRetry,
  usePolling,
  usePromiseQueue,
  type UseAsyncStateOptions,
  type UseAsyncStateReturn,
} from "./useAsync";

// ============================================================================
// 事件与交互
// ============================================================================

export {
  useClickOutside,
  useClickOutsideToggle,
  vClickOutside,
  type UseClickOutsideOptions,
} from "./useClickOutside";

export { useClipboard, type UseClipboardReturn } from "./useClipboard";

export {
  useConfirm,
  confirmDialog,
  type ConfirmOptions,
} from "./useConfirm";

export {
  useSortable,
  getDragItemClasses,
  useFileDrop,
  type UseDraggableOptions,
  type UseDraggableReturn,
  type FileDropOptions,
  type FileDropReturn,
} from "./useDraggable";

export {
  createEventBus,
  useEventBus,
  useEventListener,
  useEventValue,
  appEventBus,
  type EventBus,
} from "./useEventBus";

export {
  useFocusTrap,
  useFocusTrapWhenTrue,
  useFocusReturn,
  useAutoFocus,
  type UseFocusTrapOptions,
  type UseFocusTrapReturn,
} from "./useFocusTrap";

export {
  useKeyboard,
  type KeyboardOptions,
  type KeyboardShortcut,
} from "./useKeyboard";

// ============================================================================
// 表单与验证
// ============================================================================

export {
  useForm,
  useFormDirtyGuard,
  useFieldArray,
  type UseFormOptions,
  type UseFormReturn,
  type FormErrors,
  type FieldArrayReturn,
} from "./useForm";

// ============================================================================
// 时间与定时器
// ============================================================================

export {
  useCountdown,
  useStopwatch,
  useVerificationCode,
  useTargetDateCountdown,
  type UseCountdownOptions,
  type UseCountdownReturn,
  type UseStopwatchOptions,
  type UseStopwatchReturn,
  type UseVerificationCodeOptions,
  type UseVerificationCodeReturn,
  type UseTargetDateCountdownOptions,
  type UseTargetDateCountdownReturn,
} from "./useCountdown";

export {
  useDebounce,
  useDebouncedRef,
  useDebouncedFn,
  type UseDebounceOptions,
} from "./useDebounce";

// ============================================================================
// 视图与布局
// ============================================================================

export {
  useFullscreen,
  useDocumentFullscreen,
  useFullscreenButton,
  type UseFullscreenOptions,
  type UseFullscreenReturn,
} from "./useFullscreen";

export {
  useIntersectionObserver,
  useLazyLoad,
  useInfiniteScroll,
  useAnimateOnScroll,
  type UseIntersectionObserverOptions,
  type UseIntersectionObserverReturn,
  type UseLazyLoadOptions,
  type UseLazyLoadReturn,
  type UseInfiniteScrollOptions,
  type UseInfiniteScrollReturn,
} from "./useIntersectionObserver";

export {
  useScrollLock,
  useScrollLockWhenTrue,
  useElementScrollLock,
  useScrollPosition,
  useScrollDirection,
  type UseScrollLockOptions,
  type UseScrollLockReturn,
  type ScrollPosition,
} from "./useScrollLock";

export {
  useWindowSize,
  useMediaQuery,
  usePrefersDark,
  useElementSize,
  type UseWindowSizeOptions,
  type UseWindowSizeReturn,
  type UseMediaQueryOptions,
} from "./useWindowSize";

// ============================================================================
// 网络与存储
// ============================================================================

export {
  useNetwork,
  useOnline,
  useNetworkSpeed,
  useNetworkBanner,
  type UseNetworkOptions,
  type UseNetworkReturn,
} from "./useNetwork";

export {
  useStorage,
  useLocalStorage,
  useSessionStorage,
  type UseStorageOptions,
} from "./useStorage";

// ============================================================================
// 历史与状态管理
// ============================================================================

export {
  useHistory,
  useManualHistory,
  useTimestampedHistory,
  useSnapshot,
  type UseHistoryOptions,
  type UseHistoryReturn,
  type HistoryEntry,
} from "./useHistory";

// ============================================================================
// 定时器
// ============================================================================

export {
  useTimeout,
  useTimeoutFn,
  useInterval,
  useIntervalFn,
  useTimestamp,
  useNow,
  useRafFn,
  useDateFormat,
  useIdleCallback,
  useScheduler,
  type UseTimeoutOptions,
  type UseTimeoutReturn,
  type UseIntervalOptions,
  type UseIntervalReturn,
  type UseTimestampOptions,
  type UseTimestampReturn,
  type UseRafFnOptions,
  type UseRafFnReturn,
  type UseDateFormatOptions,
  type UseIdleCallbackOptions,
  type UseIdleCallbackReturn,
  type UseSchedulerReturn,
  type ScheduledTask,
} from "./useTimer";

// ============================================================================
// 权限
// ============================================================================

export {
  usePermission,
  useNotificationPermission,
  useClipboardPermission,
  useCameraPermission,
  useMicrophonePermission,
  useGeolocationPermission,
  useScreenWakeLock,
  isPermissionsApiSupported,
  queryPermissions,
  type PermissionName,
  type PermissionState,
  type UsePermissionOptions,
  type UsePermissionReturn,
  type UseNotificationPermissionReturn,
  type UseClipboardPermissionReturn,
  type UseCameraPermissionReturn,
  type UseMicrophonePermissionReturn,
  type UseGeolocationPermissionReturn,
  type UseScreenWakeLockReturn,
} from "./usePermission";

// ============================================================================
// 鼠标
// ============================================================================

export {
  useMouse,
  useMousePressed,
  useMouseInElement,
  useHover,
  useCursor,
  useDropZone,
  type MousePosition,
  type UseMouseOptions,
  type UseMouseReturn,
  type UseMousePressedOptions,
  type UseMousePressedReturn,
  type UseMouseInElementOptions,
  type UseMouseInElementReturn,
  type UseHoverOptions,
  type UseHoverReturn,
  type CursorType,
  type UseCursorReturn,
  type UseDropZoneOptions,
  type UseDropZoneReturn,
} from "./useMouse";

// ============================================================================
// 地理位置
// ============================================================================

export {
  useGeolocation,
  useGeolocationWatch,
  useGeolocationBounds,
  calculateDistance,
  formatDistance,
  dmsToDecimal,
  decimalToDms,
  formatCoordinates,
  getGoogleMapsUrl,
  getAppleMapsUrl,
  getNavigationUrl,
  DISTANCE_UNITS,
  type GeolocationCoordinates,
  type GeolocationPosition,
  type UseGeolocationOptions,
  type UseGeolocationReturn,
  type UseGeolocationWatchOptions,
  type GeolocationBounds,
  type UseGeolocationBoundsOptions,
  type DistanceUnit,
} from "./useGeolocation";

// ============================================================================
// 用户偏好
// ============================================================================

export {
  usePreferredDark,
  usePreferredLanguage,
  usePreferredReducedMotion,
  usePreferredContrast,
  usePreferredColorScheme,
  usePreferredTransparency,
  useDark,
  useColorMode,
  type UsePreferredDarkReturn,
  type UsePreferredLanguageReturn,
  type UsePreferredReducedMotionReturn,
  type UsePreferredContrastReturn,
  type UsePreferredColorSchemeReturn,
  type UsePreferredTransparencyReturn,
  type UseDarkOptions,
  type UseDarkReturn,
  type ColorMode,
  type UseColorModeOptions,
  type UseColorModeReturn,
} from "./usePreferences";

// ============================================================================
// 标题与文档
// ============================================================================

export {
  useTitle,
  useTitleTemplate,
  useDocumentTitle,
  useFavicon,
  usePageLeave,
  useDocumentVisibility,
  useHead,
  useScript,
  useStylesheet,
  type UseTitleOptions,
  type UseTitleReturn,
  type UseTitleTemplateOptions,
  type UseTitleTemplateReturn,
  type UseFaviconOptions,
  type UseFaviconReturn,
  type UsePageLeaveOptions,
  type UsePageLeaveReturn,
  type DocumentVisibilityState,
  type UseDocumentVisibilityReturn,
  type HeadConfig,
  type UseScriptOptions,
  type UseScriptReturn,
  type UseStylesheetOptions,
  type UseStylesheetReturn,
} from "./useTitle";

// ============================================================================
// WebSocket
// ============================================================================

export {
  useWebSocket,
  useWebSocketJSON,
  useWebSocketBinary,
  createWebSocketManager,
  type WebSocketStatus,
  type UseWebSocketOptions,
  type UseWebSocketReturn,
  type UseWebSocketJSONOptions,
  type UseWebSocketJSONReturn,
  type UseWebSocketBinaryOptions,
  type UseWebSocketBinaryReturn,
  type WebSocketManager,
} from "./useWebSocket";

// ============================================================================
// EventSource (SSE)
// ============================================================================

export {
  useEventSource,
  useEventSourceNamed,
  useServerSentEvents,
  createEventSourceManager,
  type EventSourceStatus,
  type UseEventSourceOptions,
  type UseEventSourceReturn,
  type UseEventSourceNamedOptions,
  type NamedEventData,
  type UseEventSourceNamedReturn,
  type EventSourceManager,
} from "./useEventSource";

// ============================================================================
// BroadcastChannel（跨标签页通信）
// ============================================================================

export {
  useBroadcastChannel,
  useBroadcastChannelJSON,
  useTabSync,
  useTabLeader,
  useTabMessenger,
  type UseBroadcastChannelOptions,
  type UseBroadcastChannelReturn,
  type UseBroadcastChannelJSONOptions,
  type UseTabSyncOptions,
  type UseTabSyncReturn,
  type UseTabLeaderOptions,
  type UseTabLeaderReturn,
  type UseTabMessengerReturn,
  type MessageHandler,
} from "./useBroadcastChannel";

// ============================================================================
// 图片处理
// ============================================================================

export {
  useImage,
  useImagePreload,
  useLazyImage,
  useProgressiveImage,
  validateImage,
  useImageCompression,
  type UseImageOptions,
  type UseImageReturn,
  type UseImagePreloadReturn,
  type UseLazyImageOptions,
  type UseLazyImageReturn,
  type UseProgressiveImageOptions,
  type UseProgressiveImageReturn,
  type ImageValidationResult,
  type UseImageValidationOptions,
  type UseImageCompressionOptions,
  type UseImageCompressionReturn,
} from "./useImage";

// ============================================================================
// 克隆与状态
// ============================================================================

export {
  deepClone,
  structuredClonePolyfill,
  useCloned,
  useManualClone,
  useDirtyState,
  useSnapshot,
  useSyncedRef,
  useMemoize,
  type UseClonedOptions,
  type UseClonedReturn,
  type UseManualCloneReturn,
  type UseDirtyStateReturn,
  type UseSnapshotReturn,
  type UseSyncedRefOptions,
  type UseSyncedRefReturn,
  type UseMemoizeReturn,
} from "./useClone";

// ============================================================================
// HTTP 请求
// ============================================================================

export {
  useFetch,
  createFetch,
  useLazyFetch,
  usePaginatedFetch,
  useInfiniteFetch,
  clearFetchCache,
  deleteFetchCache,
  getFetchCacheSize,
  type FetchStatus,
  type UseFetchOptions,
  type UseFetchReturn,
  type CreateFetchOptions,
  type CreateFetchReturn,
  type UsePaginatedFetchOptions,
  type UsePaginatedFetchReturn,
  type UseInfiniteFetchOptions,
  type UseInfiniteFetchReturn,
} from "./useFetch";

// ============================================================================
// v-model 绑定
// ============================================================================

export {
  useVModel,
  useVModels,
  useModelValue,
  useProxyModel,
  useControlled,
  useDebouncedVModel,
  useThrottledVModel,
  useToggle,
  useCycleList,
  type UseVModelOptions,
  type UseProxyModelOptions,
  type UseProxyModelReturn,
  type UseControlledOptions,
  type UseControlledReturn,
  type UseDebouncedVModelOptions,
  type UseThrottledVModelOptions,
  type UseToggleOptions,
  type UseToggleReturn,
  type UseCycleListReturn,
} from "./useVModel";

// ============================================================================
// SWR (Stale-While-Revalidate)
// ============================================================================

export {
  useSWR,
  useSWRMutation,
  useSWRInfinite,
  clearSWRCache,
  deleteSWRCache,
  getSWRCache,
  setSWRCache,
  revalidateSWR,
  type SWRStatus,
  type UseSWROptions,
  type UseSWRReturn,
  type UseSWRMutationOptions,
  type UseSWRMutationReturn,
  type UseSWRInfiniteOptions,
  type UseSWRInfiniteReturn,
} from "./useSWR";

// ============================================================================
// Watch 工具
// ============================================================================

export {
  watchOnce,
  watchDebounced,
  watchThrottled,
  watchPausable,
  watchIgnorable,
  watchTriggered,
  watchArray,
  watchWithFilter,
  watchAtMost,
  whenever,
  until,
  useWatchArray,
  type WatchOnceOptions,
  type WatchDebouncedOptions,
  type WatchThrottledOptions,
  type WatchPausableReturn,
  type WatchIgnorableReturn,
  type WatchTriggeredReturn,
  type WatchArrayReturn,
  type WatchAtMostReturn,
  type UntilReturn,
} from "./useWatch";

// ============================================================================
// Computed 工具
// ============================================================================

export {
  computedEager,
  computedAsync,
  computedWithControl,
  computedPick,
  computedOmit,
  toComputed,
  computedFrom,
  computedDebounced,
  computedThrottled,
  computedWithHistory,
  computedIf,
  writableComputed,
  computedArray,
  computedDelayed,
  computedDefault,
  computedObject,
  computedMap,
  computedFilter,
  computedFind,
  computedGroupBy,
  computedSort,
  type ComputedAsyncOptions,
  type ComputedAsyncReturn,
  type ComputedWithControlOptions,
  type ComputedWithControlReturn,
  type ComputedDebouncedOptions,
  type ComputedThrottledOptions,
} from "./useComputed";

// ============================================================================
// Ref 工具
// ============================================================================

export {
  refDefault,
  refDebounced,
  refThrottled,
  refHistory,
  refAutoReset,
  syncRefs,
  refWithControl,
  templateRef,
  usePrevious,
  useLatest,
  refLocked,
  useCounter,
  useBoolean,
  useObject,
  useArray,
  useSet,
  useMap,
  type RefDebouncedOptions,
  type RefThrottledOptions,
  type RefHistoryOptions,
  type RefHistoryReturn,
  type RefAutoResetOptions,
  type RefWithControlOptions,
  type RefWithControlReturn,
  type TemplateRefReturn,
} from "./useRef";

// ============================================================================
// Reactive 工具
// ============================================================================

export {
  reactiveWithOptions,
  reactiveHistory,
  reactiveForm,
  reactivePick,
  reactiveOmit,
  reactiveMerge,
  reactiveDefault,
  reactiveExtend,
  syncReactive,
  reactiveWhen,
  reactiveTransform,
  reactiveValidated,
  reactiveResettable,
  reactiveWithReadonly,
  watchReactive,
  reactiveUtils,
  type ReactiveWithOptionsConfig,
  type ReactiveHistoryOptions,
  type ReactiveHistoryReturn,
  type ReactiveFormReturn,
} from "./useReactive";

// ============================================================================
// Lifecycle 工具
// ============================================================================

export {
  useMountedState,
  useSafeMounted,
  useMounted,
  useAsyncMounted,
  useUnmounted,
  useUpdated,
  useActivated,
  useDeactivated,
  useLifecycleTracker,
  useCleanup,
  useMountedDelay,
  useMountedNextFrame,
  useMountedNextTick,
  useMountedWhen,
  useInstance,
  useRenderCount,
  useErrorCapture,
  useRenderTracking,
  useMountedOnce,
  useMountedOrActivated,
  useComponentVisible,
  useComponentAliveTime,
  useLifecycle,
  type MountedState,
  type LifecycleTrackerOptions,
  type LifecycleEvent,
  type LifecycleTrackerReturn,
  type AsyncMountedOptions,
  type CleanupFn,
} from "./useLifecycle";

// ============================================================================
// Provide/Inject 工具
// ============================================================================

export {
  createContext,
  createStateContext,
  createReactiveContext,
  createReadonlyContext,
  createEventBusContext,
  createThemeContext,
  createI18nContext,
  createStorageContext,
  createFactoryContext,
  useOptionalInject,
  useRequiredInject,
  type CreateContextReturn,
  type CreateStateContextReturn,
  type CreateReactiveContextReturn,
  type EventBusContext,
  type ThemeContext,
  type I18nContext,
} from "./useProvide";

// ============================================================================
// Component 工具
// ============================================================================

export {
  useSlotsInfo,
  useConditionalSlot,
  useSlotPass,
  useAttrsEnhanced,
  useClassMerge,
  useStyleMerge,
  useComponentInfo,
  useParentComponent,
  useRootComponent,
  useComponentRef,
  useComponentRefs,
  useExposeMethod,
  useExposeState,
  useForceUpdate,
  useComponentEmit,
  useComponentProxy,
  useComponentType,
  useCustomProperties,
  type SlotInfo,
  type ComponentInfo,
  type ComponentRefReturn,
  type ExposeOptions,
} from "./useComponent";

// ============================================================================
// Transition 工具
// ============================================================================

export {
  useTransition,
  useFade,
  useSlide,
  useScale,
  useAnimation,
  useTransitionGroup,
  useNumberTransition,
  useShake,
  usePulse,
  useTypewriter,
  type TransitionState,
  type EasingFunction,
  type TransitionConfig,
  type UseTransitionReturn,
  type AnimationConfig,
  type UseAnimationReturn,
  type FadeConfig,
  type SlideConfig,
} from "./useTransition";

// ============================================================================
// Queue 工具
// ============================================================================

export {
  useQueue,
  useStack,
  useTaskQueue,
  useNotificationQueue,
  useHistoryQueue,
  useRingBuffer,
  type QueueItem,
  type QueueConfig,
  type UseQueueReturn,
  type TaskQueueConfig,
  type TaskStatus,
  type Task,
  type UseTaskQueueReturn,
  type NotificationQueueConfig,
  type Notification,
  type UseNotificationQueueReturn,
} from "./useQueue";

// ============================================================================
// Media 工具
// ============================================================================

export {
  useMedia,
  useAudio,
  useAudioVisualizer,
  useRecorder,
  useScreenShare,
  usePictureInPicture,
  formatDuration,
  type MediaState,
  type UseMediaOptions,
  type UseMediaReturn,
  type AudioVisualizerOptions,
  type AudioVisualizerReturn,
  type RecorderOptions,
  type UseRecorderReturn,
} from "./useMedia";
