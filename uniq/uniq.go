package uniq

import (
	"bufio"
	"fmt"
	"io"
)

func Uniq(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	count := 0
	prevText := ""
	first := true
	for scanner.Scan() {
		text := scanner.Text()
		if first {
			prevText = text
			first = false
		}
		if text != prevText {
			_, err := output.Write([]byte(fmt.Sprintf("%s %d\n", prevText, count)))
			if err != nil {
				return fmt.Errorf("uniq Write %w", scanner.Err())
			}
			prevText = text
			count = 0
		}
		count++
	}
	if scanner.Err() != nil {
		return fmt.Errorf("uniq scanner.Err() %w", scanner.Err())
	}
	if prevText != "" {
		_, err := output.Write([]byte(fmt.Sprintf("%s %d\n", prevText, count)))
		if err != nil {
			return fmt.Errorf("uniq Write %w", scanner.Err())
		}
	}
	return nil

}
