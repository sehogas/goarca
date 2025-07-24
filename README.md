# GOARCA
Este servidor es una API restful que actua como proxy JSON/XML a los webservices de ARCA para los siguientes servicios:

* wgesTabRef (Consulta de Tablas de Referencia
 necesarias para los Web Services del SIM)
* wconscomunicacionembarque (Consultas para
Comunicación de Embarque)
* wgescomunicacionembarque (Comunicación de Embarque)
* wsfev1 (Facturación Electrónica Argentina)
* wsfecred (Factura Electrónica de Crédito MiPyMEs) ***Pendiente*** 
* wgesStockDepositosFiscales (Stock Depósitos Fiscales) ***Pendiente***


#### Ejecución
1. Descargar los fuentes
``git clone https://github.com/sehogas/goarca.git``

2. Configurar certificados ARCA y del servidor (vea Requisitos previos)

3. Configurar variables de entorno en el archivo .env (puede basarse en el archivo .env.example)

4. Formas de ejecución:

  ````
  $ go run ./cmd/api/.
  ````

ó

  ````
  $ make run
  ````
ó

  ````
  $ docker compose up -d --build
  ````


3. Navegar a http://localhost:4433/swagger para visualizar toda la documentación de los endpoints.


#### Requisitos previos
1. Generar clave privada RSA, crear solicitud de certificado y obtener en AFIP el certificado. 
2. Configurar variables de entorno con los datos del paso 1. Utilizar plantilla: .env.example


#### Pasos para generar clave privada RSA y Certificado AFIP
  Documentación: https://www.afip.gob.ar/ws/WSASS/html/generarcsr.html

    openssl genrsa -out MiClavePrivada 2048

    openssl req -new -key MiClavePrivada -subj "/C=AR/O=XXXX/CN=YYYY/serialNumber=CUIT 20999999992" -out misolicitud.csr

  Crear el archivo "certificado.pem" y copiar el certificado x509v2 en formato PEM generado por la página de AFIP 

#### Configurar variables de entorno 

  PRIVATE_KEY_FILE=MiClavePrivada
  CERTIFICATE_FILE=certificado.pem


### Ejemplo de creación del llamado a un servicio

````
$ gowsdl -o wgescomunicacionembarque.go -d ws/ -p wscoem  wsdl/wgescomunicacionembarque.xml
````

### Generar dentro de ./keys los certificados del servidor para https

```
$ openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650
```

### Ejemplo de Docker Compose

````
services:
  api-arca:
    image: shogas/goarca:latest
    restart: always
    ports:
      - 4433:4433
    volumes:
      - ./keys:/keys:ro
      - ./xml:/xml
    env_file:
      - .env

volumes:
  keys:
  xml:
````
---
#### Créditos
  https://github.com/hooklift/gowsdl