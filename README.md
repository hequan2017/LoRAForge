# LoRAForge (锻造)

基于 Gin-Vue-Admin 2.8.9 框架开发的云资源管理系统。

## 📌 项目简介

本项目旨在提供一套完整的云资源管理解决方案，涵盖从算力节点接入、产品规格定义、镜像管理到用户实例全生命周期管理的完整流程。

## 🚀 核心功能模块

### 1. ☁️ 云资源管理 (Cloud Resource Management)

本项目新增了 `cloud` 包，包含以下核心业务模块：

#### 🖥️ 实例管理 (Instance Management)
用户购买和使用的云主机实例管理。
*   **功能**：创建、启动、停止、重启、删除实例。
*   **关联**：
    *   1对1 关联 **镜像库** (Mirror Repository)
    *   1对1 关联 **产品规格** (Product Specification)
    *   1对1 关联 **算力节点** (Compute Node)
*   **自动化**：后端自动填充当前用户ID，自动回填 Docker 容器信息及状态。

#### 📏 产品规格 (Product Specifications)
定义云产品的硬件规格和定价。
*   **字段**：名称、显卡型号/数量、CPU核心数、内存(GB)、系统盘/数据盘容量、价格/小时。
*   **特性**：支持上架/下架管理。

#### 💻 算力节点 (Compute Nodes)
管理底层的计算资源节点。
*   **字段**：区域、IP配置(公网/内网)、SSH端口、硬件配置(CPU/内存/显卡)。
*   **安全连接**：
    *   支持 Docker TLS 连接。
    *   **证书管理**：CA证书、客户端证书、客户端私钥采用**长文本(Textarea)**输入，方便粘贴和管理。

#### 💿 镜像库 (Mirror Repository)
管理可用的操作系统或应用镜像。
*   **字段**：名称、地址、来源、描述。
*   **特性**：支持上架/下架管理。

## 🛠️ 技术特性

*   **框架**：Gin-Vue-Admin (Gin + Vue3 + Element Plus)
*   **模块化设计**：业务逻辑独立于 `cloud` 包中，不侵入系统核心代码。
*   **自动化初始化**：
    *   实现了 `SubInitializer` 接口 (`DataInserted`)。
    *   支持通过 `gowatch` 热编译自动注册 API 权限和前端菜单。
    *   解决了系统初始化时的数据依赖和接口兼容性问题。
*   **前端优化**：
    *   针对证书类长文本字段，定制了 `el-input type="textarea"` 输入组件。
    *   完善了各模块的详情查看页 (`el-descriptions`)，确保数据展示完整。

## 📂 目录结构

*   `server/source/cloud`: 后端初始化数据 (API, 菜单)
*   `server/api/v1/cloud`: 后端业务 API 接口
*   `server/service/cloud`: 后端业务逻辑层
*   `server/model/cloud`: 数据库模型定义
*   `web/src/view/cloud`: 前端 Vue 页面组件

## 🚀 快速开始

1.  **启动后端**:
    ```bash
    cd server
    go mod tidy
    go run main.go
    # 或者使用 gowatch 热加载
    gowatch
    ```

2.  **启动前端**:
    ```bash
    cd web
    npm install
    npm run dev
    ```

3.  **访问系统**:
    打开浏览器访问 `http://localhost:8080` (默认端口)。
