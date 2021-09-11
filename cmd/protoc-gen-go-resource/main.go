package main

import (
	"flag"

	"github.com/ucarion/protoc-gen-go-resource/internal/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flags flag.FlagSet
	opts := protogen.Options{ParamFunc: flags.Set}

	opts.Run(func(plugin *protogen.Plugin) error {
		return generator.Generate(plugin)

		// gen := generator.New(plugin)
		// for _, f := range plugin.Files {
		// 	if err := gen.Generate(f); err != nil {
		// 		return err
		// 	}
		// }
		//
		// return nil
	})
}
