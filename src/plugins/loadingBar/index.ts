/**
 * LoadingBar 插件
 * 为 Vue Router 提供路由切换时的加载进度条
 *
 * 使用方式：
 * ```ts
 * import { loadingBarPlugin } from "@/plugins/loadingBar";
 * app.use(loadingBarPlugin, { router });
 * ```
 */
import type { App, Plugin } from "vue";
import type { Router } from "vue-router";
import { useLoadingBar } from "./composable";
import LoadingBar from "./LoadingBar.vue";

export interface LoadingBarPluginOptions {
  router: Router;
}

/**
 * LoadingBar Vue 插件
 * 自动注册路由守卫，在路由切换时显示/隐藏进度条
 */
export const loadingBarPlugin: Plugin<LoadingBarPluginOptions> = {
  install(app: App, options: LoadingBarPluginOptions) {
    const { router } = options;
    const loadingBar = useLoadingBar();

    // 路由进入前 - 启动进度条
    router.beforeEach(() => {
      loadingBar.start();
    });

    // 路由完成后 - 完成进度条
    router.afterEach(() => {
      loadingBar.finish();
    });

    // 路由错误 - 显示失败状态
    router.onError(() => {
      loadingBar.fail();
    });

    // 全局注册组件，可在任意模板中使用 <LoadingBar />
    app.component("LoadingBar", LoadingBar);
  },
};

// 导出组件和 composable，支持手动控制进度条
export { useLoadingBar } from "./composable";
export { default as LoadingBar } from "./LoadingBar.vue";
