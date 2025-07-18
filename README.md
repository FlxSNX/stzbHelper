# stzbHelper
率土之滨攻城考勤助手  
本项目灵感来源于[团助手](http://stzbtool.cn/),它也是通过npcap抓取网络数据包来实现的功能  
目前计划只实现最基础的打城考勤,不会增加更多的功能.  如需其他功能可以参考本项目进行开发  
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
5. 如无法正常执行`build.bat`,请手动执行前后端build指令。先进入web目录执行`npm run build`,然后返回项目根目录执行`go build -tags="nomsgpack" -ldflags="-s -w" -o dist\stzbHelper-windows-amd64.exe stzbHelper`
## 开发说明
将`main.go`中的`isDebug`改为`true`再重新编译运行 就可以在打印中看到来自率土服务器类型为3的数据包  
你可以在打印中寻找一些数据包来开发其他功能  
找到需要的数据包后在`parse.go`中的`ParseData`方法通过协议号判断来解析处理你需要的数据包  
  
因为战报与同盟成员数据都是类型3的，所以本项目中只处理了类型3的数据包。   
如果你需要其他类型的数据包，我可以提供目前我所知道的数据包类型和解析方法  
数据包的类型就是数据包的第13位 也就是代码中的`buf[12]`  
类型3 zlib解压缩  
类型2 明文可以直接转为字符串  
类型5 异或解密  
