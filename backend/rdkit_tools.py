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

# ==================== MS2相似度搜索函数 ====================

def parse_ms2_data(ms2_text: str):
    """
    解析MS2文本数据，提取各能量级别的质谱峰
    
    参数:
        ms2_text: MS2数据文本
        
    返回:
        字典，键为能量级别（如'energy0'），值为(m/z, intensity)列表
        对于简单的两列格式文件（如A9_165640.txt），会返回 {'energy0': [(mz1, intensity1), ...]}
    """
    if not ms2_text:
        return {}
    
    ms2_data = {}
    current_energy = None
    lines = ms2_text.strip().split('\n')
    
    # 检查是否为简单的两列格式（没有energy标记）
    # 查看前几行是否包含'energy'关键词
    has_energy_markers = any(line.strip().startswith('energy') for line in lines[:10] if line.strip())
    
    if not has_energy_markers:
        # 简单两列格式，如A9_165640.txt
        # 创建一个默认的能量级别
        current_energy = 'energy0'
        ms2_data[current_energy] = []
        
        for line in lines:
            line = line.strip()
            if not line:
                continue
                
            # 跳过注释行
            if line.startswith('#'):
                continue
                
            # 解析两列数据：m/z intensity
            parts = line.split()
            if len(parts) < 2:
                continue
                
            try:
                mz = float(parts[0])
                intensity = float(parts[1])
                
                # 添加到默认能量级别
                ms2_data[current_energy].append((mz, intensity))
            except ValueError:
                # 跳过无法解析的行
                continue
    else:
        # 原始格式：包含energy标记
        for line in lines:
            line = line.strip()
            if not line:
                continue
                
            # 检查是否为能量级别标记
            if line.startswith('energy'):
                current_energy = line
                ms2_data[current_energy] = []
                continue
                
            # 跳过注释行
            if line.startswith('#'):
                continue
                
            # 解析质谱峰数据行
            # 格式: m/z intensity fragment_ids (additional_data...)
            parts = line.split()
            if len(parts) < 2:
                continue
                
            try:
                mz = float(parts[0])
                intensity = float(parts[1])
                
                # 添加到当前能量级别
                if current_energy:
                    ms2_data[current_energy].append((mz, intensity))
            except ValueError:
                # 跳过无法解析的行
                continue
    
    return ms2_data

def calculate_ms2_fingerprint(ms2_data, bin_size=1.0, max_mz=1000.0):
    """
    计算MS2指纹向量
    
    参数:
        ms2_data: 解析后的MS2数据
        bin_size: m/z分箱大小（Da）
        max_mz: 最大m/z值
        
    返回:
        MS2指纹向量
    """
    import numpy as np
    
    # 创建分箱
    num_bins = int(max_mz / bin_size) + 1
    fingerprint = np.zeros(num_bins)
    
    # 合并所有能量级别的峰
    all_peaks = []
    for energy_level, peaks in ms2_data.items():
        all_peaks.extend(peaks)
    
    # 将峰分配到分箱中
    for mz, intensity in all_peaks:
        bin_idx = int(mz / bin_size)
        if bin_idx < num_bins:
            # 使用强度平方根作为权重（常见于质谱相似度计算）
            fingerprint[bin_idx] += np.sqrt(intensity)
    
    # 归一化
    norm = np.linalg.norm(fingerprint)
    if norm > 0:
        fingerprint = fingerprint / norm
    
    return fingerprint

def modified_cosine_similarity(peaks1, peaks2, tolerance=0.5):
    """
    计算modified cosine相似度（完整算法）
    
    参数:
        peaks1: 第一个质谱的峰列表，每个元素为(m/z, intensity)
        peaks2: 第二个质谱的峰列表，每个元素为(m/z, intensity)
        tolerance: 质量容差（Da）
        
    返回:
        相似度分数（0-1之间）
    """
    import numpy as np
    
    if not peaks1 or not peaks2:
        return 0.0
    
    # 按m/z排序
    peaks1_sorted = sorted(peaks1, key=lambda x: x[0])
    peaks2_sorted = sorted(peaks2, key=lambda x: x[0])
    
    # 初始化指针
    i, j = 0, 0
    matched_product_sum = 0.0
    
    # 滑动窗口匹配
    while i < len(peaks1_sorted) and j < len(peaks2_sorted):
        mz1, int1 = peaks1_sorted[i]
        mz2, int2 = peaks2_sorted[j]
        
        # 计算质量差
        delta_mz = abs(mz1 - mz2)
        
        if delta_mz <= tolerance:
            # 匹配成功，计算强度乘积
            matched_product_sum += int1 * int2
            
            # 两个指针都向前移动
            i += 1
            j += 1
        elif mz1 < mz2:
            # peaks1的m/z较小，移动i
            i += 1
        else:
            # peaks2的m/z较小，移动j
            j += 1
    
    # 计算分母
    sum1 = sum(int1 * int1 for mz1, int1 in peaks1_sorted)
    sum2 = sum(int2 * int2 for mz2, int2 in peaks2_sorted)
    
    if sum1 == 0 or sum2 == 0:
        return 0.0
    
    # 计算modified cosine相似度
    similarity = matched_product_sum / (np.sqrt(sum1) * np.sqrt(sum2))
    
    return max(0.0, min(1.0, similarity))

