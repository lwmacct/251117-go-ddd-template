<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useSettings } from './composables/useSettings';

const { settings, loading, saving, errorMessage, successMessage, fetchSettings, getSetting, batchUpdateSettings, clearMessages } = useSettings();

const tabs = ref('general');

// 常规设置表单
const generalForm = ref({
  siteName: '',
  siteUrl: '',
  adminEmail: '',
  timezone: 'Asia/Shanghai',
  language: 'zh-CN',
  theme: 'light',
});

// 安全设置表单
const securityForm = ref({
  sessionTimeout: 30,
  passwordMinLength: 8,
  enableTwoFA: false,
  maxLoginAttempts: 5,
});

// 通知设置表单
const notificationForm = ref({
  enableNotifications: true,
  enableEmailNotifications: true,
  enableSMSNotifications: false,
});

// 备份设置表单
const backupForm = ref({
  enableBackup: false,
  backupFrequency: 24,
  backupRetentionDays: 30,
});

// 时区选项
const timezoneOptions = [
  { title: 'Asia/Shanghai (上海)', value: 'Asia/Shanghai' },
  { title: 'UTC', value: 'UTC' },
  { title: 'America/New_York (纽约)', value: 'America/New_York' },
  { title: 'Europe/London (伦敦)', value: 'Europe/London' },
  { title: 'Asia/Tokyo (东京)', value: 'Asia/Tokyo' },
];

onMounted(async () => {
  await fetchSettings();
  loadFormValues();
});

// 从设置列表加载表单值
const loadFormValues = () => {
  // 常规设置
  generalForm.value.siteName = getSetting('general.site_name', '');
  generalForm.value.siteUrl = getSetting('general.site_url', '');
  generalForm.value.adminEmail = getSetting('general.admin_email', '');
  generalForm.value.timezone = getSetting('general.timezone', 'Asia/Shanghai');
  generalForm.value.language = getSetting('general.language', 'zh-CN');
  generalForm.value.theme = getSetting('general.theme', 'light');

  // 安全设置
  securityForm.value.sessionTimeout = getSetting('security.session_timeout', 30);
  securityForm.value.passwordMinLength = getSetting('security.password_min_length', 8);
  securityForm.value.enableTwoFA = getSetting('security.enable_twofa', false);
  securityForm.value.maxLoginAttempts = getSetting('security.max_login_attempts', 5);

  // 通知设置
  notificationForm.value.enableNotifications = getSetting('notification.enable_notifications', true);
  notificationForm.value.enableEmailNotifications = getSetting('notification.enable_email', true);
  notificationForm.value.enableSMSNotifications = getSetting('notification.enable_sms', false);

  // 备份设置
  backupForm.value.enableBackup = getSetting('backup.enable_backup', false);
  backupForm.value.backupFrequency = getSetting('backup.backup_frequency', 24);
  backupForm.value.backupRetentionDays = getSetting('backup.retention_days', 30);
};

// 保存当前 Tab 的设置
const saveCurrentTab = async () => {
  let updates: { key: string; value: any }[] = [];

  switch (tabs.value) {
    case 'general':
      updates = [
        { key: 'general.site_name', value: generalForm.value.siteName },
        { key: 'general.site_url', value: generalForm.value.siteUrl },
        { key: 'general.admin_email', value: generalForm.value.adminEmail },
        { key: 'general.timezone', value: generalForm.value.timezone },
        { key: 'general.language', value: generalForm.value.language },
        { key: 'general.theme', value: generalForm.value.theme },
      ];
      break;
    case 'security':
      updates = [
        { key: 'security.session_timeout', value: securityForm.value.sessionTimeout },
        { key: 'security.password_min_length', value: securityForm.value.passwordMinLength },
        { key: 'security.enable_twofa', value: securityForm.value.enableTwoFA },
        { key: 'security.max_login_attempts', value: securityForm.value.maxLoginAttempts },
      ];
      break;
    case 'notification':
      updates = [
        { key: 'notification.enable_notifications', value: notificationForm.value.enableNotifications },
        { key: 'notification.enable_email', value: notificationForm.value.enableEmailNotifications },
        { key: 'notification.enable_sms', value: notificationForm.value.enableSMSNotifications },
      ];
      break;
    case 'backup':
      updates = [
        { key: 'backup.enable_backup', value: backupForm.value.enableBackup },
        { key: 'backup.backup_frequency', value: backupForm.value.backupFrequency },
        { key: 'backup.retention_days', value: backupForm.value.backupRetentionDays },
      ];
      break;
  }

  await batchUpdateSettings(updates);
};

