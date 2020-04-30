package timbrado

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// Conf contiene informacion de configuracion del webservice.
type Conf struct {
	User string
	Pass string
}

// WS interface utilizada para implementar el servicio de timbrado.
type WS interface {
	// Configure utilizada para configurar el web service.
	Configure(Conf) error
	// GenerateMessage genera el mensaje SOAP Envelope  que se enviara como body en la peticion.
	GenerateMessage(cfdi []byte) ([]byte, error)
	// URL ruta del WS.
	URL() string
	// Method de la peticion http.
	Method() string
	// ConfigureReq permite configurar datos adicionales (como headers) en la peticion.
	ConfigureReq(*http.Request) error
	// ParseResponse tranforma la respuesta de WS.
	ParseResponse(responseBody []byte) (Response, error)
}

// Response contiene la informacion de la respuesta del cfdi.
type Response struct {
	// StatusCode codigo de respuesta del webservice.
	StatusCode string
	// StatusCode mensaje de validacion del webservice.
	Message string
	CFDI    []byte
}

// TimbrarSOAP funcion que manda a llamar el webservice de timbrado.
func TimbrarSOAP(ws WS, cfdiFile, wsUser, wsPass string) (r Response, err error) {
	err = ws.Configure(Conf{User: wsUser, Pass: wsPass})
	if err != nil {
		return
	}
	b, err := ioutil.ReadFile(cfdiFile)
	if err != nil {
		return r, fmt.Errorf("can't read CFDI file: %w", err)
	}
	msg, err := ws.GenerateMessage(b)
	if err != nil {
		return
	}

	req, err := http.NewRequest(ws.Method(), ws.URL(), bytes.NewReader(msg))
	if err != nil {
		return
	}

	err = ws.ConfigureReq(req)
	if err != nil {
		return
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("codigo http de timbrado distinto de 200 (%v)\n", res.Status)
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return r, fmt.Errorf("can't read http.Response: %w", err)
	}

	name := rmvExt(cfdiFile)
	ioutil.WriteFile(name+"_response.xml", b, 0644)

	r, err = ws.ParseResponse(b)
	if err != nil {
		return
	}

	ioutil.WriteFile(name+"_timbrado.xml", r.CFDI, 0644)
	return
}

func rmvExt(cfdiFile string) string {
	extension := filepath.Ext(cfdiFile)
	name := cfdiFile[0 : len(cfdiFile)-len(extension)]
	return name
}
