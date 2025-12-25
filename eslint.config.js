import js from "@eslint/js";
import pluginVue from "eslint-plugin-vue";
import tseslint from "typescript-eslint";
import prettier from "eslint-config-prettier";
import pluginPrettier from "eslint-plugin-prettier";
import globals from "globals";
import { includeIgnoreFile } from "@eslint/compat";
import { fileURLToPath } from "node:url";

// 引用 .gitignore 文件
const gitignorePath = fileURLToPath(new URL(".gitignore", import.meta.url));

export default [
  // 从 .gitignore 导入忽略规则
  includeIgnoreFile(gitignorePath),

  // 额外的忽略规则（.gitignore 中未包含的）
  {
    ignores: [
      // TypeScript 声明文件
      "*.d.ts",
      "src/auto-imports.d.ts",
      "src/components.d.ts",
      // 独立子项目（有自己的 lint 配置）
      "docs/**",
      // Go 后端相关
      "internal/**",
      "cmd/**",
      "testing/**",
    ],
  },

  // JavaScript 基础规则
  js.configs.recommended,

  // TypeScript 规则
  ...tseslint.configs.recommended,

  // Vue 规则
  ...pluginVue.configs["flat/recommended"],

  // Prettier 集成
  prettier,

  // 全局配置
  {
    plugins: {
      prettier: pluginPrettier,
    },
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.es2021,
      },
      parserOptions: {
        parser: tseslint.parser,
        ecmaVersion: "latest",
        sourceType: "module",
      },
    },
    rules: {
      // Prettier 规则
      "prettier/prettier": "warn",

      // Vue 规则调整
      "vue/multi-word-component-names": "off",
      "vue/no-v-html": "off",
      "vue/require-default-prop": "off",
      "vue/no-unused-vars": "warn",
      // 允许 Vuetify 的动态插槽语法 (如 #item.xxx)
      "vue/valid-v-slot": ["error", { allowModifiers: true }],

      // TypeScript 规则调整
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-unused-vars": [
        "warn",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
        },
      ],
      "@typescript-eslint/explicit-function-return-type": "off",
      "@typescript-eslint/no-empty-object-type": "off",
      // 允许 this 别名（throttle/debounce 等工具函数需要）
      "@typescript-eslint/no-this-alias": "off",

      // 通用规则
      "no-console": ["warn", { allow: ["warn", "error"] }],
      "no-debugger": "warn",
    },
  },

  // Vue 文件特殊配置
  {
    files: ["**/*.vue"],
    languageOptions: {
      parserOptions: {
        parser: tseslint.parser,
      },
    },
  },
];
