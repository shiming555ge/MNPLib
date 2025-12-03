# MNPLib 前端

MNPLib 前端是一个基于 Vue 3 的现代化化学分子库管理界面，提供化合物数据浏览、化学结构查询、分子可视化等功能。

## 技术栈

### 核心框架
- **Vue 3**: 使用 Composition API 和 `<script setup>` 语法
- **Vite**: 下一代前端构建工具，提供极速的开发体验
- **Vue Router 4**: 客户端路由管理
- **Vue I18n**: 国际化支持，提供中英文界面

### UI 框架和组件
- **Bootstrap 5**: 响应式 CSS 框架
- **Bootstrap Icons**: 图标库
- **SCSS/SASS**: CSS 预处理器，用于编写可维护的样式

### 化学相关库
- **Three.js**: 3D 图形库，用于分子可视化
- **OpenChemLib**: 化学信息学库，提供分子处理功能
- **Ketcher**: 化学结构编辑器（集成在 public/Ketcher/ 目录中）

## 项目结构详解

```
frontend/
├── public/                    # 静态资源（不经过构建处理）
│   ├── Ketcher/              # 化学结构编辑器
│   │   ├── static/           # 编辑器静态资源
│   │   ├── templates/        # 化学模板文件
│   │   └── index.html        # 编辑器入口
│   ├── home_pics/            # 首页背景图片
│   ├── logo.png              # 网站Logo
│   ├── vite.svg              # Vite图标
│   └── favicon.ico           # 网站图标
├── src/                      # 源代码目录
│   ├── assets/               # 经过构建处理的静态资源
│   │   └── vue.svg           # Vue图标
│   ├── components/           # Vue组件
│   │   ├── About.vue         # 关于页面
│   │   ├── AuthModal.vue     # 认证模态框
│   │   ├── Browse.vue        # 数据浏览页面
│   │   ├── compoundCard.vue  # 化合物卡片组件
│   │   ├── CompoundDetail.vue # 化合物详情页面
│   │   ├── FilterPanel.vue   # 筛选面板组件
│   │   ├── Home.vue          # 首页
│   │   ├── MoleculeCanvas.vue # 分子画布组件
│   │   ├── MoleculeViewer.vue # 分子可视化组件
│   │   ├── Navbar.vue        # 导航栏
│   │   ├── NotFound.vue      # 404页面
│   │   ├── Query.vue         # 化学查询页面
│   │   └── SuperAdmin.vue    # 管理员面板
│   ├── composables/          # 组合式函数
│   │   └── useAuth.js        # 认证相关逻辑
│   ├── locales/              # 国际化文件
│   │   ├── en_us.json        # 英文翻译
│   │   └── zh_cn.json        # 中文翻译
│   ├── router/               # 路由配置
│   │   └── index.js          # 路由定义
│   ├── scss/                 # 样式文件
│   │   ├── navbar.scss       # 导航栏样式
│   │   └── styles.scss       # 全局样式
│   ├── App.vue               # 根组件
│   └── main.js               # 应用入口
├── index.html                # HTML入口模板
├── vite.config.js            # Vite配置
├── package.json              # 依赖配置
└── package-lock.json         # 依赖锁文件
```

## 主要功能组件

### 1. 导航栏 (Navbar.vue)
- 提供网站主要导航链接
- 用户认证状态显示
- 语言切换功能
- 响应式设计，适配移动设备

### 2. 首页 (Home.vue)
- 项目介绍和功能概览
- 特色功能展示
- 背景图片轮播

### 3. 数据浏览 (Browse.vue)
- 化合物数据表格展示
- 分页功能
- 高级筛选面板
- 实时搜索

### 4. 化学查询 (Query.vue)
- 化学结构编辑器集成
- 相似度搜索
- 子结构匹配
- SMILES 转换工具

### 5. 分子可视化 (MoleculeViewer.vue)
- 3D 分子结构展示
- 旋转、缩放、平移交互
- 原子和键的样式自定义
- 多种分子表示模式（球棍模型、空间填充等）

### 6. 认证模块 (AuthModal.vue)
- Passkey 登录
- JWT 令牌管理
- 用户状态持久化

### 7. 管理员面板 (SuperAdmin.vue)
- Passkey 管理（创建、编辑、删除）
- 用户权限管理
- 系统状态监控

## 开发指南

### 环境要求
- Node.js 18.0.0 或更高版本
- npm 9.0.0 或更高版本

### 安装依赖
```bash
cd frontend
npm install
```

### 开发服务器
```bash
npm run dev
```
开发服务器将在 `http://localhost:5173` 启动，支持热重载。

