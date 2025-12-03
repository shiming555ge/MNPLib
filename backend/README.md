# MNPLib 后端

MNPLib 后端是一个基于 Go 语言的化学分子库管理系统的服务器端，提供化合物数据管理、化学计算、用户认证等 RESTful API 服务。

## 技术栈

### 核心框架
- **Go 1.24+**: 高性能编程语言
- **Gin**: 轻量级 Web 框架，提供路由、中间件等功能
- **GORM**: ORM 框架，简化数据库操作
- **Viper**: 配置管理库，支持多种配置文件格式

### 数据库
- **MySQL 8.0+**: 关系型数据库，存储化合物数据和用户信息
- **数据库驱动**: `github.com/go-sql-driver/mysql`

### 认证和安全
- **JWT (JSON Web Tokens)**: 用户认证和授权
- **UUID**: 生成唯一的 passkey 标识符
- **加密**: 数据加密和安全性保护

### 化学计算
- **RDKit**: 化学信息学工具包（通过 Python 集成）
- **Python 3.8+**: 运行 RDKit 脚本

### 工具和工具库
- **Zerolog**: 结构化日志记录
- **Google UUID**: UUID 生成库

## 项目结构

```
backend/
├── config/                    # 配置管理
│   └── config.go             # 配置加载和初始化
├── controllers/              # 控制器层（处理 HTTP 请求）
│   ├── authController.go     # 认证相关控制器
│   ├── dataController.go     # 数据相关控制器
│   ├── passkeyController.go  # Passkey 管理控制器
│   ├── rdkitController.go    # RDKit 化学计算控制器
│   └── simple_data_controller.go # 简单数据控制器
├── database/                 # 数据库连接和操作
│   └── database.go           # 数据库初始化和连接
├── middlewares/              # 中间件
│   ├── extends_check.go      # 权限检查中间件
│   ├── jwt_auth.go           # JWT 认证中间件
│   └── validPath.go          # 路径验证中间件
├── models/                   # 数据模型（GORM 结构体）
│   ├── database.go           # 化合物数据模型
│   └── passkey.go            # Passkey 模型
├── router/                   # 路由定义
│   └── router.go             # 路由配置和注册
├── services/                 # 业务逻辑层
│   ├── initService.go        # 初始化服务
│   └── rdkitService.go       # RDKit 化学计算服务
├── static/                   # 静态文件
│   └── passkey-admin.html    # Passkey 管理页面(未使用)
├── utils/                    # 工具函数
│   ├── generate.go           # 生成工具函数
│   ├── jsonResponse.go       # JSON 响应工具
│   ├── logger.go             # 日志工具
│   ├── python-core.go        # Python 调用工具
│   └── validData.go          # 数据验证工具
├── config_example.yaml       # 配置文件示例
├── config.yaml               # 实际配置文件（需自行创建）
├── main.go                   # 应用入口点
├── rdkit_tools.py            # RDKit Python 工具脚本
└── README.md                 # 本文档
```

## 安装和运行

### 前提条件

