package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/Min-Feng/failure"
	"github.com/rs/zerolog/log"
)

func init() {
	logInit()
}

// error codes for your application.
const (
	NotFound  failure.StringCode = "NotFound"
	Forbidden failure.StringCode = "Forbidden"
)

func GetACL(projectID, userID string) (acl interface{}, e error) {
	notFound := true
	if notFound {
		return nil, failure.New(NotFound,
			failure.Context{"project_id": projectID, "user_id": userID},
		)
	}
	return nil, failure.Unexpected("unexpected error")
}

func GetProject(projectID, userID string) (project interface{}, e error) {
	_, err := GetACL(projectID, userID)
	if err != nil {
		if failure.Is(err, NotFound) {
			return nil, failure.Translate(err, Forbidden,
				failure.Message("no acl exists"),
				failure.Context{"additional_info": "hello"},
			)
		}
		return nil, failure.Wrap(err)
	}
	return nil, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	_, err := GetProject(r.FormValue("project_id"), r.FormValue("user_id"))
	if err != nil {
		HandleError(w, err)
		return
	}
}

func getHTTPStatus(err error) int {
	c, ok := failure.CodeOf(err)
	if !ok {
		return http.StatusInternalServerError
	}
	switch c {
	case NotFound:
		return http.StatusNotFound
	case Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func getMessage(err error) string {
	msg, ok := failure.MessageOf(err)
	if ok {
		return msg
	}
	return "Error"
}

func HandleError(w http.ResponseWriter, err error) {
	w.WriteHeader(getHTTPStatus(err))
	io.WriteString(w, getMessage(err))

	fmt.Println("============ JSON Format ============")
	b := failure.JSONFormat(err)
	log.Error().RawJSON("JSONFormat", b).Send()

	// {
	// 	"level": "error",
	// 	"JSONFormat": {
	// 	  "detail": [
	// 		{
	// 		  "frame": {
	// 			"func": "main.GetACL ",
	// 			"source": " /tmp/sandbox340936509/prog.go:22 "
	// 		  },
	// 		  "context": {
	// 			"project_id": "aaa",
	// 			"user_id": "111"
	// 		  },
	// 		  "code": "NotFound"
	// 		},
	// 		{
	// 		  "frame": {
	// 			"func": "main.GetProject ",
	// 			"source": " /tmp/sandbox340936509/prog.go:33 "
	// 		  },
	// 		  "context": {
	// 			"additional_info": "hello"
	// 		  },
	// 		  "message": "no acl exists",
	// 		  "code": "Forbidden"
	// 		}
	// 	  ],
	// 	  "callStack": [
	// 		{
	// 		  "func": "main.GetACL",
	// 		  "source": " /tmp/sandbox340936509/prog.go:22 "
	// 		},
	// 		{
	// 		  "func": "main.GetProject",
	// 		  "source": " /tmp/sandbox340936509/prog.go:30 "
	// 		},
	// 		{
	// 		  "func": "main.Handler",
	// 		  "source": " /tmp/sandbox340936509/prog.go:44 "
	// 		},
	// 		{
	// 		  "func": "main.main",
	// 		  "source": " /tmp/sandbox340936509/prog.go:103 "
	// 		},
	// 		{
	// 		  "func": "runtime.main",
	// 		  "source": " /usr/local/go-faketime/src/runtime/proc.go:203 "
	// 		},
	// 		{
	// 		  "func": "runtime.goexit",
	// 		  "source": " /usr/local/go-faketime/src/runtime/asm_amd64.s:1373 "
	// 		}
	// 	  ]
	// 	},
	// 	"time": "2009-11-10T23:00:00Z"
	// }

	fmt.Println()

	fmt.Println("============ Error ============")
	log.Error().Msgf("%v\n%+[1]v", err)

	fmt.Println("============ Error part ============")
	code, _ := failure.CodeOf(err)
	fmt.Printf("Code = %v\n", code)

	msg, _ := failure.MessageOf(err)
	fmt.Printf("Message = %v\n", msg)

	cs, _ := failure.CallStackOf(err)
	fmt.Printf("CallStack = %v\n", cs)

	fmt.Printf("Cause = %v\n", failure.CauseOf(err))

}

func main() {
	req := httptest.NewRequest(http.MethodGet, "/?project_id=aaa&user_id=111", nil)
	rec := httptest.NewRecorder()
	Handler(rec, req)
	log.Info().Msg("Hello world")
}
