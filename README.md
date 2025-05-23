# GOARCA
API Restful Json que sirve todos los endpoints de los siguientes webservices ARCA:

* wsaa
* wgesTabRef
* wconscomunicacionembarque
* wgescomunicacionembarque
* wsfev1  (pendiente)

#### Ejecución
1. Descargar los fuentes
``git clone https://github.com/sehogas/goarca.git``

2. Ejecutar alguna de estas opciones:

  ``go run ./cmd/api/.``

ó

  ``make``

3. Navegar a http://localhost:3000/swagger para visualizar toda la documentación de los endpoints.


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

```
$ gowsdl -o wgescomunicacionembarque.go -d ws/ -p wscoem  wsdl/wgescomunicacionembarque.xml
```

### Generar dentro de ./keys los certificados del servidor para https

```
$ openssl req -x509 -nodes -newkey rsa:2048 -keyout server.rsa.key -out server.rsa.crt -days 3650
```

---
#### Créditos
  Para la conexión soap y la generación del archivo wsaa.go se utilizó https://github.com/hooklift/gowsdl/