package handler

import (
	"fmt"
	socket "sockunx"
	"sockunx/example/helper"
	h "sockunx/pkg/helper"
	"strings"
)

// FizzBuzzRequest struct
type FizzBuzzRequest struct {
	ID   string `json:"id"`
	From int    `json:"from"`
	To   int    `json:"to"`
	Fizz string `json:"fizz"`
	Buzz string `json:"buzz"`
}

// Handlers
var (
	Index socket.Handler = func(request string) (response interface{}, e error) {
		if request == "" {
			return nil, fmt.Errorf("undefined request")
		}

		request = strings.ReplaceAll(request, "\\n", `<-|->`)
		request = strings.ReplaceAll(request, "\n", `<-|->`)
		rs := strings.Split(request, "<-|->")

		responses := []string{}
		for _, r := range rs {
			if r == "" {
				continue
			}
			data := FizzBuzzRequest{}
			e = h.FromJSONString(r, &data)
			if e != nil {
				return
			}

			output := helper.FizzBuzz(data.From, data.To, helper.FizzBuzzLabel{
				Fizz: data.Fizz,
				Buzz: data.Buzz,
			})

			responses = append(responses, h.ToJSONString(map[string]interface{}{
				data.ID: output,
			}))
		}

		response = strings.Join(responses, "\n")
		return
	}
)
