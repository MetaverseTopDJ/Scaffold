package util

import (
	"bufio"
	"bytes"
)

/*
	计算单词数量
*/

type WordCounter int

// Write 实现 fmt.Fprintf 写入方法
func (w *WordCounter) Write(b []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(b))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return -1, err
	}
	*w += WordCounter(count)
	return len(b), nil
}
