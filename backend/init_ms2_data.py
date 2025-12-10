#!/usr/bin/env python3
"""
初始化MS2中间数据并保存到数据库
用于modified cosine相似度搜索
使用MySQL数据库
"""

import os
import json
import yaml
import pymysql
from typing import List, Tuple, Dict, Any
import numpy as np

def parse_ms2_data(ms2_text: str) -> Dict[str, List[Tuple[float, float]]]:
    """
    解析MS2文本数据，提取各能量级别的质谱峰
    
    参数:
        ms2_text: MS2数据文本
        
    返回:
        字典，键为能量级别（如'energy0'），值为(m/z, intensity)列表
    """
    if not ms2_text:
        return {}
    
    ms2_data = {}
    current_energy = None
    lines = ms2_text.strip().split('\n')
    
    for line in lines:
        line = line.strip()
        if not line and current_energy == 'energy2':
            break
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

def calculate_ms2_fingerprint(ms2_data: Dict[str, List[Tuple[float, float]]], 
                             bin_size: float = 1.0, 
                             max_mz: float = 1000.0) -> np.ndarray:
    """
    计算MS2指纹向量
    
    参数:
        ms2_data: 解析后的MS2数据
        bin_size: m/z分箱大小（Da）
        max_mz: 最大m/z值
        
    返回:
        MS2指纹向量
    """
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

def modified_cosine_similarity(peaks1: List[Tuple[float, float]], 
                              peaks2: List[Tuple[float, float]], 
                              tolerance: float = 0.5) -> float:
    """
    计算modified cosine相似度（完整算法）
    
    参数:
        peaks1: 第一个质谱的峰列表，每个元素为(m/z, intensity)
        peaks2: 第二个质谱的峰列表，每个元素为(m/z, intensity)
        tolerance: 质量容差（Da）
        
    返回:
        相似度分数（0-1之间）
        
    算法描述:
        1. 对两个质谱的峰按m/z排序
        2. 使用滑动窗口匹配峰，考虑质量容差
        3. 计算匹配峰的强度乘积之和
        4. 计算modified cosine相似度
    """
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

def float_fingerprint_to_binary(fp: np.ndarray, threshold: float = 0.01) -> np.ndarray:
    """
    将浮点数指纹转换为二进制指纹
    
    参数:
        fp: 浮点数指纹数组
        threshold: 阈值，大于该值的设为1，否则为0
        
    返回:
        二进制指纹数组（0和1）
    """
    binary_fp = (fp > threshold).astype(int)
    return binary_fp


