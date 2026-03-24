package handler

import "embed"

//go:embed assets/dist-frontend assets/dist-admin
var distFS embed.FS
