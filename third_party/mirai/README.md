# Mirai 安装指南

## Step 1：下载依赖

[Download 下载](http://t.imlxy.net:64724/mirai/MiraiOK/miraiOK_windows_amd64.exe) 🎃MiraiOK 安装程序 ![Compile MiraiOK](https://github.com/LXY1226/miraiOK/workflows/Compile%20MiraiOK/badge.svg)

将下载的 `miraiOK_windows_amd64.exe` 文件保存至本文件夹 (`/third_party/mirai/`) 下

双击 `.exe` 文件运行，等待依赖下载安装

等待至 `请输入帐号密码` 提示后则安装完毕，关闭程序

## Step2：下载插件

[Download 下载](https://github.com/project-mirai/mirai-api-http/releases/download/v1.7.2/mirai-api-http-v1.7.2.jar) mirai-api-http-v1.7.2 插件

将下载的 `mirai-api-http-v1.7.2.jar` 文件保存至本文件夹 `plugins/` 下

再次运行 `.exe` ，等待加载完毕，并允许通过防火墙（如果有），关闭程序

## Step3：完成配置

进入 `plugins/MiraiAPIHTTP/` 文件夹，编辑文件夹下 `setting.yml` 文件

修改 `.yml` 文件内容为如下：

```yaml
port: 7958 # 7=> S 9=> Z 58=> KU
authKey: LOVELIVENIJIGASAKISCHOOLIDOL&NIJIGASAKIMIRAINIJI-PRODUCT # 爱与演唱会虹咲学园学园偶像&虹咲学园未来彩虹🌈-制品
enableWebsocket: true # 启用 WebSocket
```

再次运行 `.exe` ，查看输出信息，确认无误后即可运行主程序

🎉 Enjoy it !