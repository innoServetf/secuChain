import os
import requests
import lzma
import shutil
from datetime import datetime

def fetch_cve_data(year):
    # 构建 URL
    url = f"https://github.com/fkie-cad/nvd-json-data-feeds/releases/latest/download/CVE-{year}.json.xz"
    
    # 发送请求
    response = requests.get(url,verify=False)
    
    # 检查请求是否成功
    if response.status_code == 200:
         # 保存压缩文件
        with open(f'CVE-{year}.json.xz', 'wb') as file:
            file.write(response.content)
        
        # 解压文件
        with lzma.open(f'CVE-{year}.json.xz', 'rb') as f_in:
            with open(f'CVE-{year}.json', 'wb') as f_out:
                shutil.copyfileobj(f_in, f_out)
        
        # 删除压缩文件
        os.remove(f'CVE-{year}.json.xz')
        
        print(f"{year}年的 CVE 数据已下载并解压到当前文件夹。")
    else:
        print(f"请求失败，状态码: {response.status_code}")

# 获取当前年份
current_year = datetime.now().year
# 获取当前年份的 CVE 数据
fetch_cve_data(current_year)