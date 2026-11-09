// #cgo LDFLAGS: -lm -llua -ldl
package interp

import (
	_ "embed"
	"fmt"
	"os"
	"unsafe"
)

type LuaState = C.lua_State

//go:embed main.lua
var source_code string

func startLuaState() {
	fmt.Printf("[INFO] Initializing Lua State. ")
    C.gel_init()

	if C.LUA_STATE == nil {
		fmt.Printf("FAILED. LUA_STATE = NULL. Aborting.\n")
		os.Exit(1)
	}

    C.gel_exec(C.CString(source_code), C.size_t(len(source_code)))
	fmt.Printf("OK. LUA_STATE = %p\n", C.LUA_STATE)
}

func destroyLuaState() {
	fmt.Printf("[INFO] Destroying Lua State.\n")
	C.gel_finish()
}

func ExecLuaString(str string) {
	c_str := C.CString(str)
	defer C.free(unsafe.Pointer(c_str))
	c_strsize := C.size_t(len(str))
	C.gel_exec_ex(c_str, c_strsize, 0, 0);
}

func ExecLuaFunc(name string, args []Var, rets []Var) error {
	funcname := C.CString(name)
	defer C.free(unsafe.Pointer(funcname))

	c_args := make([]C.Var, 0, len(args))
	for _, arg := range args {
		switch a := arg.(type) {
			case VarInt: {
				i := a.i;
				int_var := [8]byte{};
				*(*C.int)(unsafe.Pointer(&int_var)) = C.int(i);
				c_args = append(c_args, C.Var{ typ: C.INTEGER, as: int_var });
			};
			case VarNumber: {
				return fmt.Errorf("ExecLuaFunc: Case VarNumber Not Implemented");
			};
			case VarBool: {
				return fmt.Errorf("ExecLuaFunc: Case VarBool Not Implemented");
			};
			case VarString: {
				msg := a.s;
				str_var := [8]byte{};
				str := C.CString(msg);
				defer C.free(unsafe.Pointer(str));
				*(**C.char)(unsafe.Pointer(&str_var)) = str;
				c_args = append(c_args, C.Var{ typ: C.STRING, as: str_var, });
			};
		}
	}
    var p_args *C.Var
    if len(args) > 0 { 
		p_args = (*C.Var)(unsafe.Pointer(&(c_args[0])));
	}
	args_count := len(c_args)
	p_argscount := C.int(args_count);

	c_rets := make([]C.Var, 0, len(rets))
    for _, ret := range rets {
	    switch a := ret.(type) {
	   	 case VarInt: {
	   		 i := a.i;
	   		 int_var := [8]byte{};
	   		 *(*C.int)(unsafe.Pointer(&int_var)) = C.int(i);
	   		 c_rets = append(c_rets, C.Var{ typ: C.INTEGER, as: int_var });
	   	 };
	   	 case VarNumber: {
	   		 return fmt.Errorf("ExecLuaFunc: Case VarNumber Not Implemented");
	   	 };
	   	 case VarBool: {
	   		 return fmt.Errorf("ExecLuaFunc: Case VarBool Not Implemented");
	   	 };
	   	 case VarString: {
	   		 msg := a.s;
	   		 str_var := [8]byte{};
	   		 str := C.CString(msg);
	   		 defer C.free(unsafe.Pointer(str));
	   		 *(**C.char)(unsafe.Pointer(&str_var)) = str;
	   		 c_rets = append(c_rets, C.Var{ typ: C.STRING, as: str_var, });
	   	 };
	    }
	}

	var p_rets *C.Var = nil;
	if len(rets) > 0 {
		p_rets = (*C.Var)(unsafe.Pointer(&(c_rets[0])));
	}
	rets_count := len(rets)
	p_retscount := C.int(rets_count);

	C.gel_exec_func(funcname, p_args, p_argscount, p_rets, p_retscount);

	if len(c_rets) != len(rets) {
		fmt.Printf("[ERROR] returns in C and Go are different.\n")
		fmt.Printf("[ERROR] C has %d returns.\n", len(c_rets))
		fmt.Printf("[ERROR] Go has %d returns.\n", len(rets))
		os.Exit(1)
	}

	for i := range rets {
 		var c_ret [8]byte = c_rets[i].as
		fmt.Printf("Val: %v\n", c_ret)
		switch r := rets[i].(type) {
			case VarInt: {
				r.i = *(*int)(unsafe.Pointer(&c_ret));
				rets[i] = r 
			};
			case VarNumber: {
				return fmt.Errorf("ExecLuaFunc: Case VarNumber Not Implemented");
			};
			case VarBool: {
				return fmt.Errorf("ExecLuaFunc: Case VarBool Not Implemented");
			};
			case VarString: {
				str := C.GoString(*(**C.char)(unsafe.Pointer(&c_ret[0])));
				r.s = str;
				rets[i] = r;
			};
		}
	}

	return nil;
}

//go:generate go run ../cimport/c_importer.go lua.c
