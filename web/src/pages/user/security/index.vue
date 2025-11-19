<script setup lang="ts">
import { ref } from "vue";
import TwoFactorSettings from "./components/TwoFactorSettings.vue";
import PasswordSettings from "./components/PasswordSettings.vue";

// 当前激活的 Tab
const activeTab = ref("twoFactor");
</script>

<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">
          <v-icon class="mr-2">mdi-shield-lock</v-icon>
          安全设置
        </h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-tabs v-model="activeTab" bg-color="primary">
            <v-tab value="twoFactor">
              <v-icon class="mr-2">mdi-shield-key</v-icon>
              双因素认证
            </v-tab>
            <v-tab value="password">
              <v-icon class="mr-2">mdi-lock-reset</v-icon>
              修改密码
            </v-tab>
          </v-tabs>

          <v-card-text class="pa-6">
            <v-window v-model="activeTab">
              <!-- 2FA 设置 -->
              <v-window-item value="twoFactor">
                <TwoFactorSettings />
              </v-window-item>

              <!-- 修改密码 -->
              <v-window-item value="password">
                <PasswordSettings />
              </v-window-item>
            </v-window>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
/* 可根据需要添加自定义样式 */
</style>
