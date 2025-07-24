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
	"github.com/sehogas/goarca/ws/wscoem"
)

type Wscoem struct {
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

func NewWscoem(logger *slog.Logger, environment Environment, cuit int64, printXML, saveXML bool) (*Wscoem, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	var url string
	if environment == PRODUCTION {
		url = URLWSCOEMProduction
	} else {
		url = URLWSCOEMTesting
	}

	if len(os.Getenv("WSCOEM_TIPO_AGENTE")) != 4 {
		return nil, fmt.Errorf("missing environment variable WSCOEM_TIPO_AGENTE")
	}

	if len(os.Getenv("WSCOEM_ROL")) != 4 {
		return nil, fmt.Errorf("missing or invalid environment variable WSCOEM_ROL")
	}

	return &Wscoem{
		logger:      logger,
		serviceName: "wgescomunicacionembarque",
		environment: environment,
		url:         url,
		cuit:        cuit,
		tipoAgente:  os.Getenv("WSCOEM_TIPO_AGENTE"),
		rol:         os.Getenv("WSCOEM_ROL"),
		soapTlsConfig: tls.Config{
			InsecureSkipVerify: true,
		},
		printXML: printXML,
		saveXML:  saveXML,
	}, nil
}

func (ws *Wscoem) PrintAndSaveXML(obj interface{}) {
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

func (ws *Wscoem) Dummy() (*wscoem.ResultadoEjecucionDummy, error) {
	request := &wscoem.Dummy{}

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.Dummy(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.DummyResult, nil
}

func (ws *Wscoem) RegistrarCaratula(params *wscoem.RegistrarCaratulaRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.RegistrarCaratula{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRegistrarCaratula: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RegistrarCaratula(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.RegistrarCaratulaResult, nil
}

func (ws *Wscoem) AnularCaratula(params *wscoem.AnularCaratulaRequest) (*wscoem.AnularEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.AnularCaratula{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},

		ArgAnularCaratula: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.AnularCaratula(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.AnularCaratulaResult, nil
}

func (ws *Wscoem) RectificarCaratula(params *wscoem.RectificarCaratulaRequest) (*wscoem.RectificarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.RectificarCaratula{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRectificarCaratula: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RectificarCaratula(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.RectificarCaratulaResult, nil
}

func (ws *Wscoem) RegistrarCOEM(params *wscoem.RegistrarCOEMRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.RegistrarCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRegistrarCOEM: params,
	}
	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RegistrarCOEM(request)
	if err != nil {
		return nil, err
	}
	ws.PrintAndSaveXML(response)

	return response.RegistrarCOEMResult, nil
}

func (ws *Wscoem) SolicitarCambioBuque(params *wscoem.SolicitarCambioBuqueRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarCambioBuque{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCambioBuque: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCambioBuque(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarCambioBuqueResult, nil
}

func (ws *Wscoem) SolicitarCambioFechas(params *wscoem.SolicitarCambioFechasRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarCambioFechas{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCambioFechas: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCambioFechas(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarCambioFechasResult, nil
}

func (ws *Wscoem) SolicitarCambioLOT(params *wscoem.SolicitarCambioLOTRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarCambioLOT{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCambioLOT: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCambioLOT(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarCambioLOTResult, nil
}

func (ws *Wscoem) RectificarCOEM(params *wscoem.RectificarCOEMRequest) (*wscoem.RectificarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.RectificarCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgRectificarCOEM: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.RectificarCOEM(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.RectificarCOEMResult, nil
}

func (ws *Wscoem) CerrarCOEM(params *wscoem.CerrarCOEMRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.CerrarCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgCerrarCOEM: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.CerrarCOEM(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	// var errs []error
	// result := ""
	// for _, e := range response.CerrarCOEMResult.ListaErrores.DetalleError {
	// 	if *e.Codigo != 0 {
	// 		errs = append(errs, fmt.Errorf("%d - %s", *e.Codigo, e.Descripcion))
	// 	} else {
	// 		result = strings.TrimSpace(strings.Replace(e.DescripcionAdicional, "Identificador:", "", -1))
	// 	}
	// }

	return response.CerrarCOEMResult, nil
}

func (ws *Wscoem) AnularCOEM(params *wscoem.AnularCOEMRequest) (*wscoem.AnularEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.AnularCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgAnularCOEM: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.AnularCOEM(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.AnularCOEMResult, nil
}

func (ws *Wscoem) SolicitarAnulacionCOEM(params *wscoem.SolicitarAnulacionCOEMRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarAnulacionCOEM{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarAnulacionCOEM: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarAnulacionCOEM(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarAnulacionCOEMResult, nil
}

func (ws *Wscoem) SolicitarNoABordo(params *wscoem.SolicitarNoABordoRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarNoABordo{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarNoABordo: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarNoABordo(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarNoABordoResult, nil
}

func (ws *Wscoem) SolicitarCierreCargaContoBulto(params *wscoem.SolicitarCierreCargaContoBultoRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarCierreCargaContoBulto{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCierreCargaContoBulto: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCierreCargaContoBulto(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarCierreCargaContoBultoResult, nil
}

func (ws *Wscoem) SolicitarCierreCargaGranel(params *wscoem.SolicitarCierreCargaGranelRequest) (*wscoem.RegistrarEmbarqueRta, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wscoem.SolicitarCierreCargaGranel{
		ArgWSAutenticacionEmpresa: &wscoem.WSAutenticacionEmpresa{
			WSAutenticacion: &wscoem.WSAutenticacion{
				Token: ticket.Token,
				Sign:  ticket.Sign,
			},
			CuitEmpresaConectada: ws.cuit,
			TipoAgente:           ws.tipoAgente,
			Rol:                  ws.rol,
		},
		ArgSolicitarCierreCargaGranel: params,
	}

	ws.PrintAndSaveXML(request)

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wscoem.NewWgescomunicacionembarqueSoap(client)

	response, err := service.SolicitarCierreCargaGranel(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.SolicitarCierreCargaGranelResult, nil
}
