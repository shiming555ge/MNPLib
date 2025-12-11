#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试NMR查找功能
"""

import sys
import os

# 添加当前目录到路径
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

# 测试两个版本的rdkit_tools
print("测试NMR查找功能...")
print("=" * 80)

# 首先测试rdkit_tools_f.py（可以正常运行的版本）
print("\n1. 测试 rdkit_tools_f.py（正常版本）")
try:
    # 动态导入模块
    import importlib.util
    
    # 导入rdkit_tools_f.py
    spec_f = importlib.util.spec_from_file_location("rdkit_tools_f", "rdkit_tools_f.py")
    module_f = importlib.util.module_from_spec(spec_f)
    spec_f.loader.exec_module(module_f)
    
    print("✓ 成功导入 rdkit_tools_f.py")
    
    # 测试parse_nmr_data函数
    test_nmr_text = "δ 170.5, 160.3, 140.2, 130.1, 120.5, 110.3"
    result_f = module_f.parse_nmr_data(test_nmr_text)
    print(f"  parse_nmr_data('{test_nmr_text}') = {result_f}")
    
    # 测试calculate_nmr_similarity函数
    query_nmr = "170.5, 160.3, 140.2"
    db_nmr = "170.5, 160.5, 140.0"
    similarity_f = module_f.calculate_nmr_similarity(query_nmr, db_nmr, tolerance=0.5)
    print(f"  calculate_nmr_similarity('{query_nmr}', '{db_nmr}') = {similarity_f}")
    
    # 测试nmr_search函数
    library = [
        {"id": "compound1", "nmr_13c_data": "170.5, 160.3, 140.2"},
        {"id": "compound2", "nmr_13c_data": "180.0, 150.0, 130.0"},
        {"id": "compound3", "nmr_13c_data": "170.5, 160.5, 140.0"}
    ]
    results_f = module_f.nmr_search(query_nmr, library, threshold=0.5, tolerance=0.5)
    print(f"  nmr_search('{query_nmr}', library) = {results_f}")
    
except Exception as e:
    print(f"✗ 测试 rdkit_tools_f.py 失败: {e}")
    import traceback
    traceback.print_exc()

print("\n" + "=" * 80)

# 然后测试../backend/rdkit_tools.py（更新后有问题的版本）
print("\n2. 测试 ../backend/rdkit_tools.py（更新版本）")
try:
    # 导入../backend/rdkit_tools.py
    spec_backend = importlib.util.spec_from_file_location("rdkit_tools_backend", "../backend/rdkit_tools.py")
    module_backend = importlib.util.module_from_spec(spec_backend)
    spec_backend.loader.exec_module(module_backend)
    
    print("✓ 成功导入 ../backend/rdkit_tools.py")
    
    # 测试parse_nmr_data函数
    test_nmr_text = "δ 170.5, 160.3, 140.2, 130.1, 120.5, 110.3"
    result_backend = module_backend.parse_nmr_data(test_nmr_text)
    print(f"  parse_nmr_data('{test_nmr_text}') = {result_backend}")
    
    # 测试calculate_nmr_similarity函数
    query_nmr = "170.5, 160.3, 140.2"
    db_nmr = "170.5, 160.5, 140.0"
    similarity_backend = module_backend.calculate_nmr_similarity(query_nmr, db_nmr, tolerance=0.5)
    print(f"  calculate_nmr_similarity('{query_nmr}', '{db_nmr}') = {similarity_backend}")
    
    # 测试nmr_search函数
    library = [
        {"id": "compound1", "nmr_13c_data": "170.5, 160.3, 140.2"},
        {"id": "compound2", "nmr_13c_data": "180.0, 150.0, 130.0"},
        {"id": "compound3", "nmr_13c_data": "170.5, 160.5, 140.0"}
    ]
    results_backend = module_backend.nmr_search(query_nmr, library, threshold=0.5, tolerance=0.5)
    print(f"  nmr_search('{query_nmr}', library) = {results_backend}")
    
except Exception as e:
    print(f"✗ 测试 ../backend/rdkit_tools.py 失败: {e}")
    import traceback
    traceback.print_exc()

print("\n" + "=" * 80)
print("\n测试完成！")
