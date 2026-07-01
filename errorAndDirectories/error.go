package check

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// so if err := os.Writefile() or some such
// pass check(err)

