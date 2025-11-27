from rdkit import Chem, DataStructs
from rdkit.Chem import AllChem

import base64, json

# smiles式转存pdb
def smiles_to_pdb(smiles):
    mol = Chem.MolFromSmiles(smiles)
    # mol = Chem.AddHs(mol)                       # 加氢
    AllChem.EmbedMolecule(mol)                 # 3D 构象
    AllChem.UFFOptimizeMolecule(mol)           # UFF 优化
    return Chem.MolToPDBBlock(mol)

# 创建指纹生成器
gen = Chem.rdFingerprintGenerator.GetMorganGenerator(radius=2, fpSize=2048)
# 分子指纹
def smiles_to_tanimoto(smiles):
    mol = Chem.MolFromSmiles(smiles)
    fp = gen.GetFingerprint(mol)
    return fp.ToBase64()

# 判断子结构
def is_substructure(smarts_pattern, smiles):
    patt = Chem.MolFromSmarts(smarts_pattern)
    mol = Chem.MolFromSmiles(smiles)
    return mol.HasSubstructMatch(patt)

# 批量相似度搜索
# library [{id:,fp:}]
def similarity_search(qfp, library, threshold=0.5):
    results = []
    for s in library:
        sim = DataStructs.TanimotoSimilarity(qfp, s['fp'])
        if sim >= threshold:
            results.append((s['id'], sim))
    return sorted(results, key=lambda x: -x[1])

# 批量子结构搜索
# library 现在是包含id和smiles的字典列表
def substructure_search(pattern_smarts, library):
    patt = Chem.MolFromSmarts(pattern_smarts)
    result = []
    for item in library:
        smiles = item['smiles']
        mol = Chem.MolFromSmiles(smiles)
        if mol.HasSubstructMatch(patt):
            result.append(item['id'])
    return result

# 精确匹配搜索 - 使用RDKit进行化学等价性比较
def exact_match_search(query_smiles, library):
    result = []
    query_mol = Chem.MolFromSmiles(query_smiles)
    if query_mol is None:
        return result  # 无效的查询SMILES
    
    for item in library:
        db_mol = Chem.MolFromSmiles(item['smiles'])
        if db_mol is not None and query_mol.HasSubstructMatch(db_mol) and db_mol.HasSubstructMatch(query_mol):
            result.append(item['id'])
    return result

if __name__=="__main__":
    # 模式
    while True:
        try:
            cmd = input()
            if not cmd:
                continue
                
            request = json.loads(cmd)
            msg_id = request.get("id")
            msg_content = request.get("msg")
            
            if msg_content == "init":
                response = {"id": msg_id, "reply": "initialized"}
                print(json.dumps(response))
                continue
                
            # 解析请求数据
            data = json.loads(msg_content)
            action = data.get("action")
            
            if action == "similarity_search":
                # 相似度搜索
                qfp_base64 = data.get("qfp")
                library = data.get("data")
                threshold = float(data.get('threshold', 0.5))
                
                if qfp_base64 and library:
                    # 转换查询指纹
                    fp = DataStructs.cDataStructs.ExplicitBitVect(2048)
                    fp.FromBase64(qfp_base64)
                    
                    # 转换库中的指纹
                    processed_library = []
                    for item in library:
                        if 'fp' in item:
                            lib_fp = DataStructs.cDataStructs.ExplicitBitVect(2048)
                            lib_fp.FromBase64(item['fp'])
                            processed_library.append({
                                'id': item.get('id', ''),
                                'fp': lib_fp
                            })
                    
                    # 执行相似度搜索
                    results = similarity_search(fp, processed_library, threshold)
                    response = {"id": msg_id, "reply": json.dumps(results)}
                else:
                    response = {"id": msg_id, "reply": "error: missing qfp or data parameter"}
                    
            elif action == "smiles_to_fingerprint":
                # SMILES转指纹
                smiles = data.get("smiles")
                if smiles:
                    fp_base64 = smiles_to_tanimoto(smiles)
                    response = {"id": msg_id, "reply": fp_base64}
                else:
                    response = {"id": msg_id, "reply": "error: missing smiles parameter"}
                    
            elif action == "smiles_to_pdb":
                # SMILES转PDB
                smiles = data.get("smiles")
                if smiles:
                    result = smiles_to_pdb(smiles)
                    response = {"id": msg_id, "reply": result}
                else:
                    response = {"id": msg_id, "reply": "error: missing smiles parameter"}
                    
            elif action == "is_substructure":
                # 子结构匹配
                smarts_pattern = data.get("smarts_pattern")
                smiles = data.get("smiles")
                if smarts_pattern and smiles:
                    result = is_substructure(smarts_pattern, smiles)
                    response = {"id": msg_id, "reply": str(result).lower()}
                else:
                    response = {"id": msg_id, "reply": "error: missing smarts_pattern or smiles parameter"}
                    
            elif action == "substructure_search":
                # 子结构搜索
                smarts_pattern = data.get("smarts_pattern")
                library = data.get("library")
                if smarts_pattern and library:
                    results = substructure_search(smarts_pattern, library)
                    response = {"id": msg_id, "reply": json.dumps(results)}
                else:
                    response = {"id": msg_id, "reply": "error: missing smarts_pattern or library parameter"}
                    
            elif action == "exact_match_search":
                # 精确匹配搜索
                query_smiles = data.get("smiles")
                library = data.get("library")
                if query_smiles and library:
                    results = exact_match_search(query_smiles, library)
                    response = {"id": msg_id, "reply": json.dumps(results)}
                else:
                    response = {"id": msg_id, "reply": "error: missing smiles or library parameter"}
                    
            else:
                response = {"id": msg_id, "reply": "error: unknown action"}
                
            print(json.dumps(response))
                
        except Exception as e:
            response = {"id": msg_id if 'msg_id' in locals() else "unknown", "reply": f"error: {str(e)}"}
            print(json.dumps(response))
