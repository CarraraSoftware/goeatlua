#include <assert.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>

// @NOTE: some annotations on the returns by the lua interpreter
// ---
// lua_load:
// @FROM: repo:lua/ file:manual/manual.of line:3859
// > @id{lua_load} can return
// > @Lid{LUA_OK}, @Lid{LUA_ERRSYNTAX}, or @Lid{LUA_ERRMEM}.
// ---
// lua_pcall
// @FROM: repo:lua/ file:manual/manual.of line:4042
// > The @Lid{lua_pcall} function returns one of the following status codes:
// > @Lid{LUA_OK}, @Lid{LUA_ERRRUN}, @Lid{LUA_ERRMEM}, or @Lid{LUA_ERRERR}.
//
// ---
// Values:
// @FROM: repo:lua/ file:lua.h line:48
// > #define LUA_OK		 0
// > #define LUA_YIELD	 1
// > #define LUA_ERRRUN	 2
// > #define LUA_ERRSYNTAX 3
// > #define LUA_ERRMEM	 4
// > #define LUA_ERRERR	 5

#define PRINT_ERROR(msg) fprintf(stderr, "%d:%s:%s [ERROR]: %s", __LINE__, __FILE__, __FUNCTION__, msg) 

lua_State* LUA_STATE;

typedef enum {
	INTEGER,
	FLOAT  ,
	BOOLEAN,
	STRING ,
} VarType;

typedef struct {
    // 4 bytes
    VarType typ;

    // 8 bytes
    union  {
        int         i;
        float       f;
        bool        b;
        const char* s;
    } as;

} Var;


// @NOTE: note the %s on the second line. This macro is suppposed to be
//        used together with __FUNCTION__ to say which function failed.
//        I tried putting that into it's own macro and got some seg faults.
#define MISSING_INIT_FMT                                   \
            "LUA_STATE is NULL."                           \
            " %s can't run without LUA_STATE initialized." \
            " Ensure that gel_init() is called before any other function\n"


const char* lua_err_to_str(int err)
{
    switch (err) {
        case LUA_OK       : return "LUA_OK";
        case LUA_YIELD    : return "LUA_YIELD";
        case LUA_ERRRUN   : return "LUA_ERRRUN";
        case LUA_ERRSYNTAX: return "LUA_ERRSYNTAX";
        case LUA_ERRMEM   : return "LUA_ERRMEM";
        case LUA_ERRERR   : return "LUA_ERRERR";
        default           : return "UNKNOWN_LUA_ERROR";
    }
}


void gel_init()
{
    LUA_STATE = luaL_newstate();
    luaL_openlibs(LUA_STATE);
}

bool gel_exec_ex(const char* source_code, size_t source_len, int nargs, int nret)
{

    if (LUA_STATE == NULL) {
        char errmsg[256];
        sprintf(errmsg, MISSING_INIT_FMT, __FUNCTION__);
        PRINT_ERROR(errmsg);
        return false;
    }

    int err;
    err = luaL_loadbuffer(LUA_STATE, source_code, source_len, NULL);
    if (err != LUA_OK) {
        char errmsg[256];
        sprintf(errmsg, "[ERROR] luaL_loadbuffer failed: %s\n", lua_tostring(LUA_STATE, 1));
        PRINT_ERROR(errmsg);
        return false;
    }


    err = lua_pcall (LUA_STATE, nargs, nret, 0);
    if (err != LUA_OK) {
        char errmsg[256];
        sprintf(errmsg, "[ERROR] luaL_pcall failed: %s\n", lua_tostring(LUA_STATE, 1));
        PRINT_ERROR(errmsg);
        return false;
    }

    if ( nret > 0 ) {
        for (int i = 0; i < nret; ++i)
            lua_pop(LUA_STATE, 1);
    } 
    return true;
}

int gel_exec(const char* source_code, size_t source_len)
{
    if (LUA_STATE == NULL) {
        char errmsg[256];
        sprintf(errmsg, MISSING_INIT_FMT, __FUNCTION__);
        PRINT_ERROR(errmsg);
        return false;
    }
    return gel_exec_ex(source_code, source_len, 0, 0);
}

bool gel_exec_func(const char* func_name, Var* args, int arg_count, Var* rets, int ret_count) {
    
    if (LUA_STATE == NULL) {
        char errmsg[256];
        sprintf(errmsg, MISSING_INIT_FMT, __FUNCTION__);
        PRINT_ERROR(errmsg);
        return false;
    }

    int res;
    lua_getglobal(LUA_STATE, func_name); // coloca func_name no topo da stack

    for (int i = 0; i < arg_count; ++i ) {
        Var it = args[i];
        switch (it.typ) {
            case INTEGER: {
                lua_pushinteger(LUA_STATE, args[i].as.i );
            }; break;
            case FLOAT: {
                lua_pushnumber(LUA_STATE, args[i].as.f );
            }; break;
            case BOOLEAN: {
                lua_pushboolean(LUA_STATE, args[i].as.b );
            }; break;
            case STRING: {
                const char* s = args[i].as.s;
                lua_pushstring(LUA_STATE, s);
            }; break;
            default: {
                char errmsg[256];
                sprintf(errmsg, "[ERROR]: invalid variable type: %d\n", it.typ);
                PRINT_ERROR(errmsg);
                return false;
            }; break;
        }
    }

    if (lua_pcall (LUA_STATE, arg_count, ret_count, 0) != 0) {
        char errmsg[256];
        sprintf(errmsg, "[ERROR]: %s: %s\n", func_name, lua_tostring (LUA_STATE, -1));
        PRINT_ERROR(errmsg);
        return false;
    }

    for (int i = 0; i < ret_count; ++i ) {
        Var r = rets[i];
        switch (r.typ) {
           case INTEGER: {
               rets[i].as.i = lua_tointeger(LUA_STATE, 1 );
               lua_pop(LUA_STATE, 1);
           }; break;
           case FLOAT: {
               rets[i].as.f = lua_tonumber(LUA_STATE, 1 );
               lua_pop(LUA_STATE, 1);
           }; break;
           case BOOLEAN: {
               rets[i].as.b = lua_toboolean(LUA_STATE, 1 );
               lua_pop(LUA_STATE, 1);
           }; break;
           case STRING: {
               const char* s = lua_tostring(LUA_STATE, 1 );
               size_t n = strlen(s);
               char* ret_s = malloc(sizeof(char) * n); // @TODO @LEAK
               strncpy(ret_s, s, n);
               rets[i].as.s = ret_s;
               lua_pop(LUA_STATE, 1);
           }; break;
           default: {
               char errmsg[256];
               sprintf(errmsg, "[ERROR]: invalid return type: %d\n", r.typ);
               PRINT_ERROR(errmsg);
               return false;
           }; break;
        }
    }

    return true;
}


void gel_finish()
{
	lua_close(LUA_STATE);
}
