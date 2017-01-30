package system

func Promise(f func() error) <-chan error {

	ch := make(chan error, 1)
	go func() {
		ch <- f()
	}()
	return ch
}
