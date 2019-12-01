package manager

import (
	// "fmt"
	"log"
	"os"
	// "path/filepath"

	"github.com/pelletier/go-toml"

	"github.com/lakesite/ls-config"
	"github.com/lakesite/ls-fibre"
)

// ManagerService has a toml Config property which contains 0box specific directives,
// and a pointer to the web service.
type ManagerService struct {
	Config     *toml.Tree
	WebService *fibre.WebService
}

// Init is required to initialize the manager service via a config file.
func (ms *ManagerService) Init(cfgfile string) {
	if _, err := os.Stat(cfgfile); os.IsNotExist(err) {
		log.Fatalf("File '%s' does not exist.\n", cfgfile)
	} else {
		ms.Config, _ = toml.LoadFile(cfgfile)
	}
}

// Daemonize sets up the web service and defines routes for the API.
func (ms *ManagerService) Daemonize() {
	address := config.Getenv("0BOX_HOST", "127.0.0.1") + ":" + config.Getenv("0BOX_PORT", "6999")
	ms.WebService = fibre.NewWebService("0box", address)
	ms.setupRoutes(ms.WebService)
	ms.WebService.RunWebServer()
}
