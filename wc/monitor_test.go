package wc

import (
	"testing"
)

func TestCounterMonitor(t *testing.T) {
	cm := NewCounterMonitor()
	go cm.Monitor()

	cm.Insert(&Counter{lines: 111})

	v := cm.Read()
	if v[0].lines != 111 {
		t.Errorf("wrong value stored in monitor, should %d, but got %d", 111, v[0].lines)
	}

	c := &Counter{lines: 333}
	cm.Insert(c)

	v = cm.Read()
	if len(v) != 2 {
		t.Errorf("wrong value stored in monitor, should length %d, but got %d", 2, len(v))
	}

	if v[0].lines != 111 && v[1].lines != 333 {
		t.Errorf("wrong value stored in monitor, should %d, but got %d", 111, v[0].lines)
	}

	cm.Delete(c)
	v = cm.Read()
	if len(v) != 1 {
		t.Errorf("wrong value stored in monitor, should length %d, but got %d", 1, len(v))
	}
}
