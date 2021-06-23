package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)

//通过tail包，从文件中读取msg
func main() {
	fileName := `xx.log`
	//文件读取的权限
	config := tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll: true,
	}

	//打开文件
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var (
		msg *tail.Line
		ok bool
	)

	//读取文件的每一行
	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n",
				tails.Filename)
			//读取出错，等一秒在读
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
