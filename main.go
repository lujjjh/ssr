package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/lujjjh/ssr/quickjs"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func main() {
	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{"src/entry.server.tsx"},
		Outfile:     "bundle.server.js",
		Bundle:      true,
		Format:      esbuild.FormatESModule,
	})
	if len(result.Errors) > 0 {
		for _, err := range result.Errors {
			fmt.Fprintln(os.Stderr, err.Text)
		}
		return
	}

	textEncoderClassID := quickjs.NewClassID()

	r := quickjs.NewRuntime()
	defer r.Free()

	err := quickjs.NewClass(r, textEncoderClassID, quickjs.ClassDef{
		ClassName: "TextEncoder",
	})
	if err != nil {
		panic(err)
	}

	c := r.NewContext()
	defer c.Free()

	textEncoderClassID.SetProto(c, c.GetGlobalObject())

	textEncoderEncode := c.NewGoFunction(func(c quickjs.Context, this quickjs.Value, args ...quickjs.Value) quickjs.Value {
		return c.Undefined()
	}, 0, 0)
	defer textEncoderEncode.Free()
	c.GetGlobalObject().DefinePropertyValue("encode", textEncoderEncode, quickjs.PropConfigurable|quickjs.PropWritable|quickjs.PropEnumerable)

	textEncoderCtor := c.NewGoFunction(func(c quickjs.Context, this quickjs.Value, args ...quickjs.Value) quickjs.Value {
		return c.NewClassInstance(textEncoderClassID)
	}, 0, 1)
	defer textEncoderCtor.Free()

	textEncoderCtor.SetConstructorBit(true)
	textEncoderCtor.SetConstructor(c.GetGlobalObject())

	globalObject := c.GetGlobalObject()

	c.Eval(`console = {}; void 0`, "", quickjs.EvalTypeGlobal)

	globalObject.DefinePropertyValue("TextEncoder", textEncoderCtor, quickjs.PropConfigurable)

	module := c.Eval(string(result.OutputFiles[0].Contents), result.OutputFiles[0].Path, quickjs.EvalTypeModule|quickjs.EvalFlagCompileOnly)
	defer module.Free()

	if module.IsException() {
		e := c.GetException()
		defer e.Free()
		panic(e)
	}

	value := c.EvalFunction(module)
	if value.IsException() {
		e := c.GetException()
		defer e.Free()
		log.Println(e)
	}

	entry := c.GetModuleExport(module, "default")
	defer entry.Free()

	value = entry.Call(c.Undefined())
	defer value.Free()
	if value.IsException() {
		e := c.GetException()
		defer e.Free()
		log.Println(e)
	}

	log.Println(value)
}
