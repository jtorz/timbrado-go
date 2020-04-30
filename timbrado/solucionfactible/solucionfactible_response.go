package solucionfactible

import "encoding/xml"

type soapResponse struct {
	XMLName  xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	SoapBody *soapBodyResponse
}

type soapBodyResponse struct {
	XMLName xml.Name `xml:"Body"`
	Resp    *responseBody
}
type responseBody struct {
	XMLName xml.Name `xml:"http://timbrado.ws.cfdi.solucionfactible.com timbrarResponse"`
	Return  *wsRes
}

type wsRes struct {
	XMLName xml.Name `xml:"return"`
	//CfdiTimbrado string   `xml:"http://timbrado.ws.cfdi.solucionfactible.com/xsd cfdiTimbrado"`
	Mensaje string `xml:"http://timbrado.ws.cfdi.solucionfactible.com/xsd mensaje"`
	Status  int    `xml:"http://timbrado.ws.cfdi.solucionfactible.com/xsd status"`
	Result  *resultados
}

type resultados struct {
	XMLName      xml.Name `xml:"http://timbrado.ws.cfdi.solucionfactible.com/xsd resultados"`
	Status       string   `xml:"status"`
	Mensaje      string   `xml:"mensaje"`
	CfdiTimbrado string   `xml:"cfdiTimbrado"`
	/* CadenaOriginal string   `xml:"cadenaOriginal"`
	CertificadoSAT string   `xml:"certificadoSAT"`
	FechaTimbrado  string   `xml:"fechaTimbrado"`
	QrCode         string   `xml:"qrCode"`
	SelloSAT       string   `xml:"selloSAT"`
	UUID           string   `xml:"uuid"`
	VersionTFD     string   `xml:"versionTFD"` */
}
