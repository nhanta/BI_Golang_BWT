package ReadFiles

import (
	"bufio"
	"os"
	"strings"
	"../reverseseq"
)

func ScanLines(path string) ([]string, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

// Read sequence runs analysis.
func ReadSra(path string) [][] string {

	var sra [][] string

	text, _ := ScanLines(path)
	n := len(text)
	for i := 0; i < n - 1; i += 2 {
		rText := ReverseSeq.ReverseSeq(text[i + 1])
		sra = append(sra, [] string {text[i], text[i + 1], rText})
	}
	return sra
}

// Read reference sequence analysis.
func ReadRefseq(path string) (string, string) {

	text, _ := ScanLines(path)
	title := text[0]
	seq := strings.Join(text[1:], "")

	return seq, title
}