1. **Go 1.24+**: [下载地址](https://go.dev/dl/)
2. **MySQL 8.0+**: [下载地址](https://www.mysql.com/)
3. **Python 3.8+**: [下载地址](https://www.python.org/)
4. **RDKit**: 化学信息学工具包

### 安装步骤

1. **克隆项目**:
   ```bash
   git clone https://github.com/shiming555ge/MNPLib.git
   cd MNPLib/backend
   ```

2. **安装 Go 依赖**:
   ```bash
   go mod download
   ```

3. **配置数据库**:
   - 创建 MySQL 数据库（如 `mnplib`）
   - 执行数据库初始化脚本（见下文数据库结构部分）
   - 复制 `config_example.yaml` 为 `config.yaml`
   - 修改 `config.yaml` 中的数据库配置

4. **安装 RDKit**:
   ```bash
   # 使用 conda 安装（推荐）
   conda create -n rdkit-env -c conda-forge rdkit
   conda activate rdkit-env
   
   # 或者使用 pip
   pip install rdkit-pypi
   ```

5. **运行后端服务**:
   ```bash
   go run main.go
   ```
   服务将在 `http://localhost:9090` 启动

### 构建可执行文件

```bash
# 构建
go build -o mnplib-backend main.go

# 运行
./mnplib-backend
```

## 配置说明

### 配置文件 (config.yaml)

```yaml
database:
  name: mnplib                # 数据库名称
  host: "127.0.0.1"          # 数据库主机
  port: 3306                 # 数据库端口
  user: root                 # 数据库用户名
  pass: "your_password"      # 数据库密码

enccryptokey: "32位加密密钥"  # 32位加密密钥，用于数据加密

jwt:
  secret: "your_jwt_secret"  # JWT 密钥，用于生成和验证 token

rdkit:
  python_path: "python"      # Python 解释器路径

static: false                # 是否启用静态文件服务
adress_port: ":9090"         # 服务器端口
```

### 环境变量

可以通过环境变量覆盖配置：
```bash
export DATABASE_NAME=mnplib
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USER=root
export DATABASE_PASS=your_password
export JWT_SECRET=your_jwt_secret
```

## 数据库结构

### data 表（化合物数据表）
存储所有化合物信息，包括化学结构、质谱数据、生物活性等。

```sql
CREATE TABLE `data` (
    `ID` VARCHAR(12) NOT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `Source` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `ItemName` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `ItemType` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `Formula` VARCHAR(127) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `SMILES` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `Description` ENUM('KNOWN COMPOUND','NEW NATURAL PRODUCT','NEW ANALOGS') NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `CAS_number` VARCHAR(127) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `ItemTag` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `Structure` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `MS1` VARCHAR(127) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `MS2` VARCHAR(511) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `Bioactivity` VARCHAR(511) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `NMR_13C_data` TEXT NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    `Weight` FLOAT NULL DEFAULT NULL,
    `FP` VARCHAR(255) NULL DEFAULT NULL COLLATE 'utf8mb4_uca1400_ai_ci',
    PRIMARY KEY (`ID`) USING BTREE
)
COLLATE='utf8mb4_uca1400_ai_ci'
ENGINE=InnoDB;
```

#### 字段说明：
- **ID**: 化合物唯一标识符（12位字符串，主键）
- **Source**: 数据来源
- **ItemName**: 化合物名称
- **ItemType**: 化合物类型（如生物碱、肽类等）
- **Formula**: 分子式
- **SMILES**: 简化分子线性输入规范
- **Description**: 化合物描述（已知化合物、新天然产物、新类似物）
- **CAS_number**: CAS登记号
- **ItemTag**: 化合物标签
- **Structure**: 化学结构信息
- **MS1**: 一级质谱数据
- **MS2**: 二级质谱数据（保护数据）
- **Bioactivity**: 生物活性数据（保护数据）
- **NMR_13C_data**: 碳13核磁共振数据（保护数据）
- **Weight**: 分子量
- **FP**: 分子指纹（用于相似度搜索）

### passkeys 表（用户认证表）
存储用户认证信息和权限。

```sql
CREATE TABLE `passkeys` (
    `Passkey` VARCHAR(36) NOT NULL DEFAULT uuid() COLLATE 'utf8mb3_uca1400_ai_ci',
    `Extends` TINYTEXT NOT NULL COLLATE 'utf8mb3_uca1400_ai_ci',
    `Description` VARCHAR(511) NOT NULL DEFAULT '' COLLATE 'utf8mb3_uca1400_ai_ci',
    `Operator` TINYTEXT NOT NULL DEFAULT '' COLLATE 'utf8mb3_uca1400_ai_ci',
    `Is_Active` TINYINT(1) NOT NULL DEFAULT '1',
    `Created_At` DATETIME NOT NULL DEFAULT current_timestamp()
)
COLLATE='utf8mb3_uca1400_ai_ci'
ENGINE=InnoDB;
```

#### 字段说明：
- **Passkey**: 用户唯一标识符（UUID，默认使用uuid()函数生成）
- **Extends**: 创建者信息，用于权限管理（空值表示超级管理员）
- **Description**: passkey描述
- **Operator**: 操作者名称
- **Is_Active**: 是否激活（1-激活，0-禁用）
- **Created_At**: 创建时间

### 数据关系
- `data` 表存储所有化合物数据，是系统的核心数据表
- `passkeys` 表用于用户认证和权限管理
- 保护数据字段（MS2、Bioactivity、NMR_13C_data）需要用户认证后才能访问

## API 文档

### 数据相关 API

#### 获取数据记录
- **URL**: `GET /api/data`
- **描述**: 返回指定数量的记录，支持从指定偏移量开始
- **参数**:
  - `limit` (可选): 返回的记录数量，默认为10，最大100
  - `offset` (可选): 从第几条记录开始，默认为0
- **响应**: 
  ```json
  {
    "data": [...],
    "total": 1000,
    "limit": 10,
    "offset": 0,
    "has_more": true,
    "next_offset": 10
  }
  ```

#### 根据ID获取数据
- **URL**: `GET /api/data/{id}`
- **描述**: 根据数据ID返回单条记录
- **参数**:
  - `id` (路径参数): 数据ID
- **响应**: 单条数据记录

#### 获取数据统计信息
- **URL**: `GET /api/data/statistics`
- **描述**: 返回化合物数量、物种数量、活性数据量、质谱数据量、核磁数据量等统计信息
- **响应**: 
  ```json
  {
    "total_compounds": 1000,
    "total_species": 500,
    "bioactivity_data": 300,
    "ms_data": 200,
    "nmr_data": 150
  }
  ```

#### 筛选化合物
- **URL**: `GET /api/data/filter`
- **描述**: 根据ItemType、分子量范围、Description和Source进行筛选，支持数组参数
- **参数**:
  - `limit` (可选): 返回的记录数量，默认为10，最大100
  - `offset` (可选): 从第几条记录开始，默认为0
  - `item_type` (可选): ItemType分类数组，可传入多个值。支持以下值（不区分大小写）:
    - `ALKALOID` - 生物碱类
    - `PEPTIDE` - 肽类
    - `POLYKETIDE` - 聚酮类
    - `TERPENOIDS` - 萜类
    - `CARBAZOLE` - 咔唑类
    - `INDOLE` - 吲哚类
    - `OTHERS` - 其他类别（除上述6类之外的所有化合物）
  - `min_weight` (可选): 最小分子量
  - `max_weight` (可选): 最大分子量
  - `description` (可选): Description描述数组，可传入多个值
  - `source` (可选): Source来源数组，可传入多个值
- **使用示例**:
  - 单个ItemType: `/api/data/filter?item_type=ALKALOID`
  - 多个ItemType: `/api/data/filter?item_type=ALKALOID&item_type=PEPTIDE&item_type=POLYKETIDE`
  - 包含OTHERS: `/api/data/filter?item_type=OTHERS` (返回除6个主要类别外的所有化合物)
  - 组合筛选: `/api/data/filter?item_type=ALKALOID&item_type=PEPTIDE&description=描述1&source=来源1&min_weight=100&max_weight=500`
  - 多个Source: `/api/data/filter?source=来源1&source=来源2`
- **注意**: 
  - `item_type`参数不区分大小写，前端可传入大写或小写
  - 当包含`OTHERS`时，返回除6个主要类别（ALKALOID, PEPTIDE, POLYKETIDE, TERPENOIDS, CARBAZOLE, INDOLE）之外的所有化合物
  - `description`和`source`参数支持模糊匹配（LIKE查询）
- **响应**: 
  ```json
  {
    "data": [...],
    "total": 50,
    "limit": 10,
    "offset": 0,
    "has_more": true,
    "next_offset": 10
  }
  ```

#### 获取所有ItemType分类
- **URL**: `GET /api/data/item-types`
- **描述**: 返回所有可用的ItemType分类
- **响应**: ItemType分类列表
  ```json
  ["分类1", "分类2", "分类3"]
  ```

#### 获取所有Description分类
- **URL**: `GET /api/data/descriptions`
- **描述**: 返回所有可用的Description分类
- **响应**: Description分类列表
  ```json
  ["描述1", "描述2", "描述3"]
  ```

#### 获取所有Source分类
- **URL**: `GET /api/data/sources`
- **描述**: 返回所有可用的Source分类
- **响应**: Source分类列表
  ```json
  ["来源1", "来源2", "来源3"]
  ```

#### 根据ID获取完整数据（保护数据）
- **URL**: `GET /api/data/{id}/full`
- **描述**: 根据数据ID返回MS2、Bioactivity和NMR_13C_data等保护数据，需要JWT认证
- **认证**: 需要在请求头中添加 `Authorization: Bearer <token>`
- **参数**:
  - `id` (路径参数): 数据ID
- **响应**: 保护数据字段
  ```json
  {
    "ms2": "MS2数据...",
    "bioactivity": "活性数据...",
    "nmr_13c_data": "核磁数据..."
  }
  ```

### 认证相关 API

#### 用户登录
- **URL**: `POST /api/auth/login`
- **描述**: 使用passkey进行登录，返回JWT token
- **请求体**:
  ```json
  {
    "passkey": "your_passkey"
  }
  ```
- **响应**: 
  ```json
  {
    "token": "jwt_token_here"
  }
  ```

#### 验证登录状态
- **URL**: `GET /api/auth/verify`
- **描述**: 验证JWT token是否有效（中间件已验证），返回成功状态
- **认证**: 需要在请求头中添加 `Authorization: Bearer <token>`
- **响应**: 成功状态
  ```json
  {
    "code": 200,
    "message": "success"
  }
  ```

#### 验证是否可以修改passkey
- **URL**: `GET /api/auth/verify-passkey-modifiable`
- **描述**: 验证当前用户是否有权限修改passkey（中间件已验证Extends字段为空），返回成功状态
- **认证**: 需要在请求头中添加 `Authorization: Bearer <token>`，并且用户的 Extends 字段必须为空（即
