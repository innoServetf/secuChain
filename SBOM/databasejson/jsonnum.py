from datetime import datetime
import json

year = datetime.now().year
# 指定 JSON 文件路径，使用 f-string 格式化字符串
file_path = f"CVE-{year}.json"

try:
    # 打开并读取 JSON 文件
    with open(file_path, 'r', encoding='utf-8') as file:
        data = json.load(file)  # 将 JSON 数据加载为 Python 对象

    # 获取 cve_items 数组
    cve_items = data.get('cve_items', [])

    # 计算 cve_items 数组中的条目数量
    num_cve_items = len(cve_items)

    # 输出条目数量
    print(f'cve_items 中有 {num_cve_items} 条数据。')

except json.JSONDecodeError as e:
    print(f'JSON 解析错误: {e}')
except FileNotFoundError:
    print(f'文件未找到: {file_path}')