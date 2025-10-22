# 高校课程互助与笔记分享平台

## 项目简介

本项目是一个基于Go语言和Vue.js开发的高校课程互助与笔记分享平台，采用前后端分离架构。旨在为高校学生提供一个便捷的课程资料分享和学习交流的平台。学生可以在平台上分享课程笔记、学习资料，参与课程讨论，互相帮助解决学习中的问题。

## 功能特性

### 核心功能
- **用户管理**：用户注册、登录、个人资料管理、头像上传
- **课程管理**：课程创建、加入、查询、搜索、分类管理
- **笔记分享**：笔记上传、下载、在线预览、点赞、收藏
- **评论互动**：对笔记进行评论、回复、点赞和评分
- **搜索功能**：支持课程和笔记的关键词搜索
- **权限管理**：学生、教师、管理员三级权限体系

### 管理功能
- **管理员控制台**：用户管理、课程管理、数据统计
- **教师功能**：课程创建与管理、学生管理
- **数据统计**：用户数、课程数、笔记数、评论数统计

## 页面功能介绍

### 首页 (Home)
- **功能概述**：展示平台概况、热门课程、最新笔记和推荐内容
- **主要特点**：
  - 轮播展示精选课程和笔记
  - 显示最新上传的笔记
  - 展示热门课程和推荐内容
  - 快速访问常用功能的入口

### 用户认证

#### 登录页 (Login)
- **功能概述**：用户登录系统
- **主要特点**：
  - 用户名/密码登录
  - 记住登录状态
  - 登录错误提示
  - 跳转到注册页面的链接

#### 注册页 (Register)
- **功能概述**：新用户注册
- **主要特点**：
  - 用户基本信息填写（用户名、密码、邮箱等）
  - 表单验证
  - 注册成功后自动登录
  - 展示平台主要功能和优势

### 个人中心 (Profile)
- **功能概述**：用户个人信息管理
- **主要特点**：
  - 个人资料查看和编辑
  - 头像上传和管理
  - 我的课程列表
  - 我的笔记列表
  - 我的收藏和点赞内容
  - 账户安全设置

### 课程相关

#### 课程列表 (Courses)
- **功能概述**：浏览和搜索所有课程
- **主要特点**：
  - 课程分类筛选
  - 关键词搜索
  - 分页浏览
  - 课程卡片展示（包含课程名称、教师、学生数量等信息）
  - 快速加入课程

#### 课程详情 (CourseDetail)
- **功能概述**：查看课程详细信息和相关笔记
- **主要特点**：
  - 课程基本信息展示
  - 课程相关笔记列表
  - 加入/退出课程
  - 推荐相关课程
  - 课程统计数据（学生数、笔记数等）

### 笔记相关

#### 笔记列表 (Notes)
- **功能概述**：浏览和搜索所有公开笔记
- **主要特点**：
  - 按课程、标签筛选
  - 关键词搜索
  - 分页浏览
  - 排序选项（最新、最热、最多点赞等）
  - 笔记卡片展示（标题、作者、点赞数等）

#### 我的笔记 (MyNotes)
- **功能概述**：管理用户自己创建的笔记
- **主要特点**：
  - 创建新笔记
  - 编辑/删除已有笔记
  - 查看笔记统计数据
  - 笔记状态管理（公开/私有）
  - 搜索和筛选功能

#### 笔记详情 (NoteDetail)
- **功能概述**：查看笔记详细内容和互动
- **主要特点**：
  - 笔记内容展示
  - 点赞和收藏功能
  - 评论区互动
  - 相关笔记推荐
  - 笔记下载
  - 分享功能

#### 创建/编辑笔记 (CreateNote/EditNote)
- **功能概述**：创建新笔记或编辑已有笔记
- **主要特点**：
  - 富文本编辑器
  - 文件上传
  - 标签管理
  - 关联课程选择
  - 笔记状态设置（公开/私有）

### 搜索功能 (Search)
- **功能概述**：全站内容搜索
- **主要特点**：
  - 多类型内容搜索（课程、笔记、用户）
  - 高级筛选选项
  - 搜索结果分类展示
  - 热门搜索标签
  - 搜索历史记录

