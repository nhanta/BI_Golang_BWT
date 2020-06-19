package main

import(
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"./packages/partialsuffixarrays"
	"./packages/texttobwt"

)

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Len() int{
	return len(s)
}
func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func getFirstColumn (bwt string) string {
	s := [] rune(bwt)
	sort.Sort(sortRunes(s))
	return string(s)
}

func getFirstOccurence(alphabet [] byte, firstColumn string) map[byte] int {

	m := map[byte] int {}
	for i:= range alphabet {
		id := bytes.IndexByte([] byte(firstColumn),  byte (alphabet[i]))
		m[alphabet[i]] = id
	}
	return m
}

func getCountSymbol(bwt string, symbol string, index int) int {
	return bytes.Count([] byte(bwt[0:index]), [] byte(symbol))
}

func makeRange(min, max, k int) [] byte {
	x := make([] byte, (max - min)/k + 1)
	for i := range x {
		x[i] = byte(min + i * k)
	}
	return x
}

func getCheckPointArrays (alphabet [] byte, bwt string, l int, c int) ([] byte, map[int] map[byte] int) {

	indexCheckPoint := makeRange(0, l, c)
	countDict := map[int] map[byte] int {}

	for i := range indexCheckPoint {
		counter := map[byte] int {}
		id := indexCheckPoint[i]
		instanceColumn := bwt[0:id]

		for j := range alphabet {
			key := alphabet[j]
			value := getCountSymbol(instanceColumn, string(key), int(id))
			counter[key] = value
		}
		countDict[int(id)] = counter
	}
	return indexCheckPoint, countDict
}

func uniqueGene(alphabet [] byte, str string) [] byte{
	var unique [] byte
	for i := range alphabet {
		if bytes.IndexByte([]byte(str), byte(alphabet[i])) != -1 {
			unique = append(unique, alphabet[i])
		}
	}
	return unique
}

func insertArr(arr [][] int, element [] int, index int) [][] int{
	arr = append(arr, [] int {})
	copy(arr[index + 1 :], arr[index :])
	arr[index] = element
	return arr
}

func getSubAproximateStringWithCheckPoints(bwt string, alphabet [] byte,
	                                       pattern string, d int, initialtop, initialbottom int,
	                                       partialSuffixArray map [int] int,
	                                       suffixArrayIndex [] byte,
	                                       indexCheckPoint [] byte, countDict map[int] map[byte] int,
	                                       firstOccurence map[byte] int, startPositions [] int) [] int{

	var topbottomUpdate [][] int

	topbottomUpdate = [] [] int {[] int {initialtop, initialbottom, 0, 0}, [] int {1, 1, 1, -1}}

	for reflect.DeepEqual(topbottomUpdate, [] [] int {[] int {1, 1, 1, -1}}) != true {
		tbt := topbottomUpdate[0]
		top := tbt[0]
		bottom := tbt[1]
		initialMismatch := tbt[2]
		label := tbt[3]
		newLabel := topbottomUpdate[1][3]

		aprSymbolAlphabet := uniqueGene(alphabet, bwt[top : bottom + 1])

		if pattern != "" {

			topbottomUpdate = topbottomUpdate [1:]
			lastId := len(pattern) - 1
			symbol := pattern[lastId]
			if label != newLabel {
				pattern = pattern[0 : lastId]}

			for i := range aprSymbolAlphabet {

				aprSymbol := aprSymbolAlphabet[i]
				count := label + 1
				mismatch := initialMismatch
				newTop := top
				newBottom := bottom

					if aprSymbol != symbol {
						mismatch += 1
					}

					if mismatch <= d {

						if bytes.IndexByte(indexCheckPoint, byte(newTop)) == -1 {
							newTop = firstOccurence[aprSymbol] + getCountSymbol(bwt, string(aprSymbol), newTop)
						} else {
							newTop = firstOccurence[aprSymbol] + countDict[newTop][aprSymbol]
						}

						if bytes.IndexByte(indexCheckPoint, byte(newBottom + 1)) == -1 {
							newBottom = firstOccurence[aprSymbol] + getCountSymbol(bwt, string(aprSymbol), newBottom + 1) - 1
						} else {
							newBottom = firstOccurence[aprSymbol] + countDict[newBottom + 1][aprSymbol] - 1
						}

						topbottomUpdate = insertArr(topbottomUpdate, [] int {newTop, newBottom, mismatch, count}, len(topbottomUpdate) - 1)
						fmt.Println(topbottomUpdate)
					}
			}
		} else {
			newtopbottomUpdate := topbottomUpdate[0 : len(topbottomUpdate) - 1]
			fmt.Println("new top bottom update: ", newtopbottomUpdate)
			for i := range newtopbottomUpdate {
				x := newtopbottomUpdate[i]
				top = x[0]
				bottom = x[1]

				for j := top; j < bottom + 1; j++ {
					ct := 0
					lastIndex := j
					for bytes.IndexByte(suffixArrayIndex, byte(lastIndex)) == -1 {
						ct += 1
						symbol := bwt[lastIndex]
						lastIndex = firstOccurence[symbol] + getCountSymbol(bwt,string(symbol), lastIndex)
					}
					startPositions = append(startPositions, partialSuffixArray[lastIndex] + ct)
				}
			}
			return startPositions
		}
	}
	return [] int {}
}

