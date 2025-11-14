package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"strconv"
	"stzbHelper/global"
	"stzbHelper/model"
	"sync"
	"time"
)

var databaseSelected bool = false

func main() {
	// 获取所有网络接口
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal("无法获取网络接口列表:", err)
	}

	// 如果没有找到任何接口，退出
	if len(devices) == 0 {
		log.Fatal("未找到可用的网络接口")
	}

	if global.IsDebug == true {
		// 打印所有可用的网络接口
		fmt.Println("可用的网络接口:")
		for i, device := range devices {
			fmt.Printf("%d: %s (%s)\n", i+1, device.Name, device.Description)
		}
	}

	// 使用 WaitGroup 等待所有 Goroutine 完成
	var wg sync.WaitGroup

	//model.InitDB("database")
	go StartHttpService(&wg)
	wg.Add(1)
	// 遍历所有接口并启动 Goroutine 监听
	log.Println("stzbHelper开始运行!")
	log.Println("version:", global.Version)
	//log.Println("提示：0.0.3版本开始启动软件后需要进入游戏点击自己的主公簿进行激活软件。此改动是为了之后实现多数据库与绑定游戏连接IP信息，避免出现连接到多个8001端口导致的数据错乱")
	time.Sleep(100 * time.Millisecond)
	log.Println("等待打开主公簿激活软件...")
	log.Println("未打开主公簿激活软件前软件可能会出现报错！")

	for _, device := range devices {
		wg.Add(1)
		go captureTCPPackets(device.Name, &wg)
	}

	// 等待所有 Goroutine 完成
	wg.Wait()
}

// captureTCPPackets 监听指定接口的 TCP 数据包
func captureTCPPackets(deviceName string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 打开网络接口
	handle, err := pcap.OpenLive(deviceName, 65535, true, pcap.BlockForever)
	if err != nil {
		log.Printf("无法打开接口 %s: %v\n", deviceName, err)
		return
	}
	defer handle.Close()

	// 设置过滤器，只捕获端口为 8001 的 TCP 数据包
	filter := "tcp and src port 8001"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Printf("无法在接口 %s 上设置过滤器: %v\n", deviceName, err)
		return
	}
	// 创建数据包源
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// 循环读取数据包
	if global.IsDebug == true {
		fmt.Printf("开始在接口 %s 上捕获 TCP 数据包（端口 8001）...\n", deviceName)
	}
	for packet := range packetSource.Packets() {
		handlePacket(packet)
	}
}

var fullbuf = []byte{}
var fullsize = 0
var waitbuf = false

