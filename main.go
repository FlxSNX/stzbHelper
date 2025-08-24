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
)

var isDebug bool = false
var version string = "0.0.3"

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

	if isDebug == true {
		// 打印所有可用的网络接口
		fmt.Println("可用的网络接口:")
		for i, device := range devices {
			fmt.Printf("%d: %s (%s)\n", i+1, device.Name, device.Description)
		}
	}

	// 使用 WaitGroup 等待所有 Goroutine 完成
	var wg sync.WaitGroup

	model.InitDB()
	go StartHttpService(&wg)
	wg.Add(1)
	// 遍历所有接口并启动 Goroutine 监听
	log.Println("stzbHelper开始运行!")
	log.Println("version:", version)
	//log.Println("提示：0.0.3版本开始启动软件后需要进入游戏点击自己的主公簿进行激活软件。此改动是为了之后实现多数据库与绑定游戏连接IP信息，避免出现连接到多个8001端口导致的数据错乱")
	//log.Println("等待打开主公簿激活软件...")

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
	if isDebug == true {
		fmt.Printf("开始在接口 %s 上捕获 TCP 数据包（端口 8001）...\n", deviceName)
	}
	for packet := range packetSource.Packets() {
		handlePacket(packet)
	}
}

var fullbuf = []byte{}
var fullsize = 0
var waitbuf = false
var onlySrcIp = ""
var onlyDstIp = ""

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

			if global.ExVar.BindIpInfo == true && onlySrcIp != "" && onlyDstIp != "" {
				if onlySrcIp != srcIP || onlyDstIp != dstIP {
					if isDebug == true {
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

			if isDebug == true {
				fmt.Println("")
				fmt.Println("====================================================")
				fmt.Println("")
			}
			bufread := NewBufferFrom(buf)
			bufsize := bufread.ReadInt()
			if isDebug == true {
				fmt.Println("包大小", bufsize)
			}
			cmdId := bufread.ReadInt()
			if isDebug == true {
				fmt.Println("协议号", cmdId)
			}

			if len(buf) > 14 {
				if isDebug == true {
					fmt.Println("数据类型", buf[12])
				}

				if buf[12] == 3 {
					go ParseData(cmdId, buf[17:])
				} else if buf[12] == 5 {
					if cmdId == 3686 {
						data := DecodeType5(buf[12:])
						var raw []interface{}
						err := json.Unmarshal([]byte(data), &raw)
						if err != nil {
							log.Fatal(err)
						} else {
							dataMap := raw[1].(map[string]interface{})

							if server, ok := dataMap["server"].([]interface{}); ok {
								log.Printf("服务器信息: %v\n", server)
							}

							if logData, ok := dataMap["log"].(map[string]interface{}); ok {
								roleName := logData["role_name"].(string)
								log.Printf("角色名: %s\n", roleName)
							}

							log.Println("本地IP：" + dstIP)
							log.Println("游戏服务器IP：" + srcIP)
							//log.Println("软件将绑定以上IP进行数据过滤以避免数据错乱，如果当更换了账号或网络问题导致连接IP发送变化需要在软件网页上点击刷新IP信息然后重新打开主公簿进行绑定")
							onlySrcIp = srcIP
							onlyDstIp = dstIP
						}
					}
				}
			}

			if isDebug == true {
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
