package main

import (
	"testing"
	"./packages/readfiles"
)

func benchmarkGAPM(bwt string, sra [][] string, psa map[int] int,
	               countDict map[int] map[byte] int,
	               firstOccurrence map[byte] int,
	               d int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetApproximatePatternMatchingWithCheckPointsC(bwt, sra, psa, countDict, firstOccurrence, d)
	}
}

func BenchmarkGetApproximatePatternMatchingWithCheckPointsC1(b *testing.B) {
	bwt := ReadFiles.ReadText("BWT_Ecoli.txt")
	sra := ReadFiles.ReadSra("Sra_Ecoli.fasta")
	sra = sra[0:17]
	psa1 := ReadFiles.ReadPSA("PSA_Ecoli1.json")
	cd1 := ReadFiles.ReadCountDict("Count_Ecoli1.json")
	firstColumn := GetFirstColumn(bwt)
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)
	d := 1
	benchmarkGAPM(bwt, sra, psa1, cd1, firstOccurrence, d, b)
}

func BenchmarkGetApproximatePatternMatchingWithCheckPointsC30(b *testing.B) {
	bwt := ReadFiles.ReadText("BWT_Ecoli.txt")
	sra := ReadFiles.ReadSra("Sra_Ecoli.fasta")
	sra = sra[0:17]
	psa30 := ReadFiles.ReadPSA("PSA_Ecoli30.json")
	cd30 := ReadFiles.ReadCountDict("Count_Ecoli30.json")
	firstColumn := GetFirstColumn(bwt)
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)
	d := 1
	benchmarkGAPM(bwt, sra, psa30, cd30, firstOccurrence, d, b)
}

func BenchmarkGetApproximatePatternMatchingWithCheckPointsC60(b *testing.B) {
		bwt := ReadFiles.ReadText("BWT_Ecoli.txt")
		sra := ReadFiles.ReadSra("Sra_Ecoli.fasta")
		sra = sra[0:17]
		psa60 := ReadFiles.ReadPSA("PSA_Ecoli60.json")
		cd60 := ReadFiles.ReadCountDict("Count_Ecoli60.json")
		firstColumn := GetFirstColumn(bwt)
		alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
		firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)
		d := 1
		benchmarkGAPM(bwt, sra, psa60, cd60, firstOccurrence, d, b)
}

func BenchmarkGetApproximatePatternMatchingWithCheckPointsC100(b *testing.B) {
	bwt := ReadFiles.ReadText("BWT_Ecoli.txt")
	sra := ReadFiles.ReadSra("Sra_Ecoli.fasta")
	sra = sra[0:17]
	psa100 := ReadFiles.ReadPSA("PSA_Ecoli100.json")
	cd100 := ReadFiles.ReadCountDict("Count_Ecoli100.json")
	firstColumn := GetFirstColumn(bwt)
	alphabet := []byte{'$', 'A', 'T', 'C', 'G'}
	firstOccurrence := GetFirstOccurrence(alphabet, firstColumn)
	d := 1
	benchmarkGAPM(bwt, sra, psa100, cd100, firstOccurrence, d, b)
}