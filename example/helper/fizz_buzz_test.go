package helper_test

import (
	"fmt"
	"sockunx/example/helper"
	h "sockunx/pkg/helper"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFizzBuzz(t *testing.T) {
	Convey("Subject : FizzBuzz", t, func() {
		cases := []struct {
			from   int
			to     int
			label  []helper.FizzBuzzLabel
			output []interface{}
		}{
			{from: 0, to: 5, output: []interface{}{"fizzbuzz", 1, 2, "fizz", 4, "buzz"}},
			{from: 6, to: 10, output: []interface{}{"fizz", 7, 8, "fizz", "buzz"}},
			{from: 0, to: 15, output: []interface{}{"fizzbuzz", 1, 2, "fizz", 4, "buzz", "fizz", 7, 8, "fizz", "buzz", 11, "fizz", 13, 14, "fizzbuzz"}},
			{from: 0, to: 5, output: []interface{}{"ab", 1, 2, "a", 4, "b"}, label: []helper.FizzBuzzLabel{{Fizz: "a", Buzz: "b"}}},
			{from: 6, to: 10, output: []interface{}{"fzz", 7, 8, "fzz", "bzz"}, label: []helper.FizzBuzzLabel{{Fizz: "fzz", Buzz: "bzz"}}},
			{from: 0, to: 15, output: []interface{}{"zzifzzub", 1, 2, "zzif", 4, "zzub", "zzif", 7, 8, "zzif", "zzub", 11, "zzif", 13, 14, "zzifzzub"}, label: []helper.FizzBuzzLabel{{Fizz: "zzif", Buzz: "zzub"}}},
		}

		for _, c := range cases {
			label := ""
			if len(c.label) > 0 {
				label = " with customized labeling"
			}
			Convey(fmt.Sprintf("Fizz Buzz %d to %d%s", c.to, c.from, label), func() {
				result := helper.FizzBuzz(c.from, c.to, c.label...)
				So(h.ToJSONString(result), ShouldEqual, h.ToJSONString(c.output))
			})
		}
	})
}
