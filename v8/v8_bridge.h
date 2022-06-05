#ifndef BRIDGE_H

#include <stdlib.h>

#ifdef __cplusplus
#include "libplatform/libplatform.h"
#include "v8.h"

extern "C" {
#endif

#ifdef __cplusplus
#define TYPE_PTR_ALIAS(type, alias) typedef type *alias;
#else
#define TYPE_PTR_ALIAS(type, alias) typedef void *alias;
#endif

TYPE_PTR_ALIAS(v8::Isolate, v8_isolate_ptr);
TYPE_PTR_ALIAS(v8::Persistent<v8::Context>, v8_persistent_context_ptr);
TYPE_PTR_ALIAS(v8::Persistent<v8::Module>, v8_persistent_module_ptr);
TYPE_PTR_ALIAS(v8::Persistent<v8::Value>, v8_persistent_value_ptr);

#undef TYPE_PTR_ALIAS

void v8_initialize();
void v8_dispose();

const char *v8_get_version();

v8_isolate_ptr v8_isolate_new();
void v8_isolate_dispose(v8_isolate_ptr isolate);

v8_persistent_context_ptr v8_context_new(v8_isolate_ptr isolate);
void v8_context_dispose(v8_persistent_context_ptr persistent_context);

void v8_value_dispose(v8_persistent_value_ptr persistent_value);

v8_persistent_module_ptr
v8_script_compiler_compile_module(v8_isolate_ptr isolate, const char *source,
                                  const char *resource_name);

void v8_module_dispose(v8_persistent_module_ptr persistent_module);

v8_persistent_value_ptr
v8_module_run(v8_isolate_ptr isolate,
              v8_persistent_context_ptr persistent_context,
              v8_persistent_module_ptr persistent_module);

const char *v8_value_to_string(v8_isolate_ptr isolate,
                               v8_persistent_context_ptr persistent_context,
                               v8_persistent_value_ptr persistent_value);

const char *
v8_value_to_detail_string(v8_isolate_ptr isolate,
                          v8_persistent_context_ptr persistent_context,
                          v8_persistent_value_ptr persistent_value);

v8_persistent_value_ptr
v8_function_call(v8_isolate_ptr isolate,
                 v8_persistent_context_ptr persistent_context,
                 v8_persistent_value_ptr persistent_value);

#ifdef __cplusplus
} // extern "C"
#endif

#endif // BRIDGE_H
