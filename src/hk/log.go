package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

type Log struct {
	Time           int64
	TotalBandwidth float32
	Bandwidths     []float32
}

var zone *time.Location = time.FixedZone("CST", 8*3600)
var start time.Time = time.Date(2014, 4, 1, 0, 0, 0, 0, zone)
var startUnix int64 = start.Unix()
var end time.Time = time.Date(2014, 5, 1, 0, 0, 0, 1, zone)
var endUnix int64 = end.Unix()

const interval = 5 * 60

const pointNumber = 288 * 30

func (l *Log) printTime() string {
	t := time.Unix(l.Time, 0)
	return t.In(zone).String()
}

func initTransferPoint() []Log {
	recs := make([]Log, pointNumber)
	fromUtc := startUnix + interval
	for i := 0; i < pointNumber; i++ {
		recs[i].Time = fromUtc + int64(i*interval)
	}
	return recs
}

func parseLog(line string) (utc int64, bw float32, err error) {
	rec := strings.Split(line, " ")
	timeStr := rec[0] + " " + rec[1]
	// 2014-03-06 14:37:11
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, zone)
	if err != nil {
		return
	}
	_bw, err := strconv.ParseFloat(rec[3], 32)
	if err != nil {
		return
	}
	bw = float32(_bw)
	utc = t.Unix()
	return
}

func valid(utc int64) bool {
	return utc >= startUnix && utc < endUnix
}

func addToLog(logs []Log, recs []float32) {
	pos := 0
	for i := 0; i < len(recs); i += interval {
		var bw float32
		for j := 0; j < interval; j++ {
			bw += recs[i+j]
		}
		avgBw := bw / float32(interval)
		logs[pos].TotalBandwidth += avgBw
		logs[pos].Bandwidths = append(logs[pos].Bandwidths, avgBw)
		pos++
	}
}

func readLog(path string) (recs []float32, total float32, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0700)
	if err != nil {
		return
	}
	recs = make([]float32, 3600*24*30)
	scan := bufio.NewScanner(file)
	preUtc := startUnix
	var prebw float32
	for scan.Scan() {
		line := scan.Text()
		utc, bw, err1 := parseLog(line)
		if err1 != nil {
			err = err1
			return
		}
		if !valid(utc) {
			continue
		}
		recs[utc-startUnix] = bw
		//fill hole
		total += bw
		for utc-preUtc > 1 {
			preUtc++
			recs[preUtc-startUnix] = (bw + prebw) / 2.0
			total += (bw + prebw) / 2.0
		}
		prebw = bw
		preUtc = utc
	}
	return
}
