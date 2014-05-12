package main

import (
	"flag"
	"fmt"
	hum "github.com/dustin/go-humanize"
)

func main() {
	flag.Parse()
	args := flag.Args()
	logs := initTransferPoint()
	var amount float32
	for _, v := range args {
		recs, total, err := readLog(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		addToLog(logs, recs)
		amount += total
		fmt.Println("amount", v, hum.Comma(int64(total/8.0)), "MB")
	}
	fmt.Println("total amount", hum.Comma(int64(amount/8.0)), "MB")

	LogArray(logs).Sort()

	pos := int(float32(pointNumber) * 0.95)
	fmt.Println("95%", pos, logs[pos].printTime(), logs[pos])

	pos = pointNumber - 4
	fmt.Println("4th", pos, logs[pos].printTime(), logs[pos])

	pos = pointNumber - 1
	fmt.Println("max", pos, logs[pos].printTime(), logs[pos])
	return
}
