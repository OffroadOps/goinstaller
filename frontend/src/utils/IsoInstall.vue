<template>
  <div class="iso-install">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><CdRom /></el-icon>
              <span>ISO安装功能</span>
            </div>
          </template>
          
          <div class="coming-soon">
            <el-result
              icon="info"
              title="功能开发中"
              sub-title="ISO安装功能正在开发中，敬请期待！"
            >
              <template #extra>
                <el-button type="primary" @click="goToVhd">
                  使用VHD重装
                </el-button>
              </template>
            </el-result>
            
            <div class="feature-preview">
              <h3>即将支持的功能：</h3>
              <ul>
                <li>Windows ISO镜像安装</li>
                <li>自定义分区设置</li>
                <li>驱动自动注入</li>
                <li>无人值守安装</li>
                <li>多版本系统选择</li>
              </ul>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAppStore } from '../stores/app'

const router = useRouter()
const appStore = useAppStore()

const goToVhd = () => {
  router.push({ name: 'vhd' })
}

// 添加日志
appStore.addLog('info', 'ISO安装页面已加载')
</script>

<style scoped>
.iso-install {
  padding: 20px;
  height: 100%;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.coming-soon {
  text-align: center;
  padding: 40px 20px;
}

.feature-preview {
  margin-top: 40px;
  text-align: left;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

.feature-preview h3 {
  margin-bottom: 16px;
  color: var(--el-text-color-primary);
}

.feature-preview ul {
  list-style: none;
  padding: 0;
}

.feature-preview li {
  padding: 8px 0;
  border-bottom: 1px solid var(--el-border-color-lighter);
  color: var(--el-text-color-regular);
}

.feature-preview li:before {
  content: "🔸";
  margin-right: 8px;
}
</style>