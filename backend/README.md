# Back-end

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

### RDKit 化学计算 API

#### 获取RDKit服务状态
- **URL**: `GET /api/rdkit/status`
- **描述**: 检查RDKit服务是否正常运行
- **响应**: 服务状态信息
  ```json
  {
    "initialized": true,
    "available": true,
    "status": "running"
  }
  ```

#### 相似度搜索
- **URL**: `GET /api/rdkit/similarity`
- **描述**: 根据查询指纹进行相似度搜索
- **参数**:
  - `qfp` (必需): 查询指纹
  - `threshold` (可选): 相似度阈值，默认为0.5
- **响应**: 相似度搜索结果

#### SMILES转指纹
- **URL**: `GET /api/rdkit/smiles-to-fingerprint`
- **描述**: 将SMILES字符串转换为分子指纹
- **参数**:
  - `smiles` (必需): SMILES字符串
- **响应**: 分子指纹数据

#### SMILES转PDB
- **URL**: `GET /api/rdkit/smiles-to-pdb`
- **描述**: 将SMILES字符串转换为PDB文件
- **参数**:
  - `smiles` (必需): SMILES字符串
- **响应**: PDB文件数据

#### 子结构匹配
- **URL**: `GET /api/rdkit/is-substructure`
- **描述**: 检查SMILES是否包含指定的SMARTS子结构
- **参数**:
  - `smarts_pattern` (必需): SMARTS模式
  - `smiles` (必需): SMILES字符串
- **响应**: 
  ```json
  {
    "is_substructure": true/false
  }
  ```

#### 子结构搜索
- **URL**: `GET /api/rdkit/substructure-search`
- **描述**: 根据SMARTS模式在数据库中查找所有匹配的化合物
- **参数**:
  - `smarts_pattern` (必需): SMARTS模式
- **响应**: 匹配的化合物列表

#### 精确匹配搜索
- **URL**: `GET /api/rdkit/exact-match`
- **描述**: 查找SMILES相同的结构并返回其ID
- **参数**:
  - `smiles` (必需): SMILES字符串
- **响应**: 匹配的化合物ID

## 项目结构

### config
存放用来读取配置文件的config.go（用viper）

### controllers
所有在这里的文件，负责选择调用哪个service里的函数，然后返回json数据。

### utils
此文件夹用以存放工具类，logger.go放在里头了

|文件|功能|
|---|---|
|config.go|配置文件|
|jwt.go|jwt相关|
|jsonResponse.go|格式化json输出|
|crypto.go|加密/解密，密钥放在config.yaml里|

- 后端如果有错误日志产生用zerolog搞定吧，初始化那些放在utils/logger.go里了

### middleware
此文件夹用以存放中间件

### models
此文件夹用来存放数据模型

### services
所有需要操作数据库的go文件放在这里

### routes
此文件夹用以存放路径
请放到router.go对应位置
推荐每个种api放一个文件，如user,auth，具体见apifox

### static
静态文件，如icon，css，js等

### upload
上传的文件 暂时不知道是否使用这种形式。

## To Dos

各种数据校验

注册保存手机号和邮箱

奇怪的外键问题