### 管理员功能 (Admin)
- **功能概述**：系统管理和数据统计
- **主要特点**：
  - 用户管理（查看、编辑、删除用户，修改用户角色）
  - 课程管理（查看、编辑、删除、归档课程）
  - 笔记管理（查看、编辑、删除笔记）
  - 系统数据统计（用户数、课程数、笔记数、评论数等）
  - 内容审核

### 教师功能 (Teacher)
- **功能概述**：教师课程和学生管理
- **主要特点**：
  - 课程创建和管理
  - 课程学生管理
  - 课程笔记管理
  - 教学数据统计
  - 教学资源管理

## 技术栈

### 后端技术
- **编程语言**：Go 1.20+
- **Web框架**：Gin
- **数据库**：MySQL 5.7+
- **ORM**：GORM
- **认证**：JWT
- **文件上传**：本地存储
- **中间件**：CORS、认证、错误处理

### 前端技术
- **框架**：Vue.js 3.2+
- **UI组件库**：Element Plus 2.4+
- **路由**：Vue Router 4.2+
- **HTTP客户端**：Axios 1.6+
- **构建工具**：Vue CLI 5.0+
- **PWA支持**：Service Worker

## 项目结构

```
student_shared/
├── app/                           # 后端应用程序代码
│   ├── app.go                     # 应用程序入口
│   ├── middleware/                # 中间件
│   │   ├── auth.go               # JWT认证中间件
│   │   ├── cors.go               # 跨域中间件
│   │   ├── error_handler.go      # 错误处理中间件
│   │   └── optional_auth.go      # 可选认证中间件
│   ├── model/                     # 数据模型
│   │   ├── user.go               # 用户模型
│   │   ├── course.go             # 课程模型
│   │   ├── note.go               # 笔记模型
│   │   ├── comment.go            # 评论模型
│   │   └── favorite.go           # 收藏模型
│   ├── router/                    # 路由配置
│   │   ├── router.go             # 路由注册
│   │   └── api/                  # API处理函数
│   │       ├── user.go           # 用户API
│   │       ├── course.go         # 课程API
│   │       ├── note.go           # 笔记API
│   │       ├── comment.go        # 评论API
│   │       └── admin.go          # 管理员API
│   └── utils/                     # 工具函数
│       ├── database/             # 数据库工具
│       └── jwt/                  # JWT工具
├── student_shared_view/           # 前端Vue.js应用
│   ├── public/                   # 静态资源
│   ├── src/                      # 源代码
│   │   ├── api/                  # API接口
│   │   ├── assets/               # 资源文件
│   │   ├── components/           # Vue组件
│   │   ├── router/               # 前端路由
│   │   ├── views/                # 页面组件
│   │   │   ├── Home.vue          # 首页
│   │   │   ├── Login.vue         # 登录页
│   │   │   ├── Profile.vue       # 个人资料
│   │   │   ├── NoteDetail.vue    # 笔记详情
│   │   │   ├── Admin.vue         # 管理员控制台
│   │   │   └── Teacher.vue       # 教师页面
│   │   ├── App.vue               # 根组件
│   │   └── main.js               # 入口文件
│   ├── package.json              # 前端依赖配置
│   └── vue.config.js             # Vue配置文件
├── uploads/                       # 文件上传目录
│   └── avatars/                  # 用户头像
├── go.mod                        # Go模块文件
├── go.sum                        # Go依赖校验文件
├── main.go                       # 后端主程序入口
└── README.md                     # 项目说明文档
```

## 安装和运行

### 前置条件

- Go 1.20或更高版本
- Node.js 16.0或更高版本
- MySQL 5.7或更高版本
- Git

### 后端安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd student_shared
```

2. **安装Go依赖**
```bash
go mod tidy
```

3. **配置数据库**

在MySQL中创建数据库：
```sql
CREATE DATABASE student_shared CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改数据库连接配置（在`app/utils/database/database.go`中）：
```go
dsn := "username:password@tcp(localhost:3306)/student_shared?charset=utf8mb4&parseTime=True&loc=Local"
```

