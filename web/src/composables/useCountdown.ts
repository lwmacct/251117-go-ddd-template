/**
 * 倒计时 Composable
 * 提供倒计时、计时器等时间相关功能
 */

import { ref, computed, onUnmounted, type Ref, type ComputedRef } from "vue";

// ============================================================================
// 类型定义
// ============================================================================

export interface UseCountdownOptions {
  /** 初始秒数 */
  seconds: number;
  /** 是否立即开始，默认 false */
  immediate?: boolean;
  /** 倒计时结束回调 */
  onEnd?: () => void;
  /** 每秒回调 */
  onTick?: (remaining: number) => void;
  /** 间隔毫秒数，默认 1000 */
  interval?: number;
}

export interface UseCountdownReturn {
  /** 剩余秒数 */
  remaining: Ref<number>;
  /** 是否正在运行 */
  isRunning: Ref<boolean>;
  /** 是否已结束 */
  isFinished: Ref<boolean>;
  /** 格式化时间 */
  formatted: ComputedRef<string>;
  /** 时间部分 */
  days: ComputedRef<number>;
  hours: ComputedRef<number>;
  minutes: ComputedRef<number>;
  seconds: ComputedRef<number>;
  /** 开始 */
  start: () => void;
  /** 暂停 */
  pause: () => void;
  /** 停止（重置） */
  stop: () => void;
  /** 重置 */
  reset: (newSeconds?: number) => void;
  /** 重新开始 */
  restart: () => void;
}

// ============================================================================
// 主函数
// ============================================================================

/**
 * 倒计时
 * @example
 * const { remaining, formatted, start, pause, reset } = useCountdown({
 *   seconds: 60,
 *   onEnd: () => toast.info('倒计时结束')
 * })
 *
 * start() // 开始倒计时
 */
export function useCountdown(options: UseCountdownOptions): UseCountdownReturn {
  const { seconds: initialSeconds, immediate = false, onEnd, onTick, interval = 1000 } = options;

  const remaining = ref(initialSeconds);
  const isRunning = ref(false);
  const isFinished = ref(false);

  let timer: ReturnType<typeof setInterval> | null = null;

  // 计算时间部分
  const days = computed(() => Math.floor(remaining.value / 86400));
  const hours = computed(() => Math.floor((remaining.value % 86400) / 3600));
  const minutes = computed(() => Math.floor((remaining.value % 3600) / 60));
  const seconds = computed(() => Math.floor(remaining.value % 60));

  // 格式化时间
  const formatted = computed(() => {
    const parts: string[] = [];

    if (days.value > 0) {
      parts.push(`${days.value}天`);
    }

    if (hours.value > 0 || days.value > 0) {
      parts.push(`${String(hours.value).padStart(2, "0")}:`);
    }

    parts.push(`${String(minutes.value).padStart(2, "0")}:`);
    parts.push(String(seconds.value).padStart(2, "0"));

    return parts.join("");
  });

  // 清除定时器
  const clearTimer = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  // 开始倒计时
  const start = () => {
    if (isRunning.value || remaining.value <= 0) return;

    isRunning.value = true;
    isFinished.value = false;

    timer = setInterval(() => {
      remaining.value--;
      onTick?.(remaining.value);

      if (remaining.value <= 0) {
        clearTimer();
        isRunning.value = false;
        isFinished.value = true;
        onEnd?.();
      }
    }, interval);
  };

  // 暂停
  const pause = () => {
    clearTimer();
    isRunning.value = false;
  };

  // 停止（重置到初始值）
  const stop = () => {
    clearTimer();
    isRunning.value = false;
    isFinished.value = false;
    remaining.value = initialSeconds;
  };

  // 重置
  const reset = (newSeconds?: number) => {
    clearTimer();
    isRunning.value = false;
    isFinished.value = false;
    remaining.value = newSeconds ?? initialSeconds;
  };

  // 重新开始
  const restart = () => {
    reset();
    start();
  };

  // 立即开始
  if (immediate) {
    start();
  }

  // 清理
  onUnmounted(() => {
    clearTimer();
  });

  return {
    remaining,
    isRunning,
    isFinished,
    formatted,
    days,
    hours,
    minutes,
    seconds,
    start,
    pause,
    stop,
    reset,
    restart,
  };
}

// ============================================================================
// 计时器（正计时）
// ============================================================================

export interface UseStopwatchOptions {
  /** 是否立即开始，默认 false */
  immediate?: boolean;
  /** 间隔毫秒数，默认 1000 */
  interval?: number;
  /** 每秒回调 */
  onTick?: (elapsed: number) => void;
}

/**
 * 计时器（正计时）
 * @example
 * const { elapsed, formatted, start, pause, reset } = useStopwatch()
 * start() // 开始计时
 */
