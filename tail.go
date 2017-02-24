package tail

import (
	"os"
	"time"
	"bufio"
	"io"
	"errors"
)

var Sleep time.Duration = 250 * time.Millisecond

type Line struct {
	Text string
	Pos  int64
}

func TailFile(name string) (chan Line, chan error) {
	return TailFileFromOffset(name, 0)
}

func TailFileFromOffset(name string, off int64) (chan Line, chan error) {
	lines := make(chan Line, 100)
	err := make(chan error)
	go tail(name, off, lines, err)
	return lines, err
}

func tail(name string, off int64, linesChan chan Line, errChan chan error) {
	for {
		file, err := os.Open(name)
		if err != nil {
			time.Sleep(Sleep)
			continue
		}
		_, err = file.Seek(off, os.SEEK_SET)
		if err != nil {
			select {
			case errChan <- err:
			}
			continue
		}
		curFileInfo, err := file.Stat()
		if err != nil {
			select {
			case errChan <- err:
			}
			continue
		}
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					fileInfo, err := os.Stat(name)
					if err != nil {
						select {
						case errChan <- err:
						}
						break
					}
					if !os.SameFile(fileInfo, curFileInfo) {
						select {
						case errChan <- errors.New("File moved"):
						}
						break
					}
					time.Sleep(Sleep)
				} else {
					select {
					case errChan <- err:
					}
					break
				}
			} else {
				pos, err := file.Seek(0, os.SEEK_CUR)
				if err != nil {
					select {
					case errChan <- err:
					}
					break
				}
				linesChan <- Line{Text: line, Pos: pos}
			}
		}
	}
}
