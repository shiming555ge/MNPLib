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

### RDKit 化学计算 API

#### 获取RDKit服务状态
- **URL**: `GET /api/rdkit/status`
- **描述**: 检查RDKit服务是否正常运行
- **响应**: 服务状态信息

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
  - `outputFile` (可选): 输出文件名，默认为"output.pdb"
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