func handlePacket(packet gopacket.Packet) {
	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		if appLayer := packet.ApplicationLayer(); appLayer != nil {
			PSH := tcpLayer.(*layers.TCP).PSH
			payload := appLayer.Payload()
			if len(payload) < 8 {
				return
			}
			var srcIP string
			var dstIP string
			var srcProt int
			var dstProt int
			if ipLayer := packet.NetworkLayer(); ipLayer != nil {
				switch ip := ipLayer.(type) {
				case *layers.IPv4:
					srcProt = int(tcpLayer.(*layers.TCP).SrcPort)
					dstProt = int(tcpLayer.(*layers.TCP).DstPort)
					srcIP = ip.SrcIP.String() + ":" + strconv.Itoa(srcProt)
					dstIP = ip.DstIP.String() + ":" + strconv.Itoa(dstProt)
				case *layers.IPv6:
					srcProt = int(tcpLayer.(*layers.TCP).SrcPort)
					dstProt = int(tcpLayer.(*layers.TCP).DstPort)
					srcIP = ip.SrcIP.String() + ":" + strconv.Itoa(srcProt)
					dstIP = ip.DstIP.String() + ":" + strconv.Itoa(dstProt)
				}
			}

			if global.ExVar.BindIpInfo == true && global.OnlySrcIp != "" && global.OnlyDstIp != "" {
				if global.OnlySrcIp != srcIP || global.OnlyDstIp != dstIP {
					if global.IsDebug == true {
						fmt.Println("IP信息不符合跳过数据处理")
					}
					return
				}
			}

			var buf []byte
			if PSH != true {
				waitbuf = true
				fullbuf = append(fullbuf, payload...)
				return
			} else {
				if waitbuf == true {
					waitbuf = false
					buf = append(fullbuf, payload...)
					fullbuf = []byte{}
				} else {
					buf = payload
				}
			}

			if global.IsDebug == true {
				fmt.Println("")
				fmt.Println("====================================================")
				fmt.Println("")
			}
			bufread := NewBufferFrom(buf)
			bufsize := bufread.ReadInt()
			if global.IsDebug == true {
				fmt.Println("包大小", bufsize)
			}
			cmdId := bufread.ReadInt()
			if global.IsDebug == true {
				fmt.Println("协议号", cmdId)
			}

			if len(buf) > 14 {
				if global.IsDebug == true {
					fmt.Println("数据类型", buf[12])
				}

				if buf[12] == 3 {
					//fmt.Println(len(buf), bufsize, cmdId, "-----------")
					if len(buf)-bufsize != 4 && (cmdId == 103 || cmdId == 92) {
						global.LossCmdId = cmdId
						global.LossBytes = buf
						global.PacketLoss = true
						global.NeedBufSize = bufsize
					} else {
						go ParseData(cmdId, buf[17:])
					}

				} else if buf[12] == 5 {
					//println(buf)
					if global.IsDebug == true {
						data := DecodeType5(buf[12:])
						fmt.Println(data)
					}
				} else if buf[12] == 2 {

					//if cmdId == 5028 || cmdId == 5026 {
					//	fmt.Println(string(buf[12:]))
					//}
					//
					//if cmdId == 5028 {
					//	Parse5028(buf[13:])
					//}
				} else if cmdId > 99999 && global.PacketLoss == true && (global.LossCmdId == 103 || global.LossCmdId == 92) {
					result := make([]byte, len(buf)+len(global.LossBytes))
					copy(result, global.LossBytes)
					copy(result[len(global.LossBytes):], buf)
					if len(buf)+len(global.LossBytes)-global.NeedBufSize != 4 {
						global.LossBytes = result
					} else {
						global.PacketLoss = false
						go ParseData(global.LossCmdId, result[17:])
					}

				}

				if cmdId == 3686 && databaseSelected == false {
					var data []byte
					if buf[12] == 5 {
						data = []byte(DecodeType5(buf[12:]))
					} else if buf[12] == 3 {
						data = parseZlibData(buf[17:])
					}
					var raw []interface{}
					err := json.Unmarshal([]byte(data), &raw)
					if err != nil {
						log.Fatal(err)
					} else {
						dataMap := raw[1].(map[string]interface{})
						server, ok := dataMap["server"].([]interface{})
						if ok {
							log.Printf("服务器信息: %v\n", server)
						}

						var roleName string
						if logData, ok := dataMap["log"].(map[string]interface{}); ok {
							roleName = logData["role_name"].(string)
							log.Printf("角色名: %s\n", roleName)
						}

						log.Println("本地IP：" + dstIP)
						log.Println("游戏服务器IP：" + srcIP)
						global.OnlySrcIp = srcIP
						global.OnlyDstIp = dstIP
						dabesename := roleName + "_" + server[0].(string)
						log.Println("收到主公簿数据，将打开数据库文件" + dabesename + ".db")
						model.InitDB(dabesename)
						databaseSelected = true
					}
				}
			}

			if global.IsDebug == true {
				fmt.Print("[]byte{")
				for i, b := range buf {
					if i > 0 {
						fmt.Print(", ")
					}
					fmt.Print(b)
				}
				fmt.Println("}")
				fmt.Println("")
				fmt.Println("====================================================")
				fmt.Println("")
			}
		}
	}
}

type Buffer struct {
	Byte   []byte
	pos    int
	offset int
}

func (bb *Buffer) ResetOffset() {
	bb.offset = 0
}

func NewBufferFrom(b []byte) *Buffer {
	return &Buffer{Byte: b}
}

func (bb *Buffer) ReadInt() int {
	if bb.offset+4 > len(bb.Byte) {
		return 0
	}
	value := binary.BigEndian.Uint32(bb.Byte[bb.offset : bb.offset+4])
	bb.offset += 4
	return int(value)
}

func (bb *Buffer) ReadByte() byte {
	if bb.offset+1 > len(bb.Byte) {
		return 0
	}
	value := bb.Byte[bb.offset : bb.offset+1]
	bb.offset += 1
	return value[0]
}
