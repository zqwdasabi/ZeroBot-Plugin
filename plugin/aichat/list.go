package aichat

import (
	"strings"
	"sync"

	"github.com/fumiama/deepinfra"
	"github.com/fumiama/deepinfra/model"
)

const listcap = 6

type list struct {
	mu sync.RWMutex
	m  map[int64][]string
}

func newlist() list {
	return list{
		m: make(map[int64][]string, 64),
	}
}

func (l *list) add(grp int64, txt string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	msgs, ok := l.m[grp]
	if !ok {
		msgs = make([]string, 1, listcap)
		msgs[0] = txt
		l.m[grp] = msgs
		return
	}
	if len(msgs) < cap(msgs) {
		msgs = append(msgs, txt)
		l.m[grp] = msgs
		return
	}
	copy(msgs, msgs[1:])
	msgs[len(msgs)-1] = txt
}

func (l *list) body(mn, sysp string, temp float32, grp int64) deepinfra.Model {
	m := model.NewCustom(mn, sepstr, temp, 0.9, 1024).System(sysp)
	l.mu.RLock()
	defer l.mu.RUnlock()
	sz := len(l.m[grp])
	if sz == 0 {
		return m.User("自己随机开启新话题")
	}
	return m.User(strings.Join(l.m[grp], "\n\n"))
}
