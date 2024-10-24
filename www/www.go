package www

import "embed"

// FS is our static web server content.
//go:embed *.html *.js *.css
var FS embed.FS