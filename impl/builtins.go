package main

import "fmt"

func (ev *Eval) builtin(name string, arg, funArg Value) (Value, error) {
	switch name {
	case "neg":
		switch arg.(type) {
		case int:
			return arg.(int) * -1, nil
		default:
			return nil, fmt.Errorf("type error")
		}

	case "greater":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		last := 0
		for i := 0; i < len(arr); i++ {
			n, ok := arr[i].(int)
			if !ok {
				return nil, fmt.Errorf("type error")
			}

			if last <= n && i != 0 {
				return 0, nil
			}

			last = n
		}

		return 1, nil

	case "sub":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		if len(arr) < 2 {
			return nil, fmt.Errorf("not enough arguments")
		}

		sum, ok := arr[0].(int)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		for i := 1; i < len(arr); i++ {
			n, ok := arr[i].(int)
			if !ok {
				return nil, fmt.Errorf("type error")
			}

			sum -= n
		}

		return sum, nil

	case "div":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		if len(arr) < 2 {
			return nil, fmt.Errorf("not enough arguments")
		}

		sum, ok := arr[0].(int)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		for i := 1; i < len(arr); i++ {
			n, ok := arr[i].(int)
			if !ok {
				return nil, fmt.Errorf("type error")
			}

			if n == 0 {
				return nil, fmt.Errorf("division by zero")
			}

			sum /= n
		}

		return sum, nil

	case "add":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		if len(arr) < 2 {
			return nil, fmt.Errorf("not enough arguments")
		}

		sum := 0

		for i := 0; i < len(arr); i++ {
			n, ok := arr[i].(int)
			if !ok {
				return nil, fmt.Errorf("type error")
			}

			sum += n
		}

		return sum, nil

	case "mul":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		if len(arr) < 2 {
			return nil, fmt.Errorf("not enough arguments")
		}

		sum, ok := arr[0].(int)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		for i := 1; i < len(arr); i++ {
			n, ok := arr[i].(int)
			if !ok {
				return nil, fmt.Errorf("type error")
			}

			sum *= n
		}

		return sum, nil

	case "branch":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		if len(arr) != 3 {
			return nil, fmt.Errorf("not enough arguments")
		}

		if _, ok := arr[1].(Node); !ok {
			return nil, fmt.Errorf("type error")
		}

		if _, ok := arr[2].(Node); !ok {
			return nil, fmt.Errorf("type error")
		}

		if isTruthy(arr[0]) {
			return ev.expression(arr[1].(Node), arg, funArg)
		} else {
			return ev.expression(arr[2].(Node), arg, funArg)
		}

	case "call":
		arr, ok := arg.([]Value)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		if len(arr) != 2 {
			return nil, fmt.Errorf("not enough arguments")
		}

		fun, ok := arr[0].(Node)
		if !ok {
			return nil, fmt.Errorf("type error")
		}

		return ev.expression(fun, arr[1], arr[1])

	case "print":
		fmt.Println(arg)
	}

	return nil, nil
}
