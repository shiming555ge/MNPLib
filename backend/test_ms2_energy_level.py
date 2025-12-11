#!/usr/bin/env python3
"""
测试MS2按能量级别搜索功能
"""

import json
import requests
import sys

def test_ms2_search_by_energy_level():
    """测试MS2按能量级别搜索API"""
    
    # 测试数据
    test_ms2_data = """energy0
271.06010 15.81 45 (2.9501)
273.07575 50.41 49 53 64 62 61 60 55 58 57 56 92 (7.0903 1.1007 0.33349 0.31693 0.24094 0.11569 0.11206 0.052208 0.03833 0.0080716 1.1611e-05)
energy1
55.05423 5.42 32 147 (0.86092 0.42807)
65.03858 6.17 44 (1.4683)"""
    
    # API端点
    base_url = "http://localhost:9090/api/rdkit"
    
    # 测试1: 使用energy0能量级别
    print("测试1: 使用energy0能量级别搜索")
    payload = {
        "query_ms2": test_ms2_data,
        "threshold": 0.5,
        "tolerance": 0.5,
        "prefilter_threshold": 0.3,
        "energy_level": "energy0"
    }
    
    try:
        response = requests.post(f"{base_url}/ms2-search-by-energy-level", json=payload, timeout=30)
        if response.status_code == 200:
            result = response.json()
            print(f"  成功! 返回结果: {result}")
        else:
            print(f"  失败! 状态码: {response.status_code}, 响应: {response.text}")
    except Exception as e:
        print(f"  请求失败: {e}")
    
    # 测试2: 使用energy1能量级别
    print("\n测试2: 使用energy1能量级别搜索")
    payload["energy_level"] = "energy1"
    
    try:
        response = requests.post(f"{base_url}/ms2-search-by-energy-level", json=payload, timeout=30)
        if response.status_code == 200:
            result = response.json()
            print(f"  成功! 返回结果: {result}")
        else:
            print(f"  失败! 状态码: {response.status_code}, 响应: {response.text}")
    except Exception as e:
        print(f"  请求失败: {e}")
    
    # 测试3: 使用所有能量级别（默认）
    print("\n测试3: 使用所有能量级别搜索（默认）")
    payload["energy_level"] = ""
    
    try:
        response = requests.post(f"{base_url}/ms2-search-by-energy-level", json=payload, timeout=30)
        if response.status_code == 200:
            result = response.json()
            print(f"  成功! 返回结果: {result}")
        else:
            print(f"  失败! 状态码: {response.status_code}, 响应: {response.text}")
    except Exception as e:
        print(f"  请求失败: {e}")
    
    # 测试4: 测试旧的MS2搜索API（向后兼容）
    print("\n测试4: 测试旧的MS2搜索API（向后兼容）")
    payload_old = {
        "query_ms2": test_ms2_data,
        "threshold": 0.5,
        "tolerance": 0.5,
        "prefilter_threshold": 0.3
    }
    
    try:
        response = requests.post(f"{base_url}/ms2-search", json=payload_old, timeout=30)
        if response.status_code == 200:
            result = response.json()
            print(f"  成功! 返回结果: {result}")
        else:
            print(f"  失败! 状态码: {response.status_code}, 响应: {response.text}")
    except Exception as e:
        print(f"  请求失败: {e}")

if __name__ == "__main__":
    print("开始测试MS2按能量级别搜索功能...")
    print("=" * 60)
    
    # 检查服务器是否运行
    try:
        response = requests.get("http://localhost:9090/api/rdkit/status", timeout=5)
        if response.status_code == 200:
            print("服务器正在运行")
            test_ms2_search_by_energy_level()
        else:
            print(f"服务器状态异常: {response.status_code}")
            print("请先启动服务器: go run main.go")
            sys.exit(1)
    except requests.exceptions.ConnectionError:
        print("服务器未运行，请先启动服务器: go run main.go")
        sys.exit(1)
    except Exception as e:
        print(f"检查服务器状态时出错: {e}")
        sys.exit(1)
    
    print("\n" + "=" * 60)
    print("测试完成!")
