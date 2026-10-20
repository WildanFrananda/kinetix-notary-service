package ethereum

import "sync"

type answerSequence struct {
	mu   sync.Mutex
	last sentAnswer
}
