package main

/*
#define error_urinotfound 1
#define error_handshake 2
#define error_notallowed 3
#define error_readerror 4
#define error_writeerror 5
#define error_dataerror 6
#define error_unknownerror 7
#define error_malformedurierror 8

struct pycoap_result
{
    char * result;
    int error;
};

typedef struct pycoap_result coapresult;
*/
import "C"

import (
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	coap "github.com/moroen/gocoap/v4"
)

var _debugLevel int

var _version = "0.0.2"

// Python Functions

func handleError(err error) C.int {
	if err != nil {
		if _debugLevel == 1 {
			log.Println(err.Error())
		}

		switch err {
		case coap.UriNotFound:
			return C.error_urinotfound
		case coap.ErrorHandshake:
			return C.error_handshake
		case coap.MethodNotAllowed:
			return C.error_notallowed
		case coap.ErrorWriteTimeout:
			return C.error_writeerror
		case coap.ErrorReadTimeout:
			return C.error_readerror
		case coap.ErrorBadData:
			return C.error_dataerror
		case coap.ErrorUnknownError:
			return C.error_unknownerror
		}
		return C.int(0)
	}
	return C.int(0)
}

//export coapDebugLevel
func coapDebugLevel(level C.int) C.int {

	_debugLevel = int(level)
	return C.int(0)
}

//export coapRequest
func coapRequest(gateway, uri *C.char) C.coapresult {
	var res C.coapresult

	gw := strings.Split(C.GoString(gateway), ":")
	port, err := strconv.Atoi(gw[1])
	if err != nil {
		res.result, res.error = C.CString(""), C.error_malformedurierror
		return res
	}

	params := coap.RequestParams{Host: gw[0], Port: port, Uri: C.GoString(uri)}

	response, err := coap.GetRequest(params)

	res.result, res.error = C.CString(validateResponse(response)), handleError(err)

	return res
}

//export coapRequestDTLS
func coapRequestDTLS(gateway, uri, ident, key *C.char) C.coapresult {
	var res C.coapresult

	gw := strings.Split(C.GoString(gateway), ":")
	port, err := strconv.Atoi(gw[1])
	if err != nil {
		res.result, res.error = C.CString(""), C.error_malformedurierror
		return res
	}

	params := coap.RequestParams{Host: gw[0], Port: port, Uri: C.GoString(uri), Id: C.GoString(ident), Key: C.GoString(key)}

	response, err := coap.GetRequest(params)

	res.result, res.error = C.CString(validateResponse(response)), handleError(err)

	return res
}

//export coapPutRequest
func coapPutRequest(gateway, uri, payload *C.char) C.coapresult {
	var res C.coapresult
	/*
		msg, err := PutRequest(C.GoString(gateway), C.GoString(uri), C.GoString(payload))
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		return C.CString(msg.String())
	*/
	return res
}

//export coapPutRequestDTLS
func coapPutRequestDTLS(gateway, uri, ident, key, payload *C.char) C.coapresult {
	var res C.coapresult

	gw := strings.Split(C.GoString(gateway), ":")
	port, err := strconv.Atoi(gw[1])
	if err != nil {
		res.result, res.error = C.CString(""), C.error_malformedurierror
		return res
	}

	params := coap.RequestParams{Host: gw[0], Port: port, Uri: C.GoString(uri), Id: C.GoString(ident), Key: C.GoString(key), Payload: C.GoString(payload)}

	response, err := coap.PutRequest(params)

	res.result, res.error = C.CString(validateResponse(response)), handleError(err)

	return res
}

//export coapPostRequestDTLS
func coapPostRequestDTLS(gateway, uri, ident, key, payload *C.char) C.coapresult {
	var res C.coapresult

	gw := strings.Split(C.GoString(gateway), ":")
	port, err := strconv.Atoi(gw[1])
	if err != nil {
		res.result, res.error = C.CString(""), C.error_malformedurierror
		return res
	}

	params := coap.RequestParams{Host: gw[0], Port: port, Uri: C.GoString(uri), Id: C.GoString(ident), Key: C.GoString(key), Payload: C.GoString(payload)}

	response, err := coap.PostRequest(params)

	res.result, res.error = C.CString(validateResponse(response)), handleError(err)

	return res
}

//export coapCloseConnection
func coapCloseConnection() C.int {
	err := coap.CloseDTLSConnection()
	if err != nil {
		return 1
	}
	return 0
}

//export coapVersion
func coapVersion() *C.char {
	cstr := C.CString(_version)
	return cstr
}

func validateResponse(response []byte) string {
	res := string(response)

	if !utf8.ValidString(res) {
		v := make([]rune, 0, len(res))
		for i, r := range res {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(res[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		res = string(v)
	} else {
		res = string(res)
	}

	return res
}

func main() {}
