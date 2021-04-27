package main

import (
	"bufio"

	"log"
	"net"
	"strconv"



	"github.com/pkg/errors"

	"encoding/gob"
	"flag"
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


func Open(addr string) (*bufio.ReadWriter, net.Conn, error) {
    log.Println("dial "+ addr)
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, nil, errors.Wrap(err, "Dialing "+addr+" failed")
    }

    return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), conn, nil
}
func client(ip string) error {
    testStruct := complexData{
        N: 23,
        S: "some string",
        C: &complexData{
            N: 256,
            S: "thats a recursive string",
        },
    }
    rw, conn, err := Open(ip + Port)
    defer conn.Close()
    if err != nil {
        return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
    }
    log.Println("Send the string request.")
	n, err := rw.WriteString("STRING\n")
	if err != nil {
		return errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)")
	}
	n, err = rw.WriteString("Additional data.\n")
	if err != nil {
		return errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)")
	}
	log.Println("Flush the buffer.")
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
    log.Println("Read the reply.")
	response, err := rw.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'")
	}

	log.Println("STRING request: got a response:", response)
    log.Println("Send a struct as GOB:")
	log.Printf("Outer complexData struct: \n%#v\n", testStruct)
	log.Printf("Inner complexData struct: \n%#v\n", testStruct.C)
	enc := gob.NewEncoder(rw)
	n, err = rw.WriteString("GOB\n")
	if err != nil {
		return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	}
	err = enc.Encode(testStruct)
	if err != nil {
		return errors.Wrapf(err, "Encode failed for struct: %#v", testStruct)
	}
	err = rw.Flush()
	if err != nil {
		return errors.Wrap(err, "Flush failed.")
	}
	return nil

}

func main(){
    connect := flag.String("connect", "", "IP address of process to join. If empty, go into listen mode.")
	flag.Parse()

	if *connect != "" {
		err := client(*connect)
		if err != nil {
			log.Println("Error:", errors.WithStack(err))
		}
		log.Println("Client done.")
		return
	}
}

func init(){
    log.SetFlags(log.Lshortfile)
}
