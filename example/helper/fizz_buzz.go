package helper

import "fmt"

// FizzBuzzLabel struct
type FizzBuzzLabel struct {
	Fizz string
	Buzz string
}

// FizzBuzz helper
func FizzBuzz(from, to int, label ...FizzBuzzLabel) []interface{} {
	_label := FizzBuzzLabel{
		Fizz: "fizz",
		Buzz: "buzz",
	}

	if len(label) > 0 {
		l := label[0]

		if l.Fizz != "" {
			_label.Fizz = l.Fizz
		}

		if l.Buzz != "" {
			_label.Buzz = l.Buzz
		}
	}

	result := []interface{}{}
	for i := from; i <= to; i++ {
		var value interface{}

		value = i
		if i%3 == 0 {
			value = _label.Fizz
		}

		if i%5 == 0 {
			if v, ok := value.(string); ok {
				value = fmt.Sprintf("%s%s", v, _label.Buzz)
			} else {
				value = _label.Buzz
			}
		}

		result = append(result, value)
	}

	return result
}
