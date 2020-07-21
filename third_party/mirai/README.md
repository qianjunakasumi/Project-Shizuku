<p align="center">
    <img width="160" src="https://s1.ax1x.com/2020/07/05/Up59PK.png" alt="SHIZUKU logo">
</p>
<p align="center">
    <img src="https://img.shields.io/badge/%E4%BA%A7%E5%93%81%E5%90%8D%E7%A7%B0-SHIZUKU%F0%9F%92%A7-brightgreen?style=flat-square" alt="产品名称">
    <img src="https://img.shields.io/badge/%E5%BC%80%E5%8F%91%E5%9B%A2%E9%98%9F-%E8%99%B9%E5%92%B2%E5%AD%A6%E5%9B%AD%E6%9C%AA%E6%9D%A5%E5%BD%A9%E8%99%B9%F0%9F%8C%88%E5%90%8C%E5%A5%BD%E4%BC%9A-brightgreen?style=flat-square" alt="开发团队">
    <a href="https://goreportcard.com/report/github.com/qianjunakasumi/shizuku"><img src="https://goreportcard.com/badge/github.com/qianjunakasumi/shizuku?style=flat-square" alt="Go Report Card"></a>
    <img src="https://img.shields.io/github/repo-size/qianjunakasumi/shizuku?style=flat-square" alt="repo size">
</p>

<h2 align="center">虹咲学园 SHIZUKU 未来彩虹🌈</h2>

SHIZUKU 是 虹咲学园未来彩虹🌈同好会 下的一个 AGPL许可 开源项目，项目是基于 [Mirai](https://github.com/mamoe/mirai) 提供的 HTTP API 接口开发 及 [腾讯 QQ](https://im.qq.com/) 群组中使用，以偶像为代表的辅助性功能机器人

本项目命名于 [ラブライブ！虹ヶ咲学園スクールアイドル同好会](http://lovelive-anime.jp/nijigasaki/) 中 [桜坂しずく（Osaka Shizuku）](https://lovelive-as.bushimo.jp/member/shizuku/) 的名称，项目中任何涉及包括但不限于偶像姓名、肖像、自我介绍等，其著作权均归 ©プロジェクトラブライブ！虹ヶ咲学園スクールアイドル同好会 所有

## SHIZUKU 第三方应用之 Mirai 安装指南

### Step 1：下载应用

[点击下载 🎃MiraiOK](http://t.imlxy.net:64724/mirai/MiraiOK/miraiOK_windows_amd64.exe)
![Compile MiraiOK](https://github.com/LXY1226/miraiOK/workflows/Compile%20MiraiOK/badge.svg)

将下载的 `miraiOK_windows_amd64.exe` 文件保存至本指南所在文件夹 (`/third_party/mirai/`) 下

双击 `.exe` 文件运行，等待所需环境下载安装完成

当出现 `请输入帐号密码` 提示后则代表安装完毕，请关闭程序

### Step 2：下载插件

[点击下载 mirai-api-http-v1.7.3.jar](https://github.com/project-mirai/mirai-api-http/releases/download/v1.7.2/mirai-api-http-v1.7.2.jar)

将下载的 `mirai-api-http-v1.7.3.jar` 文件保存至本指南所在文件夹的 `plugins/` 下

运行 `.exe` ，允许通过防火墙（若有），等待加载完毕后关闭程序

### Step 3：完成配置

进入 `plugins/MiraiAPIHTTP/` 文件夹，编辑文件夹下 `setting.yml` 文件

修改 `.yml` 文件内容为如下：

```yaml
port: 7958 # 7=> S 9=> Z 58=> KU
authKey: LOVELIVENIJIGASAKISCHOOLIDOL&NIJIGASAKIMIRAINIJI-PRODUCT # 爱与演唱会虹咲学园学园偶像&虹咲学园未来彩虹🌈-制品
enableWebsocket: true # 启用 WebSocket
```

再次运行 `.exe` ，检查输出信息是否与配置相符，确认无误后即可运行SHIZUKU主程序

🎉 Enjoy it !