4. **运行后端服务**
```bash
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

### 前端安装步骤

1. **进入前端目录**
```bash
cd student_shared_view
```

2. **安装依赖**
```bash
npm install
# 或使用 yarn
yarn install
```

3. **启动开发服务器**
```bash
npm run serve
# 或使用 yarn
yarn serve
```

前端应用将在 `http://localhost:8080` 启动（如果端口被占用会自动选择其他端口）

4. **构建生产版本**
```bash
npm run build
# 或使用 yarn
yarn build
```

## API文档

### 基础信息
- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: JWT Token (Header: `Authorization: Bearer <token>`)

### 用户相关 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/users/register` | 用户注册 | 否 |
| POST | `/users/login` | 用户登录 | 否 |
| GET | `/users/profile` | 获取用户资料 | 是 |
| PUT | `/users/profile` | 更新用户资料 | 是 |

### 课程相关 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/courses` | 获取课程列表 | 否 |
| GET | `/courses/:id` | 获取课程详情 | 可选 |
| POST | `/courses` | 创建课程 | 是 |
| PUT | `/courses/:id` | 更新课程 | 是 |
| DELETE | `/courses/:id` | 删除课程 | 是 |
| POST | `/courses/:id/join` | 加入课程 | 是 |

### 笔记相关 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/notes` | 获取笔记列表 | 可选 |
| GET | `/notes/:id` | 获取笔记详情 | 可选 |
| POST | `/notes` | 创建笔记 | 是 |
| PUT | `/notes/:id` | 更新笔记 | 是 |
| DELETE | `/notes/:id` | 删除笔记 | 是 |
| POST | `/notes/:id/like` | 点赞笔记 | 是 |
| DELETE | `/notes/:id/like` | 取消点赞 | 是 |
| POST | `/notes/:id/favorite` | 收藏笔记 | 是 |
| DELETE | `/notes/:id/favorite` | 取消收藏 | 是 |
| GET | `/notes/favorites` | 获取我的收藏 | 是 |
| GET | `/notes/likes` | 获取我的点赞 | 是 |
| GET | `/notes/:id/download` | 下载笔记 | 是 |

### 评论相关 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/comment/note/:noteId` | 获取笔记评论列表 | 可选 |
| POST | `/comment` | 创建评论 | 是 |
| PUT | `/comment/:id` | 更新评论 | 是 |
| DELETE | `/comment/:id` | 删除评论 | 是 |
| POST | `/comment/:id/like` | 点赞评论 | 是 |
| DELETE | `/comment/:id/like` | 取消点赞评论 | 是 |

### 搜索相关 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/search/courses` | 搜索课程 | 否 |
| GET | `/search/notes` | 搜索笔记 | 否 |

### 文件上传 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | `/upload/avatar` | 上传头像 | 是 |
| DELETE | `/upload/avatar` | 删除头像 | 是 |

### 管理员 API

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | `/admin/stats` | 获取统计数据 | 管理员 |
| GET | `/admin/users` | 获取所有用户 | 管理员 |
| PUT | `/admin/users/:id` | 管理员更新用户 | 管理员 |

## 用户角色权限

### 学生 (student)
- 浏览公开课程和笔记
- 加入课程
- 上传和管理自己的笔记
- 评论和点赞
- 收藏笔记

### 教师 (teacher)
- 学生的所有权限
- 创建和管理课程
- 管理课程内的笔记

### 管理员 (admin)
- 教师的所有权限
- 用户管理
- 系统统计数据查看
- 全局内容管理

## 部署说明

### 生产环境部署

1. **后端部署**
```bash
# 构建二进制文件
go build -o student_shared main.go

# 运行
./student_shared
```

2. **前端部署**
```bash
# 构建生产版本
cd student_shared_view
npm run build

# 将 dist 目录部署到 Web 服务器
```

