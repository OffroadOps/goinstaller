import VhdInstall from './utils/VhdInstall.vue'
import IsoInstall from './utils/IsoInstall.vue'
import LinuxInstall from './utils/LinuxInstall.vue'
import DriverManager from './utils/DriverManager.vue'
import Settings from './utils/Settings.vue'

export default [
  {
    path: '/',
    redirect: '/vhd'
  },
  {
    path: '/vhd',
    name: 'vhd',
    component: VhdInstall,
    meta: {
      title: 'VHD重装',
      icon: 'Download'
    }
  },
  {
    path: '/iso',
    name: 'iso',
    component: IsoInstall,
    meta: {
      title: 'ISO安装',
      icon: 'CdRom'
    }
  },
  {
    path: '/linux',
    name: 'linux',
    component: LinuxInstall,
    meta: {
      title: 'Linux安装',
      icon: 'Monitor'
    }
  },
  {
    path: '/driver',
    name: 'driver',
    component: DriverManager,
    meta: {
      title: '驱动管理',
      icon: 'Setting'
    }
  },
  {
    path: '/settings',
    name: 'settings',
    component: Settings,
    meta: {
      title: '系统设置',
      icon: 'Tools'
    }
  }
]