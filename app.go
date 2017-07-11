package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"flag"

	"strconv"

	"strings"

	log "github.com/Sirupsen/logrus"
	qrcode "github.com/skip2/go-qrcode"
)

var (
	SERVER_PORT int
)

func qrHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("invalid form data"))
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("invalid form data")
		return
	}

	qrSize := 100
	size := r.FormValue("size")

	if len(size) > 0 && len(strings.TrimSpace(size)) > 0 {
		if strings.Contains(size, "x") {
			qrSizeArray := strings.Split(size, "x")
			qrSize, err = strconv.Atoi(qrSizeArray[0])
			if err != nil {
				w.Write([]byte("invalid size parameter"))
				w.WriteHeader(http.StatusBadRequest)
				log.Errorf("invalid size parameter")
				return
			}
		} else {
			qrSize, err = strconv.Atoi(strings.TrimSpace(size))
			if err != nil {
				w.Write([]byte("invalid size parameter"))
				w.WriteHeader(http.StatusBadRequest)
				log.Errorf("invalid size parameter")
				return
			}
		}
	}
	data := r.FormValue("data")
	log.Infof("size: %v, data %v", qrSize, data)
	png, err := qrcode.Encode(data, qrcode.Medium, qrSize)
	if err != nil {
		w.Write([]byte("invalid form data"))
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("invalid form data: " + err.Error())
		return
	}
	w.Write(png)
}
func main() {

	log.Infof("监听的端口为 %v", SERVER_PORT)
	s := &http.Server{
		Addr: ":" + strconv.FormatInt(int64(SERVER_PORT), 10),
	}
	handlerRegister()
	go func() {
		log.Infof("%s", s.ListenAndServe())
	}()
	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Infof("http server gracefully stopping...")
	// Stop the service gracefully.
	s.Shutdown(context.Background())
	log.Infof("http server gracefully shutdown done...")
}

func handlerRegister() {
	http.HandleFunc("/", qrHandler)
}

// 自动初始化
func init() {
	initLogger()
	initCommandLineFlag()
}

func initCommandLineFlag() {
	flag.IntVar(&SERVER_PORT, "port", 9999, "-port=<N>")
	flag.Parse()
}

func initLogger() {
	file, err := os.OpenFile("qr-api.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
}
