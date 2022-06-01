#include <quickjs.h>
#include <stdlib.h>
#include <string.h>

static inline JSModuleDef *JS_GetModuleDef(JSValue v) {
  return JS_VALUE_GET_PTR(v);
}

JSValue cFunctionProxy(JSContext *ctx, JSValueConst this_val, int argc,
                       JSValueConst *argv, int magic, JSValue *func_data);

static inline int32_t JSValue_GetInt32(JSValueConst v) {
  return JS_VALUE_GET_INT(v);
}
