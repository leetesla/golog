package main

import (
	"strings"
	"fmt"
	"time"
	"os"
	"bufio"
	"io"
	"regexp"
	"log"
	"strconv"
	"net/url"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan *Message)
}


type LogProcess struct {
	readChan chan []byte
	writeChan chan *Message
	//path string//读取文件路径
	//influxDbSn string// influx data source
	reader Reader
	writer Writer
}


type ReadFromFile struct {
	path string
}

func (r *ReadFromFile)Read(rc chan []byte){
	// open file
	f,err := os.Open(r.path)

	if err != nil{
		panic(fmt.Sprintf("open file error:%s", err.Error()))
	}

	f.Seek(0,2)

	// from the end of file
	rd := bufio.NewReader(f)

	for{
		line, err:= rd.ReadBytes('\n')

		if err == io.EOF{
			time.Sleep(500 * time.Millisecond)
			continue
		}else if err != nil{
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}

		// line:= "message"
		rc <- line[:len(line)-1]
	}


}

//func(l *LogProcess) ReadFromFile(){
//	line:= "message"
//	l.readChan <- line
//}

type WriteToInfluxDB struct{
	influxDbSn string// influx data source
}

func (w *WriteToInfluxDB) Write(wc chan *Message){
	for v:= range wc{
		fmt.Println(v)
	}

}

type Message struct {
	TimeLocal time.Time
	BytesSent int
	Path,Method,Scheme,Status string
	UpstreamTime, RequetstTime float64
}

func(l *LogProcess) Process(){
	// 127.0.0.1 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854
	// echo '127.0.0.1 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 2133 "-" "KeepAliveClient" "-" 1.005 1.854' >> access.log

	// data := <-l.readChan
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"\s+\"([\d\.-]+)\"\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v:= range l.readChan{
		// l.writeChan <- strings.ToUpper(string(v))
		ret := r.FindStringSubmatch(string(v))

		fmt.Println(ret, len(ret))

		if len(ret) != 13 {
			log.Println("FindStringSubmatch fail", string(v))
		}

		message := &Message{}

		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0000", ret[4], loc)
		if err != nil{
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
		}

		message.TimeLocal = t

		bytesSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = bytesSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings.Split fail:", ret[6])
			continue
		}

		message.Method = reqSli[0]

		u, err := url.Parse(reqSli[1])

		if err != nil{
			log.Println("url parse fail:", err)
			continue
		}

		message.Path = u.Path
		message.Scheme = ret[5]
		message.Status = ret[7]

		// fmt.Println(reqSli, len(reqSli))

		upstreamTime, _ := strconv.ParseFloat(ret[11], 64)
		requestTime, _ := strconv.ParseFloat(ret[12], 64)

		message.UpstreamTime = upstreamTime
		message.RequetstTime = requestTime

		l.writeChan <- message

	}

}

//func(l *LogProcess) WriteToInfluxDB(){
//	fmt.Println(<-l.writeChan)
//}


func main()  {
	rc:= make(chan []byte)
	wc:= make(chan *Message)

	reader := &ReadFromFile{
		path:"./access.log",
	}

	writer := &WriteToInfluxDB{
		influxDbSn:"username&pwd",
	}

	lp := &LogProcess{
		readChan:rc,
		writeChan:wc,
		//path:"/tmp/access.log",
		//influxDbSn:"username&pwd",
		reader:reader,
		writer:writer,
	}

	go lp.reader.Read(lp.readChan)
	go lp.Process()
	go lp.writer.Write(lp.writeChan)

	time.Sleep(time.Millisecond * 1000 * 30)
}