// 保存所有设置
const saveAllSettings = async () => {
  const allUpdates = [
    // 常规设置
    { key: 'general.site_name', value: generalForm.value.siteName },
    { key: 'general.site_url', value: generalForm.value.siteUrl },
    { key: 'general.admin_email', value: generalForm.value.adminEmail },
    { key: 'general.timezone', value: generalForm.value.timezone },
    { key: 'general.language', value: generalForm.value.language },
    { key: 'general.theme', value: generalForm.value.theme },
    // 安全设置
    { key: 'security.session_timeout', value: securityForm.value.sessionTimeout },
    { key: 'security.password_min_length', value: securityForm.value.passwordMinLength },
    { key: 'security.enable_twofa', value: securityForm.value.enableTwoFA },
    { key: 'security.max_login_attempts', value: securityForm.value.maxLoginAttempts },
    // 通知设置
    { key: 'notification.enable_notifications', value: notificationForm.value.enableNotifications },
    { key: 'notification.enable_email', value: notificationForm.value.enableEmailNotifications },
    { key: 'notification.enable_sms', value: notificationForm.value.enableSMSNotifications },
    // 备份设置
    { key: 'backup.enable_backup', value: backupForm.value.enableBackup },
    { key: 'backup.backup_frequency', value: backupForm.value.backupFrequency },
    { key: 'backup.retention_days', value: backupForm.value.backupRetentionDays },
  ];

  await batchUpdateSettings(allUpdates);
};
</script>

