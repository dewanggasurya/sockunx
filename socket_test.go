package sockunx_test

import (
	"fmt"
	"sockunx"
	"sockunx/example/handler"
	"sockunx/pkg/helper"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// RequestResponseCase struct
type RequestResponseCase struct {
	name     string
	request  string
	response string
}

func TestSocket(t *testing.T) {
	Convey("Subject : Socket", t, func() {
		Convey("Sending & receiving", func() {
			tcs := []RequestResponseCase{
				{"Sending string 1", "hello", "hello"},
				{"Sending string 2", "hi", "hi"},
				{"Sending number 1", "123", "123"},
				{"Sending number 2", "0", "0"},
				{"Sending JSON formatted string", "{name : 'John'}", "{name : 'John'}"},
				{"Sending characters", "-%?😊", "-%?😊"},
			}

			path := "test.socket"
			server, e := sockunx.NewServer(path)
			So(e, ShouldBeNil)
			server.Run(true)
			defer server.Stop()

			client, e := sockunx.NewClient(path)
			So(e, ShouldBeNil)

			for _, tc := range tcs {
				Convey(tc.name, func() {
					response, e := client.Send(tc.request)
					So(e, ShouldBeNil)
					So(response, ShouldEqual, response)
				})
			}
		})

		Convey("Receiving with handler", func() {
			tcs := []struct {
				name    string
				handler sockunx.Handler
				tcs     []RequestResponseCase
			}{
				{
					name: "Uppercase Handler",
					handler: func(request string) (response interface{}, e error) {
						return strings.ToUpper(request), nil
					},
					tcs: []RequestResponseCase{
						{"Sending short string", "john doe", "JOHN DOE"},
						{"Sending long string", "Elit nostrud consectetur aute est ad ea.", "ELIT NOSTRUD CONSECTETUR AUTE EST AD EA."},
						{"Sending super-long string", "Amet ex incididunt cillum magna ullamco Lorem laboris dolore amet laborum veniam. Nostrud dolor eu duis ipsum reprehenderit laboris officia reprehenderit sit Lorem dolor irure esse Lorem. Ipsum excepteur ullamco esse reprehenderit esse est ad exercitation aliquip Lorem labore et ea consectetur. Quis dolor ex consequat Lorem cupidatat est magna elit sit. Nulla dolore pariatur occaecat voluptate sint enim non incididunt cupidatat incididunt. Reprehenderit ad Lorem sunt mollit ullamco.", ""},
					},
				},
				{
					name: "String to JSON Token Slice Handler",
					handler: func(request string) (response interface{}, e error) {
						return helper.ToJSONString(strings.Split(request, " ")), nil
					},
					tcs: []RequestResponseCase{
						{"Sending string 1", "john doe", helper.ToJSONString([]string{"john", "doe"})},
						{"Sending string 2", "Elit nostrud consectetur aute est ad ea.", helper.ToJSONString([]string{"Elit", "nostrud", "consectetur", "aute", "est", "ad", "ea."})},
						{"Sending string 3", "1 2 34 5 6", helper.ToJSONString([]string{"1", "2", "34", "5", "6"})},
					},
				},
				{
					name:    "Fizz Buzz Handler",
					handler: handler.Index,
					tcs: []RequestResponseCase{
						{
							name:     "Sending payload 1",
							request:  `{"id":"one","from":0,"to":15,"fizz":"zzif","buzz":"zzub"}\n`,
							response: `{"one":["zzifzzub",1,2,"zzif",4,"zzub","zzif",7,8,"zzif","zzub",11,"zzif",13,14,"zzifzzub"]}`,
						},
						{
							name:     "Sending payload 2",
							request:  `{"id":"two","from":6,"to":10,"fizz":"a","buzz":"b"}`,
							response: `{"two":["a",7,8,"a","b"]}`,
						},
						{
							name:    "Sending payload 3",
							request: `{"id":"one","from":0,"to":5,"fizz":"zzif","buzz":"zzub"}\n{"id":"two","from":6,"to":10,"fizz":"zzif2","buzz":"zzub2"}\n`,
							response: `{"one":["zzifzzub",1,2,"zzif",4,"zzub"]}
{"two":["zzif2",7,8,"zzif2","zzub2"]}`,
						},
					},
				},
			}

			path := "test.socket"

			server, e := sockunx.NewServer(path, 512)
			So(e, ShouldBeNil)
			server.Run(true)
			defer server.Stop()

			client, e := sockunx.NewClient(path, 1024)
			So(e, ShouldBeNil)

			for _, tc1 := range tcs {
				Convey(tc1.name, func() {

					server.RegisterHandler(tc1.handler)
					for _, tc2 := range tc1.tcs {
						Convey(tc2.name, func() {
							response, e := client.Send(tc2.request)

							So(e, ShouldBeNil)
							So(tc2.response, ShouldEqual, response)
						})
					}
				})
			}
		})

		Convey("Handling error", func() {
			path := "test.socket"

			server, e := sockunx.NewServer(path)
			So(e, ShouldBeNil)
			server.Run(true)
			defer server.Stop()

			server.RegisterHandler(func(request string) (response interface{}, e error) {
				if request == "" {
					return response, fmt.Errorf("well you got empty request error")
				}

				if request == "-" {
					return response, fmt.Errorf("you've come so far to find forbidden request")
				}

				return request, nil
			})

			client, e := sockunx.NewClient(path)
			So(e, ShouldBeNil)

			tcs := []struct {
				name     string
				request  string
				response string
			}{
				{"Sending empty request", "", "ERR well you got empty request error"},
				{"Sending forbidden request", "-", "ERR you've come so far to find forbidden request"},
				{"Sending acceptable value", "hi", "hi"},
			}

			for _, tc := range tcs {
				Convey(tc.name, func() {
					response, e := client.Send(tc.request)
					So(e, ShouldBeNil)
					So(response, ShouldEqual, tc.response)
				})
			}
		})

		Convey("Multiple server & client", func() {
			formatter := func(serverID int, request string) string {
				return fmt.Sprintf("From server %d : %s", serverID, request)
			}

			handler := func(serverID int) sockunx.Handler {
				return func(request string) (response interface{}, e error) {
					return formatter(serverID, request), nil
				}
			}

			Convey("Multiple servers", func() {

				servers := []*sockunx.Server{}
				clients := []*sockunx.Client{}
				for i := 0; i < 5; i++ {
					path := fmt.Sprintf("test%d.socket", i)

					server, e := sockunx.NewServer(path)
					So(e, ShouldBeNil)

					servers = append(servers, server)
					servers[i].Run(true)
					server.RegisterHandler(handler(i))

					client, e := sockunx.NewClient(path)
					So(e, ShouldBeNil)
					clients = append(clients, client)
				}

				defer func() {
					for i := range servers {
						servers[i].Stop()
					}
				}()

				for i := 0; i < 20; i++ {
					for j := 0; j < 5; j++ {
						request := fmt.Sprintf("Client %d [%d]", j, i)
						response, e := clients[j].Send(request)
						So(e, ShouldBeNil)
						So(response, ShouldEqual, formatter(j, request))
					}
				}
			})

			Convey("Multiple clients", func() {
				path := "test.socket"
				server, e := sockunx.NewServer(path)
				So(e, ShouldBeNil)
				server.Run(true)
				server.RegisterHandler(handler(1))

				clients := []*sockunx.Client{}
				for i := 0; i < 20; i++ {
					client, e := sockunx.NewClient(path)
					So(e, ShouldBeNil)
					clients = append(clients, client)
				}

				for i, client := range clients {
					request := fmt.Sprintf("Client %d", i)
					response, e := client.Send(request)
					So(e, ShouldBeNil)
					So(response, ShouldEqual, formatter(1, request))
				}
			})

		})

	})
}
