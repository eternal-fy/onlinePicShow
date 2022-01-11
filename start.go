package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"onlineProgram/fileHandle"
	"os"
	"path"
	"strconv"
)

const (
	root = "others"
)

var r *gin.Engine

type Data struct {
	Path string `uri:"path" binding:"required"`
}

func main() {
	r = gin.Default()
	r.StaticFS("/static", http.Dir("./")) //前面设置如果要访问静态页面就加前缀，后面定义静态文件的包
	first()
	second()
	third()
	forth()
	r.Run(":80")

}
func first() {
	r.GET("/", func(c *gin.Context) {
		ipContext := "  remoteAddr: " + c.Request.RemoteAddr
		MyLog(ipContext)
		c.Redirect(http.StatusMovedPermanently, "/static/config/you.html")
	})
}

func second() {
	r.GET("/put", func(c *gin.Context) {
		ipContext := "  remoteAddr: " + c.Request.RemoteAddr + " [putPictures] "
		MyLog(ipContext)
		c.Redirect(http.StatusMovedPermanently, "/static/config/addPic.html")

	})

}
func third() {
	r.POST("/upload", func(c *gin.Context) {
		r.MaxMultipartMemory = 8 << 20
		pathName := c.PostForm("pathName")
		if pathName == "" {
			c.String(500, "路径名称为空")
			return
		}
		rootPath := "./" + root
		_, err := os.Stat(rootPath)
		if err != nil {
			os.Mkdir(rootPath, 0777)
		}
		proPath := rootPath + "/" + pathName
		_, err1 := os.Stat(proPath)
		if err1 != nil {
			os.Mkdir(proPath, 0777)
		}

		fileHandle.CopyDir("./config", proPath)

		for i := 0; i < 9; i++ {
			file, err := c.FormFile("file" + strconv.Itoa(i))
			if err != nil {
				continue
			}
			filename := file.Filename
			ext := path.Ext(filename)
			if flag := fileHandle.HandleEr(c, ext, err); flag {
				continue
			}
			if i > 6 {
				c.SaveUploadedFile(file, proPath+"/pic/box/"+strconv.Itoa(i-6)+".jpg")
			} else {
				c.SaveUploadedFile(file, proPath+"/pic/"+strconv.Itoa(i)+".jpg")
			}

		}
		c.Redirect(http.StatusMovedPermanently, "/"+pathName)

	})
}
func forth() {
	r.GET("/:path", func(c *gin.Context) {
		data := new(Data)
		err := c.ShouldBindUri(&data)
		if err != nil {
			c.String(500, "路径错误")
		}
		c.Redirect(http.StatusMovedPermanently, "/static/"+root+"/"+data.Path+"/you.html")
	})

}

func MyLog(ip string) {
	log.SetPrefix("[ld-logger]")
	log.SetFlags(log.Llongfile | log.Lshortfile | log.Ltime | log.Ldate)
	logFile, err := os.OpenFile("./program.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	log.SetOutput(logFile)
	log.Println(ip)

}
