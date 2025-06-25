import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

import App from './App.vue'
import routes from './routes'

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 创建应用实例
const app = createApp(App)

// 注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 使用插件
app.use(createPinia())
app.use(router)
app.use(ElementPlus, {
  locale: zhCn,
})

// 等待DOM加载完成后再挂载
document.addEventListener('DOMContentLoaded', () => {
  // 挂载应用
  app.mount('#app')
  
  // 隐藏加载动画
  setTimeout(() => {
    const loadingContainer = document.getElementById('loading-container')
    if (loadingContainer) {
      loadingContainer.style.opacity = '0'
      setTimeout(() => {
        loadingContainer.style.display = 'none'
        // 显示Vue应用
        const appElement = document.getElementById('app')
        if (appElement) {
          appElement.style.display = 'block'
        }
      }, 300)
    }
  }, 1000) // 增加延迟确保应用完全加载
})