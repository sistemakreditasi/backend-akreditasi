package teeamai

import (
	"github.com/sistemakreditasi/backend-akreditasi/routes"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WebHook", routes.URL)
}
