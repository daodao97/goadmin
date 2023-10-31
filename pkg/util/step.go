package util

// Step .
type Step struct {
	Head int
	Tail int
}

// Steps calculates the steps.
func Steps(total, step int) (steps []Step) {
	steps = make([]Step, 0)
	//steps = append(steps, Step{Head:0, Tail:0 + step})
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
