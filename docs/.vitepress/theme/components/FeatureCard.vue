<template>
  <div class="vp-raw feature-card" :class="{ highlighted }">
    <h3 v-if="icon || title">
      <span v-if="icon" class="feature-icon">{{ icon }}</span>
      {{ title }}
    </h3>
    <p v-if="description">{{ description }}</p>
    <slot />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  description?: string;
  icon?: string;
  highlighted?: boolean;
}>();
</script>

<style scoped>
.feature-card {
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  padding: 20px;
  margin: 16px 0;
  background: var(--vp-c-bg-soft);
  transition: border-color 0.25s, box-shadow 0.25s;
}

.feature-card:hover {
  border-color: var(--vp-c-brand-1);
}

.feature-card.highlighted {
  border-color: var(--vp-c-brand-1);
  background: var(--vp-c-brand-soft);
}

.feature-icon {
  margin-right: 8px;
}
</style>
