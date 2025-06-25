# SystemReinstaller v2.0

基于 Wails + Vue 3 的现代化系统重装助手，支持 Windows 驱动备份、VHD/ISO 系统安装等功能。

## 🚀 功能特性

### 核心功能
- **VHD 系统重装**: 支持从多个服务器下载并安装 Windows VHD 镜像
- **驱动管理**: 自动备份和恢复系统驱动，支持 7z 压缩
- **ISO 安装**: Windows ISO 镜像安装（开发中）
- **Linux 安装**: Linux 系统安装（开发中）

### 技术特性
- **现代化 UI**: 基于 Vue 3 + Element Plus 的响应式界面
- **高性能**: Wails 框架提供原生性能
- **多主题**: 支持浅色/深色主题切换
- **实时日志**: 完整的操作日志记录和导出
- **进度监控**: 实时显示下载和操作进度

## 📋 系统要求

### 开发环境
- Go 1.21+
- Node.js 16+
- npm 或 yarn
- Wails v2.8.0+

### 运行环境
- Windows 10/11 (x64)
- 管理员权限（驱动操作需要）
- .NET Framework 4.7.2+

## 🛠️ 开发指南

### 快速开始

1. **克隆项目**
```bash
git clone https://github.com/your-username/SystemReinstaller.git
cd SystemReinstaller
```

2. **安装依赖**
```bash
# 安装 Go 依赖
go mod download

# 安装前端依赖
cd frontend
npm install
cd ..
```

3. **开发模式运行**
```bash
wails dev
```

4. **构建生产版本**
```bash
wails build
```

### 项目结构

```
SystemReinstaller/
├── main.go                    # 应用入口
├── app.go                     # 应用逻辑
├── wails.json                 # Wails 配置
├── go.mod                     # Go 模块定义
│
├── core/                      # 核心业务逻辑
│   ├── api_client.go          # API 客户端
│   ├── downloader.go          # 下载器
│   ├── driver_manager.go      # 驱动管理
│   ├── installer.go           # 安装器
│   └── vhd_parser.go          # VHD 解析
│
├── utils/                     # 工具函数
│   ├── logger.go              # 日志系统
│   ├── command_executor.go    # 命令执行
│   └── font.go                # 字体处理
│
├── resources/                 # 资源管理
│   └── ddlist.go              # DD 包管理
│
└── frontend/                  # 前端代码
    ├── src/
    │   ├── views/             # 页面组件
    │   ├── stores/            # 状态管理
    │   ├── utils/             # 工具函数
    │   └── styles/            # 样式文件
    └── dist/                  # 构建输出
```

### 主要技术栈

**后端 (Go)**
- Wails v2 - 跨平台应用框架
- 标准库 - 文件操作、网络请求等

**前端 (Vue 3)**
- Vue 3 - 渐进式 JavaScript 框架
- Element Plus - Vue 3 组件库
- Pinia - 状态管理
- Vue Router - 路由管理
- Vite - 构建工具

## 🔧 配置说明

### Wails 配置 (wails.json)
```json
{
  "name": "SystemReinstaller",
  "outputfilename": "SystemReinstaller",
  "frontend": {
    "dir": "./frontend",
    "install": "npm install",
    "build": "npm run build"
  }
}
```

### 前端配置 (package.json)
主要依赖：
- vue@^3.3.4
- element-plus@^2.3.8
- vue-router@^4.2.4
- pinia@^2.1.6

## 📱 功能模块

### 1. VHD 重装
- 多服务器支持
- 系统类型筛选（Windows/Windows Server/Tiny Windows）
- 语言选择（中文/英文/日文）
- 启动模式（UEFI/Legacy）
- 下载进度监控
- 安装参数配置

### 2. 驱动管理
- 系统驱动扫描和列表显示
- 一键备份所有驱动
- 支持 7z 压缩备份
- 备份历史管理
- 选择性恢复驱动
- 备份文件删除

### 3. 系统设置
- 主题切换（浅色/深色）
- 语言设置
- 下载配置（并发数、重试次数、超时时间）
- 备份配置（压缩、自动清理）
- 日志级别设置

### 4. 日志系统
- 实时日志显示
- 多级别日志（DEBUG/INFO/WARN/ERROR）
- 日志导出功能
- 日志自动清理

## 🎨 界面预览

