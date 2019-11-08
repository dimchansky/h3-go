package h3

//lint:file-ignore U1000 Ignore all unused code

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}
