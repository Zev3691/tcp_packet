package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

func main() {
	// 监听tcp连接
	listen, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 轮询接收信息
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 异步处理不会阻塞，把得到的连接放到处理函数
		go process(conn)
	}
}

func process(conn net.Conn) {
	// 当返回的时候关闭当前连接
	defer conn.Close()
	for {

		// 封包
		packet := func() error {
			// 从连接中读取数据
			reader := bufio.NewReader(conn)
			// 解码
			msg, err := decode(reader)
			if err == io.EOF {
				return errors.New("decode over")
			}
			if err != nil {
				fmt.Println("decode err ", err)
				return err
			}
			// 打印最后得到信息
			fmt.Println(msg)
			return nil
		}

		if err := packet(); err != nil {
			return
		}

		//未封包
		//noPacket := func() error {
		//	// 从连接中读取数据
		//	reader := bufio.NewReader(conn)
		//
		//	// 解码
		//	b := make([]byte, reader.Size())
		//	_, err := reader.Read(b)
		//	if err == io.EOF {
		//		return err
		//	}
		//	if err != nil {
		//		fmt.Println("decode err ", err)
		//		return err
		//	}
		//	// 打印最后得到信息
		//	fmt.Println(string(b))
		//
		//	return nil
		//}
		//if err := noPacket(); err != nil {
		//	return
		//}
	}
}

func decode(reader *bufio.Reader) (string, error) {
	// 为什么是前四个字节？
	// 四个字节是根据客户端而定的，客户端定义为4字节则服务端就去除前4个字节
	lenghtByte, _ := reader.Peek(4) // 读取包前4个字节,返回去除包前四位字节的长度
	lengthBuff := bytes.NewBuffer(lenghtByte)
	var length int16 // 二进制数据的实际长度，传入指针进去，会返回一个修改后的数据
	if err := binary.Read(lengthBuff, binary.LittleEndian, &length); err != nil {
		fmt.Println(err)
		return "", err
	}
	// 如果读取到的自己长度和计算的不一样，返回
	if int16(reader.Buffered()) < length+4 {
		return "", nil
	}
	// 真实长度加上包头的4字节
	realData := make([]byte, int(4+length))
	// 写入到realData
	_, err := reader.Read(realData)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	// 返回去掉包头的真实数据
	return string(realData[4:]), nil
}
