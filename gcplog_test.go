package gcplog

func ExamplePrint() {
	logger := New()
	logger.Print("Hello World")
	// Output:
	// {"severity":"DEFAULT","message":"Hello World"}
}

func ExamplePrintf() {
	logger := New()
	logger.Printf("%s %v", "Hello World", 12345)
	// Output:
	// {"severity":"DEFAULT","message":"Hello World 12345"}
}

func ExamplePrefixPrint() {
	logger := New()
	logger.PrefixPrint("Hello World")
	// Output:
	// {"severity":"DEFAULT","message":"DEFAULT: Hello World"}
}

func ExamplePrefixPrintf() {
	logger := New()
	logger.PrefixPrintf("%s %v", "Hello World", 12345)
	// Output:
	// {"severity":"DEFAULT","message":"DEFAULT: Hello World 12345"}
}

func ExampleSetSeverity() {
	logger := New()
	logger.SetSeverity(CRITICAL)
	logger.Print("Hello World")
	// Output:
	// {"severity":"CRITICAL","message":"Hello World"}
}
