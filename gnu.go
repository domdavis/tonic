package tonic

import "github.com/gin-gonic/gin"

// Clacks headers.
const Clacks = "X-Clacks-Overhead"

// GNU Terry Pratchett. GNU adds the X-Clacks-Overhead header.
func GNU(c *gin.Context) {
	c.Header(Clacks, "GNU Terry Pratchett, Russel Winder")
}
