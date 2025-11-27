package services

import (
	"backend/config"
	"backend/database"
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
		ID string `gorm:"column:ID" json:"id"`
		Fp string `gorm:"column:Fp" json:"fp"`
	}

	// 使用正确的表名和字段名查询
	result := database.GetDB().Table("data").Select("ID, Fp").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return "", fmt.Errorf("数据库查询失败: %v", result.Error)
	}

	// 准备发送给Python的数据
	library := make([]map[string]interface{}, len(compounds))
	for i, compound := range compounds {
		library[i] = map[string]interface{}{
			"id": compound.ID,
			"fp": compound.Fp,
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
		Smiles string `gorm:"column:Smiles" json:"smiles"`
	}

	// 使用正确的表名和字段名查询
	result := database.GetDB().Table("data").Select("ID, Smiles").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return "", fmt.Errorf("数据库查询失败: %v", result.Error)
	}

	// 准备发送给Python的数据 - 现在传递ID和SMILES
	library := make([]map[string]interface{}, len(compounds))
	for i, compound := range compounds {
		library[i] = map[string]interface{}{
			"id":     compound.ID,
			"smiles": compound.Smiles,
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
		Smiles string `gorm:"column:Smiles" json:"smiles"`
	}

	// 使用正确的表名和字段名查询
	result := database.GetDB().Table("data").Select("ID, Smiles").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return "", fmt.Errorf("数据库查询失败: %v", result.Error)
	}

	// 准备发送给Python的数据 - 传递ID和SMILES
	library := make([]map[string]interface{}, len(compounds))
	for i, compound := range compounds {
		library[i] = map[string]interface{}{
			"id":     compound.ID,
			"smiles": compound.Smiles,
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
