package leetcode

func dailyTemperatures(temperatures []int) []int {

	answers := make([]int, len(temperatures))

	stack := []dayTemp{}

	for curr_day, curr_temp := range temperatures {

		for len(stack) > 0 && top(stack).temp < curr_temp {

			prev_day := pop(&stack).idx
			answers[prev_day] = curr_day - prev_day
		}
		stack = append(stack, dayTemp{idx: curr_day, temp: curr_temp})
	}

	return answers
}

type dayTemp struct {
	idx  int
	temp int
}

func top(stack []dayTemp) dayTemp {
	return stack[len(stack)-1]
}

func pop(stack *[]dayTemp) dayTemp {
	s := *stack
	top := s[len(s)-1]
	*stack = s[:len(s)-1]
	return top
}
