package controllers

// GetSimpleDataRecords 获取简化的数据记录
// @Summary 获取简化的数据记录
// @Description 返回指定数量的简化记录，只包含id, item_name, smiles和cas_number
// @Tags simple-data
// @Accept json
// @Produce json
// @Param limit query int false "返回的记录数量，默认为10"
// @Param offset query int false "从第几条记录开始，默认为0"
// @Success 200 {object} utils.JSONResponse{data=[]models.SimpleData}