func main(){
	//bwt := "TCAAAAAAACTCCGCTGTACGGCTATCACCTATTGGTTTTCGGTGTTGTTCGGTTTGGCG$GGTCTGCTCCGCGTCCCGTCCTGCAGCCGTGTCGGATTTCTTGCTTTCGCTGCCATCCCCTCTTATTATGTAGGTAATGATCATAGCGCCACCCGGGGTTTACCTGACTCGCCTCCGATCGACTAGTGCGCTATCCTACCGTGTATTCCACATCATTCTCGCGCGCCCTCGTTTGCATCTAAACTCATGCAAATAGTTATCATAGGGAAGTACGTCCCTCCACCGACTGCCCGTGAGCGCCAAACTCCTGACTCATTTGGTTAGATCTTTTTTTTCACTGAATGTGGTATAACTTGTGATGGGGGCCACCCGGCTTCTTTTTACCGTTTCTGGGGAGCCTTAGTATTAAAAAAAAATAATAAAAAAAGGTAAGGTTAGTACGGTCAAAGTTTAAACATACATTATTTCGCCTAACATTCATAGGGATTTACTATTTTTTGTTACTTTTAGGCCCAGTGCTGATACGTCCGCCTTGGGCAAACGTGGCCTCGTTTTCCGTTTTTGCTTAGCAGGCTGGTCGAAAACAAACCTAGGATGGGGATCTACTAGGACAGTTACGGCCCCCGCTGTTTCAACGCTATGACGCGCCTGCCTAAGCTACTATTGCCCCCATCACGGTCCGTCCAGGATAAGTTCCCAGCCCCGTACTTTTGGCGCCCTATCATGTCTTTTTGGCCGGCTCCGGACCAGGAAGAGTGGGGCCGTGCCAGGTCCAGCCTGGTGGTTACCGCCCGCGATCCTAACGACGTTACCGGGGAACATCCGGGCGTAGTGTTGTTAGACTCGCAAATACATTGCCTACGAGTGACAAAAAAGAATGTGCGTGCCAGTGGTTAACAGACGCGTATCCCCAGGCAGCCCGTATCGTGCGGGCGACCCGGCCAGCGCGGTTATTGATTATCAAACCCGGACTGGGGCGCTGCTAGAGTTGCTTTCGACTCGTCCTACCTTGATAGTCTTACGGGGGCTCTCCGGTTGAAGCCCGGAGCTCTCAAACGCACTCGGACTTTTAATCCTTTCCATCACCCCTGCCTTTCAGGTCATCTCTCTTTCTACCGCCACAACTTAATCCCTGGGATCACCGATCCTTTGAGCCGATTTTTTAAAAGTGCTGGTGGCTTTCTCTTGAATCATTTCGAAAAAATGCTCTTGCACCTTTCCTCGGCCGGCTTATGCTCATCTTGCCTCTTGTCTCTGTCTTACCTAACCAGGTGTCAACTGCTTCCTGCGTCCTGAAAAACCGGCGTTTCTTGGCAGATCACGTCGTTTTCACCCCAGCACATAAAGAGAACCAGCCGCCTTTACCAGGGTTCCGACCTCTCCAGCCTTCCCTTGGTTCAGGTTGGCTCAATTGCGCTCGAACGGATCATCGCAGGTGATGTTGTGCATTGTGGCCCCCCTCTGAAATTCCCTCTTTTATGCGTGTACGAGGGGTCGCTATGTTATTGAGTAAGTAAGCAGTACATAATCACATAACCAGTGCAAGACCCCCGGGACGACGGAGTCTGTGTGTAAGGGCCTCCCGCGCGGGATGGGGACGTCCTGACGTTGGTGTGATCGAATGCCCCACTTTTTGCCCCACCCGTGCACCCTTTTGGAGCGGTCTGAAGGTCTCAACTTAATCCTGTGCACGCGTTTAACTAGAATGCCGTATTCCGGGAATCGGGCTCCCGTTTTAGGCCGCCTTTTTCCCACGCGTGTCGCTTTTCACCCTTCGACGTATCTTTCCCCTGCCCTCTATGACGGGATCTAGCGTATTGAGTCTGCGCAACGATTCTTGTCCCACACCACGAATCTGTTAGCAACGTCCCGGTGTGTGCAAGTCCCGGTTACGCGGGGTCTGGGACTTCCTGCATTGCAGTCTGCGTGAGTTACCGTGCTTGTCCCGGAGCATAGCAGACTGGAGTCTGAGCGTAGTGGAGGGACACCGCCTCGTGTGATCTTGCAGCCCTTGCGCCTTAGTTAGGCTCGGCTTACAATCTCCCTCTATTCGTCCCTCCCCGCCCTCCATCGCTTAGTTCTGTTTTCGTTTCCTCGCCCCGTCGGTCGCTGACTTTCTCGTAAAAAATCATAGGCACTGAATGATCATCGATTCTCCAGCAAGGCCCGTACATAACTTGACTGAAGCCTCAGTAAGCGTCCCTACTACGCCTGCCGGGCACATTCCGCATCCGAAGGACACTTCATATAATAGCCCCCCGGCCAGCAGGTTTCCGAGGGCTGATCGGTCCACTCGGCCCGACGAGGTCCGAGTTTAGCCACATATAGTGAGAGCCGGGCCAAGGAAAATGTATCTGAGGCGAGTGTGCAGGCACTGAATGGGATACTGTTCACATTACCCAGTCCGAGGGTCTAGGGGGGGTACCCTTAGACGATCCACCATCCTCTCTTGATGGCCGCGGGCGCGTTGTCGAACCGTGTAGATTCTTCGGCCGGGGGGGGGGCTTTGTAGGGCTCCATGAGAGGGCGTTTGAGATCGACCTGCGGCTGCTCGAAAACCGAAGAGTATGGGGCCGAAGAACACAACAGTCAGCCCATATCTCGGTGCTTAAAACCTGCATCATACTCGGGGGGAGGCACGTCCCCCTGCTTGTCGGGATGCGTGTGCCTTGACCCTTTGTGAACCTCTTTTGACCATCCGGGCGGTCAAACGCCCTTTTCATATATTTTTCTAGACCAATCGCTCGGGGGTTTCCCGAGATCCGCTCGAGCGGCTTTCTTGCGTCCGTACATATGGCCCGAATTTCATGGTGGGTGGTCCCGGCGTGGGCGTACGGAAACCGTCCCGGGGTGGTTGCTCCATCGGCGGGATAGGGCTTGAGTTATAAAATATCAGTTCACCAGGCCTAGGTACGGGCGCTCGCGAGACCGAAAAGAACGGCAGACGATTAGCAACCCACAAAGCGAGGCTGCCCTTAGGGAACGCACCGCGGGGAAACAGAAGATGAATCTACTACACCAGGACCCTACAATACCAGCGGGGATCCACTGGAACCAATCCCGCTCTCTGACGTCGTAAAGACGCACAGACGAGGCAGCGCCTCTATGACATGTACGTCAGGATTAGTTTTTTTTTTGTG"
	//bwt := "ATTC$AA" // BATATA$
	text := "ACATGCTACTTT$"
	pattern := "ATT"
	bwt := texttobwt.GetBWT(text)
	l := len(bwt)
	d := 1
	initialtop := 0
	initialbottom := len(bwt) - 1
	fmt.Println(bwt)
	alphabet := [] byte {'$', 'A', 'T', 'C', 'G'}
	firstColumn := getFirstColumn(bwt)
	// fmt.Println(firstColumn)
	firstOccurence := getFirstOccurence(alphabet, firstColumn)
	// fmt.Println(firstOccurence)
	// fmt.Println(getCountSymbol(bwt, "A", 6))
	// fmt.Println(makeRange(2, 23, 3))

	indexCheckPoint, countDict := getCheckPointArrays (alphabet, bwt, l, 3)
	//fmt.Println("index check point: ", indexCheckPoint)
	//fmt.Println(countDict)
	// fmt.Println(uniqueGene(alphabet, bwt))

	psa := partialsuffixarrays.ConstructPartialSuffixArrayNoneConcurrency(text, 2)
	//fmt.Println(psa)
	var sai [] byte
	for k := range psa {
		sai = append(sai, byte(k))
	}
	// fmt.Println("Suffix Array Index: ", sai)
	positions := getSubAproximateStringWithCheckPoints(bwt, alphabet, pattern, d, initialtop, initialbottom,
		                                               psa, sai, indexCheckPoint, countDict,
		                                               firstOccurence , [] int {})
	fmt.Println(positions)
}
