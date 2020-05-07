package config

import (
	"fmt"
	"os"
	"runtime"
)

// ServerPort puerto de escucha del servidor web.
var ServerPort = os.Getenv("TIMBRADO_PORT")

// RootDir ruta absoluta del directorio donde se encuentra el proyecto.
var RootDir = os.Getenv("TIMBRADO_ROOT")

// XMLOutDir ruta del directorio donde se almacenan los archivos xml de timbrado.
var XMLOutDir = os.Getenv("TIMBRADO_OUT")

// XSLTPath ruta donde se encuentra el archivo xslt de la cadena original.
var XSLTPath = RootDir + "/assets/xslt/cadenaoriginal_3_3.xslt"

// ClientDist ruta del directorio donde se encuentra los archivos compilados del front-end.
var ClientDist = RootDir + "/client-dist"

// CMDxsltproc comando xsltproc a ejecutar.
var CMDxsltproc string

func init() {
	if runtime.GOOS == "windows" {
		CMDxsltproc = RootDir + "/tools/xslt/bin/xsltproc.exe"
	} else {
		CMDxsltproc = RootDir + "/tools/xslt/bin/xsltproc"
	}

	fmt.Printf("ServerPort: %s\n", ServerPort)
	fmt.Printf("RootDir: %s\n", RootDir)
	fmt.Printf("XMLOutDir: %s\n", XMLOutDir)
	fmt.Printf("XSLTPath: %s\n", XSLTPath)
	fmt.Printf("ClientDist: %s\n", ClientDist)
	fmt.Printf("CMDxsltproc: %s\n", CMDxsltproc)
}