export function useStopwatch(options: UseStopwatchOptions = {}) {
  const { immediate = false, interval = 1000, onTick } = options;

  const elapsed = ref(0);
  const isRunning = ref(false);

  let timer: ReturnType<typeof setInterval> | null = null;

  // 计算时间部分
  const hours = computed(() => Math.floor(elapsed.value / 3600));
  const minutes = computed(() => Math.floor((elapsed.value % 3600) / 60));
  const seconds = computed(() => Math.floor(elapsed.value % 60));

  // 格式化时间
  const formatted = computed(() => {
    const h = String(hours.value).padStart(2, "0");
    const m = String(minutes.value).padStart(2, "0");
    const s = String(seconds.value).padStart(2, "0");
    return `${h}:${m}:${s}`;
  });

  // 清除定时器
  const clearTimer = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  // 开始
  const start = () => {
    if (isRunning.value) return;

    isRunning.value = true;

    timer = setInterval(() => {
      elapsed.value++;
      onTick?.(elapsed.value);
    }, interval);
  };

  // 暂停
  const pause = () => {
    clearTimer();
    isRunning.value = false;
  };

  // 停止（重置）
  const stop = () => {
    clearTimer();
    isRunning.value = false;
    elapsed.value = 0;
  };

  // 重置
  const reset = () => {
    clearTimer();
    isRunning.value = false;
    elapsed.value = 0;
  };

  // 重新开始
  const restart = () => {
    reset();
    start();
  };

  // 立即开始
  if (immediate) {
    start();
  }

  // 清理
  onUnmounted(() => {
    clearTimer();
  });

  return {
    elapsed,
    isRunning,
    formatted,
    hours,
    minutes,
    seconds,
    start,
    pause,
    stop,
    reset,
    restart,
  };
}

// ============================================================================
// 验证码倒计时
// ============================================================================

export interface UseVerificationCodeOptions {
  /** 倒计时秒数，默认 60 */
  seconds?: number;
  /** 发送验证码函数 */
  onSend?: () => Promise<void>;
  /** 发送成功回调 */
  onSendSuccess?: () => void;
  /** 发送失败回调 */
  onSendError?: (error: Error) => void;
}

/**
 * 验证码倒计时
 * @example
 * const { remaining, isRunning, isSending, send, buttonText } = useVerificationCode({
 *   onSend: () => api.sendVerificationCode(phone)
 * })
 */
export function useVerificationCode(options: UseVerificationCodeOptions = {}) {
  const { seconds = 60, onSend, onSendSuccess, onSendError } = options;

  const isSending = ref(false);

  const countdown = useCountdown({
    seconds,
    immediate: false,
  });

  // 按钮文本
  const buttonText = computed(() => {
    if (isSending.value) return "发送中...";
    if (countdown.isRunning.value) return `${countdown.remaining.value}秒后重发`;
    return "发送验证码";
  });

  // 是否禁用
  const isDisabled = computed(() => {
    return isSending.value || countdown.isRunning.value;
  });

  // 发送验证码
  const send = async () => {
    if (isDisabled.value) return;

    if (!onSend) {
      countdown.start();
      return;
    }

    isSending.value = true;

    try {
      await onSend();
      countdown.start();
      onSendSuccess?.();
    } catch (e) {
      const error = e instanceof Error ? e : new Error(String(e));
      onSendError?.(error);
    } finally {
      isSending.value = false;
    }
  };

  return {
    remaining: countdown.remaining,
    isRunning: countdown.isRunning,
    isSending,
    buttonText,
    isDisabled,
    send,
    reset: countdown.reset,
  };
}

// ============================================================================
// 目标日期倒计时
// ============================================================================

export interface UseTargetDateCountdownOptions {
  /** 目标日期 */
  targetDate: Date | string | number;
  /** 倒计时结束回调 */
  onEnd?: () => void;
  /** 间隔毫秒数，默认 1000 */
  interval?: number;
}

/**
 * 目标日期倒计时
 * @example
 * const { days, hours, minutes, seconds, isFinished } = useTargetDateCountdown({
 *   targetDate: '2024-12-31T23:59:59'
 * })
 */
export function useTargetDateCountdown(options: UseTargetDateCountdownOptions) {
  const { targetDate, onEnd, interval = 1000 } = options;

  const target = new Date(targetDate).getTime();
  const remaining = ref(0);
  const isFinished = ref(false);

  let timer: ReturnType<typeof setInterval> | null = null;

  // 计算剩余时间
  const updateRemaining = () => {
    const now = Date.now();
    const diff = Math.max(0, Math.floor((target - now) / 1000));
    remaining.value = diff;

    if (diff === 0 && !isFinished.value) {
      isFinished.value = true;
      if (timer) {
        clearInterval(timer);
        timer = null;
      }
      onEnd?.();
    }
  };

  // 计算时间部分
  const days = computed(() => Math.floor(remaining.value / 86400));
  const hours = computed(() => Math.floor((remaining.value % 86400) / 3600));
  const minutes = computed(() => Math.floor((remaining.value % 3600) / 60));
  const seconds = computed(() => Math.floor(remaining.value % 60));

  // 格式化
  const formatted = computed(() => {
    return `${days.value}天 ${String(hours.value).padStart(2, "0")}:${String(minutes.value).padStart(2, "0")}:${String(seconds.value).padStart(2, "0")}`;
  });

  // 初始化
  updateRemaining();

  if (!isFinished.value) {
    timer = setInterval(updateRemaining, interval);
  }

  // 清理
  onUnmounted(() => {
    if (timer) {
      clearInterval(timer);
    }
  });

  return {
    remaining,
    isFinished,
    days,
    hours,
    minutes,
    seconds,
    formatted,
  };
}
