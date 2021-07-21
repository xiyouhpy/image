module github.com/xiyouhpy/image

go 1.16

replace github.com/xiyouhpy/image => ./

require (
	github.com/garyburd/redigo v1.6.2
	github.com/gin-gonic/gin v1.7.2
	github.com/go-sql-driver/mysql v1.6.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d // indirect
)
