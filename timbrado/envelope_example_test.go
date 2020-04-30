package timbrado

import "encoding/xml"

type soapResponse struct {
	XMLName  xml.Name          `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *soapBodyResponse `xml:"Body"`
}

type soapBodyResponse struct {
	Resp *responseBody `xml:"Timbrado"`
}

type responseBody struct {
	Message string `xml:"message"`
	Code    string `xml:"code"`
	XML     string `xml:"XML"`
}
