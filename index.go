package main
import (
	"fmt"
	"./packages/texttobwt"
	"./packages/partialsuffixarrays"
	"./packages/checkpoint"
	"./packages/memusage"
	"./packages/readfiles"
	"runtime"
	"strconv"
	"time"
)

func main () {
	seq, title := ReadFiles.ReadRefseq("Refseq_Ecoli.fasta")
	fmt.Println(title)
	fmt.Println("")
	fmt.Println("Length of genome: ", len(seq))
	seq += "$"
	l := len(seq)
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}

	c := [] int {1, 30, 60, 100}
	k := [] int {1, 30, 60, 100}

	stBWT := time.Now()
	bwt := texttobwt.GetBWT(seq)
	fmt.Println("Runtime of BWT: ------", time.Since(stBWT), "------")
	fmt.Println("Memory of bwt: ")
	MemUsage.PrintMemUsage()
	ReadFiles.WriteText(bwt, "BWT_Ecoli.txt")

	bwt = ""
	runtime.GC()

	for i := 0; i < len(k); i++ {
		str := strconv.Itoa(k[i])
		stPSA := time.Now()
		psa := partialsuffixarrays.ConstructPartialSuffixArrayNoneConcurrency(seq, k[i])
		fmt.Println("Runtime of PSA ", str, "------", time.Since(stPSA), "------")
		fmt.Println("Memory of PSA ", str, ":")
		MemUsage.PrintMemUsage()
		path := "PSA_Ecoli" + str + ".json"
		ReadFiles.WritePSA(psa, path)
		psa = map[int]int{}
		runtime.GC()
	}

	bwt = ReadFiles.ReadText("BWT_Ecoli.txt")
	for i := 0; i < len(c); i++ {
		str := strconv.Itoa(c[i])
		stC := time.Now()
		countDict := CheckPoint.GetCheckPointArrays (alphabet, bwt, l, c[i])
		fmt.Println("Runtime of C ", str, "------", time.Since(stC), "------")
		fmt.Println("Memory of C ", str, ":")
		MemUsage.PrintMemUsage()
		path := "Count_Ecoli" + str + ".json"
		ReadFiles.WriteCountDict(countDict, path)
		countDict = map[int] map[byte] int {}
		runtime.GC()
	}
}

