package main

import (
	"bufio"
    "fmt"
	"log"
	"net"
	"strconv"
    "os"
    "encoding/gob"
    "image/jpeg"
    "image"
    "strings"
	"github.com/pkg/errors"
    "time"
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

func check(e error) {
    if e != nil {
        log.Println(e)
    }
}
func Open(addr string) (*bufio.ReadWriter, net.Conn, error) {
    log.Println("dial "+ addr)
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return nil, nil, errors.Wrap(err, "Dialing "+addr+" failed")
    }

    return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), conn, nil
}
func client(ip string) error {
    rw, conn, err := Open(ip + Port)
    if err != nil {
        return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
    }

    defer conn.Close()

    consoleTextReader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("-> ")
        text, _ := consoleTextReader.ReadString('\n')
        text = strings.Replace(text, "\r\n", "", -1)
        if strings.Compare(text, "close") == 0 {
            fmt.Println("shutting")
            return nil

        } else if strings.Compare(text, "rec") == 0 {
            fmt.Println("working")
            n, err := rw.WriteString("rec\n")
        	if err != nil {
        		fmt.Println(err, "Could not send the STRING ("+strconv.Itoa(n)+" bytes written)")
                continue
        	}

            log.Println("Flush the buffer.")
        	err = rw.Flush()
        	if err != nil {
        		fmt.Println(errors.Wrap(err, "Flush failed."))
                continue
        	}
            var data image.RGBA
        	dec := gob.NewDecoder(rw)
        	err = dec.Decode(&data)
        	if err != nil {
        		log.Println("Error decoding GOB data:", err)
        		continue
        	}
        	log.Printf("received the image")

            s := time.Now().Format(time.Stamp)
            s = strings.Replace(s, ":", "_", -1)
            log.Println(s)
            f, err := os.Create("C:\\Users\\ASUS\\Pictures\\gofiles\\sneek_"+s+".jpeg")
            check(err)
            defer f.Close()
            w := bufio.NewWriter(f)
            err = jpeg.Encode(w, data.SubImage(image.Rect(0, 0, 1920, 1080)), nil)
            check(err)
            w.Flush()
            fmt.Printf("wrote the jpeg\n")

        } else {
            fmt.Println("can't understand you")

        }
    }

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
