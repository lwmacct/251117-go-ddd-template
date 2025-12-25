<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useTwoFactor } from "../composables/useTwoFactor";

// 使用 2FA Composable
const {
  loading,
  enabled,
  recoveryCodesCount,
  showDisableDialog,
  qrcodeImage,
  secret,
  verifyCode,
  recoveryCodes,
  setupStep,
  errorMessage,
  successMessage,
  fetchStatus,
  startSetup,
  verifyAndEnable,
  disable2FA,
  closeSetupDialog,
  copyToClipboard,
  downloadRecoveryCodes,
} = useTwoFactor();

const statusText = computed(() => (enabled.value ? "已启用" : "未启用"));
const statusColor = computed(() => (enabled.value ? "success" : "warning"));

const resetToStatus = () => {
  closeSetupDialog();
};

// 组件挂载时获取状态
onMounted(() => {
  fetchStatus();
});
</script>

<template>
  <div>
    <v-progress-linear v-if="loading" indeterminate color="primary" class="mb-4" />

    <v-alert type="info" variant="tonal" class="mb-4">
      <div class="text-body-2">
        <strong>什么是双因素认证？</strong>
        <p class="mb-0 mt-1">启用 2FA 后，登录时需要密码 + 手机验证器中的 6 位验证码，能显著提升账号安全性。</p>
      </div>
    </v-alert>

    <v-row class="mb-4">
      <v-col cols="12" md="7">
        <v-card>
          <v-card-title class="d-flex align-center justify-space-between">
            <div>
              <div class="text-subtitle-1">双因素认证</div>
              <div class="text-caption text-medium-emphasis">当前状态：{{ statusText }}</div>
            </div>
            <v-chip :color="statusColor" variant="flat" size="small">{{ statusText }}</v-chip>
          </v-card-title>
          <v-card-text class="pt-0">
            <div v-if="enabled" class="text-body-2 mb-2">恢复码剩余 {{ recoveryCodesCount }} 个</div>
            <div class="d-flex gap-2">
              <v-btn
                v-if="!enabled || setupStep !== 'status'"
                color="primary"
                prepend-icon="mdi-shield-plus"
                :loading="loading"
                :disabled="loading"
                @click="startSetup"
              >
                立即启用
              </v-btn>
              <v-btn
                v-else
                color="error"
                variant="outlined"
                prepend-icon="mdi-shield-off"
                :loading="loading"
                :disabled="loading"
                @click="showDisableDialog = true"
              >
                禁用 2FA
              </v-btn>
              <v-btn
                v-if="setupStep !== 'status' && !loading"
                variant="text"
                prepend-icon="mdi-arrow-left"
                @click="resetToStatus"
              >
                返回状态
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 成功/错误消息 -->
    <v-alert v-if="successMessage" type="success" density="compact" class="mb-4" closable @click:close="successMessage = ''">
      {{ successMessage }}
    </v-alert>

    <v-alert v-if="errorMessage" type="error" density="compact" class="mb-4" closable @click:close="errorMessage = ''">
      {{ errorMessage }}
    </v-alert>

    <!-- 设置步骤：二维码与验证 -->
    <div v-if="setupStep === 'setup' || setupStep === 'verify'" class="mb-6">
      <v-alert type="warning" variant="tonal" class="mb-4">
        <div class="text-body-2">
          <strong>步骤 {{ setupStep === "setup" ? "1" : "2" }}：</strong>
          <span v-if="setupStep === 'setup'">使用 Google/Microsoft Authenticator 扫描二维码或手动输入密钥。</span>
          <span v-else>输入手机验证器显示的 6 位验证码完成绑定。</span>
        </div>
      </v-alert>

      <v-row>
        <v-col cols="12" md="6">
          <v-card class="mb-4">
            <v-card-text class="text-center">
              <div class="text-subtitle-2 mb-3">扫描二维码</div>
              <div v-if="qrcodeImage" class="d-flex justify-center">
                <v-img :src="qrcodeImage" width="220" height="220" contain alt="2FA 二维码" />
              </div>
              <div v-else class="d-flex justify-center align-center" style="width: 220px; height: 220px">
                <v-progress-circular indeterminate color="primary" />
              </div>
              <div class="text-caption text-medium-emphasis mt-3">使用 Authenticator App 扫码</div>
            </v-card-text>
          </v-card>
        </v-col>
        <v-col cols="12" md="6">
          <v-card class="mb-4">
            <v-card-text>
              <div class="text-subtitle-2 mb-3">手动输入密钥</div>
              <v-text-field
                :value="secret"
                readonly
                variant="outlined"
                density="compact"
                append-inner-icon="mdi-content-copy"
                @click:append-inner="copyToClipboard(secret)"
              />
              <div class="text-caption text-medium-emphasis mt-2">无法扫码时可输入此密钥</div>
            </v-card-text>
          </v-card>

          <v-card>
            <v-card-text>
              <v-text-field
                v-model="verifyCode"
                label="输入验证码"
                placeholder="请输入6位验证码"
                variant="outlined"
                density="compact"
                maxlength="6"
                class="mb-4"
                :rules="[
                  (v) => !!v || '请输入验证码',
                  (v) => (v && v.length === 6) || '验证码必须为6位数字',
                  (v) => /^\d+$/.test(v) || '验证码只能包含数字',
                ]"
              />

              <div class="d-flex justify-between">
                <v-btn variant="text" @click="setupStep = 'setup'"> 上一步 </v-btn>
                <div class="d-flex gap-2">
                  <v-btn variant="outlined" @click="resetToStatus"> 取消 </v-btn>
                  <v-btn
                    color="primary"
                    :loading="loading"
                    :disabled="!verifyCode || verifyCode.length !== 6"
                    @click="verifyAndEnable"
                  >
                    验证并启用
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </div>

    <!-- 设置步骤：显示恢复码 -->
    <div v-if="setupStep === 'codes'" class="mb-6">
      <v-alert type="warning" variant="tonal" class="mb-4">
        <div class="text-body-2">
          <strong>⚠️ 请立即保存恢复码！</strong>
          <p class="mb-0 mt-1">若丢失验证器设备，可用恢复码登录。恢复码只显示一次，请妥善保存。</p>
        </div>
      </v-alert>

      <v-card>
        <v-card-text>
          <div class="text-subtitle-2 mb-3">恢复码（共 {{ recoveryCodes.length }} 个）</div>
          <div class="d-flex flex-wrap gap-2 mb-4">
            <v-chip v-for="(code, index) in recoveryCodes" :key="index" variant="outlined" size="small">
              {{ code }}
            </v-chip>
          </div>

          <div class="d-flex flex-wrap gap-2 mb-4">
            <v-btn
              variant="outlined"
              size="small"
              prepend-icon="mdi-content-copy"
              @click="copyToClipboard(recoveryCodes.join('\\n'))"
            >
              复制所有
            </v-btn>
            <v-btn variant="outlined" size="small" prepend-icon="mdi-download" @click="downloadRecoveryCodes"> 下载保存 </v-btn>
          </div>

          <v-btn color="primary" @click="resetToStatus"> 完成 </v-btn>
        </v-card-text>
      </v-card>
    </div>

    <!-- 禁用 2FA 确认对话框 -->
    <v-dialog v-model="showDisableDialog" max-width="500">
      <v-card>
        <v-card-title class="text-h5">
          <v-icon class="mr-2" color="error">mdi-alert</v-icon>
          确认禁用 2FA
        </v-card-title>

        <v-card-text>
          <v-alert type="warning" variant="tonal" class="mb-4"> 禁用双因素认证将降低您的账号安全性。确定要继续吗？ </v-alert>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn variant="text" @click="showDisableDialog = false"> 取消 </v-btn>
          <v-btn color="error" :loading="loading" @click="disable2FA"> 确认禁用 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </div>
</template>

<style scoped>
.gap-2 {
  gap: 8px;
}

.recovery-codes code {
  font-size: 14px;
  padding: 4px 8px;
  background-color: #f5f5f5;
  border-radius: 4px;
}
</style>
