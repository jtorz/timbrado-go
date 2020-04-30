package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jtorz/timbrado-golang/cfdi"
	"github.com/jtorz/timbrado-golang/timbrado"
	"github.com/jtorz/timbrado-golang/timbrado/facturehoy"
	"github.com/jtorz/timbrado-golang/timbrado/solucionfactible"
	"github.com/jtorz/timbrado-golang/timbrado/timbox"
)

func main() {
	keyPath := flag.String("key", "assets/certificate.key", "ruta donde se encuentra el archivo .key del certificado del SAT")
	certPath := flag.String("cert", "assets/certificate.cer", "ruta donde se encuentra el archivo .cer del certificado del SAT")
	certpass := flag.String("certpass", "12345678a", "password de la llave del certificado del sat")
	webServ := flag.String("ws", "facturehoy", "web service a utilizar: facturehoy,solucionfactible,timbox")
	file := flag.String("file", "test/original.xml", "archivo cfdi a sellar y timbrar")

	wsUser := flag.String("wsuser", "", "usuario de autenticacion en ws")
	wsPass := flag.String("wsPass", "", "password de autenticacion en ws")
	flag.Parse()

	err := cfdi.LoadCert(*certPath, *keyPath, []byte(*certpass))
	if err != nil {
		fmt.Printf("No se puedo cargar certificado %s\n\t %s\n", *certPath, err)
		os.Exit(1)
	}

	cfdiFile, err := cfdi.Sellar(*file)
	if err != nil {
		fmt.Printf("No se pudo sellar xml\n\t %s\n", err)
		os.Exit(1)
	}
	var ws timbrado.WS = timbox.WS{}
	switch *webServ {
	case "facturehoy":
		fmt.Println("Utilizando ws de facturehoy")
		ws = facturehoy.WS{}
	case "solucionfactible":
		fmt.Println("Utilizando ws de solucionfactible")
		ws = solucionfactible.WS{}
	default:
		fmt.Println("Utilizando ws de timbox")
	}
	r, err := timbrado.TimbrarSOAP(ws, cfdiFile, *wsUser, *wsPass)
	if err != nil {
		fmt.Printf("No se pudo timbrar\n\tEstado (%s) %s\n", r.StatusCode, err)
		os.Exit(1)
	}
	fmt.Printf("Terminado.\nEstado: %s\n%s", r.StatusCode, r.Message)
}
