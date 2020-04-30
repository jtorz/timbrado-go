package facturehoy

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net/http"
	"text/template"

	"github.com/jtorz/timbrado-golang/timbrado"
)

const url = "http://pruebasclientes.facturehoy.com:8080/CFDI33/WsEmisionTimbrado33"
const soapAction = "urn:EmitirTimbrar"
const servicioID = "5906390"

var tmpl *template.Template
var usr = "AAA010101AAA.Test.User"
var pass = "Prueba$1"

//https://www3.facturehoy.com/home/ambiente-de-pruebas-cfdi33/
/*
 Usuario = AAA010101AAA.Test.User
 Contraseña = Prueba$1
 Id Servicio = 5906390 ->   Facturehoy solo timbrado
 Id Servicio = 36424534 ->  Facturehoy sella y timbra.
 Id Servicio = 36424640 ->  Facturehoy sella y timbra. Facturehoy sella y timbra. (Sólo plantilla TXT. Da clic y descarga Plantilla y Ejemplos TXT)
                            https://www.dropbox.com/s/6zd9af521q446c3/Ejemplos%20y%20Plantilla%20TXT%20CFDI%203.3.zip?dl=0
*/
func init() {
	var err error
	tmpl, err = template.New("soapMsg").Parse(`<?xml version="1.0" encoding="UTF-8" standalone="no"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" xmlns:wsp="http://www.w3.org/ns/ws-policy" xmlns:wsp1_2="http://schemas.xmlsoap.org/ws/2004/09/policy" xmlns:wsam="http://www.w3.org/2007/05/addressing/metadata" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:tns="http://cfdi.ws2.facturehoy.certus.com/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"><SOAP-ENV:Body><mns1:EmitirTimbrar xmlns:mns1="http://cfdi.ws2.facturehoy.certus.com/"><usuario>{{.usuario}}</usuario><contrasenia>{{.contrasenia}}</contrasenia><idServicio>{{.idServicio}}</idServicio><xml>{{.xml}}</xml></mns1:EmitirTimbrar></SOAP-ENV:Body></SOAP-ENV:Envelope>`)
	if err != nil {
		panic(fmt.Errorf("could parse soap message template %w", err))
	}
}

// WS implementa la interface SoapWS
type WS struct{}

// Configure utilizada para configurar el web service.
func (ws WS) Configure(c timbrado.Conf) error {
	if c.User == "" {
		fmt.Println("no se proporciono el usuario utilizando usuario de pruebas")
	} else {
		usr = c.User
	}
	if c.Pass == "" {
		fmt.Println("no se proporciono la contraseña utilizando contraseña de pruebas")
	} else {
		pass = c.Pass
	}
	return nil
}

// GenerateMessage genera el mensaje SOAP Envelope  que se enviara como body en la peticion.
func (ws WS) GenerateMessage(cfdi []byte) ([]byte, error) {
	sb := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(sb, "soapMsg", map[string]string{
		"usuario":     usr,
		"contrasenia": pass,
		"idServicio":  servicioID,
		"xml":         base64.StdEncoding.EncodeToString(cfdi),
	})
	if err != nil {
		return nil, err
	}
	return sb.Bytes(), nil
}

// URL ruta del WS.
func (ws WS) URL() string {
	return url
}

// Method de la peticion http.
func (ws WS) Method() string {
	return "POST"
}

// ConfigureReq permite configurar datos adicionales (como headers) en la peticion.
func (ws WS) ConfigureReq(req *http.Request) error {
	req.Header.Set("SOAPAction", soapAction)
	return nil
}

// ParseResponse tranforma la respuesta de WS.
func (ws WS) ParseResponse(responseBody []byte) (timbrado.Response, error) {
	wsRes := &soapResponse{}
	err := xml.Unmarshal(responseBody, &wsRes)
	if err != nil {
		return timbrado.Response{}, err
	}

	cfdi, err := base64.StdEncoding.DecodeString(wsRes.SoapBody.Resp.Return.XML)
	if err != nil {
		return timbrado.Response{}, err
	}

	return timbrado.Response{
		CFDI:       cfdi,
		Message:    wsRes.SoapBody.Resp.Return.Message,
		StatusCode: wsRes.SoapBody.Resp.Return.CodigoError,
	}, nil
}
