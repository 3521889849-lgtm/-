# 酒店住宿管理后台系统 - Web 前端

## 技术栈

- **React 18** - UI 框架
- **Ant Design 5** - UI 组件库
- **React Router 6** - 路由管理
- **Axios** - HTTP 请求
- **Vite** - 构建工具
- **Day.js** - 日期处理

## 项目结构

```
web/
├── src/
│   ├── components/        # 公共组件
│   │   ├── Layout.jsx    # 布局组件（顶部导航、左侧菜单）
│   │   └── RoomForm.jsx  # 房源表单组件
│   ├── pages/            # 页面组件
│   │   ├── RoomStatus.jsx              # 房态/房源管理页面
│   │   ├── RoomTypeManagement.jsx      # 房型管理页面
│   │   ├── FacilityManagement.jsx      # 设施管理页面
│   │   └── CancellationPolicyManagement.jsx  # 退订政策管理页面
│   ├── api/             # API 接口封装
│   │   └── room.js      # 房源相关 API
│   ├── App.jsx          # 应用入口
│   ├── main.jsx         # 入口文件
│   └── index.css        # 全局样式
├── index.html           # HTML 模板
├── package.json         # 依赖配置
├── vite.config.js       # Vite 配置
└── README.md           # 说明文档
```

## 安装依赖

```bash
cd web
npm install
```

## 启动开发服务器

```bash
npm run dev
```

前端服务将在 `http://localhost:3000` 启动。

## 构建生产版本

```bash
npm run build
```

## 功能特性

### 1. 房源管理（房态页面）
- ✅ 房源列表展示
- ✅ 状态筛选（启用/停用/维修）
- ✅ 房源名称搜索
- ✅ 添加时间筛选
- ✅ 添加/修改/删除房源
- ✅ 状态切换（启用/停用）
- ✅ 关联房源功能
- ✅ 分页展示

### 2. 房型管理
- ✅ 房型列表展示
- ✅ 添加/修改/删除房型
- ✅ 状态筛选和搜索

### 3. 设施管理
- ✅ 设施列表展示
- ✅ 添加/修改/删除设施
- ✅ 在房源表单中勾选设施

### 4. 退订政策管理
- ✅ 退订政策列表展示
- ✅ 添加/修改/删除退订政策
- ✅ 支持自定义规则描述和违约金比例
- ✅ 房型筛选

## 界面说明

### 布局结构
- **顶部导航栏**：包含 logo、分店切换、主要标签（房态/订单/报表）、用户信息
- **左侧菜单**：设置管理、基础设置、会员管理、财务管理、系统设置等
- **主内容区**：根据路由显示不同页面

### 房源管理页面
- 顶部：添加房源按钮
- 搜索区：状态、房源名称、添加时间筛选
- 同步区：同步到途游按钮
- 表格：房源列表，包含操作按钮（查看/修改/删除/启用/停用/关联房源）
- 分页：支持分页和快速跳转

## API 对接

前端通过 `/api/v1` 前缀调用后端 API，Vite 开发服务器已配置代理到 `http://localhost:8080`。

## 注意事项

1. 确保后端 API 服务已启动（`http://localhost:8080`）
2. 图片上传功能需要后端支持 multipart/form-data
3. 所有 API 调用都包含错误处理
4. 表单验证使用 Ant Design 的 Form 组件
