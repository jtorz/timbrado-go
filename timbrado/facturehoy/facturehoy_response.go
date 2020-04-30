package facturehoy

import "encoding/xml"

type soapResponse struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *soapBodyResponse
}

type soapBodyResponse struct {
	XMLName xml.Name `xml:"Body"`
	Resp    *responseBody
}
type responseBody struct {
	XMLName xml.Name `xml:"http://cfdi.ws2.facturehoy.certus.com/ EmitirTimbrarResponse"`
	Return  *wsRes
}

type wsRes struct {
	XMLName     xml.Name `xml:"return"`
	Message     string   `xml:"message"`
	XML         string   `xml:"XML"`
	CodigoError string   `xml:"codigoError"`
	//IsError string   `xml:"isError"`
	//CadenaOriginal        string `xml:"cadenaOriginal"`
	//CadenaOriginalTimbre  string `xml:"cadenaOriginalTimbre"`
	//FechaHoraTimbrado     string `xml:"fechaHoraTimbrado"`
	//FolioUDDI             string `xml:"folioUDDI"`
	//PDF                   string `xml:"PDF"`
	//RutaDescargaPDF       string `xml:"rutaDescargaPDF"`
	//RutaDescargaXML       string `xml:"rutaDescargaXML"`
	//SelloDigitalEmisor    string `xml:"selloDigitalEmisor"`
	//SelloDigitalTimbreSAT string `xml:"selloDigitalTimbreSAT"`
}
