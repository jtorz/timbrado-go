package timbox

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"net/http"
	"text/template"

	"github.com/jtorz/timbrado-golang/timbrado"
)

const url = "https://staging.ws.timbox.com.mx/timbrado_cfdi33/action"
const soapAction = "timbrar_cfdi"

var tmpl *template.Template
var usr string
var pass string

func init() {
	var err error
	tmpl, err = template.New("soapMsg").Parse(`<?xml version="1.0" encoding="UTF-8" standalone="no"?><SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="urn:WashOut" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:soap-enc="http://schemas.xmlsoap.org/soap/encoding/" ><SOAP-ENV:Body><mns:timbrar_cfdi xmlns:mns="urn:WashOut" SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><username xsi:type="xsd:string">{{.user}}</username><password xsi:type="xsd:string">{{.pass}}</password><sxml xsi:type="xsd:string">{{.cfdi}}</sxml></mns:timbrar_cfdi></SOAP-ENV:Body></SOAP-ENV:Envelope>`)
	if err != nil {
		panic(fmt.Errorf("could parse soap message template %w", err))
	}
}

// WS implementa la interface SoapWS
type WS struct {
	User string
	Pass string
}

// Configure utilizada para configurar el web service.
func (ws WS) Configure(c timbrado.Conf) error {
	if c.User == "" {
		return errors.New("no se proporciono el usuario")
	}
	if c.Pass == "" {
		return errors.New("no se proporciono la contraseña")
	}
	pass = c.Pass
	usr = c.User
	return nil
}

// GenerateMessage genera el mensaje SOAP Envelope  que se enviara como body en la peticion.
func (ws WS) GenerateMessage(cfdi []byte) ([]byte, error) {
	sb := &bytes.Buffer{}
	err := tmpl.ExecuteTemplate(sb, "soapMsg", map[string]string{
		"user": usr,
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
	env := &soapResponse{}
	err := xml.Unmarshal(responseBody, &env)
	if err != nil {
		return timbrado.Response{}, err
	}
	if env.SoapBody.Fault != nil {
		return timbrado.Response{
			CFDI:       nil,
			Message:    env.SoapBody.Fault.Msg,
			StatusCode: env.SoapBody.Fault.Code,
		}, errors.New(env.SoapBody.Fault.Msg)
	}
	cfdi := html.UnescapeString(env.SoapBody.Resp.Res.XML)

	return timbrado.Response{
		CFDI:       []byte(cfdi),
		Message:    "",
		StatusCode: "",
	}, nil
}
