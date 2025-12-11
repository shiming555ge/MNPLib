package controllers

import (
	"backend/services"
	"backend/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SimilaritySearch(c *gin.Context) {
	// 绑定请求参数
	qfp := c.Query("qfp")
	threshold := c.DefaultQuery("threshold", "0.5")
	// 验证参数
	if qfp == "" {
		utils.JsonErrorResponse(c, 200400, "查询指纹qfp不能为空")
		return
	}

	result, err := services.SimilaritySearch(qfp, threshold)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("相似度搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// GetRdkitStatus 获取RDKit服务状态
func GetRdkitStatus(c *gin.Context) {
	status := services.GetRdkitStatus()
	utils.JsonSuccessResponse(c, status)
}

// SmilesToFingerprint SMILES转指纹
func SmilesToFingerprint(c *gin.Context) {
	smiles := c.Query("smiles")
	if smiles == "" {
		utils.JsonErrorResponse(c, 200400, "SMILES字符串不能为空")
		return
	}

	result, err := services.SmilesToFingerprint(smiles)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("SMILES转指纹失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// SmilesToPDB SMILES转PDB文件
func SmilesToPDB(c *gin.Context) {
	smiles := c.Query("smiles")
	if smiles == "" {
		utils.JsonErrorResponse(c, 200400, "SMILES字符串不能为空")
		return
	}

	result, err := services.SmilesToPDB(smiles)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("SMILES转PDB失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// IsSubstructure 子结构匹配
func IsSubstructure(c *gin.Context) {
	smartsPattern := c.Query("smarts_pattern")
	smiles := c.Query("smiles")
	if smartsPattern == "" || smiles == "" {
		utils.JsonErrorResponse(c, 200400, "SMARTS模式和SMILES字符串都不能为空")
		return
	}

	result, err := services.IsSubstructure(smartsPattern, smiles)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("子结构匹配失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, map[string]interface{}{
		"is_substructure": result,
	})
}

// SubstructureSearch 子结构搜索 - 根据SMARTS模式在数据库中查找所有匹配的化合物
func SubstructureSearch(c *gin.Context) {
	smartsPattern := c.Query("smarts_pattern")
	if smartsPattern == "" {
		utils.JsonErrorResponse(c, 200400, "SMARTS模式不能为空")
		return
	}

	result, err := services.SubstructureSearch(smartsPattern)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("子结构搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// ExactMatchSearch 精确匹配搜索 - 查找SMILES相同的结构并返回其ID
func ExactMatchSearch(c *gin.Context) {
	smiles := c.Query("smiles")
	if smiles == "" {
		utils.JsonErrorResponse(c, 200400, "SMILES字符串不能为空")
		return
	}

	result, err := services.ExactMatchSearch(smiles)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("精确匹配搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// NMRSearch 核磁谱搜索
func NMRSearch(c *gin.Context) {
	// 绑定请求参数
	queryNMR := c.Query("query_nmr")
	threshold := c.DefaultQuery("threshold", "0.5")
	tolerance := c.DefaultQuery("tolerance", "0.5")

	// 验证参数
	if queryNMR == "" {
		utils.JsonErrorResponse(c, 200400, "查询核磁谱数据不能为空")
		return
	}

	result, err := services.NMRSearch(queryNMR, threshold, tolerance)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("核磁谱搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// MS2Search MS2相似度搜索（接收原始MS2文本）
func MS2Search(c *gin.Context) {
	// 定义请求结构体
	type MS2SearchRequest struct {
		QueryMS2           string  `json:"query_ms2" binding:"required"`
		Threshold          float64 `json:"threshold"`
		Tolerance          float64 `json:"tolerance"`
		PrefilterThreshold float64 `json:"prefilter_threshold"`
	}

	// 绑定请求参数
	var req MS2SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 200400, fmt.Sprintf("请求参数错误: %v", err))
		return
	}

	// 设置默认值
	if req.Threshold == 0 {
		req.Threshold = 0.5
	}
	if req.Tolerance == 0 {
		req.Tolerance = 0.5
	}
	if req.PrefilterThreshold == 0 {
		req.PrefilterThreshold = 0.3
	}

	// 验证参数
	if req.QueryMS2 == "" {
		utils.JsonErrorResponse(c, 200400, "查询MS2数据不能为空")
		return
	}

	// 转换参数为字符串
	thresholdStr := fmt.Sprintf("%.2f", req.Threshold)
	toleranceStr := fmt.Sprintf("%.2f", req.Tolerance)
	prefilterThresholdStr := fmt.Sprintf("%.2f", req.PrefilterThreshold)

	result, err := services.MS2Search(req.QueryMS2, thresholdStr, toleranceStr, prefilterThresholdStr)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("MS2相似度搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// MS2SearchByFingerprint MS2相似度搜索（接收指纹数据）
func MS2SearchByFingerprint(c *gin.Context) {
	// 定义请求结构体
	type MS2FingerprintSearchRequest struct {
		FingerprintJson    string  `json:"fingerprint_json" binding:"required"`
		Threshold          float64 `json:"threshold"`
		Tolerance          float64 `json:"tolerance"`
		PrefilterThreshold float64 `json:"prefilter_threshold"`
	}

	// 绑定请求参数
	var req MS2FingerprintSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 200400, fmt.Sprintf("请求参数错误: %v", err))
		return
	}

	// 设置默认值
	if req.Threshold == 0 {
		req.Threshold = 0.5
	}
	if req.Tolerance == 0 {
		req.Tolerance = 0.5
	}
	if req.PrefilterThreshold == 0 {
		req.PrefilterThreshold = 0.3
	}

	// 验证参数
	if req.FingerprintJson == "" {
		utils.JsonErrorResponse(c, 200400, "指纹JSON数据不能为空")
		return
	}

	// 转换参数为字符串
	thresholdStr := fmt.Sprintf("%.2f", req.Threshold)
	toleranceStr := fmt.Sprintf("%.2f", req.Tolerance)
	prefilterThresholdStr := fmt.Sprintf("%.2f", req.PrefilterThreshold)

	result, err := services.MS2SearchByFingerprint(req.FingerprintJson, thresholdStr, toleranceStr, prefilterThresholdStr)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("MS2指纹搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}

// MS2SearchByEnergyLevel MS2相似度搜索（按能量级别分别比对）
func MS2SearchByEnergyLevel(c *gin.Context) {
	// 定义请求结构体
	type MS2EnergyLevelSearchRequest struct {
		QueryMS2           string  `json:"query_ms2" binding:"required"`
		Threshold          float64 `json:"threshold"`
		Tolerance          float64 `json:"tolerance"`
		PrefilterThreshold float64 `json:"prefilter_threshold"`
		EnergyLevel        string  `json:"energy_level"`
	}

	// 绑定请求参数
	var req MS2EnergyLevelSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JsonErrorResponse(c, 200400, fmt.Sprintf("请求参数错误: %v", err))
		return
	}

	// 设置默认值
	if req.Threshold == 0 {
		req.Threshold = 0.5
	}
	if req.Tolerance == 0 {
		req.Tolerance = 0.5
	}
	if req.PrefilterThreshold == 0 {
		req.PrefilterThreshold = 0.3
	}
	if req.EnergyLevel == "" {
		req.EnergyLevel = "energy0" // 默认使用energy0
	}

	// 验证参数
	if req.QueryMS2 == "" {
		utils.JsonErrorResponse(c, 200400, "查询MS2数据不能为空")
		return
	}

	// 验证能量级别参数
	validEnergyLevels := []string{"energy0", "energy1", "energy2"}
	valid := false
	for _, level := range validEnergyLevels {
		if req.EnergyLevel == level {
			valid = true
			break
		}
	}
	if !valid {
		utils.JsonErrorResponse(c, 200400, "能量级别参数无效，必须是energy0、energy1或energy2")
		return
	}

	// 转换参数为字符串
	thresholdStr := fmt.Sprintf("%.2f", req.Threshold)
	toleranceStr := fmt.Sprintf("%.2f", req.Tolerance)
	prefilterThresholdStr := fmt.Sprintf("%.2f", req.PrefilterThreshold)

	result, err := services.MS2SearchByEnergyLevel(req.QueryMS2, thresholdStr, toleranceStr, prefilterThresholdStr, req.EnergyLevel)
	if err != nil {
		utils.JsonErrorResponse(c, 200500, fmt.Sprintf("MS2按能量级别搜索失败: %v", err))
		return
	}
	utils.JsonSuccessResponse(c, result)
}
