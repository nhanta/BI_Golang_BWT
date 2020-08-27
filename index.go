package main

import (
	"./packages/checkpoint"
	"./packages/memusage"
	"./packages/readfiles"
	"./packages/texttobwt"
	"./packages/partialsuffixarrays"
	"fmt"
	"runtime"
	"strconv"
	"time"
)

// function, which takes a string as
// argument and return the reverse of string.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func main () {
	seq, title := ReadFiles.ReadRefseq("Ref_SARs_CoV_2.fa")
	rSeq := Reverse(seq)
	fmt.Println(title)
	fmt.Println("")
	fmt.Println("Length of genome: ", len(seq))
	seq += "$"
	l := len(seq)
	rSeq += "$"

	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}

	c := [] int {1, 30, 60, 100}
	k := [] int {1, 30, 60, 100}

	stBWT := time.Now()
	bwt := texttobwt.GetBWT(seq)
	fmt.Println("Runtime of BWT: ------", time.Since(stBWT), "------")
	fmt.Println("Memory of bwt: ")
	MemUsage.PrintMemUsage()
	ReadFiles.WriteText(bwt, "BWT_SARs_CoV_2.txt")

	rBWT := texttobwt.GetBWT(rSeq)
	ReadFiles.WriteText(rBWT, "rBWT_SARs_CoV_2.txt")

	bwt = ""
	rBWT = ""
	runtime.GC()

	for i := 0; i < len(k); i++ {
		str := strconv.Itoa(k[i])
		stPSA := time.Now()
		psa := partialsuffixarrays.ConstructPartialSuffixArrayNoneConcurrency(seq, k[i])
		fmt.Println("Runtime of PSA ", str, "------", time.Since(stPSA), "------")
		fmt.Println("Memory of PSA ", str, ":")
		MemUsage.PrintMemUsage()
		path := "PSA_SARs_CoV_2_" + str + ".json"
		ReadFiles.WritePSA(psa, path)
		psa = map[int]int{}
		runtime.GC()
	}

	bwt = ReadFiles.ReadText("BWT_SARs_CoV_2.txt")
	for i := 0; i < len(c); i++ {
		str := strconv.Itoa(c[i])
		stC := time.Now()
		countDict := CheckPoint.GetCheckPointArrays (alphabet, bwt, l, c[i])
		fmt.Println("Runtime of C ", str, "------", time.Since(stC), "------")
		fmt.Println("Memory of C ", str, ":")
		MemUsage.PrintMemUsage()
		path := "Count_SARs_CoV_2_" + str + ".json"
		ReadFiles.WriteCountDict(countDict, path)
		countDict = map[int] map[byte] int {}
		runtime.GC()
	}

	rBWT = ReadFiles.ReadText("rBWT_SARs_CoV_2.txt")
	for i := 0; i < len(c); i++ {
		str := strconv.Itoa(c[i])
		stC := time.Now()
		countDict := CheckPoint.GetCheckPointArrays (alphabet, rBWT, l, c[i])
		fmt.Println("Runtime of C ", str, "------", time.Since(stC), "------")
		fmt.Println("Memory of C ", str, ":")
		MemUsage.PrintMemUsage()
		path := "rCount_SARs_CoV_2_" + str + ".json"
		ReadFiles.WriteCountDict(countDict, path)
		countDict = map[int] map[byte] int {}
		runtime.GC()
	}
}

