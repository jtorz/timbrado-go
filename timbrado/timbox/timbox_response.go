package timbox

import "encoding/xml"

/*
<soap:Envelope>
  <soap:Body>
    <tns:timbrar_cfdi_response>
      <timbrar_cfdi_result>
        <xml>

<soap:Envelope>
  <soap:Body>
    <soap:Fault>
      <faultcode>
      <faultstring>
*/

type soapResponse struct {
	XMLName  xml.Name  `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	SoapBody *soapBody `xml:"Body"`
}

type soapBody struct {
	Resp  *response `xml:"urn:WashOut timbrar_cfdi_response"`
	Fault *fault    `xml:"Fault"`
}

type fault struct {
	Code string `xml:"faultcode"`
	Msg  string `xml:"faultstring"`
}

type response struct {
	Res *result `xml:"timbrar_cfdi_result"`
}

type result struct {
	XML string `xml:"xml"`
}
