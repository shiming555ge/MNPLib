package services

import (
	"backend/config"
	"backend/database"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var p *utils.PythonProcess

// InitRdkit 初始化RDKit Python进程
func InitRdkit() error {
	pwd, err := os.Getwd()
	if err != nil {
		utils.LogError(err)
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	path := config.Config.GetString("rdkit.python_path")
	if path == "" {
		errMsg := "RDkit初始化失败: python_path配置为空"
		utils.Log(errMsg)
		return fmt.Errorf(errMsg)
	}

	// 使用filepath处理路径，确保跨平台兼容性
	pythonScriptPath := filepath.Join(pwd, "rdkit_tools.py")

	var initErr error
	p, initErr = utils.NewPythonProcess(path, pythonScriptPath)
	if initErr != nil {
		utils.LogError(initErr)
		utils.Log(fmt.Sprintf("RDkit进程启动失败: Python路径=%v, 脚本路径=%v", path, pythonScriptPath))
		return fmt.Errorf("RDkit进程启动失败: %v", initErr)
	}

	res, err := p.SendAndWait("init")
	if err != nil {
		utils.LogError(err)
		utils.Log("RDkit初始化通信失败")
		return fmt.Errorf("RDkit初始化通信失败: %v", err)
	}
	if res != "initialized" {
		errMsg := fmt.Sprintf("RDkit初始化失败: 响应不正确, 期望='initialized', 实际='%s'", res)
		utils.Log(errMsg)
		return fmt.Errorf(errMsg)
	}
	utils.Log("RDkit初始化成功")
	return nil
}

// GetRdkitStatus 获取RDKit服务状态
func GetRdkitStatus() map[string]interface{} {
	status := map[string]interface{}{
		"initialized": p != nil,
		"available":   p != nil,
	}

	if p != nil {
		status["status"] = "running"
	} else {
		status["status"] = "not_initialized"
	}

	return status
}

// SimilaritySearch 相似度搜索
func SimilaritySearch(fp string, threshold string) (string, error) {
	// 检查Python进程是否已初始化
	if p == nil {
		errMsg := "RDkit进程未初始化"
		utils.Log(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 获取数据库中的所有指纹数据
	var compounds []struct {
		ID     string `gorm:"column:ID" json:"id"`
		SMILES string `gorm:"column:SMILES;type:TEXT" json:"smiles,omitempty"`
		FP     string `gorm:"column:FP" json:"fp"`
	}

	// 使用正确的表名和字段名查询
	result := database.GetDB().Table("data").Select("ID, FP").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return "", fmt.Errorf("数据库查询失败: %v", result.Error)
	}
	// 检查fp是否存在
	for i, compound := range compounds {
		// 初始化FP
		if compound.FP == "" {
			var err error
			compound.FP, err = SmilesToFingerprint(compound.SMILES)
			if err != nil {
				return "", fmt.Errorf("初始化化合物(%v)FP失败： %v", compound.ID, result.Error)
			}
			database.GetDB().Table("data").Where("ID = ?", compound.ID).Update("FP", compound.FP)
			compounds[i] = compound
		}
	}

	// 准备发送给Python的数据
	library := make([]map[string]interface{}, len(compounds))
	for i, compound := range compounds {
		library[i] = map[string]interface{}{
			"id": compound.ID,
			"fp": compound.FP,
		}
	}

	// 发送相似度搜索请求
	requestData := map[string]interface{}{
		"action":    "similarity_search",
		"qfp":       fp,
		"threshold": threshold,
		"data":      library,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("相似度搜索失败: %v", err)
	}

	return res, nil
}

// SmilesToFingerprint SMILES转指纹
func SmilesToFingerprint(smiles string) (string, error) {
	if p == nil {
		errMsg := "RDkit进程未初始化"
		utils.Log(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	requestData := map[string]interface{}{
		"action": "smiles_to_fingerprint",
		"smiles": smiles,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("SMILES转指纹失败: %v", err)
	}

	return res, nil
}

// SmilesToPDB SMILES转PDB
func SmilesToPDB(smiles string) (string, error) {
	if p == nil {
		errMsg := "RDkit进程未初始化"
		utils.Log(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	requestData := map[string]interface{}{
		"action": "smiles_to_pdb",
		"smiles": smiles,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("SMILES转PDB失败: %v", err)
	}

	return res, nil
}

// IsSubstructure 子结构匹配
func IsSubstructure(smartsPattern string, smiles string) (bool, error) {
	if p == nil {
		errMsg := "RDkit进程未初始化"
		utils.Log(errMsg)
		return false, fmt.Errorf(errMsg)
	}

	requestData := map[string]interface{}{
		"action":         "is_substructure",
		"smarts_pattern": smartsPattern,
		"smiles":         smiles,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		utils.LogError(err)
		return false, fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		utils.LogError(err)
		return false, fmt.Errorf("子结构匹配失败: %v", err)
	}

	// 解析结果为布尔值
	result := strings.ToLower(strings.TrimSpace(res))
	return result == "true" || result == "1", nil
}

// SubstructureSearch 子结构搜索 - 根据SMARTS模式在数据库中查找所有匹配的化合物
func SubstructureSearch(smartsPattern string) (string, error) {
	// 检查Python进程是否已初始化
	if p == nil {
		errMsg := "RDkit进程未初始化"
		utils.Log(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 获取数据库中的所有SMILES数据
	var compounds []struct {
		ID     string `gorm:"column:ID" json:"id"`
		SMILES string `gorm:"column:SMILES" json:"smiles"`
	}

	// 使用正确的表名和字段名查询
	result := database.GetDB().Table("data").Select("ID, SMILES").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return "", fmt.Errorf("数据库查询失败: %v", result.Error)
	}

	// 准备发送给Python的数据 - 现在传递ID和SMILES
	library := make([]map[string]interface{}, len(compounds))
	for i, compound := range compounds {
		library[i] = map[string]interface{}{
			"id":     compound.ID,
			"smiles": compound.SMILES,
		}
	}

	// 发送子结构搜索请求
	requestData := map[string]interface{}{
		"action":         "substructure_search",
		"smarts_pattern": smartsPattern,
		"library":        library,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("子结构搜索失败: %v", err)
	}

	return res, nil
}

// ExactMatchSearch 精确匹配搜索 - 查找SMILES相同的结构并返回其ID
func ExactMatchSearch(smiles string) (string, error) {
	// 检查Python进程是否已初始化
	if p == nil {
		errMsg := "RDkit进程未初始化"
		utils.Log(errMsg)
		return "", fmt.Errorf(errMsg)
	}

	// 获取数据库中的所有SMILES数据
	var compounds []struct {
		ID     string `gorm:"column:ID" json:"id"`
		SMILES string `gorm:"column:SMILES" json:"smiles"`
	}

	// 使用正确的表名和字段名查询
	result := database.GetDB().Table("data").Select("ID, SMILES").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return "", fmt.Errorf("数据库查询失败: %v", result.Error)
	}

	// 准备发送给Python的数据 - 传递ID和SMILES
	library := make([]map[string]interface{}, len(compounds))
	for i, compound := range compounds {
		library[i] = map[string]interface{}{
			"id":     compound.ID,
			"smiles": compound.SMILES,
		}
	}

	// 发送精确匹配搜索请求
	requestData := map[string]interface{}{
		"action":  "exact_match_search",
		"smiles":  smiles,
		"library": library,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		utils.LogError(err)
		return "", fmt.Errorf("精确匹配搜索失败: %v", err)
	}

	return res, nil
}

// FilterCompounds 筛选化合物 - 根据ItemType、分子量范围、Description和Source进行筛选，支持数组参数
func FilterCompounds(itemTypes []string, minWeight, maxWeight float64, descriptions []string, sources []string, limit, offset int) ([]models.PublicData, int64, error) {
	// 构建查询条件
	query := database.GetDB().Table("data")

	// ItemType筛选 - 支持数组
	if len(itemTypes) > 0 {
		// 定义主要的6个类别（小写形式）
		mainCategories := []string{"alkaloid", "peptide", "polyketide", "terpenoids", "carbazole", "indole"}

		// 将前端传过来的itemTypes转换为小写
		lowerItemTypes := make([]string, len(itemTypes))
		hasOthers := false
		specificTypes := []string{}

		for i, itemType := range itemTypes {
			lowerType := strings.ToLower(itemType)
			lowerItemTypes[i] = lowerType

			// 检查是否包含"others"
			if lowerType == "others" {
				hasOthers = true
			} else {
				specificTypes = append(specificTypes, lowerType)
			}
		}

		if hasOthers {
			// 如果包含"others"，查询除了6个主要类别之外的所有化合物
			// 但也要包含用户选择的其他具体类别（如果有的话）
			if len(specificTypes) > 0 {
				// 用户既选择了"others"又选择了其他具体类别
				// 查询：不属于6个主要类别，但属于用户选择的具体类别
				query = query.Where("LOWER(ItemType) NOT IN (?) AND LOWER(ItemType) IN (?)", mainCategories, specificTypes)
			} else {
				// 用户只选择了"others"
				// 查询：不属于6个主要类别的所有化合物
				query = query.Where("LOWER(ItemType) NOT IN (?)", mainCategories)
			}
		} else {
			// 如果不包含"others"，正常查询用户选择的类别
			query = query.Where("LOWER(ItemType) IN (?)", lowerItemTypes)
		}
	}

	// Description筛选 - 支持数组
	if len(descriptions) > 0 {
		// 为每个description创建LIKE条件
		likeConditions := make([]string, len(descriptions))
		likeArgs := make([]interface{}, len(descriptions))
		for i, desc := range descriptions {
			likeConditions[i] = "Description LIKE ?"
			likeArgs[i] = "%" + desc + "%"
		}
		// 使用OR连接多个LIKE条件
		query = query.Where(strings.Join(likeConditions, " OR "), likeArgs...)
	}

	// Source筛选 - 支持数组
	if len(sources) > 0 {
		// 为每个source创建LIKE条件
		likeConditions := make([]string, len(sources))
		likeArgs := make([]interface{}, len(sources))
		for i, source := range sources {
			likeConditions[i] = "Source LIKE ?"
			likeArgs[i] = "%" + source + "%"
		}
		// 使用OR连接多个LIKE条件
		query = query.Where(strings.Join(likeConditions, " OR "), likeArgs...)
	}

	// 分子量筛选 - 使用Weight字段
	if minWeight > 0 || maxWeight > 0 {
		if minWeight > 0 {
			query = query.Where("Weight >= ?", minWeight)
		}
		if maxWeight > 0 {
			query = query.Where("Weight <= ?", maxWeight)
		}
	}

	// 获取总记录数
	var totalCount int64
	countQuery := query
	result := countQuery.Count(&totalCount)
	if result.Error != nil {
		utils.LogError(result.Error)
		return nil, 0, fmt.Errorf("获取总记录数失败: %v", result.Error)
	}

	// 应用分页并获取数据
	var compounds []models.PublicData
	result = query.Offset(offset).Limit(limit).Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return nil, 0, fmt.Errorf("数据库查询失败: %v", result.Error)
	}

	return compounds, totalCount, nil
}

// GetItemTypes 获取所有ItemType分类
func GetItemTypes() ([]string, error) {
	var itemTypes []string
	result := database.GetDB().Table("data").
		Where("ItemType IS NOT NULL AND ItemType != ''").
		Distinct("ItemType").
		Pluck("ItemType", &itemTypes)

	if result.Error != nil {
		utils.LogError(result.Error)
		return nil, fmt.Errorf("获取ItemType分类失败: %v", result.Error)
	}

	return itemTypes, nil
}

// GetDescriptions 获取所有Description分类
func GetDescriptions() ([]string, error) {
	var descriptions []string
	result := database.GetDB().Table("data").
		Where("Description IS NOT NULL AND Description != ''").
		Distinct("Description").
		Pluck("Description", &descriptions)

	if result.Error != nil {
		utils.LogError(result.Error)
		return nil, fmt.Errorf("获取Description分类失败: %v", result.Error)
	}

	return descriptions, nil
}

// GetSources 获取所有Source分类
func GetSources() ([]string, error) {
	var sources []string
	result := database.GetDB().Table("data").
		Where("Source IS NOT NULL AND Source != ''").
		Distinct("Source").
		Pluck("Source", &sources)

	if result.Error != nil {
		utils.LogError(result.Error)
		return nil, fmt.Errorf("获取Source分类失败: %v", result.Error)
	}

	return sources, nil
}
