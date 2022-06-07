package main

import (
	"context"
	"log"
	"runtime"
	"time"

	"github.com/jackc/puddle/puddleg"
	_ "go.uber.org/automaxprocs"

	v8 "github.com/lujjjh/ssr/v8"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func main() {
	bundleContent := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{"src/entry.server.tsx"},
		Outfile:     "dist/entry.server.mjs",
		Bundle:      true,
		Write:       true,
		Target:      esbuild.ES2017,
		Format:      esbuild.FormatESModule,
	}).OutputFiles[0].Contents

	v8.Initialize()
	defer v8.Dispose()

	type entry struct {
		c *v8.Context
		f *v8.Value
	}

	pool := puddleg.NewPool(func(ctx context.Context) (res entry, err error) {
		isolate := v8.NewIsolate()
		c := isolate.NewContext()

		polyfillModule, _ := isolate.CompileModule(`
			globalThis.TextEncoder = function TextEncoder() {};
			TextEncoder.prototype.encode = function () {};
	
			const console = globalThis.console = {};
		`, "polyfill.js")
		defer polyfillModule.Dispose()
		polyfillModule.Run(c)

		module, _ := isolate.CompileModule(string(bundleContent), "main.js")
		defer module.Dispose()

		f := module.Run(c)

		return entry{c, f}, nil
	}, func(e entry) {
		e.f.Dispose()
		e.c.Dispose()
		e.c.Isolate().Dispose()
	}, int32(runtime.GOMAXPROCS(0)))

	ln, err := reuseport.Listen("tcp4", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on :3000")

	log.Fatal(fasthttp.Serve(ln, func(c *fasthttp.RequestCtx) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resource, err := pool.Acquire(ctx)
		if err != nil {
			c.Error("Failed to acquire SSR runtime", fasthttp.StatusInternalServerError)
			return
		}
		defer resource.Release()

		e := resource.Value()
		c.SuccessString("text/html", e.f.Call().String())
	}))
}
