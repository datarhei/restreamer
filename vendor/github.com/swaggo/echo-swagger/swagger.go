package echoSwagger

import (
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo/v4"
	swaggerFiles "github.com/swaggo/files/v2"
	"github.com/swaggo/swag"
)

// Config stores echoSwagger configuration variables.
type Config struct {
	// The url pointing to API definition (normally swagger.json or swagger.yaml). Default is `mockedSwag.json`.
	URLs                 []string
	DocExpansion         string
	DomID                string
	InstanceName         string
	DeepLinking          bool
	PersistAuthorization bool
	SyntaxHighlight      bool

	// The information for OAuth2 integration, if any.
	OAuth *OAuthConfig
}

// OAuthConfig stores configuration for Swagger UI OAuth2 integration. See
// https://swagger.io/docs/open-source-tools/swagger-ui/usage/oauth2/ for further details.
type OAuthConfig struct {
	// The ID of the client sent to the OAuth2 IAM provider.
	ClientId string

	// The OAuth2 realm that the client should operate in. If not applicable, use empty string.
	Realm string

	// The name to display for the application in the authentication popup.
	AppName string
}

// URL presents the url pointing to API definition (normally swagger.json or swagger.yaml).
func URL(url string) func(*Config) {
	return func(c *Config) {
		c.URLs = append(c.URLs, url)
	}
}

// DeepLinking true, false.
func DeepLinking(deepLinking bool) func(*Config) {
	return func(c *Config) {
		c.DeepLinking = deepLinking
	}
}

// SyntaxHighlight true, false.
func SyntaxHighlight(syntaxHighlight bool) func(*Config) {
	return func(c *Config) {
		c.SyntaxHighlight = syntaxHighlight
	}
}

// DocExpansion list, full, none.
func DocExpansion(docExpansion string) func(*Config) {
	return func(c *Config) {
		c.DocExpansion = docExpansion
	}
}

// DomID #swagger-ui.
func DomID(domID string) func(*Config) {
	return func(c *Config) {
		c.DomID = domID
	}
}

// InstanceName specified swag instance name.
func InstanceName(instanceName string) func(*Config) {
	return func(c *Config) {
		c.InstanceName = instanceName
	}
}

// PersistAuthorization Persist authorization information over browser close/refresh.
// Defaults to false.
func PersistAuthorization(persistAuthorization bool) func(*Config) {
	return func(c *Config) {
		c.PersistAuthorization = persistAuthorization
	}
}

func OAuth(config *OAuthConfig) func(*Config) {
	return func(c *Config) {
		c.OAuth = config
	}
}

func newConfig(configFns ...func(*Config)) *Config {
	config := Config{
		URLs:                 []string{"doc.json", "doc.yaml"},
		DocExpansion:         "list",
		DomID:                "swagger-ui",
		InstanceName:         "swagger",
		DeepLinking:          true,
		PersistAuthorization: false,
		SyntaxHighlight:      true,
	}

	for _, fn := range configFns {
		fn(&config)
	}

	if config.InstanceName == "" {
		config.InstanceName = swag.Name
	}

	return &config
}

// WrapHandler wraps swaggerFiles.Handler and returns echo.HandlerFunc
var WrapHandler = EchoWrapHandler()

