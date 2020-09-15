package urlshorten

import (
	"log"
	"os"
)

// 被导入时自动执行：main()中导入"urlShorten/urlshorten"时会执行包下面全部的init()
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, _ := os.OpenFile("urlshorten.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	log.SetOutput(f)
}
