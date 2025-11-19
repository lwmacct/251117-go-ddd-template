<script setup lang="ts">
import { onMounted } from "vue";
import { useTwoFactor } from "../composables/useTwoFactor";

// 使用 2FA Composable
const {
  loading,
  enabled,
  recoveryCodesCount,
  showSetupDialog,
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

// 组件挂载时获取状态
onMounted(() => {
  fetchStatus();
});
</script>

<template>
  <div>
    <!-- 状态显示 -->
    <v-alert v-if="enabled" type="success" variant="tonal" class="mb-4">
      <div class="d-flex align-center justify-space-between">
        <div>
          <div class="font-weight-bold">2FA 已启用</div>
          <div class="text-caption">剩余恢复码: {{ recoveryCodesCount }} 个</div>
        </div>
        <v-icon size="large">mdi-check-circle</v-icon>
      </div>
    </v-alert>

    <v-alert v-else type="warning" variant="tonal" class="mb-4">
      <div class="font-weight-bold">2FA 未启用</div>
      <div class="text-caption">建议启用双因素认证以提高账号安全性</div>
    </v-alert>

    <!-- 说明 -->
    <div class="mb-4">
      <h3 class="text-h6 mb-2">什么是双因素认证？</h3>
      <p class="text-body-2">
        双因素认证 (2FA) 是一种额外的安全层，除了密码外，还需要输入您手机应用（如 Google Authenticator、Microsoft Authenticator 等）生成的 6 位验证码才能登录。
      </p>
    </div>

    <!-- 操作按钮 -->
    <div class="d-flex gap-2">
      <v-btn v-if="!enabled" color="primary" prepend-icon="mdi-shield-plus" :loading="loading" @click="startSetup"> 启用 2FA </v-btn>

      <v-btn v-else color="error" variant="outlined" prepend-icon="mdi-shield-off" :loading="loading" @click="showDisableDialog = true"> 禁用 2FA </v-btn>
    </div>

    <!-- 成功/错误消息 -->
    <v-alert v-if="successMessage" type="success" density="compact" class="mt-4" closable @click:close="successMessage = ''">
      {{ successMessage }}
    </v-alert>

    <v-alert v-if="errorMessage" type="error" density="compact" class="mt-4" closable @click:close="errorMessage = ''">
      {{ errorMessage }}
    </v-alert>

    <!-- 设置 2FA 对话框 -->
    <v-dialog v-model="showSetupDialog" max-width="600" persistent>
      <v-card>
        <v-card-title class="text-h5">
          <v-icon class="mr-2">mdi-shield-lock-outline</v-icon>
          设置双因素认证
        </v-card-title>

        <v-card-text>
          <!-- 步骤 1: 扫描二维码 -->
          <div v-if="setupStep === 'qrcode'">
            <v-stepper-header class="elevation-0 mb-4">
              <v-stepper-item :complete="false" :value="1" title="扫描二维码"></v-stepper-item>
              <v-divider></v-divider>
              <v-stepper-item :value="2" title="验证"></v-stepper-item>
              <v-divider></v-divider>
              <v-stepper-item :value="3" title="保存恢复码"></v-stepper-item>
            </v-stepper-header>

            <v-alert type="info" variant="tonal" density="compact" class="mb-4">
              <div class="text-body-2">
                <div class="font-weight-bold mb-2">步骤 1: 使用认证器应用扫描二维码</div>
                <div>1. 在手机上安装 Google Authenticator 或 Microsoft Authenticator</div>
                <div>2. 打开应用，选择"添加账户"或扫描二维码</div>
                <div>3. 扫描下方的二维码</div>
              </div>
            </v-alert>

            <div class="text-center mb-4">
              <img :src="qrcodeImage" alt="2FA QR Code" style="max-width: 256px; border: 1px solid #e0e0e0; border-radius: 8px" />
            </div>

            <v-alert type="info" variant="tonal" density="compact" class="mb-4">
              <div class="text-body-2">
                <div class="font-weight-bold mb-2">或手动输入密钥：</div>
                <div class="d-flex align-center gap-2">
                  <code class="flex-1">{{ secret }}</code>
                  <v-btn icon="mdi-content-copy" variant="text" size="small" @click="copyToClipboard(secret)"></v-btn>
                </div>
              </div>
            </v-alert>

            <v-btn color="primary" block @click="setupStep = 'verify'"> 下一步：验证 </v-btn>
          </div>

          <!-- 步骤 2: 验证 -->
          <div v-if="setupStep === 'verify'">
            <v-stepper-header class="elevation-0 mb-4">
              <v-stepper-item :complete="true" :value="1" title="扫描二维码"></v-stepper-item>
              <v-divider></v-divider>
              <v-stepper-item :complete="false" :value="2" title="验证"></v-stepper-item>
              <v-divider></v-divider>
              <v-stepper-item :value="3" title="保存恢复码"></v-stepper-item>
            </v-stepper-header>

            <v-alert type="info" variant="tonal" density="compact" class="mb-4">
              <div class="text-body-2">
                <div class="font-weight-bold mb-2">步骤 2: 输入验证码</div>
                <div>请输入认证器应用中显示的 6 位验证码</div>
              </div>
            </v-alert>

            <v-text-field v-model="verifyCode" label="6 位验证码" type="text" maxlength="6" variant="outlined" autofocus :error-messages="errorMessage" class="mb-4"></v-text-field>

            <div class="d-flex gap-2">
              <v-btn variant="outlined" @click="setupStep = 'qrcode'"> 返回 </v-btn>
              <v-btn color="primary" :loading="loading" :disabled="verifyCode.length !== 6" @click="verifyAndEnable"> 验证并启用 </v-btn>
            </div>
          </div>

          <!-- 步骤 3: 保存恢复码 -->
          <div v-if="setupStep === 'codes'">
            <v-stepper-header class="elevation-0 mb-4">
              <v-stepper-item :complete="true" :value="1" title="扫描二维码"></v-stepper-item>
              <v-divider></v-divider>
              <v-stepper-item :complete="true" :value="2" title="验证"></v-stepper-item>
              <v-divider></v-divider>
              <v-stepper-item :complete="false" :value="3" title="保存恢复码"></v-stepper-item>
            </v-stepper-header>

            <v-alert type="success" variant="tonal" density="compact" class="mb-4">
              <div class="text-body-2">
                <div class="font-weight-bold mb-2">✓ 2FA 已成功启用！</div>
              </div>
            </v-alert>

            <v-alert type="warning" variant="tonal" class="mb-4">
              <div class="text-body-2">
                <div class="font-weight-bold mb-2">重要：请保存以下恢复码</div>
                <div>如果您无法访问认证器应用，可以使用这些恢复码登录。每个恢复码只能使用一次。</div>
              </div>
            </v-alert>

            <v-card variant="outlined" class="mb-4">
              <v-card-text>
                <div class="recovery-codes">
                  <code v-for="(code, index) in recoveryCodes" :key="index" class="d-block mb-2">{{ code }}</code>
                </div>
              </v-card-text>
            </v-card>

            <div class="d-flex gap-2 mb-4">
              <v-btn prepend-icon="mdi-download" variant="outlined" @click="downloadRecoveryCodes"> 下载恢复码 </v-btn>
            </div>

            <v-btn color="primary" block @click="closeSetupDialog"> 完成 </v-btn>
          </div>
        </v-card-text>

        <v-card-actions>
          <v-btn v-if="setupStep !== 'codes'" variant="text" @click="closeSetupDialog"> 取消 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

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
