<template>
  <ol class="vp-raw steps-guide">
    <li v-for="(step, index) in steps" :key="index">
      <strong>{{ step.title }}</strong>
      <p>
        {{ step.description }}
      </p>
    </li>
  </ol>
</template>

<script setup lang="ts">
interface Step {
  title: string;
  description: string;
}

defineProps<{
  steps: Step[];
}>();
</script>

<style scoped>
.steps-guide {
  counter-reset: step-counter;
  list-style: none;
  padding: 0;
  margin: 20px 0;
}

.steps-guide li {
  position: relative;
  padding: 16px 16px 16px 48px;
  margin-bottom: 12px;
  border-left: 2px solid var(--vp-c-divider);
  background: var(--vp-c-bg-soft);
  border-radius: 4px;
  counter-increment: step-counter;
}

.steps-guide li::before {
  content: counter(step-counter);
  position: absolute;
  left: 12px;
  top: 16px;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--vp-c-brand-1);
  color: white;
  border-radius: 50%;
  font-weight: 600;
  font-size: 12px;
}

.steps-guide li:hover {
  border-left-color: var(--vp-c-brand-1);
}

.steps-guide strong {
  display: block;
  margin-bottom: 4px;
  color: var(--vp-c-text-1);
}

.steps-guide p {
  margin: 0;
  color: var(--vp-c-text-2);
}
</style>
