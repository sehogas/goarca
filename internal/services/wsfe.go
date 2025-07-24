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
	"github.com/sehogas/goarca/ws/wsfe"
)

type Wsfe struct {
	logger        *slog.Logger
	serviceName   string
	environment   Environment
	url           string
	cuit          int64
	soapTlsConfig tls.Config
	printXML      bool
	saveXML       bool
}

func NewWsfe(logger *slog.Logger, environment Environment, cuit int64, printXML, saveXML bool) (*Wsfe, error) {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	var url string
	if environment == PRODUCTION {
		url = URLWSFEProduction
	} else {
		url = URLWSFETesting
	}

	return &Wsfe{
		logger:      logger,
		serviceName: "wsfe",
		environment: environment,
		url:         url,
		cuit:        cuit,
		soapTlsConfig: tls.Config{
			InsecureSkipVerify: true,
		},
		printXML: printXML,
		saveXML:  saveXML,
	}, nil
}

func (ws *Wsfe) PrintAndSaveXML(obj interface{}) {
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

func (ws *Wsfe) FEDummy() (*wsfe.DummyResponse, error) {
	request := &wsfe.FEDummy{}

	client := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(client)

	response, err := service.FEDummy(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEDummyResult, nil
}

func (ws *Wsfe) FEUltimoComprobanteEmitido(ptoVta int32, cbteTipo int32) (*wsfe.FERecuperaLastCbteResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wsfe.FECompUltimoAutorizado{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		PtoVta:   ptoVta,
		CbteTipo: cbteTipo,
	}

	ws.PrintAndSaveXML(request)

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECompUltimoAutorizado(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FECompUltimoAutorizadoResult, nil
}

func (ws *Wsfe) FECAESolicitar(cab *wsfe.FECabRequest, det []*wsfe.FECAEDetRequest) (*wsfe.FECAEResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wsfe.FECAESolicitar{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		FeCAEReq: &wsfe.FECAERequest{
			FeCabReq: &wsfe.FECAECabRequest{
				FECabRequest: cab,
			},
			FeDetReq: &wsfe.ArrayOfFECAEDetRequest{
				FECAEDetRequest: det,
			},
		},
	}

	ws.PrintAndSaveXML(request)

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECAESolicitar(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FECAESolicitarResult, nil
}

func (ws *Wsfe) FEParamGetTiposCbte() (*wsfe.CbteTipoResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposCbte(&wsfe.FEParamGetTiposCbte{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposCbteResult, nil
}

func (ws *Wsfe) FEParamGetTiposConcepto() (*wsfe.ConceptoTipoResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposConcepto(&wsfe.FEParamGetTiposConcepto{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposConceptoResult, nil
}

func (ws *Wsfe) FEParamGetTiposDoc() (*wsfe.DocTipoResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposDoc(&wsfe.FEParamGetTiposDoc{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposDocResult, nil
}

func (ws *Wsfe) FEParamGetTiposIva() (*wsfe.IvaTipoResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposIva(&wsfe.FEParamGetTiposIva{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposIvaResult, nil
}

func (ws *Wsfe) FEParamGetTiposMonedas() (*wsfe.MonedaResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposMonedas(&wsfe.FEParamGetTiposMonedas{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposMonedasResult, nil
}

func (ws *Wsfe) FEParamGetTiposOpcional() (*wsfe.OpcionalTipoResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposOpcional(&wsfe.FEParamGetTiposOpcional{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposOpcionalResult, nil
}

func (ws *Wsfe) FEParamGetTiposTributos() (*wsfe.FETributoResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposTributos(&wsfe.FEParamGetTiposTributos{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposTributosResult, nil
}

func (ws *Wsfe) FEParamGetPtosVenta() (*wsfe.FEPtoVentaResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetPtosVenta(&wsfe.FEParamGetPtosVenta{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetPtosVentaResult, nil
}

func (ws *Wsfe) FEParamGetCotizacion(monId, fchCotiz string) (*wsfe.FECotizacionResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	request := &wsfe.FEParamGetCotizacion{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		MonId:    monId,
		FchCotiz: fchCotiz,
	}

	ws.PrintAndSaveXML(request)

	response, err := service.FEParamGetCotizacion(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetCotizacionResult, nil
}

func (ws *Wsfe) FECompTotXRequest() (*wsfe.FERegXReqResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECompTotXRequest(&wsfe.FECompTotXRequest{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FECompTotXRequestResult, nil
}

func (ws *Wsfe) FECAEARegInformativo(cab *wsfe.FECabRequest, det []*wsfe.FECAEADetRequest) (*wsfe.FECAEAResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	request := &wsfe.FECAEARegInformativo{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		FeCAEARegInfReq: &wsfe.FECAEARequest{
			FeCabReq: &wsfe.FECAEACabRequest{
				FECabRequest: cab,
			},
			FeDetReq: &wsfe.ArrayOfFECAEADetRequest{
				FECAEADetRequest: det,
			},
		},
	}

	ws.PrintAndSaveXML(request)

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FECAEARegInformativo(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FECAEARegInformativoResult, nil
}

func (ws *Wsfe) FECAEASinMovimientoConsultar(caea string, ptoVta int32) (*wsfe.FECAEASinMovConsResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	request := &wsfe.FECAEASinMovimientoConsultar{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		CAEA:   caea,
		PtoVta: ptoVta,
	}

	ws.PrintAndSaveXML(request)

	response, err := service.FECAEASinMovimientoConsultar(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FECAEASinMovimientoConsultarResult, nil
}

func (ws *Wsfe) FECompConsultar(ptoVta, cbteTipo int32, cbteNro int64) (*wsfe.FECompConsultaResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	request := &wsfe.FECompConsultar{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		FeCompConsReq: &wsfe.FECompConsultaReq{
			PtoVta:   ptoVta,
			CbteTipo: cbteTipo,
			CbteNro:  cbteNro,
		},
	}

	ws.PrintAndSaveXML(request)

	response, err := service.FECompConsultar(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FECompConsultarResult, nil
}

func (ws *Wsfe) FEParamGetTiposPaises() (*wsfe.FEPaisResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetTiposPaises(&wsfe.FEParamGetTiposPaises{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetTiposPaisesResult, nil
}

func (ws *Wsfe) FEParamGetActividades() (*wsfe.FEActividadesResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	response, err := service.FEParamGetActividades(&wsfe.FEParamGetActividades{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
	})
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetActividadesResult, nil
}

func (ws *Wsfe) FEParamGetCondicionIvaReceptor(claseCmp string) (*wsfe.CondicionIvaReceptorResponse, error) {
	ticket, err := GetTA(ws.environment, ws.serviceName, ws.cuit)
	if err != nil {
		return nil, err
	}

	conexion := soap.NewClient(ws.url, soap.WithTLS(&ws.soapTlsConfig))
	service := wsfe.NewServiceSoap(conexion)

	request := &wsfe.FEParamGetCondicionIvaReceptor{
		Auth: &wsfe.FEAuthRequest{
			Token: ticket.Token,
			Sign:  ticket.Sign,
			Cuit:  ticket.Cuit},
		ClaseCmp: claseCmp,
	}

	ws.PrintAndSaveXML(request)

	response, err := service.FEParamGetCondicionIvaReceptor(request)
	if err != nil {
		return nil, err
	}

	ws.PrintAndSaveXML(response)

	return response.FEParamGetCondicionIvaReceptorResult, nil
}