// EchoWrapHandler wraps `http.Handler` into `echo.HandlerFunc`.
func EchoWrapHandler(options ...func(*Config)) echo.HandlerFunc {
	config := newConfig(options...)

	// create a template with name
	index, _ := template.New("swagger_index.html").Parse(indexTemplate)

	var re = regexp.MustCompile(`^(.*/)([^?].*)?[?|.]*$`)

	return func(c echo.Context) error {
		if c.Request().Method != http.MethodGet {
			return echo.NewHTTPError(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		}

		matches := re.FindStringSubmatch(c.Request().RequestURI)
		path := matches[2]

		switch filepath.Ext(path) {
		case ".html":
			c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			c.Response().Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			c.Response().Header().Set("Content-Type", "application/javascript")
		case ".json":
			c.Response().Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".yaml":
			c.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")
		case ".png":
			c.Response().Header().Set("Content-Type", "image/png")
		}

		response := c.Response()
		// This check fixes an error introduced here: https://github.com/labstack/echo/blob/8da8e161380fd926d4341721f0328f1e94d6d0a2/response.go#L86-L88
		if _, ok := response.Writer.(http.Flusher); ok {
			defer response.Flush()
		}

		switch path {
		case "":
			_ = c.Redirect(http.StatusMovedPermanently, matches[1]+"/"+"index.html")
		case "index.html":
			_ = index.Execute(c.Response().Writer, config)
		case "doc.json":
			doc, err := swag.ReadDoc(config.InstanceName)
			if err != nil {
				c.Error(err)

				return nil
			}

			_, _ = c.Response().Writer.Write([]byte(doc))
		case "doc.yaml":
			jsonString, err := swag.ReadDoc(config.InstanceName)
			if err != nil {
				c.Error(err)

				return nil
			}
			doc, err := yaml.JSONToYAML([]byte(jsonString))
			if err != nil {
				c.Error(err)

				return nil
			}
			_, _ = c.Response().Writer.Write(doc)
		default:
			c.Request().URL.Path = matches[2]
			http.FileServer(http.FS(swaggerFiles.FS)).ServeHTTP(c.Response(), c.Request())
		}

		return nil
	}
}

const indexTemplate = `<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css" >
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
  <style>
    html
    {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
    }
    *,
    *:before,
    *:after
    {
        box-sizing: inherit;
    }

    body {
      margin:0;
      background: #fafafa;
    }
  </style>
</head>

<body>

<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" style="position:absolute;width:0;height:0">
  <defs>
    <symbol viewBox="0 0 20 20" id="unlocked">
          <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V6h2v-.801C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8z"></path>
    </symbol>

    <symbol viewBox="0 0 20 20" id="locked">
      <path d="M15.8 8H14V5.6C14 2.703 12.665 1 10 1 7.334 1 6 2.703 6 5.6V8H4c-.553 0-1 .646-1 1.199V17c0 .549.428 1.139.951 1.307l1.197.387C5.672 18.861 6.55 19 7.1 19h5.8c.549 0 1.428-.139 1.951-.307l1.196-.387c.524-.167.953-.757.953-1.306V9.199C17 8.646 16.352 8 15.8 8zM12 8H8V5.199C8 3.754 8.797 3 10 3c1.203 0 2 .754 2 2.199V8z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="close">
      <path d="M14.348 14.849c-.469.469-1.229.469-1.697 0L10 11.819l-2.651 3.029c-.469.469-1.229.469-1.697 0-.469-.469-.469-1.229 0-1.697l2.758-3.15-2.759-3.152c-.469-.469-.469-1.228 0-1.697.469-.469 1.228-.469 1.697 0L10 8.183l2.651-3.031c.469-.469 1.228-.469 1.697 0 .469.469.469 1.229 0 1.697l-2.758 3.152 2.758 3.15c.469.469.469 1.229 0 1.698z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow">
      <path d="M13.25 10L6.109 2.58c-.268-.27-.268-.707 0-.979.268-.27.701-.27.969 0l7.83 7.908c.268.271.268.709 0 .979l-7.83 7.908c-.268.271-.701.27-.969 0-.268-.269-.268-.707 0-.979L13.25 10z"/>
    </symbol>

    <symbol viewBox="0 0 20 20" id="large-arrow-down">
      <path d="M17.418 6.109c.272-.268.709-.268.979 0s.271.701 0 .969l-7.908 7.83c-.27.268-.707.268-.979 0l-7.908-7.83c-.27-.268-.27-.701 0-.969.271-.268.709-.268.979 0L10 13.25l7.418-7.141z"/>
    </symbol>


    <symbol viewBox="0 0 24 24" id="jump-to">
      <path d="M19 7v4H5.83l3.58-3.59L8 6l-6 6 6 6 1.41-1.41L5.83 13H21V7z"/>
    </symbol>

    <symbol viewBox="0 0 24 24" id="expand">
      <path d="M10 18h4v-2h-4v2zM3 6v2h18V6H3zm3 7h12v-2H6v2z"/>
    </symbol>

  </defs>
</svg>

<div id="{{.DomID}}"></div>

<script src="./swagger-ui-bundle.js"> </script>
<script src="./swagger-ui-standalone-preset.js"> </script>
<script>
window.onload = function() {
  // Build a system
  const ui = SwaggerUIBundle({
	urls: [
	{{range $index, $url := .URLs}}
		{
			name: "{{$url}}",
			url: "{{$url}}",
		},
	{{end}}
	],
    syntaxHighlight: {{.SyntaxHighlight}},
    deepLinking: {{.DeepLinking}},
    docExpansion: "{{.DocExpansion}}",
    persistAuthorization: {{.PersistAuthorization}},
    dom_id: "#{{.DomID}}",
    validatorUrl: null,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  })

  {{if .OAuth}}
  ui.initOAuth({
    clientId: "{{.OAuth.ClientId}}",
    realm: "{{.OAuth.Realm}}",
    appName: "{{.OAuth.AppName}}"
  })
  {{end}}

  window.ui = ui
}
</script>
</body>

</html>
`
