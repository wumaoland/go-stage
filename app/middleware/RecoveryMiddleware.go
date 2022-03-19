package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
	"stage/app/common"
	"stage/config/conf"
	"strings"
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var message string
			if _, ok := err.(int); ok {
				message = conf.Message(err.(int))
			} else {
				message = fmt.Sprintf("%s", err)
			}
			for _, v := range writer() {
				log.New(v, "", log.LstdFlags).Printf("%s\n\n", trace(message))
				log.New(v, "", log.LstdFlags).Printf("%s", "----------------------------------------------------------------------")
			}
			common.NewDisplay(c).Show(err)
		}
	}()
	c.Next()
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
