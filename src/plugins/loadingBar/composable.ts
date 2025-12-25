/**
 * 路由加载进度条 Composable
 * 使用 VueUse useTransition 实现平滑的进度动画
 */
import { ref, computed } from "vue";
import { useTransition, TransitionPresets } from "@vueuse/core";

// 全局单例状态（模块级别，确保整个应用共享同一个进度条）
const isLoading = ref(false);
const targetProgress = ref(0);
const color = ref<"primary" | "error">("primary");

// 使用 VueUse useTransition 实现平滑过渡
// easeOutCubic: 快速启动、缓慢结束，给人"即将完成"的感觉
const progress = useTransition(targetProgress, {
  duration: 200,
  transition: TransitionPresets.easeOutCubic,
});

// 定时器引用（模块级别，避免多实例冲突）
let timer: ReturnType<typeof setInterval> | null = null;

/**
 * 路由加载进度条 Hook
 * 提供 start/finish/fail 方法控制进度条
 */
export function useLoadingBar() {
  /**
   * 停止进度递增
   */
  const stop = () => {
    if (timer) {
      clearInterval(timer);
      timer = null;
    }
  };

  /**
   * 开始加载
   * 启动进度条，模拟进度递增直到 90%
   */
  const start = () => {
    stop();
    isLoading.value = true;
    targetProgress.value = 0;
    color.value = "primary";

    // 每 200ms 递增 5-15%，到 90% 后停止（等待真正完成）
    timer = setInterval(() => {
      if (targetProgress.value < 90) {
        targetProgress.value += Math.random() * 10 + 5;
      }
    }, 200);
  };

  /**
   * 完成加载
   * 快速推进到 100%，延迟 300ms 后隐藏
   */
  const finish = () => {
    stop();
    targetProgress.value = 100;

    // 延迟隐藏，让用户看到完成动画
    setTimeout(() => {
      isLoading.value = false;
      targetProgress.value = 0;
    }, 300);
  };

  /**
   * 加载失败
   * 显示红色错误状态，延迟 500ms 后隐藏
   */
  const fail = () => {
    stop();
    color.value = "error";
    targetProgress.value = 100;

    setTimeout(() => {
      isLoading.value = false;
      targetProgress.value = 0;
      color.value = "primary";
    }, 500);
  };

  return {
    // 响应式状态（只读）
    isLoading: computed(() => isLoading.value),
    progress,
    color: computed(() => color.value),

    // 控制方法
    start,
    finish,
    fail,
  };
}
