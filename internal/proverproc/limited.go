package proverproc

import "bytes"

type limited struct {
	buffer *bytes.Buffer
	limit  int
}

func (l *limited) Write(p []byte) (int, error) {
	room := l.limit - l.buffer.Len()
	if room > 0 {
		if len(p) > room {
			l.buffer.Write(p[:room])
		} else {
			l.buffer.Write(p)
		}
	}

	return len(p), nil
}
