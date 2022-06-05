#include "v8_bridge.h"

static auto platform = v8::platform::NewDefaultPlatform();
static auto default_allocator =
    v8::ArrayBuffer::Allocator::NewDefaultAllocator();

void v8_initialize() {
  v8::V8::InitializePlatform(platform.get());
  v8::V8::Initialize();
}

void v8_dispose() {
  v8::V8::Dispose();
  v8::V8::DisposePlatform();
}

const char *v8_get_version() { return v8::V8::GetVersion(); }

v8_isolate_ptr v8_isolate_new() {
  v8::Isolate::CreateParams create_params;
  create_params.array_buffer_allocator = default_allocator;
  return v8::Isolate::New(create_params);
}

void v8_isolate_dispose(v8_isolate_ptr isolate) { isolate->Dispose(); }

v8_persistent_context_ptr v8_context_new(v8_isolate_ptr isolate) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);

  auto context = v8::Context::New(isolate);
  auto persistent_context = new v8::Persistent<v8::Context>(isolate, context);

  return persistent_context;
}

void v8_context_dispose(v8_persistent_context_ptr persistent_context) {
  persistent_context->Reset();
  delete persistent_context;
}

void v8_value_dispose(v8_persistent_value_ptr persistent_value) {
  persistent_value->Reset();
  delete persistent_value;
}

v8_persistent_module_ptr
v8_script_compiler_compile_module(v8_isolate_ptr isolate, const char *source,
                                  const char *resource_name) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);
  v8::TryCatch try_catch(isolate);

  // TODO: Handle errors.
  auto source_string =
      v8::String::NewFromUtf8(isolate, source).ToLocalChecked();

  // TODO: Handle errors.
  auto resource_name_string =
      v8::String::NewFromUtf8(isolate, resource_name).ToLocalChecked();

  v8::ScriptOrigin script_origin(isolate, resource_name_string, 0, 0, false, -1,
                                 v8::Local<v8::Value>(), false, false, true);

  v8::ScriptCompiler::Source module_source(source_string, script_origin);

  // TODO: Handle errors.
  auto module = v8::ScriptCompiler::CompileModule(isolate, &module_source)
                    .ToLocalChecked();

  auto persistent_module = new v8::Persistent<v8::Module>(isolate, module);

  return persistent_module;
}

void v8_module_dispose(v8_persistent_module_ptr persistent_module) {
  persistent_module->Reset();
  delete persistent_module;
}

bool v8_module_instantiate(v8_isolate_ptr isolate,
                           v8_persistent_context_ptr persistent_context,
                           v8_persistent_module_ptr petsistent_module) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);

  // auto persistent_context = context->Get(isolate);
  // auto persistent_module = module->Get(isolate);

  return true;
}

static v8::MaybeLocal<v8::Module>
resolve_module_callback(v8::Local<v8::Context> context,
                        v8::Local<v8::String> specifier,
                        v8::Local<v8::FixedArray> import_assertions,
                        v8::Local<v8::Module> referrer) {
  return v8::MaybeLocal<v8::Module>();
}

v8_persistent_value_ptr
v8_module_run(v8_isolate_ptr isolate,
              v8_persistent_context_ptr persistent_context,
              v8_persistent_module_ptr persistent_module) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);

  auto context = persistent_context->Get(isolate);
  auto module = persistent_module->Get(isolate);

  v8::Context::Scope context_scope(context);

  module->InstantiateModule(context, resolve_module_callback).ToChecked();

  auto promise = module->Evaluate(context).ToLocalChecked();

  if (module->IsGraphAsync()) {
    auto persistent_value = new v8::Persistent<v8::Value>(isolate, promise);
    return persistent_value;
  }

  auto value = module->GetModuleNamespace();

  auto default_value =
      v8::Local<v8::Object>::Cast(value)
          ->Get(context,
                v8::String::NewFromUtf8(isolate, "default",
                                        v8::NewStringType::kInternalized)
                    .ToLocalChecked())
          .ToLocalChecked();

  auto persistent_value = new v8::Persistent<v8::Value>(isolate, default_value);

  return persistent_value;
}

const char *v8_value_to_string(v8_isolate_ptr isolate,
                               v8_persistent_context_ptr persistent_context,
                               v8_persistent_value_ptr persistent_value) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);

  auto context = persistent_context->Get(isolate);
  auto value = persistent_value->Get(isolate);

  // TODO: Handle errors.
  auto str = value->ToString(context).ToLocalChecked();
  v8::String::Utf8Value utf8_value(isolate, str);

  // TODO: Handle errors.
  return strdup(*utf8_value);
}

const char *
v8_value_to_detail_string(v8_isolate_ptr isolate,
                          v8_persistent_context_ptr persistent_context,
                          v8_persistent_value_ptr persistent_value) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);

  auto context = persistent_context->Get(isolate);
  auto value = persistent_value->Get(isolate);

  // TODO: Handle errors.
  auto str = value->ToDetailString(context).ToLocalChecked();
  v8::String::Utf8Value utf8_value(isolate, str);

  // TODO: Handle errors.
  return strdup(*utf8_value);
}

v8_persistent_value_ptr
v8_function_call(v8_isolate_ptr isolate,
                 v8_persistent_context_ptr persistent_context,
                 v8_persistent_value_ptr persistent_value) {
  v8::Locker locker(isolate);
  v8::Isolate::Scope isolate_scope(isolate);
  v8::HandleScope handle_scope(isolate);

  auto context = persistent_context->Get(isolate);
  auto value = persistent_value->Get(isolate);

  v8::Context::Scope context_scope(context);

  auto result_value = v8::Local<v8::Function>::Cast(value)
                          ->Call(context, v8::Undefined(isolate), 0, nullptr)
                          .ToLocalChecked();

  auto persistent_result_value =
      new v8::Persistent<v8::Value>(isolate, result_value);

  return persistent_result_value;
}
