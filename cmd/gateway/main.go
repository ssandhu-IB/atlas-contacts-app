package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/Infoblox-CTO/ngp.api.toolkit/gw"
	"github.com/infobloxopen/atlas-contacts-app/cmd/config"
)

var (
	ServerAddr  string
	GatewayAddr string
	SwaggerDir  string
)

func main() {
	// create HTTP handler for gateway
	errHandler := runtime.WithProtoErrorHandler(gw.ProtoMessageErrorHandler)
	opHandler := runtime.WithMetadata(gw.MetadataAnnotator)
	serverHandler, err := NewAtlasContactsAppHandler(context.Background(), ServerAddr, errHandler, opHandler)
	// strip all but trailing "/" on incoming requests
	serverHandler = http.StripPrefix(
		config.GATEWAY_URL[:len(config.GATEWAY_URL)-1],
		serverHandler,
	)
	if err != nil {
		log.Fatalln(err)
	}
	// map HTTP endpoints to handlers
	mux := http.NewServeMux()
	mux.Handle("/atlas-contacts-app/v1/", serverHandler)
	mux.HandleFunc("/swagger/", SwaggerHandler)
	// serve handlers on the gateway address
	http.ListenAndServe(GatewayAddr, mux)
}

func init() {
	// default gateway values; optionally configured via command-line flags
	flag.StringVar(&ServerAddr, "server", config.SERVER_ADDRESS, "address of the gRPC server")
	flag.StringVar(&GatewayAddr, "gateway", config.GATEWAY_ADDRESS, "address of the gateway server")
	flag.StringVar(&SwaggerDir, "swagger-dir", config.SWAGGER_DIR, "directory of the swagger.json file")
	flag.Parse()
}
