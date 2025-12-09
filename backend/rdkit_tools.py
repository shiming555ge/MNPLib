from rdkit import Chem, DataStructs
from rdkit.Chem import AllChem, Descriptors

import base64, json

# smiles式转存pdb
def smiles_to_pdb(smiles):
    mol = Chem.MolFromSmiles(smiles)
    mol = Chem.AddHs(mol)                       # 加氢
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
    # 首先尝试直接匹配
    patt = Chem.MolFromSmarts(smarts_pattern)
    mol = Chem.MolFromSmiles(smiles)
    if patt is None or mol is None:
        return False
    
    if mol.HasSubstructMatch(patt):
        return True
    
    return False

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
    # 首先尝试直接解析SMARTS模式
    patt = Chem.MolFromSmarts(pattern_smarts)
    if patt is None:
        return []  # 无效的SMARTS模式
    
    result = []
    for item in library:
        smiles = item['smiles']
        if not smiles:
            continue
            
        mol = Chem.MolFromSmiles(smiles)
        if mol is None:
            continue
            
        # 尝试匹配
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

# 计算分子量
def calculate_molecular_weight(smiles):
    mol = Chem.MolFromSmiles(smiles)
    if mol is None:
        return None
    return Descriptors.MolWt(mol)

# 解析核磁谱数据
def parse_nmr_data(nmr_text):
    """
    解析核磁谱文本数据，提取化学位移值
    格式示例: "δ 170.5, 160.3, 140.2, 130.1, 120.5, 110.3"
    或: "170.5, 160.3, 140.2, 130.1, 120.5, 110.3"
    """
    if not nmr_text:
        return []
    
    # 移除常见的核磁谱标记
    text = nmr_text.replace('δ', '').replace('ppm', '').strip()
    
    # 分割并提取数字
    values = []
    for part in text.split(','):
        part = part.strip()
        # 尝试提取数字
        try:
            # 移除非数字字符
            num_str = ''
            for char in part:
                if char.isdigit() or char == '.' or char == '-':
                    num_str += char
            if num_str:
                value = float(num_str)
                values.append(value)
        except:
            continue
    
    return sorted(values)

# 计算核磁谱相似度
def calculate_nmr_similarity(query_nmr, db_nmr, tolerance=0.5):
    """
    计算两个核磁谱的相似度
    使用简单的峰值匹配算法，允许一定的容差
    """
    if not query_nmr or not db_nmr:
        return 0.0
    
    query_peaks = parse_nmr_data(query_nmr)
    db_peaks = parse_nmr_data(db_nmr)
    
    if not query_peaks or not db_peaks:
        return 0.0
    
    matched_count = 0
    query_matched = [False] * len(query_peaks)
    db_matched = [False] * len(db_peaks)
    
    # 尝试匹配每个查询峰值
    for i, q_peak in enumerate(query_peaks):
        for j, d_peak in enumerate(db_peaks):
            if not db_matched[j] and abs(q_peak - d_peak) <= tolerance:
                matched_count += 1
                query_matched[i] = True
                db_matched[j] = True
                break
    
    # 计算相似度分数
    # 使用F1分数：2 * precision * recall / (precision + recall)
    precision = matched_count / len(query_peaks) if query_peaks else 0
    recall = matched_count / len(db_peaks) if db_peaks else 0
    
    if precision + recall == 0:
        return 0.0
    
    f1_score = 2 * precision * recall / (precision + recall)
    return f1_score

# 核磁谱搜索
def nmr_search(query_nmr, library, threshold=0.5, tolerance=0.5):
    """
    在数据库中搜索相似的核磁谱
    library: 包含id和nmr_13c_data的字典列表
    """
    results = []
    
    for item in library:
        db_nmr = item.get('nmr_13c_data', '')
        if not db_nmr:
            continue
            
        similarity = calculate_nmr_similarity(query_nmr, db_nmr, tolerance)
        if similarity >= threshold:
            results.append((item['id'], similarity))
    
    # 按相似度降序排序
    return sorted(results, key=lambda x: -x[1])

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
                    
            elif action == "calculate_molecular_weight":
                # 计算分子量
                smiles = data.get("smiles")
                if smiles:
                    weight = calculate_molecular_weight(smiles)
                    if weight is not None:
                        response = {"id": msg_id, "reply": str(weight)}
                    else:
                        response = {"id": msg_id, "reply": "error: invalid smiles"}
                else:
                    response = {"id": msg_id, "reply": "error: missing smiles parameter"}
                    
            elif action == "nmr_search":
                # 核磁谱搜索
                query_nmr = data.get("query_nmr")
                library = data.get("library")
                threshold = float(data.get('threshold', 0.5))
                tolerance = float(data.get('tolerance', 0.5))
                
                if query_nmr and library:
                    results = nmr_search(query_nmr, library, threshold, tolerance)
                    response = {"id": msg_id, "reply": json.dumps(results)}
                else:
                    response = {"id": msg_id, "reply": "error: missing query_nmr or library parameter"}
                    
            else:
                response = {"id": msg_id, "reply": "error: unknown action"}
                
            print(json.dumps(response))
                
        except Exception as e:
            response = {"id": msg_id if 'msg_id' in locals() else "unknown", "reply": f"error: {str(e)}"}
            print(json.dumps(response))
