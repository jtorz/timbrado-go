package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jtorz/timbrado-golang/cfdi"
	"github.com/jtorz/timbrado-golang/timbrado"
	"github.com/jtorz/timbrado-golang/timbrado/facturehoy"
	"github.com/jtorz/timbrado-golang/timbrado/solucionfactible"
	"github.com/jtorz/timbrado-golang/timbrado/timbox"
)

const carpetaTimbrado = "c:/timbrado"

func init() {
	_, err := os.Stat(carpetaTimbrado)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(carpetaTimbrado, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}
}
func main() {

	keyPath := flag.String("key", "assets/certificate.key", "ruta donde se encuentra el archivo .key del certificado del SAT")
	certPath := flag.String("cert", "assets/certificate.cer", "ruta donde se encuentra el archivo .cer del certificado del SAT")
	certpass := flag.String("certpass", "12345678a", "password de la llave del certificado del sat")
	err := cfdi.LoadCert(*certPath, *keyPath, []byte(*certpass))
	if err != nil {
		fmt.Printf("No se puedo cargar certificado %s\n\t %s\n", *certPath, err)
	}

	r := gin.Default()

	r.Use(static.Serve("/", static.LocalFile("./client/dist", true)))
	api(r.Group("/api"))
	r.Run(":3000")
}

func api(r *gin.RouterGroup) {
	r.GET("/cert", getCert)
	r.POST("/cert", setCert)
	r.POST("/cert/key", setKey)
	r.GET("/webservices", catalogoWS)
	r.GET("/webservices/ws", getWS)
	r.POST("/webservices/ws", setWS)
	r.POST("/timbrar", timbrar)
}

func getCert(c *gin.Context) {
	if cfdi.X509Cert == nil {
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"serialNumber": cfdi.NoCert,
		"emisor": gin.H{
			"o":  cfdi.X509Cert.Issuer.Organization,
			"ou": cfdi.X509Cert.Issuer.OrganizationalUnit,
			"cn": cfdi.X509Cert.Issuer.CommonName,
		},
		"receptor": gin.H{
			"o":  cfdi.X509Cert.Subject.Organization,
			"ou": cfdi.X509Cert.Subject.OrganizationalUnit,
			"cn": cfdi.X509Cert.Subject.CommonName,
		},
	})
}

func setCert(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer src.Close()
	cf, err := ioutil.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = cfdi.SetCert(cf)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	getCert(c)
}

func setKey(c *gin.Context) {
	file, err := c.FormFile("file")
	pass := c.PostForm("pass")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer src.Close()
	cf, err := ioutil.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = cfdi.SetKey(cf, []byte(pass))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

type ws struct {
	ID   string `json:"id" binding:"required"`
	Name string `json:"name"`
}

type wsauth struct {
	WS       ws     `json:"ws" binding:"required"`
	Usuario  string `json:"usuario"`
	Password string `json:"password"`
}

var catWS = []ws{
	{ID: "timbox", Name: "Timbox"},
	{ID: "facturehoy", Name: "Facturehoy"},
	{ID: "solucionfactible", Name: "Solucionfactible"},
}

var currentws = wsauth{WS: catWS[1], Usuario: "", Password: ""}

func catalogoWS(c *gin.Context) {
	c.JSON(http.StatusOK, catWS)
}

func getWS(c *gin.Context) {
	c.JSON(http.StatusOK, currentws)
}

func setWS(c *gin.Context) {
	wsauth := wsauth{}
	err := c.Bind(&wsauth)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	w, err := findWS(wsauth.WS.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	wsauth.WS = w
	currentws = wsauth
	time.Sleep(2 * time.Second)
	c.Status(http.StatusOK)
}

func findWS(w string) (ws, error) {
	for i := range catWS {
		if w == catWS[i].ID {
			return catWS[i], nil
		}
	}
	return ws{}, fmt.Errorf("ws service %s no encontrado", w)
}

func timbrar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer src.Close()
	cfdiPath := carpetaTimbrado + "/" + time.Now().Format("2006_01_02T15_04_05")

	f, err := os.OpenFile(cfdiPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	_, err = io.Copy(f, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	err = f.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ws := selectWS()
	cfdiFile, err := cfdi.Sellar(cfdiPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error)
		return
	}

	r, err := timbrado.TimbrarSOAP(ws, cfdiFile, currentws.Usuario, currentws.Password)
	if err != nil {
		if r.Message != "" {
			c.JSON(http.StatusInternalServerError, r)
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"StatusCode": r.StatusCode,
		"Message":    r.Message,
		"CFDI":       string(r.CFDI),
	})
}

func selectWS() timbrado.WS {
	switch currentws.WS.ID {
	case "facturehoy":
		return facturehoy.WS{}
	case "solucionfactible":
		return solucionfactible.WS{}
	case "timbox":
		return timbox.WS{}
	default:
		panic(fmt.Sprintf("web service %s desconocido", currentws.WS.ID))
	}
}
