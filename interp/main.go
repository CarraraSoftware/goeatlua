package interp

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func example_1_call_function_declared_in_lua_file() {
    var err error

	// which function
	name := "fancy_add"

	// args
	msg  := VarString{ "Uma Outra String" };
	p1   := VarInt{ 303 };
	p2   := VarInt{ 201 };

	// returns
	rets := []Var{  VarInt{} };

	// callit
	err = ExecLuaFunc( name, []Var{ msg, p1, p2 }, rets)
	if err != nil { panic(fmt.Sprintf("fancy_add: é isso aí: %s\n",err)) }
	fmt.Printf("Resultado de Fancy Add: add(%d, %d) = %d\n", p1, p2, rets[0]);
}

func example_2_declare_and_call_function() {
    var err error

	// define the function with a string
	var func_src_code strings.Builder
	func_src_code.WriteString("function concat(a, b)\n")
	func_src_code.WriteString("    return a .. \" +\" .. b\n")
	func_src_code.WriteString("end\n")
	ExecLuaString(func_src_code.String())

	// which function
	name := "concat"

	// args
	a := VarString{ "Uma String" };
	b := VarString{ " Outra String"};

	// returns
	r := []Var{ VarString{} };

	// callit
	err = ExecLuaFunc(name, []Var{ a,b }, r);
    if err != nil { panic(fmt.Sprintf("concat string não rolou: %s\n",err)) }
	fmt.Printf("Resultado de Concat String: %s\n", r[0].(VarString).s);
}


func Init() {
 	startLuaState()
}

func Deinit() {
    destroyLuaState()
}

func REPL() {
	var err error
	var str string
	scanner := bufio.NewScanner(os.Stdin)
    for {
		fmt.Printf("> ")
		if !scanner.Scan() {
			break
		}

		str = scanner.Text()
        if str == "exit" {
			break;
		}

		ExecLuaString(str)
	}
	err = scanner.Err()
	if err != nil {
		fmt.Printf("[ERROR] Bufio Scanner failed: %s\n", err)
	}
}