### 构建生产版本
```bash
npm run build
```
构建结果将输出到 `dist/` 目录。

### 预览生产版本
```bash
npm run preview
```
预览构建结果，确保生产环境运行正常。

## 配置说明

### Vite 配置 (vite.config.js)
```javascript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:9090', // 后端API地址
        changeOrigin: true
      }
    }
  }
})
```

### 环境变量
创建 `.env` 文件来配置环境变量：
```env
VITE_API_BASE_URL=http://localhost:9090
VITE_APP_TITLE=MNPLib
```

在代码中使用：
```javascript
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL
```

## 与后端集成

### API 调用
前端通过代理配置与后端 API 通信：
- 开发环境：Vite 代理到 `http://localhost:9090`
- 生产环境：需要配置反向代理（如 Nginx）

### 认证流程
1. 用户输入 passkey 登录
2. 前端调用 `/api/auth/login` 获取 JWT token
3. 将 token 存储在 localStorage 中
4. 后续请求在 Authorization 头中携带 token
5. 使用 `useAuth` composable 管理认证状态

### 数据流
1. 组件通过 composable 或直接调用 API 获取数据
2. 使用 Vue 的响应式系统更新 UI
3. 错误处理通过 try-catch 和全局错误处理器

## 样式指南

### SCSS 结构
- `styles.scss`: 全局样式和变量定义
- `navbar.scss`: 导航栏特定样式
- 组件样式使用 `<style scoped>` 在组件内定义

### Bootstrap 定制
通过覆盖 Bootstrap 变量来自定义主题：
```scss
// 在 styles.scss 中
$primary: #007bff;
$secondary: #6c757d;
$font-family-sans-serif: 'Arial', sans-serif;

@import 'bootstrap/scss/bootstrap';
```

### 响应式设计
- 使用 Bootstrap 的栅格系统
- 移动设备优先的设计原则
- 断点：xs (<576px), sm (≥576px), md (≥768px), lg (≥992px), xl (≥1200px)

## 国际化

### 添加新语言
1. 在 `src/locales/` 目录下创建新的 JSON 文件，如 `fr_fr.json`
2. 在 `main.js` 中注册新语言：
   ```javascript
   import frFR from './locales/fr_fr.json'
   
   const i18n = createI18n({
     locale: 'fr',
     messages: {
       fr: frFR
     }
   })
   ```

### 在组件中使用
```vue
<template>
  <h1>{{ $t('welcome.title') }}</h1>
</template>

<script setup>
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
console.log(t('welcome.message'))
</script>
```

## 性能优化

### 代码分割
Vite 自动进行代码分割，按需加载路由组件。

### 图片优化
- 使用适当的图片格式（WebP、AVIF）
- 实现懒加载
- 使用 CDN 或图片压缩工具

### 构建优化
```bash
# 分析构建大小
npm run build -- --report
```

## 测试

### 单元测试
```bash
# 安装测试依赖
npm install -D vitest @vue/test-utils

# 运行测试
npm test
```

### 端到端测试
```bash
# 安装 Cypress
npm install -D cypress

# 打开 Cypress
npx cypress open
```

## 部署

### 静态文件部署
1. 构建项目：`npm run build`
2. 将 `dist/` 目录内容上传到静态文件服务器
3. 配置服务器重写规则，将所有请求重定向到 `index.html`

### Docker 部署
创建 `Dockerfile`：
```dockerfile
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 常见问题

### 1. 开发服务器无法启动
- 检查端口 5173 是否被占用
- 确保 Node.js 版本符合要求
- 清除 node_modules 并重新安装依赖

### 2. 无法连接到后端 API
- 检查后端服务是否运行在 `http://localhost:9090`
- 确认 Vite 代理配置正确
- 查看浏览器控制台网络错误

### 3. 样式不生效
- 检查 SCSS 文件是否正确导入
- 确认 Bootstrap 变量覆盖顺序
- 查看组件样式是否使用了 `scoped` 属性

### 4. 国际化文本不显示
- 检查语言文件路径和格式
- 确认在 `main.js` 中正确注册了语言
- 查看翻译键名是否正确

## 贡献指南

1. 遵循 Vue 3 组合式 API 最佳实践
2. 使用 TypeScript 或 JSDoc 进行类型提示
3. 为新功能添加相应的测试
4. 更新相关文档
5. 遵循项目的代码风格和提交规范

## 许可证

本项目采用 MIT 许可证 - 查看根目录的 [LICENSE](../LICENSE) 文件了解详情。

---

**MNPLib 前端** - 为化学研究提供现代化的用户界面体验
