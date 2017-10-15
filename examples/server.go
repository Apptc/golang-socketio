package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"


	"github.com/geneva-lake/golang-socketio"
	"github.com/geneva-lake/golang-socketio/transport"

)

type ChannelInner struct {
	Channel string `json:"channel"`
}

type MessageInner struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

var currentRoot string
var assetsDir http.FileSystem

func main() {
	currentRoot, _ := os.Getwd()
	assetsDir = http.Dir(currentRoot + "\\examples" + "\\assets")
	fmt.Println(assetsDir)

	server := gosocketio.NewServer(transport.GetDefaultPollingTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Connected")

		//c.Emit("/message", MessageInner{10, "main", "using emit"})

		//c.Join("test")
		//c.BroadcastTo("test", "/message", MessageInner{10, "main", "using broadcast"})
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	server.On("/join", func(c *gosocketio.Channel, channel ChannelInner) string {
		time.Sleep(2 * time.Second)
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	//fs := http.FileServer(http.Dir("assets"))
	//fmt.Println(fs)
	serveMux.HandleFunc("/", AssetsFileHandler)

	log.Println("Starting server...")
	log.Panic(http.ListenAndServe(":3811", serveMux))
}

func AssetsFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return
	}
	var file string = r.URL.Path
	f, err := assetsDir.Open(file)
	if err != nil {
		log.Println("can not open file ", file, " ", err)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatalln("can not open file ", file, " ", err)
	}
	http.ServeContent(w, r, file, fi.ModTime(), f)
}