### 主界面
- 现代化卡片式布局
- 响应式设计，支持不同屏幕尺寸
- 直观的操作流程
- 实时状态反馈

### 主题支持
- **浅色主题**: 清爽简洁的浅色界面
- **深色主题**: 护眼的深色界面
- 自动保存主题选择

## 🚦 使用指南

### 驱动备份
1. 打开"驱动管理"页面
2. 设置备份保存位置
3. 选择是否启用压缩
4. 点击"开始备份"
5. 等待备份完成

### VHD 重装
1. 打开"VHD重装"页面
2. 选择服务器
3. 设置筛选条件
4. 选择要安装的 VHD 版本
5. 配置高级选项
6. 下载并安装

### 驱动恢复
1. 在"驱动管理"页面的备份历史中选择备份
2. 点击"恢复驱动"
3. 确认恢复操作
4. 重启系统（推荐）

## 🔒 权限说明

本应用需要管理员权限来执行以下操作：
- 备份系统驱动文件
- 安装/恢复驱动
- 访问系统关键目录
- 执行系统级命令

## 🐛 故障排除

### 常见问题

**Q: 驱动备份失败，提示权限不足**
A: 请确保以管理员身份运行应用程序

**Q: 下载速度很慢**
A: 可以尝试切换其他服务器，或检查网络连接

**Q: 7z 压缩失败**
A: 确保系统已安装 7-Zip，或将 7z.exe 放置在 bin 目录下

**Q: 界面显示异常**
A: 尝试切换主题或重启应用

### 日志查看
应用会自动生成日志文件，保存在 `logs/` 目录下：
- 文件名格式：`app_YYYYMMDD.log`
- 包含详细的操作记录和错误信息

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发规范
1. 使用 Go fmt 格式化 Go 代码
2. 使用 ESLint 和 Prettier 格式化前端代码
3. 提交前确保所有测试通过
4. 遵循 [Conventional Commits](https://conventionalcommits.org/) 规范

### 提交类型
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🙏 致谢

感谢以下开源项目：
- [Wails](https://wails.io/) - 跨平台应用框架
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Element Plus](https://element-plus.org/) - Vue 3 组件库
- [7-Zip](https://www.7-zip.org/) - 文件压缩工具

## 📞 联系方式

- 项目主页: https://github.com/your-username/SystemReinstaller
- 问题反馈: https://github.com/your-username/SystemReinstaller/issues
- 邮箱: support@systemreinstaller.com

## 🗺️ 开发路线图

### v2.1.0 (计划中)
- [ ] ISO 安装功能完善
- [ ] Linux 系统安装支持
- [ ] 多语言界面支持
- [ ] 自动更新功能

### v2.2.0 (计划中)
- [ ] 网络安装支持
- [ ] 批量操作功能
- [ ] 云端备份同步
- [ ] 移动端适配

### v3.0.0 (远期规划)
- [ ] 插件系统
- [ ] 自定义脚本支持
- [ ] 企业版功能
- [ ] API 接口开放

## 📊 性能指标

- 应用启动时间: < 3秒
- 内存占用: < 100MB
- 驱动备份速度: 取决于驱动数量和磁盘性能
- 下载速度: 取决于网络环境和服务器性能

## 🔐 安全说明

### 数据安全
- 本地操作，不上传个人数据
- 备份文件仅保存在本地
- 支持备份文件加密（计划中）

### 网络安全
- 仅连接可信服务器
- 支持 HTTPS 加密传输
- 文件完整性校验

## 🌐 国际化

目前支持的语言：
- 简体中文 (zh-CN) ✅
- English (en-US) 🚧

计划支持的语言：
- 繁体中文 (zh-TW)
- 日本語 (ja-JP)
- 한국어 (ko-KR)
- Deutsch (de-DE)
- Français (fr-FR)

## 📈 版本历史

### v2.0.0 (当前版本)
- 🎉 全新 Wails + Vue 3 架构
- ✨ 现代化用户界面
- 🚀 性能大幅提升
- 🔧 完善的驱动管理功能
- 📱 响应式设计

### v1.x (Legacy)
- 基于 Fyne 的传统界面
- 基础驱动备份功能
- 简单的 DD 重装功能

---

**Happy Coding! 🎉**

如果这个项目对您有帮助，请考虑给我们一个 ⭐️ Star！"# goinstaller" 
