package embedded

import "embed"

//go:embed public/* views/*
var BoosterFiles embed.FS
