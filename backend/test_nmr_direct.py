#!/usr/bin/env python3
"""
直接测试rdkit_tools.py中的NMR函数
"""

import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

# 导入rdkit_tools中的函数
from rdkit_tools import parse_nmr_data, calculate_nmr_similarity, nmr_search

def test_nmr_functions():
    """测试NMR相关函数"""
    
    print("测试parse_nmr_data函数...")
    
    # 测试数据
    test_nmr_text1 = "170.5, 160.3, 140.2, 130.1, 120.5, 110.3"
    test_nmr_text2 = "δ 170.5, 160.3, 140.2, 130.1, 120.5, 110.3 ppm"
    test_nmr_text3 = "170.5 160.3 140.2"
    
    print(f"测试1: {test_nmr_text1}")
    result1 = parse_nmr_data(test_nmr_text1)
    print(f"  结果: {result1}")
    
    print(f"\n测试2: {test_nmr_text2}")
    result2 = parse_nmr_data(test_nmr_text2)
    print(f"  结果: {result2}")
    
    print(f"\n测试3: {test_nmr_text3}")
    result3 = parse_nmr_data(test_nmr_text3)
    print(f"  结果: {result3}")
    
    print("\n测试calculate_nmr_similarity函数...")
    similarity = calculate_nmr_similarity(test_nmr_text1, test_nmr_text2, tolerance=0.5)
    print(f"  相似度: {similarity}")
    
    print("\n测试nmr_search函数...")
    library = [
        {"id": "MNP001", "nmr_13c_data": "170.5, 160.3, 140.2, 130.1, 120.5, 110.3"},
        {"id": "MNP002", "nmr_13c_data": "180.5, 170.3, 150.2, 140.1, 130.5, 120.3"},
        {"id": "MNP003", "nmr_13c_data": "190.5, 180.3, 160.2, 150.1, 140.5, 130.3"},
    ]
    
    results = nmr_search(test_nmr_text1, library, threshold=0.5, tolerance=0.5)
    print(f"  搜索结果: {results}")
    
    print("\n测试完成!")

if __name__ == "__main__":
    print("开始直接测试NMR函数...")
    print("=" * 60)
    
    try:
        test_nmr_functions()
    except Exception as e:
        print(f"测试失败: {e}")
        import traceback
        traceback.print_exc()
    
    print("\n" + "=" * 60)
    print("测试完成!")
