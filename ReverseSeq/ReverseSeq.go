package ReverseSeq

func ReverseSeq(seq string) string {
	l := len(seq) - 1
	var rs string
	for i := l; i >= 0; i-- {
		var r string
		sym := seq[i]
		if sym == byte('A') {
			r = "T"
		} else if sym == byte('C') {
			r = "G"
		} else if sym == byte('G') {
			r = "C"
		} else {
			r = "A"
		}

		rs += r
	}
	return rs
}
