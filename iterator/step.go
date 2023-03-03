package iterator

type Step struct {
	Head int
	Tail int
}

// Steps calculates the steps.
//
//	for _, step := range pkg.Steps(len(ids), 10) {
//				cur := ids[step.Head:step.Tail]
//		}
func Steps(total, step int) (steps []Step) {
	steps = make([]Step, 0)
	for i := 0; i < total; i++ {
		if i%step == 0 {
			head := i
			tail := head + step
			if tail > total {
				tail = total
			}
			steps = append(steps, Step{Head: head, Tail: tail})
		}
	}
	return steps
}
