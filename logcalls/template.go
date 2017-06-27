package {{$.PackageName}}

// #cgo pkg-config: libcryptsetup
// #include <libcryptsetup.h>
// #include <stdlib.h>
// #include "gocryptsetup.h"
// #include "{{$.FileName "h"}}"
// void {{$.Ns}}_log_default(int, char *, void *);
import "C"
import (
	"unsafe"
	"log"
)

// Setup a default logging function that uses the standard logger.
func init() {
	C.crypt_set_log_callback(
		nil,
		(*[0]byte)(C.{{$.Ns}}_log_default),
		nil,
	)
}
//export {{$.Ns}}_log_default
func {{$.Ns}}_log_default(level C.int, msg *C.char, usrptr unsafe.Pointer) {
	log.Print(C.GoString(msg))
}

{{define "declparam"}}{{end}}
{{define "returnvalue"}}{{end}}

{{range .Methods}}
func (d Device) {{.GoName}}({{range $k, $v := .DeclParams}}{{if $k}}, {{end}}{{$v.Name}} {{$v.GoType}}{{end}}) ({{with .Return}}out {{.}}, {{end}}err error) {
	arglist := (*C.struct_{{$.Ns}}_logstack)(nil)
	{{range .Params}}
	{{if eq "string" (.GoType)}}
	_{{.Name}} := C.CString({{.Value}})
	defer C.free(unsafe.Pointer(_{{.Name}}))
	{{else if eq "[]byte" (.GoType)}}
	_{{.Name}} := unsafe.Pointer(nil)
	if {{.Value}} != nil {
		_{{.Name}} = C.CBytes({{.Value}})
		defer C.free(_{{.Name}})
	}
	{{else}}
	_{{.Name}} := C.{{.Type}}({{.Value}})
	{{end}}
	{{end}}
	
	ival := C.{{$.Ns}}_{{.Name}}(
		&arglist,
		d.cd,
		{{range .Params}}
		_{{.Name}},
		{{end}})
	defer C.{{$.Ns}}_logstack_free(arglist)
	
	msg := ""
	for i := arglist; i != nil; i = i.prev {
		msg = C.GoString(i.message) + "\n" + msg
	}
	if ival < 0 {
		err = newError(int(ival), msg)
		return
	}
	
	{{with .Return}}out = ({{.}})(ival){{end}}
	return
}
{{end}}
