package cfdi

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jtorz/timbrado-golang/config"
	"github.com/youmark/pkcs8"
)

// Cert cadena del certificado del SAT.
var Cert string

// NoCert cadena del numero de certificado del SAT.
var NoCert string

// X509Cert certificado del SAT.
var X509Cert *x509.Certificate

// Key llave del certificado del SAT.
var Key *rsa.PrivateKey

// LoadCert carga el certificado.
func LoadCert(certFile, keyFile string, pass []byte) error {
	cf, err := ioutil.ReadFile(certFile)
	if err != nil {
		return err
	}

	kf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return err
	}

	SetCert(cf)
	if err != nil {
		return err
	}
	return SetKey(kf, pass)
}

// SetCert configura la informacion del certificado.
func SetCert(cf []byte) (err error) {
	X509Cert, err = x509.ParseCertificate(cf)
	if err != nil {
		return err
	}

	Cert = base64.StdEncoding.EncodeToString(cf)
	decodeNoCert()
	return nil
}

// SetKey configura la llave del certificado.
func SetKey(kf []byte, pass []byte) (err error) {
	v, err := pkcs8.ParsePKCS8PrivateKey(kf, pass)
	if err != nil {
		return err
	}
	var ok bool
	Key, ok = v.(*rsa.PrivateKey)
	if !ok {
		return errors.New("not a rsa.privatekey")
	}
	return nil
}

func decodeNoCert() {
	s := fmt.Sprintf("%x", X509Cert.SerialNumber)
	sb := strings.Builder{}
	sb.Grow(len(s) / 2)
	for i, d := range s {
		if i%2 == 1 {
			sb.WriteRune(d)
		}
	}
	NoCert = sb.String()
}

// Sellar genera e inyecta el Sello, Certificado y NoCertificado en el archivo cfdi.
func Sellar(cfdiFile string) (cfdiFileSellado string, err error) {
	cfdiFileSellado = getCfdiFileName(cfdiFile)
	cfdiRaw, err := ioutil.ReadFile(cfdiFile)
	if err != nil {
		return
	}

	cfdiRaw = bytes.Replace(cfdiRaw,
		[]byte("<cfdi:Comprobante"),
		[]byte(`<cfdi:Comprobante Certificado="`+Cert+`" NoCertificado="`+NoCert+`"`),
		1,
	)
	err = ioutil.WriteFile(cfdiFileSellado, cfdiRaw, 0644)

	sello, err := generarSello(cfdiFileSellado)
	if err != nil {
		return
	}
	//sello = ""
	cfdiRaw = bytes.Replace(cfdiRaw,
		[]byte("<cfdi:Comprobante"),
		[]byte(`<cfdi:Comprobante Sello="`+sello+`"`),
		1,
	)
	err = ioutil.WriteFile(cfdiFileSellado, cfdiRaw, 0644)
	return
}

func getCfdiFileName(cfdiFile string) string {
	extension := filepath.Ext(cfdiFile)
	name := cfdiFile[0 : len(cfdiFile)-len(extension)]
	return name + "_sellado" + extension
}

func generarSello(cfdiFile string) (string, error) {
	co, err := GeneraCadenaOriginal(cfdiFile)
	if err != nil {
		return "", fmt.Errorf("no se puedo generar la cadena original: %w", err)
	}
	s, err := Digest(co)
	log.Printf("Cadena original: `%s`", string(co))
	if err != nil {
		return "", fmt.Errorf("no se pudo digerir %w", err)
	}
	return base64.StdEncoding.EncodeToString(s), nil
}

// GeneraCadenaOriginal genera la cadena original con el archivo xslt.
func GeneraCadenaOriginal(cfdiFile string) ([]byte, error) {
	cmd := exec.Cmd{
		Args: []string{"xsltproc", config.XSLTPath, cfdiFile},
		Env:  os.Environ(),
		Path: config.CMDxsltproc,
	}
	return cmd.Output()
}

// Digest genera el hash(sha 256) de la cadena original y lo encripta
func Digest(cadenaOriginal []byte) ([]byte, error) {
	h := sha256.Sum256(cadenaOriginal)
	return rsa.SignPKCS1v15(rand.Reader, Key, crypto.SHA256, h[:])
}
