package main

// this generates, and provides acces to a bunch of generated metrics/packets
// it does this using a single backing array to save memory allocations

import (
	"bytes"
	"fmt"
)

// so that whatever the timestamp is, if we add it to this, it will always use same amount
// of chars. in other words, don't use that many DP's that this statement would no longer be true!
const tsTpl = 1000000000

type dummyPackets struct {
	key       string
	amount    int
	packetLen int
	scratch   *bytes.Buffer
}

func NewDummyPackets(key string, amount int) *dummyPackets {
	tpl := "%s.dummyPacket 123 %d"
	packetLen := 17 + 10 + len(key)
	scratchBuf := make([]byte, 0, packetLen*amount)
	scratch := bytes.NewBuffer(scratchBuf)
	for i := 1; i <= amount; i++ {
		ts := tsTpl + i
		l, err := fmt.Fprintf(scratch, tpl, key, ts)
		if err != nil {
			panic(err)
		}
		if packetLen != l {
			panic(fmt.Sprintf("bad packet length (or bad write) at index %d.  supposed len: %d, real len: %d", i, packetLen, l))
		}
	}
	return &dummyPackets{key, amount, packetLen, scratch}
}

func (dp *dummyPackets) Get(i int) []byte {
	if i >= dp.amount {
		panic("can't ask for higher index then what we have in dummyPackets")
	}
	sliceFull := dp.scratch.Bytes()
	return sliceFull[dp.packetLen*i : dp.packetLen*(i+1)]
}

func (dp *dummyPackets) All() chan []byte {
	ret := make(chan []byte, 10000) // pretty arbitrary, but seems to help perf
	go func(dp *dummyPackets, ret chan []byte) {
		sliceFull := dp.scratch.Bytes()
		for i := 0; i < dp.amount; i++ {
			ret <- sliceFull[dp.packetLen*i : dp.packetLen*(i+1)]
		}
		close(ret)
	}(dp, ret)
	return ret
}

func mergeAll(in ...chan []byte) chan []byte {
	ret := make(chan []byte, 10000) // pretty arbitrary, but seems to help perf
	go func() {
		for _, inChan := range in {
			for val := range inChan {
				ret <- val
			}
		}
		close(ret)
	}()
	return ret
}
