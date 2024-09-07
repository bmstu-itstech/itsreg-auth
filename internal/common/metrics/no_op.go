package metrics

type NoOp struct{}

func (m NoOp) Inc(key string, value int) {}
