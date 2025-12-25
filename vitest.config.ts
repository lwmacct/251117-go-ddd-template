import { fileURLToPath } from "node:url";
import { mergeConfig, defineConfig, configDefaults } from "vitest/config";
import viteConfig from "./vite.config";

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      environment: "happy-dom",
      exclude: [...configDefaults.exclude, "e2e/**"],
      root: fileURLToPath(new URL("./", import.meta.url)),
      coverage: {
        provider: "v8",
        reporter: ["text", "html", "lcov"],
        exclude: ["node_modules/**", "dist/**", "**/*.d.ts", "**/*.config.*", "**/types/**"],
      },
      // 全局 API 可用
      globals: true,
    },
  }),
);
