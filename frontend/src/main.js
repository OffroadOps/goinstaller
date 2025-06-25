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

// 挂载应用
app.mount('#app')

// 隐藏加载动画 - 修复元素选择器
setTimeout(() => {
  // 移除整个初始加载容器
  const container = document.querySelector('.container')
  if (container) {
    container.style.opacity = '0'
    setTimeout(() => {
      container.style.display = 'none'
      // 显示Vue应用
      document.getElementById('app').style.display = 'block'
    }, 300)
  }
}, 300)