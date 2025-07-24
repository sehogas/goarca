package services

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/goarca/ws/wgestabref"
)

type Wsgestabref struct {
	logger        *slog.Logger
	serviceName   string
	environment   Environment
	url           string
	cuit          int64
	tipoAgente    string
	rol           string
	soapTlsConfig tls.Config
	printXML      bool
	saveXML       bool
}

func Newgestabref(logger *slog.Logger, environment Environment, cuit int64, printXML, saveXML bool) (*Wsgestabref, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	var url string
	if environment == PRODUCTION {
		url = URLWGESTABREFProduction
	} else {
		url = URLWGESTABREFTesting
	}

	if len(os.Getenv("WGESTABREF_TIPO_AGENTE")) != 4 {
		return nil, fmt.Errorf("missing environment variable WGESTABREF_TIPO_AGENTE")
	}

	if len(os.Getenv("WGESTABREF_ROL")) != 4 {
		return nil, fmt.Errorf("missing or invalid environment variable WGESTABREF_ROL")
	}

	return &Wsgestabref{
		logger:      logger,
		serviceName: "wGesTabRef",
		environment: environment,
		url:         url,
		cuit:        cuit,
		tipoAgente:  os.Getenv("WGESTABREF_TIPO_AGENTE"),
		rol:         os.Getenv("WGESTABREF_ROL"),
		soapTlsConfig: tls.Config{
			InsecureSkipVerify: true,
		},
		printXML: printXML,
		saveXML:  saveXML,
	}, nil
}

func (ws *Wsgestabref) PrintAndSaveXML(obj interface{}) {
	if ws.printXML || ws.saveXML {
		data, err := xml.MarshalIndent(obj, " ", "  ")
		if err == nil {
			if ws.printXML {
				ws.logger.Info("printXML", "XML", string(data))
			}
			if ws.saveXML {
				path := fmt.Sprintf("xml/%s", ws.serviceName)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					if err := os.Mkdir(path, os.ModePerm); err != nil {
						ws.logger.Error(fmt.Sprintf("Error creating directory [ %s ]", path), "err", err.Error())
						return
					}
				}
				fileName := fmt.Sprintf("%s/%s_%s.xml", path, strings.ReplaceAll(strings.ReplaceAll(reflect.TypeOf(obj).String(), "*", ""), " ", "_"), strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", ""), "-", ""), "T", ""))
				f, err := os.Create(fileName)
				if err != nil {
					ws.logger.Error(fmt.Sprintf("Error creating file [ %s ]", fileName), "err", err.Error())
				}
				defer f.Close()
				_, err = f.WriteString(string(data))
				if err != nil {
					ws.logger.Error(fmt.Sprintf("Error writing file [ %s ]", fileName), "err", err.Error())
				}
				ws.logger.Info("saveXML", "file", fileName)
			}
		}
	}
}

func (ws *Wsgestabref) Dummy() (*wgestabref.WsDummyResponse, error) {
	request := &wgestabref.Dummy{}

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return nil, err
	}

	return response.DummyResult, nil
}

func (ws *Wsgestabref) ConsultarFechaUltAct(idReferencia string) (*wgestabref.FechaUltAct, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ConsultarFechaUltAct{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ConsultarFechaUltAct(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ConsultarFechaUltActResult, nil
}

func (ws *Wsgestabref) ListaArancel(idReferencia string) (*wgestabref.Opciones, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaArancel{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaArancel(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaArancelResult, nil
}

func (ws *Wsgestabref) ListaDescripcion(idReferencia string) (*wgestabref.Descripciones, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaDescripcion{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaDescripcion(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaDescripcionResult, nil
}

func (ws *Wsgestabref) ListaDescripcionDecodificacion(idReferencia string) (*wgestabref.DescripcionesCodificaciones, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaDescripcionDecodificacion{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaDescripcionDecodificacion(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaDescripcionDecodificacionResult, nil
}

func (ws *Wsgestabref) ListaEmpresas(idReferencia string) (*wgestabref.Empresas, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaEmpresas{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaEmpresas(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaEmpresasResult, nil
}

func (ws *Wsgestabref) ListaLugaresOperativos(idReferencia string) (*wgestabref.LugaresOperativos, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaLugaresOperativos{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaLugaresOperativos(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaLugaresOperativosResult, nil
}

func (ws *Wsgestabref) ListaPaisesAduanas(idReferencia string) (*wgestabref.PaisesAduanas, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaPaisesAduanas{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaPaisesAduanas(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaPaisesAduanasResult, nil
}

func (ws *Wsgestabref) ListaTablasReferencia() (*wgestabref.TablasReferencia, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaTablasReferencia{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaTablasReferencia(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaTablasReferenciaResult, nil
}

func (ws *Wsgestabref) ListaVigencias(idReferencia string) (*wgestabref.Vigencias, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaVigencias{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaVigencias(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaVigenciasResult, nil
}

func (ws *Wsgestabref) ListaDatoComplementario(idReferencia string) (*wgestabref.DatosComplementarios, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wgestabref.ListaDatoComplementario{
		Autentica: &wgestabref.Autenticacion{
			Cuit:       strconv.FormatInt(ws.cuit, 10),
			TipoAgente: ws.tipoAgente,
			Rol:        ws.rol,
			Token:      &ticket.Token,
			Sign:       &ticket.Sign,
		},
		IdReferencia: idReferencia,
	}

	//ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wgestabref.NewWgesTabRefSoap(client)

	response, err := service.ListaDatoComplementario(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ListaDatoComplementarioResult, nil
}
