package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	chanIn := in
	for i := 0; i < len(stages); i++ {
		outStage := make(Bi)
		go func(i int, chanInL In, chanOutL Bi) {
			defer close(chanOutL)
			chanOut := stages[i](chanInL)
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
		}(i, chanIn, outStage)
		chanIn = outStage
	}
	return chanIn
}
