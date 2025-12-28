<script setup lang="ts">
import { useResponsiveTabs, type TabItem } from "@/composables";
import ResponsiveTabs from "@/components/ResponsiveTabs.vue";
import TwoFactorSettings from "./components/TwoFactorSettings.vue";
import PasswordSettings from "./components/PasswordSettings.vue";
import DeleteAccountSettings from "./components/DeleteAccountSettings.vue";

/** 安全设置 Tab 列表 */
const securityTabs: TabItem[] = [
  { value: "twoFactor", label: "双因素认证", icon: "mdi-shield-key" },
  { value: "password", label: "修改密码", icon: "mdi-lock-reset" },
  { value: "deleteAccount", label: "删除账户", icon: "mdi-account-remove" },
];

const { currentTab, isVertical, handleTabChange } = useResponsiveTabs({
  defaultTab: "twoFactor",
});
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
          <ResponsiveTabs
            :model-value="currentTab"
            :tabs="securityTabs"
            :vertical="isVertical"
            @update:model-value="handleTabChange"
          >
            <template #twoFactor>
              <TwoFactorSettings />
            </template>
            <template #password>
              <PasswordSettings />
            </template>
            <template #deleteAccount>
              <DeleteAccountSettings />
            </template>
          </ResponsiveTabs>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
