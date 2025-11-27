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
	outputFile := c.DefaultQuery("outputFile", "output.pdb")
	if smiles == "" {
		utils.JsonErrorResponse(c, 200400, "SMILES字符串不能为空")
		return
	}

	result, err := services.SmilesToPDB(smiles, outputFile)
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
