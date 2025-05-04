package main

import (
	"encoding/binary"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"stzbHelper/model"
	"sync"
)

var isDebug bool = false

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

	fullbuf := []byte{}
	fullsize := 0

	// 循环读取数据包
	if isDebug == true {
		fmt.Printf("开始在接口 %s 上捕获 TCP 数据包（端口 8001）...\n", deviceName)
	}
	for packet := range packetSource.Packets() {
		// 解析 TCP 层
		if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
			// 获取 TCP 数据层（Payload）
			if appLayer := packet.ApplicationLayer(); appLayer != nil {
				payload := appLayer.Payload()
				if len(payload) < 4 {
					continue
				}
				var buf []byte
				if fullsize > 0 { // 需要组合数据包时
					if len(payload)+len(fullbuf) < fullsize {
						fullbuf = append(fullbuf, payload...)
						continue
					} else {
						buf = append(fullbuf, payload...)
						fullsize = 0
						fullbuf = []byte{}
					}
				} else {
					if len(payload) < int(binary.BigEndian.Uint32(payload[:4])) {
						fullbuf = append(fullbuf, payload...)
						fullsize = int(binary.BigEndian.Uint32(payload[:4]))
						continue
					}
					buf = payload
				}
				if isDebug == true {
					fmt.Println("")
					fmt.Println("====================================================")
					fmt.Println("")
				}
				bufread := NewBufferFrom(buf)
				cmdId2 := bufread.ReadInt()
				if isDebug == true {
					fmt.Println("字节", cmdId2)
				}
				cmdId := bufread.ReadInt()
				if isDebug == true {
					fmt.Println("协议号", cmdId)
				}
				if len(buf) > 14 {
					if isDebug == true {
						fmt.Println("数据类型", buf[12])
					}

					// 只处理类型3的数据
					if buf[12] == 3 {
						ParseData(cmdId, buf[17:])
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
