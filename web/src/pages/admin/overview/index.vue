<script setup lang="ts">
import { ref, onMounted } from 'vue'

/**
 * 数据概览页面
 * 用于查看和管理 CDN 压测节点列表
 */

interface NodeItem {
  name: string
  status: string
  ip: string
  updatedAt: string
}

// 示例统计数据
const statistics = ref([
  { title: '总节点数', value: 0, icon: 'mdi-server', color: 'primary' },
  { title: '在线节点', value: 0, icon: 'mdi-server-network', color: 'success' },
  { title: '离线节点', value: 0, icon: 'mdi-server-off', color: 'error' },
  { title: '告警数量', value: 0, icon: 'mdi-alert', color: 'warning' },
])

const nodes = ref<NodeItem[]>([])

// 加载数据
const loadData = async () => {
  // TODO: 实现数据加载逻辑
  console.log('Loading overview data...')
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="overview-page">
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-6">数据概览</h1>
      </v-col>
    </v-row>

    <!-- 统计卡片 -->
    <v-row>
      <v-col v-for="stat in statistics" :key="stat.title" cols="12" sm="6" md="3">
        <v-card>
          <v-card-text>
            <div class="d-flex align-center">
              <v-icon :color="stat.color" size="48" class="mr-4">{{ stat.icon }}</v-icon>
              <div>
                <div class="text-h5">{{ stat.value }}</div>
                <div class="text-body-2 text-medium-emphasis">{{ stat.title }}</div>
              </div>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- 节点列表 -->
    <v-row class="mt-4">
      <v-col cols="12">
        <v-card>
          <v-card-title>CDN 压测节点列表</v-card-title>
          <v-card-text>
            <v-data-table
              :headers="[
                { title: '节点名称', key: 'name' },
                { title: '状态', key: 'status' },
                { title: 'IP 地址', key: 'ip' },
                { title: '最后更新', key: 'updatedAt' },
                { title: '操作', key: 'actions', sortable: false },
              ]"
              :items="nodes"
              no-data-text="暂无数据"
            >
              <template #item.status="{ item }">
                <v-chip :color="item.status === 'online' ? 'success' : 'error'" size="small">
                  {{ item.status === 'online' ? '在线' : '离线' }}
                </v-chip>
              </template>
              <template #item.actions>
                <v-btn icon="mdi-pencil" size="small" variant="text"></v-btn>
                <v-btn icon="mdi-delete" size="small" variant="text"></v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<style scoped>
.overview-page {
  width: 100%;
}
</style>