def cosine_similarity_fingerprint(fp1, fp2):
    """
    计算指纹的cosine相似度（用于快速预筛选）
    
    参数:
        fp1: 第一个指纹向量（可以是浮点数数组或二进制数组）
        fp2: 第二个指纹向量（可以是浮点数数组或二进制数组）
        
    返回:
        相似度分数（0-1之间）
    """
    import numpy as np
    
    fp1_arr = np.array(fp1)
    fp2_arr = np.array(fp2)
    
    dot_product = np.dot(fp1_arr, fp2_arr)
    norm1 = np.linalg.norm(fp1_arr)
    norm2 = np.linalg.norm(fp2_arr)
    
    if norm1 == 0 or norm2 == 0:
        return 0.0
    
    return dot_product / (norm1 * norm2)


def base64_to_binary_fingerprint(base64_str: str, original_length: int = 1001):
    """
    将base64编码转换回二进制指纹
    
    参数:
        base64_str: base64编码的字符串
        original_length: 原始指纹长度
        
    返回:
        二进制指纹数组（0和1）
    """
    import numpy as np
    
    # 解码base64
    bytes_data = base64.b64decode(base64_str)
    
    # 转换为二进制数组
    binary_bytes = np.frombuffer(bytes_data, dtype=np.uint8)
    binary_bits = np.unpackbits(binary_bytes)
    
    # 截取原始长度
    binary_fp = binary_bits[:original_length]
    
    return binary_fp


def binary_fingerprint_similarity(binary_fp1, binary_fp2):
    """
    计算二进制指纹的相似度（Tanimoto/Jaccard相似度）
    
    参数:
        binary_fp1: 第一个二进制指纹数组（0和1）
        binary_fp2: 第二个二进制指纹数组（0和1）
        
    返回:
        相似度分数（0-1之间）
    """
    import numpy as np
    
    fp1_arr = np.array(binary_fp1, dtype=np.uint8)
    fp2_arr = np.array(binary_fp2, dtype=np.uint8)
    
    # 计算交集和并集
    intersection = np.sum(fp1_arr & fp2_arr)
    union = np.sum(fp1_arr | fp2_arr)
    
    if union == 0:
        return 0.0
    
    return intersection / union

def ms2_search(query_ms2_text, library, threshold=0.5, tolerance=0.5, prefilter_threshold=0.3):
    """
    MS2相似度搜索（两阶段策略）
    
    参数:
        query_ms2_text: 查询MS2文本数据
        library: 数据库中的化合物列表，每个元素包含id和MS2_fingerprint
        threshold: 最终相似度阈值
        tolerance: 质量容差（Da）
        prefilter_threshold: 预筛选阈值
        
    返回:
        相似度结果列表，按相似度降序排序
    """
    import numpy as np
    import json
    
    # 1. 解析查询MS2数据
    query_ms2_data = parse_ms2_data(query_ms2_text)
    if not query_ms2_data:
        return []
    
    # 合并所有能量级别的峰
    query_peaks = []
    for peaks in query_ms2_data.values():
        query_peaks.extend(peaks)
    
    # 计算查询指纹
    query_fingerprint = calculate_ms2_fingerprint(query_ms2_data)
    
    results = []
    
    # 2. 第一阶段：指纹预筛选
    prefiltered_items = []
    for item in library:
        if 'MS2_fingerprint' not in item or not item['MS2_fingerprint']:
            continue
            
        try:
            # 解析指纹数据
            fingerprint_data = json.loads(item['MS2_fingerprint'])
            db_fingerprint = fingerprint_data.get('fingerprint', [])
            
            if not db_fingerprint:
                continue
                
            # 计算指纹相似度
            fingerprint_sim = cosine_similarity_fingerprint(query_fingerprint, db_fingerprint)
            
            if fingerprint_sim >= prefilter_threshold:
                prefiltered_items.append({
                    'id': item['id'],
                    'fingerprint_sim': fingerprint_sim,
                    'fingerprint_data': fingerprint_data
                })
        except:
            continue
    
    print(f"预筛选后保留 {len(prefiltered_items)} 个化合物")
    
    # 3. 第二阶段：精确modified cosine计算
    for item in prefiltered_items:
        try:
            # 获取数据库化合物的峰数据
            db_peaks = item['fingerprint_data'].get('peaks', [])
            
            if not db_peaks:
                continue
                
            # 计算modified cosine相似度
            similarity = modified_cosine_similarity(query_peaks, db_peaks, tolerance)
            
            if similarity >= threshold:
                results.append((item['id'], similarity))
        except:
            continue
    
    # 按相似度降序排序
    return sorted(results, key=lambda x: -x[1])


