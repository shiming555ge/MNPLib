package services

import (
	"backend/database"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

// InitializeCompoundData 初始化化合物数据，计算缺失的FP、Structure和Weight
func InitializeCompoundData() error {
	// 检查Python进程是否已初始化
	if p == nil {
		errMsg := "RDkit进程未初始化，无法计算化合物数据"
		utils.Log(errMsg)
		return fmt.Errorf(errMsg)
	}

	// 获取所有化合物数据
	var compounds []models.Data
	result := database.GetDB().Table("data").Find(&compounds)
	if result.Error != nil {
		utils.LogError(result.Error)
		return fmt.Errorf("获取化合物数据失败: %v", result.Error)
	}

	utils.Log(fmt.Sprintf("开始初始化 %d 个化合物的数据...", len(compounds)))

	// 统计需要计算的数据
	compoundsToUpdate := []struct {
		ID            string
		SMILES        string
		NeedFP        bool
		NeedStructure bool
		NeedWeight    bool
	}{}

	for _, compound := range compounds {
		needFP := compound.FP == nil || *compound.FP == ""
		needStructure := compound.Structure == nil || *compound.Structure == ""
		needWeight := compound.Weight == nil

		// 检查SMILES是否存在
		var smiles string
		if compound.SMILES != nil {
			smiles = *compound.SMILES
		} else {
			// 如果没有SMILES，跳过这个化合物
			continue
		}

		if needFP || needStructure || needWeight {
			compoundsToUpdate = append(compoundsToUpdate, struct {
				ID            string
				SMILES        string
				NeedFP        bool
				NeedStructure bool
				NeedWeight    bool
			}{
				ID:            compound.ID,
				SMILES:        smiles,
				NeedFP:        needFP,
				NeedStructure: needStructure,
				NeedWeight:    needWeight,
			})
		}
	}

	if len(compoundsToUpdate) == 0 {
		utils.Log("所有化合物数据已初始化完成")
		return nil
	}

	utils.Log(fmt.Sprintf("发现 %d 个化合物需要初始化数据", len(compoundsToUpdate)))

	// 逐个计算缺失的数据
	for _, compound := range compoundsToUpdate {
		utils.Log(fmt.Sprintf("正在初始化化合物 %s 的数据...", compound.ID))

		updates := make(map[string]interface{})

		// 计算FP
		if compound.NeedFP {
			fp, err := calculateFingerprint(compound.SMILES)
			if err != nil {
				utils.LogError(err)
				utils.Log(fmt.Sprintf("计算FP失败: ID=%s", compound.ID))
			} else {
				updates["FP"] = fp
				utils.Log(fmt.Sprintf("成功计算FP: ID=%s", compound.ID))
			}
		}

		// 计算Structure
		if compound.NeedStructure {
			structure, err := calculateStructure(compound.SMILES)
			if err != nil {
				utils.LogError(err)
				utils.Log(fmt.Sprintf("计算Structure失败: ID=%s", compound.ID))
			} else {
				updates["Structure"] = structure
				utils.Log(fmt.Sprintf("成功计算Structure: ID=%s", compound.ID))
			}
		}

		// 计算Weight
		if compound.NeedWeight {
			weight, err := calculateMolecularWeight(compound.SMILES)
			if err != nil {
				utils.LogError(err)
				utils.Log(fmt.Sprintf("计算Weight失败: ID=%s", compound.ID))
			} else {
				updates["Weight"] = weight
				utils.Log(fmt.Sprintf("成功计算Weight: ID=%s", compound.ID))
			}
		}

		// 更新数据库
		if len(updates) > 0 {
			updateResult := database.GetDB().Table("data").
				Where("ID = ?", compound.ID).
				Updates(updates)
			if updateResult.Error != nil {
				utils.LogError(updateResult.Error)
				utils.Log(fmt.Sprintf("更新化合物数据失败: ID=%s", compound.ID))
			} else {
				utils.Log(fmt.Sprintf("成功更新化合物数据: ID=%s", compound.ID))
			}
		}
	}

	utils.Log("化合物数据初始化完成")
	return nil
}

// calculateFingerprint 计算指纹
func calculateFingerprint(smiles string) (string, error) {
	requestData := map[string]interface{}{
		"action": "smiles_to_fingerprint",
		"smiles": smiles,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		return "", fmt.Errorf("计算指纹失败: %v", err)
	}

	return res, nil
}

// calculateStructure 计算结构
func calculateStructure(smiles string) (string, error) {
	requestData := map[string]interface{}{
		"action": "smiles_to_pdb",
		"smiles": smiles,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return "", fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		return "", fmt.Errorf("计算结构失败: %v", err)
	}

	return res, nil
}

// calculateMolecularWeight 计算分子量
func calculateMolecularWeight(smiles string) (float32, error) {
	requestData := map[string]interface{}{
		"action": "calculate_molecular_weight",
		"smiles": smiles,
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return 0, fmt.Errorf("请求数据序列化失败: %v", err)
	}

	res, err := p.SendAndWait(string(requestJSON))
	if err != nil {
		return 0, fmt.Errorf("计算分子量失败: %v", err)
	}

	weight, err := strconv.ParseFloat(res, 32)
	if err != nil {
		return 0, fmt.Errorf("解析分子量失败: %v", err)
	}

	return float32(weight), nil
}
