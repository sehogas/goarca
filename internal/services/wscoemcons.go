package services

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/hooklift/gowsdl/soap"
	"github.com/sehogas/goarca/ws/wscoemcons"
)

type Wscoemcons struct {
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

func NewWscoemcons(logger *slog.Logger, environment Environment, cuit int64, printXML, saveXML bool) (*Wscoemcons, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	var url string
	if environment == PRODUCTION {
		url = URLWSCOEMConsProduction
	} else {
		url = URLWSCOEMConsTesting
	}

	if len(os.Getenv("WSCOEM_TIPO_AGENTE")) != 4 {
		return nil, fmt.Errorf("missing environment variable WSCOEM_TIPO_AGENTE")
	}

	if len(os.Getenv("WSCOEM_ROL")) != 4 {
		return nil, fmt.Errorf("missing or invalid environment variable WSCOEM_ROL")
	}

	return &Wscoemcons{
		logger:      logger,
		serviceName: "wconscomunicacionembarque",
		environment: environment,
		cuit:        cuit,
		url:         url,
		tipoAgente:  os.Getenv("WSCOEM_TIPO_AGENTE"),
		rol:         os.Getenv("WSCOEM_ROL"),
		soapTlsConfig: tls.Config{
			InsecureSkipVerify: true,
		},
		printXML: printXML,
		saveXML:  saveXML,
	}, nil
}

func (ws *Wscoemcons) PrintAndSaveXML(obj interface{}) {
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

func (ws *Wscoemcons) Dummy() (*wscoemcons.ResultadoEjecucionOfDummyOutput, error) {
	request := &wscoemcons.Dummy{}

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.DummyResult, nil
}

func (ws *Wscoemcons) ObtenerConsultaEstadosCOEM(identificadorCaratula string) (*wscoemcons.ResultadoEjecucionOfResultadoEstadoProceso, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoemcons.ObtenerConsultaEstadosCOEM{
		ArgWSAutenticacionEmpresa: &wscoemcons.WSAutenticacionEmpresa{
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
			WSAutenticacion: &wscoemcons.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
		},
		IdentificadorCabecera: identificadorCaratula,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.ObtenerConsultaEstadosCOEM(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ObtenerConsultaEstadosCOEMResult, nil
}

func (ws *Wscoemcons) ObtenerConsultaNoAbordo(identificadorCaratula string) (*wscoemcons.ResultadoEjecucionOfResultadoNoAbordoProceso, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoemcons.ObtenerConsultaNoAbordo{
		ArgWSAutenticacionEmpresa: &wscoemcons.WSAutenticacionEmpresa{
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
			WSAutenticacion: &wscoemcons.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
		},
		IdentificadorCabecera: identificadorCaratula,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.ObtenerConsultaNoAbordo(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ObtenerConsultaNoAbordoResult, nil
}

func (ws *Wscoemcons) ObtenerConsultaSolicitudes(identificadorCaratula string) (*wscoemcons.ResultadoEjecucionOfResultadoSolicitudProceso, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoemcons.ObtenerConsultaSolicitudes{
		ArgWSAutenticacionEmpresa: &wscoemcons.WSAutenticacionEmpresa{
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
			WSAutenticacion: &wscoemcons.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
		},
		IdentificadorCabecera: identificadorCaratula,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoemcons.NewWconscomunicacionembarqueSoap(client)

	response, err := service.ObtenerConsultaSolicitudes(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.ObtenerConsultaSolicitudesResult, nil
}