<template>
  <div class="settings-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-2">系统设置</h1>
        <p class="text-body-2 text-medium-emphasis mb-6">配置系统参数和偏好设置</p>
      </v-col>
    </v-row>

    <v-row v-if="errorMessage || successMessage">
      <v-col cols="12">
        <v-alert v-if="errorMessage" type="error" closable @click:close="clearMessages">
          {{ errorMessage }}
        </v-alert>
        <v-alert v-if="successMessage" type="success" closable @click:close="clearMessages">
          {{ successMessage }}
        </v-alert>
      </v-col>
    </v-row>

    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4"></v-progress-linear>

    <v-row>
      <v-col cols="12">
        <v-card>
          <v-tabs v-model="tabs" bg-color="primary">
            <v-tab value="general">
              <v-icon start>mdi-cog</v-icon>
              常规设置
            </v-tab>
            <v-tab value="security">
              <v-icon start>mdi-shield-lock</v-icon>
              安全设置
            </v-tab>
            <v-tab value="notification">
              <v-icon start>mdi-bell</v-icon>
              通知设置
            </v-tab>
            <v-tab value="backup">
              <v-icon start>mdi-backup-restore</v-icon>
              备份设置
            </v-tab>
          </v-tabs>

          <v-card-text class="pa-6">
            <v-tabs-window v-model="tabs">
              <!-- 常规设置 -->
              <v-tabs-window-item value="general">
                <v-form>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field v-model="generalForm.siteName" label="站点名称" variant="outlined" hint="网站显示名称"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model="generalForm.siteUrl" label="站点 URL" variant="outlined" hint="网站访问地址" placeholder="https://example.com"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model="generalForm.adminEmail" label="管理员邮箱" type="email" variant="outlined" hint="接收系统通知的邮箱"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select v-model="generalForm.timezone" label="时区" :items="timezoneOptions" variant="outlined" hint="系统默认时区"></v-select>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select
                        v-model="generalForm.language"
                        label="语言"
                        :items="[
                          { title: '简体中文', value: 'zh-CN' },
                          { title: 'English', value: 'en-US' },
                        ]"
                        variant="outlined"
                        hint="系统界面语言"
                      ></v-select>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-select
                        v-model="generalForm.theme"
                        label="主题"
                        :items="[
                          { title: '浅色', value: 'light' },
                          { title: '深色', value: 'dark' },
                          { title: '自动', value: 'auto' },
                        ]"
                        variant="outlined"
                        hint="系统界面主题"
                      ></v-select>
                    </v-col>
                  </v-row>
                </v-form>
              </v-tabs-window-item>

              <!-- 安全设置 -->
              <v-tabs-window-item value="security">
                <v-alert type="info" variant="tonal" class="mb-4">
                  <v-icon start>mdi-information</v-icon>
                  配置系统安全相关参数，这些设置将影响所有用户
                </v-alert>
                <v-form>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field v-model.number="securityForm.sessionTimeout" label="会话超时时间 (分钟)" type="number" variant="outlined" hint="用户无活动后自动登出的时间" :min="5" :max="1440"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model.number="securityForm.passwordMinLength" label="密码最小长度" type="number" variant="outlined" hint="用户密码的最小字符数" :min="6" :max="32"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model.number="securityForm.maxLoginAttempts" label="最大登录尝试次数" type="number" variant="outlined" hint="超过此次数后锁定账户" :min="3" :max="10"></v-text-field>
                    </v-col>
                    <v-col cols="12">
                      <v-switch v-model="securityForm.enableTwoFA" label="强制启用两步验证" color="primary" hint="要求所有用户启用 2FA" persistent-hint></v-switch>
                    </v-col>
                  </v-row>
                </v-form>
              </v-tabs-window-item>

              <!-- 通知设置 -->
              <v-tabs-window-item value="notification">
                <v-form>
                  <v-row>
                    <v-col cols="12">
                      <v-switch v-model="notificationForm.enableNotifications" label="启用系统通知" color="primary" hint="开启/关闭所有系统通知" persistent-hint class="mb-4"></v-switch>
                    </v-col>
                    <v-col cols="12">
                      <v-switch v-model="notificationForm.enableEmailNotifications" label="启用邮件通知" color="primary" :disabled="!notificationForm.enableNotifications" hint="通过邮件发送通知" persistent-hint class="mb-4"></v-switch>
                    </v-col>
                    <v-col cols="12">
                      <v-switch v-model="notificationForm.enableSMSNotifications" label="启用短信通知" color="primary" :disabled="!notificationForm.enableNotifications" hint="通过短信发送通知" persistent-hint class="mb-4"></v-switch>
                    </v-col>
                    <v-col cols="12">
                      <v-alert type="info" variant="tonal">
                        <v-icon start>mdi-information</v-icon>
                        更多通知设置（如通知模板、推送配置）将在后续版本中添加
                      </v-alert>
                    </v-col>
                  </v-row>
                </v-form>
              </v-tabs-window-item>

              <!-- 备份设置 -->
              <v-tabs-window-item value="backup">
                <v-form>
                  <v-row>
                    <v-col cols="12">
                      <v-switch v-model="backupForm.enableBackup" label="启用自动备份" color="primary" hint="定期自动备份系统数据" persistent-hint class="mb-4"></v-switch>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model.number="backupForm.backupFrequency" label="备份频率 (小时)" type="number" variant="outlined" :disabled="!backupForm.enableBackup" hint="自动备份的时间间隔" :min="1" :max="168"></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field v-model.number="backupForm.backupRetentionDays" label="备份保留天数" type="number" variant="outlined" :disabled="!backupForm.enableBackup" hint="备份文件保留时长" :min="1" :max="365"></v-text-field>
                    </v-col>
                    <v-col cols="12" v-if="backupForm.enableBackup">
                      <v-alert type="warning" variant="tonal">
                        <v-icon start>mdi-alert</v-icon>
                        <div>
                          <div class="text-subtitle-2 mb-1">备份注意事项：</div>
                          <ul class="pl-4">
                            <li>确保有足够的存储空间</li>
                            <li>定期测试备份恢复流程</li>
                            <li>备份文件应存储在安全位置</li>
                          </ul>
                        </div>
                      </v-alert>
                    </v-col>
                  </v-row>
                </v-form>
              </v-tabs-window-item>
            </v-tabs-window>
          </v-card-text>

          <v-card-actions class="pa-6">
            <v-spacer></v-spacer>
            <v-btn variant="text" @click="loadFormValues" :disabled="saving">重置</v-btn>
            <v-btn color="secondary" variant="tonal" @click="saveCurrentTab" :loading="saving">保存当前页</v-btn>
            <v-btn color="primary" @click="saveAllSettings" :loading="saving">保存所有设置</v-btn>
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
