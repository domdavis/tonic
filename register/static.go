package register

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Static content will be served from the content root against the router.
func Static(contentRoot string, router *gin.Engine) error {
	files, err := os.ReadDir(contentRoot)

	if err != nil {
		return fmt.Errorf("failed to read %s: %w", contentRoot, err)
	}

	for _, file := range files {
		route := fmt.Sprintf("/%s", file.Name())
		root := fmt.Sprintf("%s/%s", contentRoot, file.Name())

		if file.IsDir() {
			router.Static(route, root)
		} else {
			router.StaticFile(route, root)
		}
	}

	return nil
}
