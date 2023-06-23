package sort

import (
	"bufio"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"math/rand"
	"os"
	"path"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func Sort(input io.Reader, maxRows int, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	rows := make([]string, 0, maxRows)
	dir := os.TempDir()
	sortingJobName := randStringRunes(10)
	chunkNum := 0

	var chunks []string

	dumpChunk := func() error {
		chunkFileName := path.Join(dir, fmt.Sprintf("%s_%d", sortingJobName, chunkNum))
		chunks = append(chunks, chunkFileName)
		file, err := os.Create(chunkFileName)
		if err != nil {
			return err
		}
		sort.Strings(rows)
		for _, r := range rows {
			_, err := file.Write(append([]byte(r), '\n'))
			if err != nil {
				return err
			}
		}
		chunkNum++
		return file.Close()
	}

	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		rows = append(rows, scanner.Text())
		if len(rows) >= maxRows {
			err := dumpChunk()
			if err != nil {
				return err
			}
			rows = rows[:0]
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	err := dumpChunk()
	if err != nil {
		return err
	}

	err = merge(chunks, output)
	if err != nil {
		return err
	}

	return nil
}

func merge(chunks []string, output io.Writer) error {
	if len(chunks) == 0 {
		return nil
	}
	if len(chunks) == 1 {
		file, err := os.Open(chunks[0])
		if err != nil {
			return fmt.Errorf("merge %w", err)
		}
		_, err = io.Copy(output, file)
		if err != nil {
			return fmt.Errorf("merge Copy %w", err)
		}
		err = file.Close()
		if err != nil {
			return fmt.Errorf("merge close %w", err)
		}
		return nil
	}
	if len(chunks) == 2 {
		chunk1, err := os.Open(chunks[0])
		if err != nil {
			return fmt.Errorf("merge %w", err)
		}
		chunk2, err := os.Open(chunks[1])
		if err != nil {
			return fmt.Errorf("merge %w", err)
		}
		scanner1 := bufio.NewScanner(chunk1)
		scanner2 := bufio.NewScanner(chunk2)
		move := 0
		for {
			ok1 := func() bool {
				if move == 1 || move == 0 {
					return scanner1.Scan()
				}
				return true
			}()
			if scanner1.Err() != nil {
				return scanner1.Err()
			}
			ok2 := func() bool {
				if move == 2 || move == 0 {
					return scanner2.Scan()
				}
				return true
			}()
			if scanner2.Err() != nil {
				return scanner2.Err()
			}
			if ok1 && ok2 {
				text1 := scanner1.Text()
				text2 := scanner2.Text()
				if text1 < text2 {
					_, err = output.Write(append([]byte(text1), '\n'))
					if err != nil {
						return fmt.Errorf("merge Write %w", err)
					}
					move = 1
				} else {
					_, err = output.Write(append([]byte(text2), '\n'))
					if err != nil {
						return fmt.Errorf("merge Write %w", err)
					}
					move = 2
				}
				continue
			}
			if ok1 {
				_, err = output.Write(append([]byte(scanner1.Text()), '\n'))
				if err != nil {
					return fmt.Errorf("merge Write %w", err)
				}
				move = 0
				continue
			}
			if ok2 {
				_, err = output.Write(append([]byte(scanner2.Text()), '\n'))
				if err != nil {
					return fmt.Errorf("merge Write %w", err)
				}
				move = 0
				continue
			}
			break
		}
		if scanner1.Err() != nil {
			return fmt.Errorf("merge scanner.Err() %w", scanner1.Err())
		}
		if scanner2.Err() != nil {
			return fmt.Errorf("merge scanner.Err() %w", scanner2.Err())
		}
		return nil
	}

	mid := len(chunks) / 2
	dir := os.TempDir()
	name := "part_" + randStringRunes(10)
	part1, err := os.Create(path.Join(dir, fmt.Sprintf("%s_%d", name, 0)))
	if err != nil {
		return fmt.Errorf("merge %w", err)
	}
	part2, err := os.Create(path.Join(dir, fmt.Sprintf("%s_%d", name, 1)))
	if err != nil {
		return fmt.Errorf("merge %w", err)
	}
	group := errgroup.Group{}
	group.Go(func() error {
		return merge(chunks[:mid], part1)
	})
	group.Go(func() error {
		return merge(chunks[mid:], part2)
	})
	err = group.Wait()
	if err != nil {
		return fmt.Errorf("merge group.Wait %w", err)
	}
	err = part1.Close()
	if err != nil {
		return fmt.Errorf("merge close %w", err)
	}
	err = part2.Close()
	if err != nil {
		return fmt.Errorf("merge close %w", err)
	}

	err = merge([]string{part1.Name(), part2.Name()}, output)
	if err != nil {
		return fmt.Errorf("merge merge %w", err)
	}

	return err
}
