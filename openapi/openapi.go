package openapi

import (
	_ "embed"

	"github.com/gin-gonic/gin"
)

// -------------------------------------------------------------------------------------
// Config
// -------------------------------------------------------------------------------------

var (
	//go:embed openapi.yaml
	openapiYAML []byte

	swaggerUI = []byte(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta
      name="description"
      content="SwaggerUI"
    />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js" crossorigin></script>
  <script src="https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js" crossorigin></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: window.location.pathname + '/openapi.yaml',
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        layout: "StandaloneLayout",
      });
    };
  </script>
  </body>
</html>
	`)
)

// -------------------------------------------------------------------------------------
// Handlers
// -------------------------------------------------------------------------------------

func HandleSwaggerUI(c *gin.Context) {
	c.Header("content-type", "text/html")
	_, err := c.Writer.Write(swaggerUI)
	if err != nil {
		c.Error(err)
	}
}

func HandleSpecYAML(c *gin.Context) {
	c.Header("content-type", "application/yaml")
	_, err := c.Writer.Write(openapiYAML)
	if err != nil {
		c.Error(err)
	}
}
