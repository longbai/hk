package main

import (
	"sort"
)

type LogArray []Log

func (p LogArray) Len() int           { return len(p) }
func (p LogArray) Less(i, j int) bool { return p[i].TotalBandwidth < p[j].TotalBandwidth }
func (p LogArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p LogArray) Sort()              { sort.Sort(p) }
