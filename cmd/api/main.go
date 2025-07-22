package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sehogas/goarca/afip"
	"github.com/sehogas/goarca/cmd/api/docs"
	"github.com/sehogas/goarca/internal/middleware"
	"github.com/sehogas/goarca/internal/util"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	Version string = "development"

	Wscoem      *afip.Wscoem
	Wscoemcons  *afip.Wscoemcons
	Wsgestabref *afip.Wsgestabref
	Wsfe        *afip.Wsfe

	validate *validator.Validate
)

//	@title			API proxy a los webservices de ARCA
//	@version		1.0
//	@description	Esta API Json Rest actua como proxy SOAP a los servicios web de ARCA.
//	@termsOfService	http://swagger.io/terms/

// @contact.name	Sebastian Hogas
// @contact.email	sehogas@gmail.com
func main() {
	godotenv.Load()

	environment := afip.TESTING
	if strings.ToLower(strings.TrimSpace(os.Getenv("PROD"))) == "true" {
		environment = afip.PRODUCTION
	}

	cuit, err := strconv.ParseInt(os.Getenv("CUIT"), 10, 64)
	if err != nil {
		log.Fatalln("variable de entorno CUIT faltante o no numérica")
	}

	if os.Getenv("PRIVATE_KEY_FILE") == "" {
		log.Fatalln("variable de entorno PRIVATE_KEY_FILE faltante")
	}

	if os.Getenv("CERTIFICATE_FILE") == "" {
		log.Fatalln("falta variable de entorno CERTIFICATE_FILE")
	}

	if len(os.Getenv("WSCOEM_TIPO_AGENTE")) != 4 {
		log.Fatalln("falta variable de entorno WSCOEM_TIPO_AGENTE")
	}

	if len(os.Getenv("WSCOEM_ROL")) != 4 {
		log.Fatalln("variable de entorno WSCOEM_ROL faltante o inválida")
	}

	if len(os.Getenv("WGESTABREF_TIPO_AGENTE")) != 4 {
		log.Fatalln("falta variable de entorno WGESTABREF_TIPO_AGENTE")
	}

	if len(os.Getenv("WGESTABREF_ROL")) != 4 {
		log.Fatalln("variable de entorno WGESTABREF_ROL faltante o inválida")
	}

	if len(os.Getenv("KEYS_FILE")) == 0 {
		log.Fatalln("Falta variable de entorno KEYS_FILE")
	}

	port := 4433
	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatalln("variable de entorno PORT no numérica. ", err)
		}
	}

	Wscoem, err = afip.NewWscoem(environment, cuit, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	Wscoemcons, err = afip.NewWscoemcons(environment, cuit, os.Getenv("WSCOEM_TIPO_AGENTE"), os.Getenv("WSCOEM_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	Wsgestabref, err = afip.Newgestabref(environment, cuit, os.Getenv("WGESTABREF_TIPO_AGENTE"), os.Getenv("WGESTABREF_ROL"))
	if err != nil {
		log.Fatalln(err)
	}

	Wsfe, err = afip.NewWsfe(environment, cuit)
	if err != nil {
		log.Fatalln(err)
	}

	/* API Rest */
	validate = validator.New(validator.WithRequiredStructEnabled())

	middlewareCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		Debug:            false,
	})

	middlewareApiKey, err := middleware.NewApiKeyMiddleware(os.Getenv("KEYS_FILE"))
	if err != nil {
		log.Fatalln(err)
	}

	router := http.NewServeMux()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"https"}

	router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	coem := http.NewServeMux()
	coem.HandleFunc("/dummy", DummyCoemHandler)
	coem.HandleFunc("POST /registrar-caratula", RegistrarCaratulaHandler)
	coem.HandleFunc("PUT /rectificar-caratula", RectificarCaratulaHandler)
	coem.HandleFunc("DELETE /anular-caratula", AnularCaratulaHandler)
	coem.HandleFunc("PUT /solicitar-cambio-buque", SolicitarCambioBuqueHandler)
	coem.HandleFunc("PUT /solicitar-cambio-fechas", SolicitarCambioFechasHandler)
	coem.HandleFunc("PUT /solicitar-cambio-lot", SolicitarCambioLOTHandler)
	coem.HandleFunc("POST /registrar-coem", RegistrarCOEMHandler)
	coem.HandleFunc("PUT /rectificar-coem", RectificarCOEMHandler)
	coem.HandleFunc("POST /cerrar-coem", CerrarCOEMHandler)
	coem.HandleFunc("DELETE /anular-coem", AnularCOEMHandler)
	coem.HandleFunc("POST /solicitar-anulacion-coem", SolicitarAnulacionCOEMHandler)
	coem.HandleFunc("POST /solicitar-no-abordo", SolicitarNoABordoHandler)
	coem.HandleFunc("POST /solicitar-cierre-carga-conto-bulto", SolicitarCierreCargaContoBultoHandler)
	coem.HandleFunc("POST /solicitar-cierre-carga-granel", SolicitarCierreCargaGranelHandler)

	coemcons := http.NewServeMux()
	coemcons.HandleFunc("/dummy", DummyCoemconsHandler)
	coemcons.HandleFunc("/obtener-consulta-estados-coem", ObtenerConsultaEstadosCOEMHandler)
	coemcons.HandleFunc("/obtener-consulta-no-abordo", ObtenerConsultaNoAbordoHandler)
	coemcons.HandleFunc("/obtener-consulta-solicitudes", ObtenerConsultaSolicitudesHandler)

	gestabref := http.NewServeMux()
	gestabref.HandleFunc("/dummy", DummyGesTabRefHandler)
	gestabref.HandleFunc("/consultar-fecha-ult-act", ConsultarFechaUltActHandler)
	gestabref.HandleFunc("/lista-arancel", ListaArancelHandler)
	gestabref.HandleFunc("/lista-descripcion", ListaDescripcionHandler)
	gestabref.HandleFunc("/lista-descripcion-decodificacion", ListaDescripcionDecodificacionHandler)
	gestabref.HandleFunc("/lista-empresas", ListaEmpresasHandler)
	gestabref.HandleFunc("/lista-lugares-operativos", ListaLugaresOperativosHandler)
	gestabref.HandleFunc("/lista-paises-aduanas", ListaPaisesAduanasHandler)
	gestabref.HandleFunc("/lista-tablas-referencia", ListaTablasReferenciaHandler)

	fe := http.NewServeMux()
	fe.HandleFunc("/Dummy", FEDummyHandler)
	fe.HandleFunc("/CompUltimoAutorizado", FECompUltimoAutorizadoHandler)
	fe.HandleFunc("/GetTiposCbte", FEParamGetTiposCbteHandler)
	fe.HandleFunc("/GetTiposConcepto", FEParamGetTiposConceptoHandler)
	fe.HandleFunc("/GetTiposDoc", FEParamGetTiposDocHandler)
	fe.HandleFunc("/GetTiposIva", FEParamGetTiposIvaHandler)
	fe.HandleFunc("/GetTiposMonedas", FEParamGetTiposMonedasHandler)
	fe.HandleFunc("/GetTiposOpcional", FEParamGetTiposOpcionalHandler)
	fe.HandleFunc("/GetTiposTributos", FEParamGetTiposTributosHandler)
	fe.HandleFunc("POST /GetPtosVenta", FEParamGetPtosVentaHandler)
	fe.HandleFunc("POST /CAESolicitar", FECAESolicitarHandler)

	v1 := http.NewServeMux()
	v1.HandleFunc("/info", InfoHandler)
	v1.Handle("/coem/", http.StripPrefix("/coem", coem))
	v1.Handle("/coemcons/", http.StripPrefix("/coemcons", coemcons))
	v1.Handle("/gestabref/", http.StripPrefix("/gestabref", gestabref))
	v1.Handle("/fe/", http.StripPrefix("/fe", fe))
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))

	stack := middleware.CreateStack(
		middlewareCors.Handler,
		middleware.Recovery,
		middleware.Logging,
		middlewareApiKey.Handler,
	)

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      stack(router),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	go util.GracefulShutdown(server)

	log.Printf("Starting server on port %v", server.Addr)
	err = server.ListenAndServeTLS("keys/server.crt", "keys/server.key")
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %s", err)
	}
}
