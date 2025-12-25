<script setup lang="ts">
/**
 * Auth 共享布局
 *
 * 为 login 和 register 提供统一的左右分屏布局：
 * - 左侧：表单区（通过 router-view 渲染子路由）
 * - 右侧：品牌区（固定不变）
 */
</script>

<template>
  <div class="auth-page">
    <!-- 左侧表单区 -->
    <div class="form-panel">
      <div class="form-container">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>

    <!-- 右侧品牌区 -->
    <div class="brand-panel">
      <div class="brand-content">
        <!-- Logo -->
        <div class="brand-logo">
          <v-icon size="64" color="primary">mdi-shield-account</v-icon>
        </div>

        <!-- 标语 -->
        <h1 class="brand-title">DDD Template</h1>
        <p class="brand-subtitle">基于领域驱动设计的现代化应用模板</p>

        <!-- 特性列表 -->
        <div class="brand-features">
          <div class="feature-item">
            <v-icon color="primary" size="20" class="mr-3">mdi-check-circle</v-icon>
            <span>安全可靠的认证系统</span>
          </div>
          <div class="feature-item">
            <v-icon color="primary" size="20" class="mr-3">mdi-check-circle</v-icon>
            <span>RBAC 权限管理</span>
          </div>
          <div class="feature-item">
            <v-icon color="primary" size="20" class="mr-3">mdi-check-circle</v-icon>
            <span>企业级 DDD 架构</span>
          </div>
        </div>
      </div>

      <!-- 底部版权信息 -->
      <div class="brand-footer">
        <span>© 2024 DDD Template. All rights reserved.</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  display: flex;
  min-height: 100vh;
  width: 100%;
  background: #f5f5f5;
}

/* 左侧表单区 - 40% */
.form-panel {
  flex: 4;
  background: #ffffff;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 48px;
  min-width: 400px;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
}

.form-container {
  width: 100%;
  max-width: 380px;
}

/* 右侧品牌区 - 60% */
.brand-panel {
  flex: 6;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 40px;
  position: relative;
}

.brand-content {
  text-align: center;
  max-width: 420px;
}

.brand-logo {
  margin-bottom: 24px;
}

.brand-title {
  font-size: 2rem;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
  letter-spacing: -0.5px;
}

.brand-subtitle {
  font-size: 1rem;
  color: #666;
  margin-bottom: 48px;
  line-height: 1.6;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 16px;
  text-align: left;
}

.feature-item {
  display: flex;
  align-items: center;
  font-size: 0.95rem;
  color: #555;
}

.brand-footer {
  position: absolute;
  bottom: 24px;
  left: 0;
  right: 0;
  text-align: center;
  color: #999;
  font-size: 0.85rem;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 响应式设计 - 移动端隐藏右侧品牌区 */
@media (max-width: 960px) {
  .brand-panel {
    display: none;
  }

  .form-panel {
    flex: 1;
    min-width: unset;
    padding: 24px 16px;
    box-shadow: none;
  }

  .form-container {
    max-width: 400px;
    margin: 0 auto;
  }
}
</style>
