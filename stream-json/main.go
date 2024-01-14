package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func main() {
	// https://pkg.go.dev/encoding/json#Decoder.Token
	// https://stackoverflow.com/questions/31794355/decode-large-stream-json

	data1 := `
[
	{"Name": "Ed", "Text": "Knock knock."},
	{"Name": "Sam", "Text": "Who's there?"},
	{"Name": "Ed", "Text": "Go fmt."},
	{"Name": "Sam", "Text": "Go fmt who?"},
	{"Name": "Ed", "Text": "Go fmt yourself!"}
]
`
	reader1 := strings.NewReader(data1)
	err1 := decodeJson(reader1, AnalysisHandler1)
	if err1 != nil {
		fmt.Printf("data1 err=%v\n", err1)
	}

	fmt.Println()

	// 	data2 := `
	// {
	//   "List":[
	//     {
	//       "Name":"Ed",
	//       "Text":"Knock knock."
	//     },
	//     {
	//       "Name":"Sam",
	//       "Text":"Who's there?"
	//     },
	//     {
	//       "Name":"Ed",
	//       "Text":"Go fmt."
	//     },
	//     {
	//       "Name":"Sam",
	//       "Text":"Go fmt who?"
	//     },
	//     {
	//       "Name":"Ed",
	//       "Text":"Go fmt yourself!"
	//     }
	//   ]
	// }
	// `
	// 	reader2 := strings.NewReader(data2)
	// 	err2 := decodeJson(reader2, AnalysisHandler2)
	// 	fmt.Printf("data1 err=%v\n", err2)
}

func decodeJson(reader io.Reader, handler func(decoder *json.Decoder) error) (err error) {
	decoder := json.NewDecoder(reader)
	for decoder.More() {
		fmt.Println("offset=", decoder.InputOffset())

		token, TokenErr := decoder.Token()
		if TokenErr != nil {
			if TokenErr == io.EOF {
				return nil
			}
			return TokenErr
		}

		// fmt.Println("token=", token)
		// fmt.Println("offset=", decoder.InputOffset())
		switch token {
		case json.Delim('['):
			for decoder.More() {
				err = handler(decoder)
				if err != nil {
					return err
				}
			}

		// 單一函數沒辦法同時解析 array obj
		// 因為只要讀取 token
		// cursor 就會前進
		//
		case json.Delim('{'):
			err = handler(decoder)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func AnalysisHandler1(decoder *json.Decoder) error {
	type Message struct {
		Name, Text string
	}

	var payload Message
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Println("decode fail:", err)
		return err
	}

	fmt.Printf("payload=%#v\n", payload)
	return nil
}

func AnalysisHandler2(decoder *json.Decoder) error {
	type Item struct {
		Name, Text string
	}

	type Message struct {
		List []Item
	}

	var payload Message
	err := decoder.Decode(&payload)
	if err != nil {
		fmt.Println("decode fail:", err)
		return err
	}

	fmt.Printf("payload=%#v\n", payload)
	return nil
}