3. **Nginx 配置示例**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /path/to/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 文件上传目录
    location /uploads/ {
        root /path/to/backend;
    }
}
```

## 开发指南

### 代码规范
- Go代码遵循官方代码规范
- Vue.js代码使用ESLint规范
- 提交信息使用约定式提交格式

### 开发流程
1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建Pull Request

### 数据库迁移
项目使用GORM的自动迁移功能，首次运行时会自动创建所需的数据表。

## 常见问题

### Q: 如何重置管理员密码？
A: 可以直接在数据库中修改用户表的密码字段（需要使用bcrypt加密）。

### Q: 如何修改文件上传大小限制？
A: 在Gin路由中添加`MaxMultipartMemory`配置。

### Q: 前端如何配置API基础URL？
A: 在`src/api/index.js`中修改`baseURL`配置。

### Q: 删除课程后为什么在数据库中仍然存在？
A: 系统采用软删除机制，删除课程时只是将状态设置为`inactive`，而不是真正从数据库中删除记录。这样设计是为了数据安全和可恢复性考虑。

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件

## 贡献者

感谢所有为这个项目做出贡献的开发者！

## 联系方式

如有任何问题或建议，请通过以下方式联系：

- 项目Issues：[GitHub Issues](https://github.com/yourusername/student_shared/issues)
- 邮箱：your-email@example.com

---

**注意**：这是一个教育项目，仅供学习和研究使用。在生产环境中使用前，请确保进行充分的安全测试和性能优化。

## AI 摘要集成
# 高校课程互助与笔记分享平台

## 后端配置（统一 conf 文件）

- 配置文件路径：`app/conf/congfig.conf`
- 格式：`KEY=VALUE`（无分区），支持注释行（以 `#` 或 `;` 开头）。
- 环境变量覆盖：若存在同名环境变量，优先使用环境变量值。

示例内容：

```
# 数据库配置
MYSQL_USER=root
MYSQL_PASS=abc123456
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=shared_student

# 服务端口
SERVER_PORT=8080

# AI配置
OPENAI_API_KEY= # 在此填入真实密钥或通过环境变量注入
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-4o-mini
OPENAI_TIMEOUT_SECONDS=30
```

## 生效范围

- 数据库：`app/utils/database/database.go` 读取上述 MySQL 配置。
- 服务端口：`app/app.go` 读取 `SERVER_PORT`。
- AI 摘要：`app/utils/ai/provider_openai.go` 读取 `OPENAI_*` 配置（密钥为空时不调用真实AI，自动降级到本地算法）。

## 使用说明

- 开发环境可直接编辑 `congfig.conf`；生产环境建议使用环境变量覆盖敏感项（如 `OPENAI_API_KEY`）。
- 修改后无需重启即可被读取（当前实现为进程启动时加载一次，如需热更新可后续扩展）。

## 前端说明（摘要再生成）

- 笔记详情页新增“重新生成”按钮（作者/管理员可见），触发后端 `POST /api/v1/ai/summarize`，成功后刷新 `GET /api/v1/ai/notes/:id/meta`。具体见 `student_shared_view/src/views/NoteDetail.vue`。
- 搜索说明：`GET /api/v1/notes/search?keyword=xxx` 已支持在 AI 关键词（NoteAIMeta.keywords）中匹配，提升相关性。

## 语义检索（Semantic Search）

- 端点：`GET /api/v1/search/semantic`
- 查询参数：
  - `q`：搜索关键词（必填）
  - `type`：`notes` | `courses` | `all`（默认 `all`）
  - `page`：页码（默认 `1`）
  - `page_size`：每页数量（默认 `10`，最大 `100`）
- 响应字段：
  - 通用：`page`, `page_size`
  - 笔记：`notes`（数组），每项包含：`id`, `title`, `excerpt`, `similarity`, `highlighted_title`, `highlighted_excerpt`, `user_id`, `course_id`, `course_name`, `status`, `created_at`；以及 `total_notes`
  - 课程：`courses`（数组），每项包含：`id`, `name`, `description_excerpt`, `similarity`, `highlighted_name`, `highlighted_description`, `school`, `department`, `semester`, `teacher`, `created_at`；以及 `total_courses`
- 排序与高亮：后端按语义相似度自动排序，并返回 HTML 高亮字段（`highlighted_*`），前端直接用 `v-html` 渲染即可。
- 过滤支持：目前支持按类型（笔记/课程）与常规关键词预筛选，后续可扩展课程维度（学校/院系/学期）与笔记维度（作者/课程/状态）参数。