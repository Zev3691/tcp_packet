package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"reflect"
	"strconv"
)

func main() {
	// 开启一个tcp连接
	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for i := 0; i < 30; i++ {
		// 示例数据
		msg := fmt.Sprintf("{\"Id\":%s,\"Name\":\"golang\",\"Message\":\"message\"}", strconv.Itoa(i))

		// 封包
		packet := func() {
			data, err := encode(msg)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("data ", string(data))
			_, err = conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		packet()

		// 未封包
		//noPacket := func() {
		//	data, err := encode(msg)
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//	fmt.Println("data ", string(data))
		//	_, err = conn.Write(data)
		//	if err != nil {
		//		fmt.Println(err)
		//		return
		//	}
		//}
		//noPacket()

	}
	fmt.Println("send data over... ")
}

func encode(msg string) ([]byte, error) {
	var length = int32(len(msg)) // 此处定义头部字节长度
	fmt.Println(reflect.TypeOf(length))
	var nb = new(bytes.Buffer)

	// 如果第三参数data是数字类型，则设置长度
	if err := binary.Write(nb, binary.LittleEndian, length); err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 如果第三参数data非数字类型，则写入数据
	if err := binary.Write(nb, binary.LittleEndian, []byte(msg)); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return nb.Bytes(), nil
}
