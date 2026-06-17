package main

import _ "embed"

// openapiSpec is the generated OpenAPI document, embedded so the binary can
// serve its own contract (e.g. for SDK regeneration or a docs endpoint).
//
//go:embed docs/swagger.yaml
var openapiSpec []byte
