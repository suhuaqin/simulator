package handler

import (
	"errors"
	"math/rand"
	"time"
)

type transferHelp struct {
	remainder uint32
	// 小于该值丢弃
	discardLE int64
}

func newTransferHelp() *transferHelp {
	rand.Seed(time.Now().Unix())
	return &transferHelp{
		remainder: 1,
		discardLE: -1,
	}
}

func (t *transferHelp) isDiscard() bool {
	return int64(rand.Uint32()%t.remainder) < t.discardLE
}

func (t *transferHelp) setConfig(remainder uint32, discardLE int64) error {
	if remainder == 0 {
		return errors.New("remainder can not zero")
	}
	t.remainder = remainder
	t.discardLE = discardLE
	return nil
}
