package main

import (
	"bufio"
	"io"
	"log"
	"net"

	"strings"
	"sync"

	"github.com/pkg/errors"

	"encoding/gob"

)
type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}
const (
    Port = ":6100"
)
type HandleFunc func(*bufio.ReadWriter)

type EndPoint struct {
    listener net.Listener
    handler map[string]HandleFunc

    m sync.RWMutex
}
func NewEndPoint() *EndPoint {
    return &EndPoint{
        handler: map[string]HandleFunc{},
    }

}


func (e *EndPoint) AddHandle (name string, f HandleFunc) {
    e.m.Lock()
    e.handler[name] = f
    e.m.Unlock()
}

func (e *EndPoint) Listen() error {
    var err error
    e.listener, err = net.Listen("tcp", Port)
    if err!= nil {
        return errors.Wrapf(err, "Unableto listen on %s", Port)
    }
    log.Println("Listen on", e.listener.Addr().String())
    for {
        log.Println("Accept an incoming connection request")
        conn, err := e.listener.Accept()
        if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
        log.Println("Handle incoming messages")
        go e.handleMessages(conn)
    }
}

func (e *EndPoint) handleMessages (conn net.Conn) {
    rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
    defer conn.Close()
    for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
        cmd = strings.Trim(cmd, "\n")
        log.Println(cmd+"'")

        e.m.RLock()
        handleCommand, ok := e.handler[cmd]
        e.m.RUnlock()
        if !ok {
            log.Println("Command "+cmd+ " is not registered")
            return
        }
        handleCommand(rw)
    }

}
func handleStrings(rw *bufio.ReadWriter){
    log.Print("handling STRING ")
    s, err := rw.ReadString('\n')
    if err != nil {
        log.Println("cannot read from connection, ", err)
    }
    s = strings.Trim(s, "\n")
    log.Println(s)
    _, err = rw.WriteString("thank you\n")
    if err != nil {
        log.Println("ERror writing in reply")
    }
    err = rw.Flush()
    if err != nil {
        log.Println("Flush failed ", err)
    }


}
func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data:")
	var data complexData
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	log.Printf("Outer complexData struct: \n%#v\n", data)
	log.Printf("Inner complexData struct: \n%#v\n", data.C)
}


func server() error {
    endp := NewEndPoint()
    endp.AddHandle("STRING", handleStrings)
    endp.AddHandle("GOB", handleGob)
    return endp.Listen()
}


func main(){
    err := server()
	if err != nil {
		log.Println("Error:", errors.WithStack(err))
	}

	log.Println("Server done.")
}
func init() {
	log.SetFlags(log.Lshortfile)
}
