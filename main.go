package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/big"
	"net/http"
	"log"
)

var _n = "t76eqou0a5rdnij8mixvut88ca8lk9y513x5xxj6dcihqar4buvofcyu7ymbt84rllt9awioc7b1ra3ridjxt4ofkzm1i2iriaz"
var _d = "ekwnx0sfrsk6fy9vsgvgp9yqz1r4vzcnnmb6efvn1klz22e614qfwh41r8sch3fvwmp9yc2t2zadz4ip7eodiagcffxovccch5t"
var pkey *rsa.PrivateKey

func init() {
	n, _ := new(big.Int).SetString(_n, 36)
	d, _ := new(big.Int).SetString(_d, 36)
	p, _ := new(big.Int).SetString("526cqefo0m5h5qydu25d301g8ntrhmbvq7om4ldki3kgmcnhrd", 36)
	q, _ := new(big.Int).SetString("5rq25sorip2lewym841g6h2qx68oq38n8pdxvpb88tpiq4b1mb", 36)
	pkey = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: 65537,
		},
		D:      d,
		Primes: []*big.Int{p, q},
	}
}
func main() {
	listen := ":8089"
	http.HandleFunc("/rpc/obtainTicket.action", ActivateIdea)
	http.HandleFunc("/rpc/releaseTicket.action", ActivateIdea)
	http.HandleFunc("/rpc/prolongTicket.action", ActivateIdea)
	http.HandleFunc("//rpc/obtainTicket.action", ActivateIdea)
	http.HandleFunc("//rpc/releaseTicket.action", ActivateIdea)
	http.HandleFunc("//rpc/prolongTicket.action", ActivateIdea)
	log.Println("ListenAndServe:"+listen, "server start")
	err := http.ListenAndServe(listen, nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func ActivateIdea(writer http.ResponseWriter, request *http.Request) {

	ip := request.Header.Get("Remote_addr")
	if ip == "" {
		ip = request.RemoteAddr
	}
	log.Println(request.Method, ip, request.URL)
	salt := request.FormValue("salt")
	userName := request.FormValue("userName")
	ticket := fmt.Sprintf(`<ObtainTicketResponse>
	<message></message>
	<prolongationPeriod>86400000</prolongationPeriod>
	<responseCode>OK</responseCode>
	<salt>%s</salt>
	<ticketId>1</ticketId>
	<ticketProperties>licensee=%s	licenseType=0	</ticketProperties>
</ObtainTicketResponse>`, salt, userName)
	hashed := md5.Sum([]byte(ticket))
	sig, _ := rsa.SignPKCS1v15(rand.Reader, pkey, crypto.MD5, hashed[:])
	resp := fmt.Sprintf("<!-- %x -->\n%s", sig, ticket)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(writer, resp)
}
