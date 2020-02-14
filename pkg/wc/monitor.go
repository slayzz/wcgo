package wc

type CounterMonitor struct {
	readValue   chan CounterList
	writeValue  chan *Counter
	deleteValue chan CounterList
}

func NewCounterMonitor() CounterMonitor {
	return CounterMonitor{
		readValue:   make(chan CounterList),
		writeValue:  make(chan *Counter),
		deleteValue: make(chan CounterList),
	}
}

func (cl *CounterMonitor) Insert(val *Counter) {
	cl.writeValue <- val
}

func (cl *CounterMonitor) Read() CounterList {
	return <-cl.readValue
}

func (cl *CounterMonitor) Delete(c *Counter) {
	counters := <-cl.readValue
	for i, cv := range counters {
		if cv == c {
			counters[i] = counters[len(counters)-1]
			counters[len(counters)-1] = nil
			counters = counters[:len(counters)-1]
			break
		}
	}
	cl.deleteValue <- counters
}

func (cl *CounterMonitor) Monitor() {
	var counterList CounterList
	for {
		select {
		case c := <-cl.writeValue:
			counterList = append(counterList, c)
		case cl.readValue <- counterList:
		case newList := <-cl.deleteValue:
			counterList = newList
		}
	}
}
