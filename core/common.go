package core

import (
	"io"
	"log"
	"net"
	"sync"
	"unsafe"
)

type Header struct {
	Address   NetAddress
	Timestamp int64
}

func connCopy(source, target net.Conn, wg *sync.WaitGroup) {
	_, err := io.Copy(source, target)
	if err != nil {
		log.Println("Connection interrupted", err.Error())
	}
	_ = source.Close()
	wg.Done()
}

func forward(conn1, conn2 net.Conn) {
	log.Printf("Forward channel [%s/%s] <-> [%s/%s]\n",
		conn1.RemoteAddr(), conn1.LocalAddr(), conn2.RemoteAddr(), conn2.LocalAddr())

	var wg sync.WaitGroup
	// wait tow goroutines
	wg.Add(2)
	go connCopy(conn1, conn2, &wg)
	go connCopy(conn2, conn1, &wg)
	//blocking when the wg is locked
	wg.Wait()
}

func readHeader(conn net.Conn) NetAddress {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Fail to read local addresses", err.Error())
		return NetAddress{}
	}
	address := ParseNetAddress(string(buffer[:n]))
	log.Println("proxy address", address)
	return address
}

// 回写
func writeHeader(conn net.Conn, address NetAddress) bool {
	_, err := conn.Write([]byte(address.String()))
	if err != nil {
		log.Println("Fail to write response header")
		return false
	}
	return true
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func Serialize(data *Header) []byte {
	length := unsafe.Sizeof(*data)
	bytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(data)),
		cap:  int(length),
		len:  int(length),
	}
	return *(*[]byte)(unsafe.Pointer(bytes))
}

func Deserialize(bytes *[]byte) *Header {
	return *(**Header)(unsafe.Pointer(bytes))
}
