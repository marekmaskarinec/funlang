package main

type Value interface {
}

func isTruthy(val Value) bool {
	switch val.(type) {
	case int:
		return val.(int) != 0
	case []Value, Node:
		return true
	}

	return false
}

type Eval struct {
	defs map[string]Node
}

func (ev *Eval) expression(ast Node, arg, funArg Value) (Value, error) {
	switch ast.value.kind {
	case kindArrayStart:
		arr := make([]Value, len(ast.children))
		for i := 0; i < len(ast.children); i++ {
			val, err := ev.expression(ast.children[i], arg, funArg)
			if err != nil {
				return ast, err
			}

			arr[i] = val
		}

		return arr, nil

	case kindArgSubst:
		return funArg, nil

	case kindChainOperator:
		val, err := ev.expression(ast.children[0], arg, funArg)
		if err != nil {
			return val, err
		}

		return ev.expression(ast.children[1], val, funArg)

	case kindNumber:
		return ast.value.toInt()

	case kindIdent:
		if ast.value.value == "fun" {
			return ast.children[0], nil
		} else if node, exists := ev.defs[ast.value.value]; exists {
			return ev.expression(node, arg, arg)
		} else {
			val, err := ev.builtin(ast.value.value, arg, funArg)
			if err != nil {
				err = ast.value.errorf("%w", err)
			}

			return val, err
		} 

	case kindRef:
		arr, err := ev.expression(ast.children[0], arg, funArg)
		if err != nil {
			return arg, err
		}

		if _, ok := arr.([]Value); !ok {
			return arg, ast.children[0].value.errorf("type error: array expected")
		}

		index, err := ev.expression(ast.children[1], arg, funArg)
		if err != nil {
			return arg, err
		}

		if _, ok := index.(int); !ok {
			return arg, ast.children[1].value.errorf("type error: number expected")
		}

		return arr.([]Value)[index.(int)], nil

	default:
		for i := 0; i < len(ast.children); i++ {
			val, err := ev.expression(ast.children[i], arg, funArg)
			if err != nil {
				return val, err
			}

			if i == len(ast.children)-1 {
				return val, nil
			}
		}
	}

	return arg, nil
}
