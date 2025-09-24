package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/sehogas/goarca/cmd/api/docs"
	"github.com/sehogas/goarca/internal/middleware"
	"github.com/sehogas/goarca/internal/services"
	"github.com/sehogas/goarca/internal/util"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var (
	Version string = "development"

	Wscoem      *services.Wscoem
	Wscoemcons  *services.Wscoemcons
	Wsgestabref *services.Wsgestabref
	Wsfe        *services.Wsfe
)

//	@title			API proxy a los webservices de ARCA
//	@version		1.0
//	@description	Esta API Json Rest actua como proxy SOAP a los servicios web de ARCA.
//	@termsOfService	http://swagger.io/terms/

// @contact.name	Sebastian Hogas
// @contact.email	sehogas@gmail.com
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: util.GetLogLevelFromEnv(),
	}))

	// Create a log.Logger that uses the SlogWriter
	errorLogWriter := util.NewSlogWriter(logger, slog.LevelError)
	serverErrorLogger := log.New(errorLogWriter, "", 0)

	slog.SetDefault(logger)

	godotenv.Load()

	environment := services.TESTING
	if strings.ToLower(strings.TrimSpace(os.Getenv("PROD"))) == "true" {
		environment = services.PRODUCTION
	}

	printXML, err := strconv.ParseBool(os.Getenv("PRINT_XML"))
	if err != nil {
		printXML = (environment == services.TESTING)
	}

	saveXML, err := strconv.ParseBool(os.Getenv("SAVE_XML"))
	if err != nil {
		saveXML = (environment == services.TESTING)
	}

	logger.Debug("Mode", "PRODUCTION", (environment == services.PRODUCTION), "Log XML", printXML, "Save XML", saveXML)

	cuit, err := strconv.ParseInt(os.Getenv("CUIT"), 10, 64)
	if err != nil {
		logger.Error("missing or invalid environment variable CUIT")
		os.Exit(1)
	}
	logger.Debug("", "CUIT", cuit)

	if os.Getenv("PRIVATE_KEY_FILE") == "" {
		logger.Error("missing environment variable PRIVATE_KEY_FILE")
		os.Exit(1)
	}

	if os.Getenv("CERTIFICATE_FILE") == "" {
		logger.Error("missing environment variable CERTIFICATE_FILE")
		os.Exit(1)
	}

	if len(os.Getenv("KEYS_FILE")) == 0 {
		logger.Error("missing environment variable KEYS_FILE")
		os.Exit(1)
	}

	port := 4433
	if os.Getenv("PORT") != "" {
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			logger.Error("environment variable PORT no numeric.")
			os.Exit(1)
		}
	}

	Wscoem, err = services.NewWscoem(logger, environment, cuit, printXML, saveXML)
	if err != nil {
		logger.Error("NewWscoem()", "err", err.Error())
		os.Exit(1)
	}

	Wscoemcons, err = services.NewWscoemcons(logger, environment, cuit, printXML, saveXML)
	if err != nil {
		logger.Error("NewWscoemcons()", "err", err.Error())
		os.Exit(1)
	}

	Wsgestabref, err = services.Newgestabref(logger, environment, cuit, printXML, saveXML)
	if err != nil {
		logger.Error("Newgestabref()", "err", err.Error())
		os.Exit(1)
	}

	Wsfe, err = services.NewWsfe(logger, environment, cuit, printXML, saveXML)
	if err != nil {
		logger.Error("NewWsfe()", "err", err.Error())
		os.Exit(1)
	}

	/* API Rest */

	middlewareCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		Debug:            false,
	})

	middlewareApiKey, err := middleware.NewApiKeyMiddleware(os.Getenv("KEYS_FILE"))
	if err != nil {
		logger.Error("NewApiKeyMiddleware()", "err", err.Error())
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
	coem.HandleFunc("/Dummy", DummyCoemHandler)
	coem.HandleFunc("POST /RegistrarCaratula", RegistrarCaratulaHandler)
	coem.HandleFunc("DELETE /AnularCaratula", AnularCaratulaHandler)
	coem.HandleFunc("PUT /RectificarCaratula", RectificarCaratulaHandler)
	coem.HandleFunc("POST /RegistrarCOEM", RegistrarCOEMHandler)
	coem.HandleFunc("PUT /SolicitarCambioBuque", SolicitarCambioBuqueHandler)
	coem.HandleFunc("PUT /SolicitarCambioFechas", SolicitarCambioFechasHandler)
	coem.HandleFunc("PUT /SolicitarCambioLOT", SolicitarCambioLOTHandler)
	coem.HandleFunc("PUT /RectificarCOEM", RectificarCOEMHandler)
	coem.HandleFunc("POST /CerrarCOEM", CerrarCOEMHandler)
	coem.HandleFunc("DELETE /AnularCOEM", AnularCOEMHandler)
	coem.HandleFunc("POST /SolicitarAnulacionCOEM", SolicitarAnulacionCOEMHandler)
	coem.HandleFunc("POST /SolicitarNoABordo", SolicitarNoABordoHandler)
	coem.HandleFunc("POST /SolicitarCierreCargaContoBulto", SolicitarCierreCargaContoBultoHandler)
	coem.HandleFunc("POST /SolicitarCierreCargaGranel", SolicitarCierreCargaGranelHandler)

	coemcons := http.NewServeMux()
	coemcons.HandleFunc("/Dummy", DummyCoemconsHandler)
	coemcons.HandleFunc("/ObtenerConsultaEstadosCOEM", ObtenerConsultaEstadosCOEMHandler)
	coemcons.HandleFunc("/ObtenerConsultaNoAbordo", ObtenerConsultaNoAbordoHandler)
	coemcons.HandleFunc("/ObtenerConsultaSolicitudes", ObtenerConsultaSolicitudesHandler)

	gestabref := http.NewServeMux()
	gestabref.HandleFunc("/Dummy", DummyGesTabRefHandler)
	gestabref.HandleFunc("/ConsultarFechaUltAct", ConsultarFechaUltActHandler)
	gestabref.HandleFunc("/ListaArancel", ListaArancelHandler)
	gestabref.HandleFunc("/ListaDescripcion", ListaDescripcionHandler)
	gestabref.HandleFunc("/ListaDescripcionDecodificacion", ListaDescripcionDecodificacionHandler)
	gestabref.HandleFunc("/ListaEmpresas", ListaEmpresasHandler)
	gestabref.HandleFunc("/ListaLugaresOperativos", ListaLugaresOperativosHandler)
	gestabref.HandleFunc("/ListaPaisesAduanas", ListaPaisesAduanasHandler)
	gestabref.HandleFunc("/ListaVigencias", ListaVigenciasHandler)
	gestabref.HandleFunc("/ListaTablasReferencia", ListaTablasReferenciaHandler)
	gestabref.HandleFunc("/ListaDatoComplementario", ListaDatoComplementarioHandler)

	fe := http.NewServeMux()
	fe.HandleFunc("/FEDummy", FEDummyHandler)
	fe.HandleFunc("/FECompUltimoAutorizado", FECompUltimoAutorizadoHandler)
	fe.HandleFunc("/FEParamGetTiposCbte", FEParamGetTiposCbteHandler)
	fe.HandleFunc("/FEParamGetTiposConcepto", FEParamGetTiposConceptoHandler)
	fe.HandleFunc("/FEParamGetTiposDoc", FEParamGetTiposDocHandler)
	fe.HandleFunc("/FEParamGetTiposIva", FEParamGetTiposIvaHandler)
	fe.HandleFunc("/FEParamGetTiposMonedas", FEParamGetTiposMonedasHandler)
	fe.HandleFunc("/FEParamGetTiposOpcional", FEParamGetTiposOpcionalHandler)
	fe.HandleFunc("/FEParamGetTiposTributos", FEParamGetTiposTributosHandler)
	fe.HandleFunc("/FEParamGetPtosVenta", FEParamGetPtosVentaHandler)
	fe.HandleFunc("/FEParamGetCotizacion", FEParamGetCotizacionHandler)
	fe.HandleFunc("/FECompTotXRequest", FECompTotXRequestHandler)
	fe.HandleFunc("/FECAEASinMovimientoConsultar", FECAEASinMovimientoConsultarHandler)
	fe.HandleFunc("/FECompConsultar", FECompConsultarHandler)
	fe.HandleFunc("/FEParamGetTiposPaises", FEParamGetTiposPaisesHandler)
	fe.HandleFunc("/FEParamGetActividades", FEParamGetActividadesHandler)
	fe.HandleFunc("/FEParamGetCondicionIvaReceptor", FEParamGetCondicionIvaReceptorHandler)
	fe.HandleFunc("POST /FECAESolicitar", FECAESolicitarHandler)
	fe.HandleFunc("POST /FECAEARegInformativo", FECAEARegInformativoHandler)

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
		ErrorLog:     serverErrorLogger,
	}

	go util.GracefulShutdown(server)

	logger.Info("Starting server", "PORT", server.Addr)
	err = server.ListenAndServeTLS("keys/server.crt", "keys/server.key")
	if err != nil && err != http.ErrServerClosed {
		logger.Error("http server error", "err", err.Error())
	}
}
