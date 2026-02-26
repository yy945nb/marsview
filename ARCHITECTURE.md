# Marsview 项目架构文档

[English](#architecture-overview) | 中文

## 项目概述

Marsview 是一款面向中后台管理系统的低代码可视化搭建平台。开发者可在平台上创建项目、页面和组件，支持事件交互、接口调用、数据联动和逻辑编排，并可通过微前端框架 microApp 集成到已有业务系统中。

## 技术栈

| 类别       | 技术                                    |
| ---------- | --------------------------------------- |
| 框架       | React 18                                |
| 语言       | TypeScript                              |
| 构建工具   | Vite 5                                  |
| UI 组件库  | Ant Design 5                            |
| 状态管理   | Zustand                                 |
| 路由       | React Router v6                         |
| 样式       | Less + styled-components                |
| 拖拽       | react-dnd                               |
| 代码编辑器 | Monaco Editor                           |
| 流程编排   | @xyflow/react                           |
| 图表       | @ant-design/plots                       |
| HTTP 请求  | Axios                                   |
| 包管理器   | pnpm (workspace monorepo)               |
| 代码规范   | ESLint + Prettier + lint-staged + yorkie |

## 仓库结构

```
marsview/
├── packages/                  # 前端 monorepo 工作区
│   ├── editor/                # 编辑器端 —— 可视化页面搭建
│   ├── admin/                 # 访问端 —— 项目管理和页面访问
│   └── materials/             # 组件物料库 —— 可复用 UI 组件
├── package.json               # 根 package.json（公共依赖 + 脚本）
├── pnpm-workspace.yaml        # pnpm workspace 配置
├── README.md                  # 中文说明
├── README.en-US.md            # 英文说明
└── CHANGELOG.md               # 版本变更日志
```

## 核心 Package 说明

### 1. `@marsview/editor`（编辑器端）

编辑器是平台的核心，提供可视化拖拽搭建、事件流配置、接口管理和逻辑编排能力。

- **开发端口**：8080
- **构建输出**：`dist/editor`

#### 目录结构

```
packages/editor/src/
├── api/                # 接口层（页面、项目、用户等 API）
├── components/         # 公共 UI 组件
├── config/             # 应用配置
├── layout/             # 布局组件
├── packages/           # 动态组件包（与 materials 结构对应）
│   ├── Advanced/       #   高级组件
│   ├── Basic/          #   基础组件
│   ├── Container/      #   容器组件
│   ├── EChart/         #   图表组件
│   ├── FeedBack/       #   反馈组件
│   ├── FormItems/      #   表单组件
│   ├── Functional/     #   功能组件
│   ├── Layout/         #   布局组件
│   ├── MarsRender/     #   渲染引擎
│   ├── Other/          #   其他组件
│   ├── Page/           #   页面组件
│   └── Scene/          #   场景组件
├── pages/              # 路由页面
│   ├── admin/          #   后台管理页面
│   ├── editor/         #   编辑器画布
│   ├── home/           #   首页
│   ├── login/          #   登录
│   └── welcome/        #   欢迎页
├── router/             # React Router 路由定义
├── stores/             # Zustand 状态管理（pageStore）
└── utils/              # 工具函数
```

#### 关键依赖

- `react-dnd` — 拖拽交互
- `@monaco-editor/react` — 在线代码编辑
- `@xyflow/react` — 流程 / 逻辑编排图
- `bytemd` — Markdown 编辑器
- `react-infinite-viewer` — 画布缩放与拖拽

### 2. `@marsview/admin`（访问端）

访问端是项目和页面管理入口，支持以微前端子应用形式嵌入到其他系统中。

- **开发端口**：8090
- **构建输出**：`dist/admin`

#### 目录结构

```
packages/admin/src/
├── api/                # 接口层（含 Mock 数据）
├── components/         # 公共组件（Header、Menu、Tab 等）
├── config/             # 配置
├── layout/             # 布局
├── pages/              # 路由页面
│   ├── console/        #   控制台
│   ├── login/          #   登录
│   ├── page/           #   页面管理
│   ├── project/        #   项目管理
│   └── welcome/        #   欢迎页
├── router/             # 路由定义
├── stores/             # Zustand 状态管理（projectStore）
└── utils/              # 工具函数
```

#### 关键特性

- 支持微前端集成（`window.mount` / `window.unmount`）
- 内置 Mock API 便于独立开发
- 依赖 `@marsview/materials` 组件库进行页面渲染

### 3. `@marsview/materials`（组件物料库）

物料库提供平台所有可复用的 UI 组件，按类型分组，通过 Vite 的 `import.meta.glob()` 实现动态加载。

#### 目录结构

```
packages/materials/
├── index.tsx           # 动态组件注册入口
├── Basic/              # 基础组件（Avatar, Icon, Image, Link, Text, Title, Statistic）
├── Container/          # 容器组件
├── EChart/             # 图表组件
├── FeedBack/           # 反馈组件（Alert, Modal 等）
├── FormItems/          # 表单组件
├── Functional/         # 功能性组件
├── Layout/             # 布局组件
├── MarsRender/         # 通用渲染引擎（处理组件渲染、事件、API 调用）
├── Other/              # 其他组件
├── Page/               # 页面级组件
├── Scene/              # 场景组件
├── stores/             # 共享状态管理（pageStore）
├── types/              # 类型定义
└── utils/              # 工具库（请求封装、存储、Action 处理等）
```

## 架构设计

### 整体架构

```
┌─────────────────────────────────────────────────┐
│                   用户浏览器                      │
│  ┌──────────────────┐  ┌──────────────────────┐  │
│  │   Editor (8080)  │  │    Admin (8090)       │  │
│  │   可视化搭建编辑器 │  │    项目 / 页面管理    │  │
│  └────────┬─────────┘  └──────────┬───────────┘  │
│           │                       │               │
│           └───────────┬───────────┘               │
│                       │                           │
│              ┌────────▼────────┐                  │
│              │   Materials     │                  │
│              │   组件物料库     │                  │
│              └─────────────────┘                  │
└─────────────────────────────────────────────────┘
                        │
                        ▼
               ┌─────────────────┐
               │   后端 API       │
               │  (Node.js/MySQL) │
               └─────────────────┘
```

> **注意**：后端代码不包含在本开源仓库中。本仓库仅包含前端部分，默认连接线上 API。

### 状态管理

- 使用 **Zustand** 进行轻量级状态管理
- Editor 端维护 `pageStore`（页面数据、组件树、事件流）
- Admin 端维护 `projectStore`（项目配置、菜单、权限）
- Materials 共享 `pageStore` 供渲染时使用

### 组件动态加载

Materials 通过 `import.meta.glob()` 按需加载组件：

```typescript
// packages/materials/index.tsx
const modules = import.meta.glob('./**/*.tsx');
```

组件按分类目录组织，在编辑器中通过拖拽添加到画布，在访问端通过 MarsRender 进行统一渲染。

### 构建与部署

| 命令                 | 说明              | 输出目录      |
| -------------------- | ----------------- | ------------- |
| `pnpm start:editor`  | 启动编辑器开发    | —             |
| `pnpm start:admin`   | 启动访问端开发    | —             |
| `pnpm build:editor`  | 构建编辑器        | `dist/editor` |
| `pnpm build:admin`   | 构建访问端        | `dist/admin`  |
| `pnpm build`         | 构建全部          | `dist/`       |

### 环境要求

- Node.js >= 18.0.0
- 包管理器：pnpm（通过 `preinstall` 脚本强制使用）

---

## Architecture Overview

Marsview is a low-code visual development platform for building mid-to-back-office admin systems. It is organized as a **pnpm workspace monorepo** with three frontend packages:

| Package                | Purpose                             | Dev Port |
| ---------------------- | ----------------------------------- | -------- |
| `@marsview/editor`     | Visual drag-and-drop page builder   | 8080     |
| `@marsview/admin`      | Project and page management portal  | 8090     |
| `@marsview/materials`  | Shared component library            | —        |

### Tech Stack

- **Framework**: React 18 + TypeScript
- **Build**: Vite 5
- **UI**: Ant Design 5
- **State**: Zustand
- **Routing**: React Router v6
- **Drag & Drop**: react-dnd
- **Code Editor**: Monaco Editor
- **Flow Editor**: @xyflow/react
- **Charts**: @ant-design/plots
- **Styling**: Less + styled-components

### Key Design Decisions

1. **Monorepo with pnpm workspaces** — shared dependencies and cross-package references (`workspace:^`).
2. **Dynamic component loading** — Materials uses `import.meta.glob()` for on-demand component registration.
3. **Zustand for state management** — lightweight stores (`pageStore`, `projectStore`) instead of Redux.
4. **Micro-frontend ready** — Admin package supports `window.mount`/`window.unmount` for microApp integration.
5. **Vite-based builds** — fast development server and optimized production builds.
