# crawlergo

![chromedp](https://img.shields.io/badge/chromedp-v0.5.2-brightgreen.svg)![Chromium version](https://img.shields.io/badge/chromium-79.0.3945.0-important.svg)![SKP](https://img.shields.io/badge/Project-360天相-blue.svg)

> A powerful dynamic crawler for web vulnerability scanners

crawlergo是一个使用`chrome headless`模式进行URL入口收集的**动态爬虫**。 使用Golang语言开发，基于[chromedp](https://github.com/chromedp/chromedp) 进行一些定制化开发后操纵CDP协议，对整个页面关键点进行HOOK，灵活表单填充提交，完整的事件触发，尽可能的收集网站暴露出的入口。同时，依靠智能URL去重模块，在过滤掉了大多数伪静态URL之后，仍然确保不遗漏关键入口链接，大幅减少重复任务。

crawlergo 目前支持以下特性：

* 原生浏览器环境，协程池调度任务
* 表单智能填充、自动化提交
* 完整DOM事件收集，自动化触发
* 智能URL去重，去掉大部分的重复页面
* 全面分析收集，包括javascript文件内容、页面注释、robots.txt文件和常见路径Fuzz
* 支持Host绑定，自动添加Referer。

目前开放编译好的程序给大家使用，该项目属于商业化产品的一部分，代码暂无法开源。

## 运行截图

![](D:\go_projects\crawlergo\imgs\demo.gif)

## 安装

**安装使用之前，请仔细阅读并确认[免责声明](./Disclaimer.md)。**

1. crawlergo 只依赖或chrome运行即可，前往[下载](https://www.chromium.org/getting-involved/download-chromium)新版本的chromium，或者直接[点击下载Linux79版本](https://storage.googleapis.com/chromium-browser-snapshots/Linux_x64/706915/chrome-linux.zip)。

2. 前往[页面下载](https://github.com/0Kee-Team/crawlergo/releases)最新版本的crawlergo解压到任意目录，如果是linux或者macOS系统，请赋予crawlergo**可执行权限(+x)**。

## Quick Start

### Go！

假设你的chromium安装在 `/tmp/chromium/` ，开启最大20标签页，爬取AWVS靶场：

```shell
./crawlergo -c /tmp/chromium/chrome -t 20 http://testphp.vulnweb.com/
```

### 系统调用

默认打印当前域名请求，但多数情况我们希望调用crawlergo返回的结果，所以设置输出模式为 `json`，使用python调用并收集结果的示例如下：

```python
#!/usr/bin/python3
# coding: utf-8

import simplejson
import subprocess


def main():
    target = "http://testphp.vulnweb.com/"
    cmd = ["./crawlergo", "-c", "/tmp/chromium/chrome", "-o", "json", target]
    rsp = subprocess.Popen(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    output, error = rsp.communicate()
	#  "--[Mission Complete]--"  是任务结束的分隔字符串
    result = simplejson.loads(output.decode().split("--[Mission Complete]--")[1])
    req_list = result["req_list"]
    print(req_list[0])


if __name__ == '__main__':
    main()
```

### 返回结果

当设置输出模式为 `json`时，返回的结果反序列化之后包含四个部分：

* `all_req_list`： 本次爬取任务过程中发现的所有请求，包含其他域名的任何资源类型。
* `req_list`：本次爬取任务的**同域名结果**，经过伪静态去重，不包含静态资源链接。理论上是 `all_req_list `的子集
* `all_domain_list`：发现的所有域名列表。
* `sub_domain_list`：发现的任务目标的子域名列表。

## 参数说明

crawlergo 拥有灵活的参数配置，以下是详细的选项说明：

* `--chromium-path Path, -c Path`    chrome的可执行程序路径
* `--custom-headers Headers`   自定义HTTP头，使用传入json序列化之后的数据，这个是全局定义，将被用于所有请求
* `--post-data PostData, -d PostData`   提供POST数据，目标使用POST请求方法
* `--max-crawled-count Number, -m Number`   爬虫最大任务数量，避免因伪静态造成长时间无意义抓取。
* `--filter-mode Mode, -f Mode`   过滤模式，简单：只过滤静态资源和完全重复的请求。智能：拥有过滤伪静态的能力。严格：更加严格的伪静态过滤规则。
* `--output-mode value, -o value`   结果输出模式，console：打印当前域名结果。json：打印所有结果的json序列化字符串，可直接被反序列化解析。
* `--incognito-context, -i`   浏览器启动隐身模式
* `--max-tab-count Number, -t Number`   爬虫同时开启最大标签页，即同时爬取的页面数量。
* `--fuzz-path`  使用常见路径Fuzz目标，获取更多入口。
* `--robots-path` 从robots.txt 文件中解析路径，获取更多入口。
* `--tab-run-timeout Timeout`   单个Tab标签页的最大运行超时。
* `--wait-dom-content-loaded-timeout Timeout`  爬虫等待页面加载完毕的最大超时。

## Bypass headless detect

https://intoli.com/blog/not-possible-to-block-chrome-headless/chrome-headless-test.html

![](D:\go_projects\crawlergo\imgs\bypass.png)

## 关于360天相

crawlergo是[**360天相**](https://skp.360.cn/)的子模块，天相是360自研的**资产管理与威胁探测系统**，主打强大的资产识别能力和全方位分析体系，拥有高效率的扫描能力，核心技术由 [360 0KeeTeam](https://0kee.360.cn/) 和 **360 RedTeam** 提供支持。

![](D:\go_projects\crawlergo\imgs\skp.png)

详情请访问：[https://skp.360.cn/](https://skp.360.cn/)

## 推荐用法

crawlergo 返回了全量的请求和URL信息，可以有多种使用方法：

* 联动其它的开源被动扫描器  example
* 子域名收集  example
* 旁站入口收集  example
* 结合celery实现分布式扫描
* Host绑定设置  example
* 带Cookie扫描  example

## // TODO

* 支持不同Host的目标输入
* 支持从文件中读取请求作为输入
* 输出结果到消息队列

## Trouble Shooting

* 'Fetch.enable' wasn't found

  Fetch是新版chrome支持的功能，如果出现此错误，说明你的版本较低，请升级chrome到最新版即可。

## Follow me

如果你有关于动态爬虫的想法，欢迎和我交流。

微博：[@9ian1i](https://weibo.com/u/5242748339) Github: [@9ian1i](https://github.com/Qianlitp) 

相关文章：[漏扫动态爬虫实践](https://www.anquanke.com/post/id/178339)