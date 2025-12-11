#!/usr/bin/env python3
"""
测试NMR功能
"""

import json
import requests
import sys

def test_nmr_search():
    """测试NMR搜索API"""
    
    # 测试数据
    test_nmr_data = "170.5, 160.3, 140.2, 130.1, 120.5, 110.3"
    
    # API端点
    base_url = "http://localhost:9090/api/rdkit"
    
    print("测试NMR搜索功能...")
    
    # 测试NMR搜索
    params = {
        "query_nmr": test_nmr_data,
        "threshold": "0.5",
        "tolerance": "0.5"
    }
    
    try:
        response = requests.get(f"{base_url}/nmr-search", params=params, timeout=30)
        print(f"状态码: {response.status_code}")
        print(f"响应头: {response.headers}")
        print(f"响应内容: {response.text}")
        
        if response.status_code == 200:
            result = response.json()
            print(f"成功! 返回结果: {result}")
        else:
            print(f"失败! 状态码: {response.status_code}")
    except Exception as e:
        print(f"请求失败: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    print("开始测试NMR搜索功能...")
    print("=" * 60)
    
    # 检查服务器是否运行
    try:
        response = requests.get("http://localhost:9090/api/rdkit/status", timeout=5)
        if response.status_code == 200:
            print("服务器正在运行")
            test_nmr_search()
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
