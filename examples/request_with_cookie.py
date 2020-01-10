#!/usr/bin/python3
# coding: utf-8

import simplejson
import subprocess
"""
    添加Cookie扫描示例
    
    命令行调用时：
    ./crawlergo -c /home/test/chrome-linux/chrome -o json --ignore-url-keywords quit,exit,zhuxiao --custom-headers "{\"Cookie\": \"crawlergo=Cool\"}"

    使用 --ignore-url-keywords 添加你想要的排除的关键字，避免访问注销请求
"""


def main():
    target = "http://testphp.vulnweb.com/"
    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) "
                      "Chrome/74.0.3945.0 Safari/537.36",
        "Cookie": "crawlergo=Cool"
    }
    cmd = ["./crawlergo", "-c", "/home/test/chrome-linux/chrome",
           "-o", "json", "--ignore-url-keywords", "quit,exit,zhuxiao", "--custom-headers", simplejson.dumps(headers),
           target]

    rsp = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = rsp.communicate()

    result = simplejson.loads(output.decode().split("--[Mission Complete]--")[1])
    req_list = result["req_list"]
    for each in req_list:
        print(each)


if __name__ == '__main__':
    main()





