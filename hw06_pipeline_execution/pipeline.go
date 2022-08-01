package hw06pipelineexecution

import "time"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 && done == nil {
		c := make(Bi)
		close(c)
		return c
	}
	if len(stages) == 0 && done != nil {
		c := make(Bi)
		go func() {
			for {
				select {
				case <-done:
					close(c)
					return
				default:
					time.Sleep(1 * time.Millisecond)
				}
			}
		}()
		return c
	}
	chanIn := in
	for _, st := range stages {
		outStage := make(Bi)
		go func(stage Stage, chanInL In, chanOutL Bi) {
			defer close(chanOutL)
			chanOut := stage(chanInL)
			for {
				select {
				case d, ok := <-chanOut:
					if !ok {
						return
					}
					chanOutL <- d
				case <-done:
					return
				}
			}
		}(st, chanIn, outStage)
		chanIn = outStage
	}
	return chanIn
}
