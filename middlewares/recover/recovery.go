package recover

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"zweb"
)

func Recovery() zweb.HandlerFunc {
	return func(ctx *zweb.Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				ctx.Fail(http.StatusInternalServerError, "Internal server error")
			}
		}()

		ctx.Next()
	}
}

func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller Callers 用来返回调用栈的程序计数器, 第 0 个 Caller 是 Callers 本身
	// ，第 1 个是上一层 trace，第 2 个是再上一层的 defer func。因此，为了日志简洁一点，我们跳过了前 3 个 Caller。

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")

	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)   // runtime.FuncForPC(pc) 获取对应的函数
		file, line := fn.FileLine(pc) //通过 fn.FileLine(pc) 获取到调用该函数的文件名和行号，打印在日志中。
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
