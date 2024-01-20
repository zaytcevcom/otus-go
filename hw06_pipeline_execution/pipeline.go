package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		stageCh := make(Bi)
		go process(stageCh, out, done)
		out = stage(stageCh)
	}

	return out
}

func process(in Bi, out Out, done In) {
	defer close(in)

	for {
		select {
		case <-done:
			return
		case result, ok := <-out:
			if !ok {
				return
			}
			in <- result
		}
	}
}
