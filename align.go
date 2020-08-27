package main
import (
	"./packages/readfiles"
	"runtime"
	"strconv"
	"sync"
	"fmt"
)

// Read alignment.
func AlignRead (seq , r string, di, lc, nm int) map[string] [] int {
	l := len(r)
	val := map[string] [] int {}
	var str string
	var iR int
	for i := 0; i < l; i++ {
		id := lc + i

		if r[i] != seq[id] {
			str += strconv.Itoa(i - iR) + string(r[i])
			iR = i + 1
		}
	}
	str += strconv.Itoa(l - iR)
	val[str] = [] int {di, lc, nm}
	return val
}

// Sequence Alignment
func Align(seq string, sra [][]string, data map[string] [] int) map[string] map[string] [] int {

	l := len(data)
	N := len(sra)

	aln := map[string] map[string] [] int {}
	mapSra := map[string] [] string {}
	var keys [] string
	var wg sync.WaitGroup
	var mu sync.Mutex

	goroutines := runtime.NumCPU()
	lastGoroutine := goroutines - 1
	stride := l/goroutines
	wg.Add(goroutines)

	for k, _ := range data {
		keys = append(keys, k)
	}

	for i := 0; i < N; i++ {
		read := sra[i]
		mapSra[read[0]] = [] string {read[1], read[2]}
	}

	for g := 0; g < goroutines; g++ {
		go func(g int) {
			start := g * stride
			end := start + stride
			if g == lastGoroutine {
				end = l
			}
			for i := start; i < end; i++{

				name := keys[i]
				rd := mapSra[name]
				di := data[name][0]
				lc := data[name][1]
				nm := data[name][2]
				if di == 0 {
					r := rd[0]
					value := AlignRead (seq, r, di, lc, nm)
					mu.Lock()
					aln[name] = value
					mu.Unlock()
				} else {
					r := rd[1]
					value := AlignRead (seq, r, di, lc, nm)
					mu.Lock()
					aln[name] = value
					mu.Unlock()
				}
			}
			wg.Done()
		}(g)
	}
	wg.Wait()
	return aln
}

func main(){

	seq, title := ReadFiles.ReadRefseq("Ref_SARs_CoV_2.fa")
	fmt.Println(title)
	sra := ReadFiles.ReadSra("Sra_SARs_CoV_2.fasta")

	m0 := ReadFiles.ReadPos("m0.json")
	m1 := ReadFiles.ReadPos("m1.json")
	m2 := ReadFiles.ReadPos("m2.json")
	m3 := ReadFiles.ReadPos("m3.json")

	aln0 := Align (seq , sra, m0)
	aln1 := Align (seq , sra, m1)
	aln2 := Align (seq , sra, m2)
	aln3 := Align (seq , sra, m3)

	ReadFiles.WriteAln(aln0, "aln_0.json")
	ReadFiles.WriteAln(aln1, "aln_1.json")
	ReadFiles.WriteAln(aln2, "aln_2.json")
	ReadFiles.WriteAln(aln3, "aln_3.json")
}

