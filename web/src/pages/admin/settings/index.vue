<script setup lang="ts">
import { ref } from "vue";

/**
 * 系统设置页面
 * 用于配置系统参数和偏好设置
 */

const tabs = ref("general");

// 设置表单
const settings = ref({
  siteName: "",
  siteUrl: "",
  email: "",
  timezone: "Asia/Shanghai",
  language: "zh-CN",
  theme: "light",
  enableNotifications: true,
  enableBackup: false,
});

const saveSettings = () => {
  console.log("Saving settings...", settings.value);
  // TODO: 实现保存逻辑
};
</script>

<template>
  <div class="settings-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">系统设置</h1>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-tabs v-model="tabs" bg-color="primary">
            <v-tab value="general">常规设置</v-tab>
            <v-tab value="security">安全设置</v-tab>
            <v-tab value="notification">通知设置</v-tab>
            <v-tab value="backup">备份设置</v-tab>
          </v-tabs>

          <v-card-text class="pa-6">
            <v-tabs-window v-model="tabs">
              <!-- 常规设置 -->
              <v-tabs-window-item value="general">
                <v-form>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field v-model="settings.siteName" label="站点名称" variant="outlined"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model="settings.siteUrl" label="站点 URL" variant="outlined"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model="settings.email" label="管理员邮箱" type="email" variant="outlined"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select v-model="settings.timezone" label="时区" :items="['Asia/Shanghai', 'UTC', 'America/New_York']" variant="outlined"></v-select>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select
                        v-model="settings.language"
                        label="语言"
                        :items="[
                          {
                            title: '简体中文',
                            value: 'zh-CN',
                          },
                          {
                            title: 'English',
                            value: 'en-US',
                          },
                        ]"
                        variant="outlined"
                      ></v-select>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select
                        v-model="settings.theme"
                        label="主题"
                        :items="[
                          {
                            title: '浅色',
                            value: 'light',
                          },
                          {
                            title: '深色',
                            value: 'dark',
                          },
                        ]"
                        variant="outlined"
                      ></v-select>
                    </v-col>
                  </v-row>
                </v-form>
              </v-tabs-window-item>

              <!-- 安全设置 -->
              <v-tabs-window-item value="security">
                <v-alert type="info" variant="tonal" class="mb-4"> 配置系统安全相关参数 </v-alert>
                <v-form>
                  <v-text-field label="会话超时时间 (分钟) " type="number" variant="outlined" class="mb-4"></v-text-field>
                  <v-text-field label="密码最小长度" type="number" variant="outlined" class="mb-4"></v-text-field>
                  <v-switch label="启用两步验证" color="primary" class="mb-4"></v-switch>
                </v-form>
              </v-tabs-window-item>

              <!-- 通知设置 -->
              <v-tabs-window-item value="notification">
                <v-switch v-model="settings.enableNotifications" label="启用系统通知" color="primary" class="mb-4"></v-switch>
                <v-alert type="info" variant="tonal"> 更多通知设置将在后续版本中添加 </v-alert>
              </v-tabs-window-item>

              <!-- 备份设置 -->
              <v-tabs-window-item value="backup">
                <v-switch v-model="settings.enableBackup" label="启用自动备份" color="primary" class="mb-4"></v-switch>
                <v-text-field label="备份频率 (小时) " type="number" variant="outlined" :disabled="!settings.enableBackup"></v-text-field>
              </v-tabs-window-item>
            </v-tabs-window>
          </v-card-text>

          <v-card-actions class="pa-6">
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="saveSettings">保存设置</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.settings-page {
  width: 100%;
}
</style>
