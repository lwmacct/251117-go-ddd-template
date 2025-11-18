<template>
  <div class="vp-raw api-endpoint">
    <p>
      <code class="api-method" :class="`method-${method.toLowerCase()}`">{{ method }}</code>
      <code>{{ path }}</code>
      <span v-if="version" class="vp-badge tip">{{ version }}</span>
      <span v-if="deprecated" class="vp-badge warning">已废弃</span>
    </p>
    <p v-if="description" class="api-description">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  path: string;
  description?: string;
  version?: string;
  deprecated?: boolean;
}>();
</script>

<style scoped>
.api-endpoint {
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 16px;
  margin: 16px 0;
  background: var(--vp-c-bg-soft);
}

.api-method {
  padding: 2px 6px;
  margin-right: 8px;
  font-weight: 600;
  font-size: 11px;
  text-transform: uppercase;
  border-radius: 4px;
}

.method-get { background: #61affe; color: white; }
.method-post { background: #49cc90; color: white; }
.method-put { background: #fca130; color: white; }
.method-patch { background: #50e3c2; color: white; }
.method-delete { background: #f93e3e; color: white; }

.api-description {
  color: var(--vp-c-text-2);
  margin-top: 8px;
}

.vp-badge {
  display: inline-block;
  padding: 0 10px;
  margin-left: 6px;
  font-size: 12px;
  font-weight: 500;
  line-height: 18px;
  border-radius: 10px;
  border: 1px solid transparent;
}

.vp-badge.tip {
  background: var(--vp-c-tip-soft);
  color: var(--vp-c-tip-1);
}

.vp-badge.warning {
  background: var(--vp-c-warning-soft);
  color: var(--vp-c-warning-1);
}
</style>
