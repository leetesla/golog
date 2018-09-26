package mains

import (
	"strings"
	"fmt"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}


type LogProcess struct {
	readChan chan string
	writeChan chan string
	//path string//读取文件路径
	//influxDbSn string// influx data source
}


type ReadFromFile struct {
	path string
}

func (r *ReadFromFile)Read(rc chan string){
	line:= "message"
	rc <- line
}

//func(l *LogProcess) ReadFromFile(){
//	line:= "message"
//	l.readChan <- line
//}

type WriteToInfluxDB struct{
	influxDbSn string// influx data source
}

func (w *WriteToInfluxDB) Write(wc chan string){
	fmt.Println(<-wc)
}


func(l *LogProcess) Process(){
	data := <-l.readChan
	l.writeChan <- strings.ToUpper(data)
}

//func(l *LogProcess) WriteToInfluxDB(){
//	fmt.Println(<-l.writeChan)
//}


func main()  {
	rc:= make(chan string)
	wc:= make(chan string)

	lp := &LogProcess{
		readChan:rc,
		writeChan:wc,
		path:"/tmp/access.log",
		influxDbSn:"username&pwd",
	}

	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteToInfluxDB()

	time.Sleep(time.Millisecond * 1000)
}