def ms2_search_by_fingerprint(fingerprint_json, library, threshold=0.5, tolerance=0.5, prefilter_threshold=0.3):
    """
    MS2相似度搜索（使用预计算的指纹数据）
    
    参数:
        fingerprint_json: 查询指纹JSON数据，包含peaks和fingerprint_base64
        library: 数据库中的化合物列表，每个元素包含id和MS2_fingerprint
        threshold: 最终相似度阈值
        tolerance: 质量容差（Da）
        prefilter_threshold: 预筛选阈值
        
    返回:
        相似度结果列表，按相似度降序排序
    """
    import numpy as np
    import json
    
    # 1. 解析查询指纹数据
    try:
        query_data = json.loads(fingerprint_json)
        query_peaks = query_data.get('peaks', [])
        query_fingerprint_base64 = query_data.get('fingerprint_base64', '')
        
        if not query_peaks or not query_fingerprint_base64:
            return []
            
        # 解码base64指纹为二进制指纹
        query_binary_fp = base64_to_binary_fingerprint(query_fingerprint_base64)
        
    except Exception as e:
        print(f"解析查询指纹数据失败: {e}")
        return []
    
    results = []
    
    # 2. 第一阶段：指纹预筛选（使用二进制指纹）
    prefiltered_items = []
    for item in library:
        if 'MS2_fingerprint' not in item or not item['MS2_fingerprint']:
            continue
            
        try:
            # 解析指纹数据
            fingerprint_data = json.loads(item['MS2_fingerprint'])
            db_fingerprint_base64 = fingerprint_data.get('fingerprint_base64', '')
            db_peaks = fingerprint_data.get('peaks', [])
            
            if not db_fingerprint_base64 or not db_peaks:
                continue
                
            # 解码base64指纹为二进制指纹
            db_binary_fp = base64_to_binary_fingerprint(db_fingerprint_base64)
            
            # 计算二进制指纹相似度（Tanimoto/Jaccard）
            fingerprint_sim = binary_fingerprint_similarity(query_binary_fp, db_binary_fp)
            
            if fingerprint_sim >= prefilter_threshold:
                prefiltered_items.append({
                    'id': item['id'],
                    'fingerprint_sim': fingerprint_sim,
                    'fingerprint_data': fingerprint_data
                })
        except Exception as e:
            print(f"处理数据库化合物 {item.get('id', 'unknown')} 失败: {e}")
            continue
    
    print(f"预筛选后保留 {len(prefiltered_items)} 个化合物")
    
    # 3. 第二阶段：精确modified cosine计算
    for item in prefiltered_items:
        try:
            # 获取数据库化合物的峰数据
            db_peaks = item['fingerprint_data'].get('peaks', [])
            
            if not db_peaks:
                continue
                
            # 计算modified cosine相似度
            similarity = modified_cosine_similarity(query_peaks, db_peaks, tolerance)
            
            if similarity >= threshold:
                results.append((item['id'], similarity))
        except Exception as e:
            print(f"计算modified cosine相似度失败: {e}")
            continue
    
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
                    
            elif action == "ms2_search":
                # MS2相似度搜索
                query_ms2 = data.get("query_ms2")
                library = data.get("library")
                threshold = float(data.get('threshold', 0.5))
                tolerance = float(data.get('tolerance', 0.5))
                prefilter_threshold = float(data.get('prefilter_threshold', 0.3))
                
                if query_ms2 and library:
                    results = ms2_search(query_ms2, library, threshold, tolerance, prefilter_threshold)
                    response = {"id": msg_id, "reply": json.dumps(results)}
                else:
                    response = {"id": msg_id, "reply": "error: missing query_ms2 or library parameter"}
                    
            elif action == "ms2_search_by_fingerprint":
                # MS2指纹搜索
                fingerprint_json = data.get("fingerprint_json")
                library = data.get("library")
                threshold = float(data.get('threshold', 0.5))
                tolerance = float(data.get('tolerance', 0.5))
                prefilter_threshold = float(data.get('prefilter_threshold', 0.3))
                
                if fingerprint_json and library:
                    results = ms2_search_by_fingerprint(fingerprint_json, library, threshold, tolerance, prefilter_threshold)
                    response = {"id": msg_id, "reply": json.dumps(results)}
                else:
                    response = {"id": msg_id, "reply": "error: missing fingerprint_json or library parameter"}
                    
            else:
                response = {"id": msg_id, "reply": "error: unknown action"}
                
            print(json.dumps(response))
                
        except Exception as e:
            response = {"id": msg_id if 'msg_id' in locals() else "unknown", "reply": f"error: {str(e)}"}
            print(json.dumps(response))
