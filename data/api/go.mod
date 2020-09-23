module fullRssAPI

go 1.15

replace gopkg.in/urfave/cli.v2 => github.com/urfave/cli/v2 v2.2.0

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/lib/pq v1.8.0
	github.com/gin-contrib/cors v1.3.1
	github.com/satori/go.uuid v1.2.0
	github.com/antchfx/htmlquery v1.2.3
	github.com/gorilla/feeds v1.1.1
	github.com/mmcdole/gofeed v1.1.0
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd
)
