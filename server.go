package squareroot

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	readTimeout  = time.Second * 15
	writeTimeout = time.Second * 15
)

func Boot(cert tls.Certificate, httpPort, httpsPort int, handler http.Handler) error {
	srvr := &http.Server{
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         fmt.Sprintf(":%v", httpsPort),
	}

	srvr.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	srvr.Handler = handler

	err := srvr.ListenAndServeTLS("", "")

	if err != nil {
		return err
	}

	//consider go func to tls redirect
	return http.ListenAndServe(fmt.Sprintf(":%v", httpPort), http.HandlerFunc(redirectTLS))
}

func GetCertified(pubPath, privPath string) tls.Certificate {
	pubfile, err := ioutil.ReadFile(pubPath)

	if err != nil {
		panic(err)
	}

	privatefile, err := ioutil.ReadFile(privPath)

	if err != nil {
		panic(err)
	}

	cert, err := tls.X509KeyPair(pubfile, privatefile)

	if err != nil {
		panic(err)
	}

	return cert
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	moveURL := fmt.Sprintf("https://%s%s", r.Host, r.RequestURI)
	http.Redirect(w, r, moveURL, http.StatusPermanentRedirect)
}
