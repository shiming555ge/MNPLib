package services

import (
	"backend/config"
	"backend/utils"
	"os"
)

var p *utils.PythonProcess

func InitRdkit() {
	pwd, err := os.Getwd()
	if err != nil {
		utils.LogError(err)
		return
	}
	p = utils.NewPythonProcess(config.Config.GetString("rdkit.pytho_path"), pwd+"/rdkit_tools.py")
	utils.Log("RDkit初始化成功")
}

func SimilaritySearch(smiles string) {
	_, err := p.SendAndWait("SimSearch")
	if err != nil {
		utils.LogError(err)
	}
	//p.SendAndWait(database.GetDB().Where("")))
}