def binary_fingerprint_to_base64(binary_fp: np.ndarray) -> str:
    """
    将二进制指纹转换为base64编码
    
    参数:
        binary_fp: 二进制指纹数组（0和1）
        
    返回:
        base64编码的字符串
    """
    import base64
    
    # 将二进制数组转换为字节
    # 每8位转换为一个字节
    # 确保长度是8的倍数
    padded_length = ((len(binary_fp) + 7) // 8) * 8
    padded = np.zeros(padded_length, dtype=np.uint8)
    padded[:len(binary_fp)] = binary_fp
    
    # 转换为字节
    bytes_data = np.packbits(padded).tobytes()
    
    # 转换为base64
    return base64.b64encode(bytes_data).decode('ascii')


def base64_to_binary_fingerprint(base64_str: str, original_length: int = 1001) -> np.ndarray:
    """
    将base64编码转换回二进制指纹
    
    参数:
        base64_str: base64编码的字符串
        original_length: 原始指纹长度
        
    返回:
        二进制指纹数组（0和1）
    """
    import base64
    
    # 解码base64
    bytes_data = base64.b64decode(base64_str)
    
    # 转换为二进制数组
    binary_bytes = np.frombuffer(bytes_data, dtype=np.uint8)
    binary_bits = np.unpackbits(binary_bytes)
    
    # 截取原始长度
    binary_fp = binary_bits[:original_length]
    
    return binary_fp


def calculate_ms2_fingerprint_modified(ms2_data: Dict[str, List[Tuple[float, float]]]) -> Dict[str, Any]:
    """
    计算MS2指纹（modified cosine版本）
    
    参数:
        ms2_data: 解析后的MS2数据
        
    返回:
        包含峰列表和base64编码的二进制指纹的字典
    """
    # 合并所有能量级别的峰
    all_peaks = []
    for energy_level, peaks in ms2_data.items():
        all_peaks.extend(peaks)
    
    # 计算浮点数指纹向量（用于快速预筛选）
    fingerprint_float = calculate_ms2_fingerprint(ms2_data)
    
    # 转换为二进制指纹
    fingerprint_binary = float_fingerprint_to_binary(fingerprint_float, threshold=0.01)
    
    # 转换为base64编码
    fingerprint_base64 = binary_fingerprint_to_base64(fingerprint_binary)
    
    return {
        'peaks': all_peaks,
        'fingerprint_base64': fingerprint_base64  # 只存储base64编码的二进制指纹
    }

def load_config(config_path: str = 'config.yaml') -> Dict[str, Any]:
    """
    加载配置文件
    
    参数:
        config_path: 配置文件路径
        
    返回:
        配置字典
    """
    if not os.path.exists(config_path):
        # 尝试使用示例配置文件
        config_path = 'config_example.yaml'
        if not os.path.exists(config_path):
            raise FileNotFoundError(f"配置文件 {config_path} 不存在")
    
    with open(config_path, 'r', encoding='utf-8') as f:
        config = yaml.safe_load(f)
    
    return config

def init_database_ms2_data(config_path: str = 'config.yaml'):
    """
    初始化数据库中的MS2数据
    
    参数:
        config_path: 配置文件路径
    """
    # 加载配置
    config = load_config(config_path)
    
    db_config = config.get('database', {})
    host = db_config.get('host', '127.0.0.1')
    port = db_config.get('port', 3306)
    user = db_config.get('user', 'root')
    password = db_config.get('pass', '')
    database = db_config.get('name', 'test')
    
    print(f"连接数据库: {user}@{host}:{port}/{database}")
    
    try:
        # 连接数据库
        conn = pymysql.connect(
            host=host,
            port=port,
            user=user,
            password=password,
            database=database,
            charset='utf8mb4',
            cursorclass=pymysql.cursors.DictCursor
        )
        
        with conn.cursor() as cursor:
            # 检查MS2_fingerprint列是否存在
            cursor.execute("SHOW COLUMNS FROM data LIKE 'MS2_fingerprint'")
            column_exists = cursor.fetchone()
            
            if not column_exists:
                print("添加MS2_fingerprint列到data表...")
                cursor.execute("ALTER TABLE data ADD COLUMN MS2_fingerprint TEXT")
                conn.commit()
            
            # 获取所有有MS2数据的化合物
            cursor.execute("SELECT ID, MS2_full FROM data WHERE MS2_full IS NOT NULL AND MS2_full != ''")
            compounds = cursor.fetchall()
            
            print(f"找到 {len(compounds)} 个有MS2数据的化合物")
            
            for i, compound in enumerate(compounds):
                compound_id = compound['ID']
                ms2_text = compound['MS2_full']
                
                if i % 10 == 0:
                    print(f"处理第 {i+1}/{len(compounds)} 个化合物: {compound_id}")
                
                try:
                    # 解析MS2数据
                    ms2_data = parse_ms2_data(ms2_text)
                    
                    if not ms2_data:
                        print(f"  化合物 {compound_id}: 无法解析MS2数据")
                        continue
                    
                    # 计算MS2指纹（modified cosine版本）
                    fingerprint_data = calculate_ms2_fingerprint_modified(ms2_data)
                    
                    # 将指纹转换为JSON字符串保存
                    fingerprint_json = json.dumps(fingerprint_data)
                    
                    # 更新数据库
                    cursor.execute(
                        "UPDATE data SET MS2_fingerprint = %s WHERE ID = %s",
                        (fingerprint_json, compound_id)
                    )
                    
                except Exception as e:
                    print(f"  化合物 {compound_id} 处理失败: {str(e)}")
            
            conn.commit()
            print("MS2数据初始化完成")
            
    except pymysql.Error as e:
        print(f"数据库连接失败: {e}")
    finally:
        if 'conn' in locals() and conn.open:
            conn.close()

def test_ms2_similarity():
    """测试MS2相似度计算"""
    # 示例MS2数据
    ms2_text1 = """energy0
271.06010 15.81 45 (2.9501)
273.07575 50.41 49 53 64 62 61 60 55 58 57 56 92 (7.0903 1.1007 0.33349 0.31693 0.24094 0.11569 0.11206 0.052208 0.03833 0.0080716 1.1611e-05)
energy1
55.05423 5.42 32 147 (0.86092 0.42807)
65.03858 6.17 44 (1.4683)"""
    
    ms2_text2 = """energy0
271.06010 20.81 45 (2.9501)
273.07575 55.41 49 53 64 62 61 60 55 58 57 56 92 (7.0903 1.1007 0.33349 0.31693 0.24094 0.11569 0.11206 0.052208 0.03833 0.0080716 1.1611e-05)
energy1
55.05423 8.42 32 147 (0.86092 0.42807)
65.03858 9.17 44 (1.4683)"""
    
    print("测试MS2相似度计算...")
    
    # 解析MS2数据
    ms2_data1 = parse_ms2_data(ms2_text1)
    ms2_data2 = parse_ms2_data(ms2_text2)
    
    print(f"MS2数据1能量级别: {list(ms2_data1.keys())}")
    print(f"MS2数据2能量级别: {list(ms2_data2.keys())}")
    
    # 合并所有能量级别的峰
    all_peaks1 = []
    for peaks in ms2_data1.values():
        all_peaks1.extend(peaks)
    
    all_peaks2 = []
    for peaks in ms2_data2.values():
        all_peaks2.extend(peaks)
    
    print(f"质谱1峰数量: {len(all_peaks1)}")
    print(f"质谱2峰数量: {len(all_peaks2)}")
    
    # 计算modified cosine相似度
    similarity = modified_cosine_similarity(all_peaks1, all_peaks2, tolerance=0.5)
    print(f"Modified cosine相似度: {similarity:.4f}")
    
    # 测试指纹计算
    fingerprint_data1 = calculate_ms2_fingerprint_modified(ms2_data1)
    fingerprint_data2 = calculate_ms2_fingerprint_modified(ms2_data2)
    
    print(f"指纹1长度: {len(fingerprint_data1['fingerprint'])}, 峰数量: {len(fingerprint_data1['peaks'])}")
    print(f"指纹2长度: {len(fingerprint_data2['fingerprint'])}, 峰数量: {len(fingerprint_data2['peaks'])}")

if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description='初始化MS2中间数据')
    parser.add_argument('--config', default='config.yaml', help='配置文件路径')
    parser.add_argument('--test', action='store_true', help='运行测试')
    
    args = parser.parse_args()
    
    if args.test:
        test_ms2_similarity()
    else:
        init_database_ms2_data(args.config)
