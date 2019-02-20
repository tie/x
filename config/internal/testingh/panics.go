package testingh

func DoesNotPanic(f func()) bool {
	defer func() { recover() }()
	f()
	return true
}
