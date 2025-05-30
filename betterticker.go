package betterticker

import "time"

func NewBetterTicker(duration time.Duration) *BetterTicker {
	return &BetterTicker{
		Duration:     duration,
		C:            make(chan time.Time),
		timeLastTick: time.Now(),
		internalC:    make(chan time.Time),
		stop:         make(chan time.Time),
		internalStop: make(chan time.Time),
	}
}

type BetterTicker struct {
	Duration     time.Duration
	C            chan (time.Time)
	timeLastTick time.Time
	internalC    chan (time.Time)
	internalStop chan (time.Time)
	stop         chan (time.Time)
}

func (bt *BetterTicker) Start() {
	go runCore(bt.timeLastTick, bt.Duration, bt.internalC, bt.internalStop)
	go bt.tick()
}

func (bt *BetterTicker) Stop() {
	bt.stop <- time.Now()
}

func (bt *BetterTicker) tick() {
	t := time.Now()
	bt.timeLastTick = t
	bt.C <- t
	for {
		select {
		case t = <-bt.internalC:
			bt.timeLastTick = t
			go runCore(bt.timeLastTick, bt.Duration, bt.internalC, bt.internalStop)
			bt.C <- t

		case <-bt.stop:
			bt.internalStop <- time.Now()
			time.Sleep(1)
			close(bt.stop)
			close(bt.internalC)
			close(bt.C)
			close(bt.internalStop)

			return
		}
	}
}

func runCore(timeLastTick time.Time, duration time.Duration, internalC chan (time.Time), exitC chan (time.Time)) {
	timeNow := time.Now()
	for timeLastTick.Add(duration).After(timeNow) {
		timeNow = time.Now()
		select {
		case <-exitC:
			return
		default:
			continue
		}
	}

	internalC <- timeNow
}
