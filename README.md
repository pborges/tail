Simple naive tail -f
-----

```
func cleanup(name string, offset *int64) {
	ioutil.WriteFile(name, []byte(strconv.FormatInt(*offset, 10)), 0755)
}
func main() {
	filename := "test.log"

	var cur int64 = 0
	offsetFilename := filename + ".offset"
	if b, err := ioutil.ReadFile(offsetFilename); err == nil {
		cur, _ = strconv.ParseInt(string(b), 10, 64)
	}

	c := make(chan os.Signal, 5)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGINT)
	signal.Notify(c, os.Interrupt, syscall.SIGKILL)
	go func() {
		<-c
		cleanup(offsetFilename, &cur)
		os.Exit(1)
	}()

	l, e := tail.TailFileFromOffset(filename, cur)
	for {
		select {
		case line := <-l:
			cur = line.Pos
			fmt.Print("[", line.Pos, "] ", line.Text)
		case err := <-e:
			log.Println("ERROR", err)
		}
	}
}
```