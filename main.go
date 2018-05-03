package main

import (
	"fmt"
	"os"

	"github.com/drmjo/heimdall/jsonUtil"
	"github.com/drmjo/heimdall/proxy"
	"github.com/drmjo/heimdall/spec"
)

var apiJson []byte

func main() {

	apiJson = jsonUtil.ResolveRefs(flags.RootJson)

	if flags.GenerateSpecs {
		jsonUtil.WriteToFile(apiJson, flags.RootJson, flags.PrettyPrint)
		os.Exit(0)
	}

	api := spec.NewApi(apiJson)

	fmt.Printf("OpenApi Version: %s\n", api.OpenApi)

	if flags.Serve {
		config := &proxy.Config{
			flags.NoGzip,
			flags.ListenOn}

		proxy.Serve(api, config)

	}
}
