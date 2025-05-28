# stzbHelper
率土之滨攻城考勤助手  
本项目灵感来源于[团助手](http://stzbtool.cn/),它也是通过npcap抓取网络数据包来实现的功能  
目前计划只实现最基础的打城考勤,不会增加更多的功能.  如需其他功能可以参考本项目进行开发  
已差不多可以使用 [演示视频](https://www.bilibili.com/video/BV1ABVqzFERV)
## 使用说明
本程序依赖于 [Npcap](https://npcap.com/#download) 抓取网络数据包来获取战报与同盟成员信息, 所以你在使用前需要先安装Npcap(https://npcap.com/dist/npcap-1.82.exe)  
## 支持情况
- PC客户端
- 模拟器移动端客户端
- 移动端客户端（使用移动端设备时需运行本程序的主机带有无线网卡，并打开移动热点给移动端设备连接）
## 功能
- 攻城任务考勤（统计目标成员的主力，拆迁队伍数量和攻城次数)
## 构建
1. 构建前需确保已安装 golang 1.24及以上版本、nodejs  
2. 执行 `go mod tidy`或者`go mod download` 安装依赖
3. 到项目的web目录下 执行 `npm install` 安装依赖
4. 执行 `build.bat`
