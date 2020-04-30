package solucionfactible

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/jtorz/timbrado-golang/timbrado"
)

const url = "https://testing.solucionfactible.com/ws/services/Timbrado"
const soapAction = "urn:timbrar"

var tmpl *template.Template
var usr = "testing@solucionfactible.com"
var pass = "timbrado.SF.16672"

// https://solucionfactible.com/sfic/capitulos/timbrado/ws-timbrado.jsp
/*
 Metodos
 timbrar: Recibe y timbra uno o varios CFDI.
 cancelar: Recibe el folio fiscal (UUID) de uno o varios CFDI y los cancela.
*/

func init() {
	var err error
	tmpl, err = template.New("soapMsg").Parse(`<?xml version="1.0" encoding="UTF-8"?>
	<env:Envelope xmlns:env="http://www.w3.org/2003/05/soap-envelope"><env:Body><ns:timbrar xmlns:ns="http://timbrado.ws.cfdi.solucionfactible.com"><ns:usuario>{{.usr}}</ns:usuario><ns:password>{{.pass}}</ns:password><ns:cfdi>{{.cfdi}}</ns:cfdi><ns:zip>false</ns:zip></ns:timbrar></env:Body></env:Envelope>`)
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
		"usr":  usr,
		"pass": pass,
		"cfdi": base64.StdEncoding.EncodeToString(cfdi),
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

	if wsRes.SoapBody.Resp.Return.Status != 200 {
		return timbrado.Response{}, errors.New(wsRes.SoapBody.Resp.Return.Mensaje)
	}

	cfdi, err := base64.StdEncoding.DecodeString(wsRes.SoapBody.Resp.Return.Result.CfdiTimbrado)
	if err != nil {
		return timbrado.Response{}, err
	}

	return timbrado.Response{
		CFDI:       cfdi,
		Message:    wsRes.SoapBody.Resp.Return.Result.Mensaje,
		StatusCode: wsRes.SoapBody.Resp.Return.Result.Status,
	}, nil
}
