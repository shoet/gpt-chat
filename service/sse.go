package service

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/shoet/gpt-chat/interfaces"
)

type SSEClient struct {
	client interfaces.Client
}

func (c *SSEClient) Request(req *http.Request, chunkSep string, chunkHandler func(b []byte) error) error {
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	chEve := make(chan string)
	chErr := make(chan error)
	sseScanner := NewSSEScanner(res.Body, chunkSep)
	go func() {
		for sseScanner.Scan() {
			b := sseScanner.Bytes()
			chEve <- string(b)
		}
		if err := sseScanner.Err(); err != nil {
			chErr <- err
			return
		}
		chErr <- io.EOF
	}()

	for {
		select {
		case err := <-chErr:
			if err == io.EOF {
				return nil
			}
			return err
		case event := <-chEve:
			if err := chunkHandler([]byte(event)); err != nil {
				return err
			}
		}
	}
}

type SSEScanner struct {
	scanner *bufio.Scanner
}

func NewSSEScanner(r io.Reader, sseSep string) *SSEScanner {
	// TODO: バッファ周りリファクタリング
	scanner := bufio.NewScanner(r)
	initBufferSize := 1024
	maxBufferSize := 4096
	scanner.Buffer(make([]byte, initBufferSize), maxBufferSize)

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		beforeSep := bytes.Index(data, []byte(sseSep)) // 最初sepの直前
		if beforeSep >= 0 {
			// 最初のsepの位置, dataのsepの直前までのスライス, nil
			return beforeSep + len(sseSep), data[0:beforeSep], nil
		}
		if atEOF {
			// 残りのすべて
			return len(data), data, nil
		}
		return 0, nil, nil
	}
	scanner.Split(split)

	return &SSEScanner{
		scanner: scanner,
	}
}

func (s *SSEScanner) Scan() bool {
	return s.scanner.Scan()
}

func (s *SSEScanner) Bytes() []byte {
	return s.scanner.Bytes()
}

func (s *SSEScanner) Err() error {
	return s.scanner.Err()
}
