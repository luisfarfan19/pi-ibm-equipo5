# Watson Service

El servicio tiene como funcionalidad recibir API requests para realizar diversos llamados a IBM cloud, asi como logica interna necesaria para el proyecto de **Automatización de procesos de gestión de canal al menudeo mediante un modelo de visión**.

## Instalación

- GO 1.24.1
- Gorilla mux
- Docker (opcional)

## Ejecutar el servicio
Se puede ejecutar por medio de Docker o directamente con GO.

Directamente desde go

```go
go mod tidy
go get github.com/gorilla/mux

go build .
go run main.go
```

Utilizando Docker
```Docker
docker build -t wantson-service .
docker run -d -p 8080:8080 wantson-service
```

Para validar que se ha ejecutado correctamente podemos llamar el api /health
```
curl http://localhost:8080/health
```
y recibir el status UP
```
{"status":"UP"}
```

## Ejecutando el api de validacion de gondolas
Para ejecutar el api v1 de la validacion de gondolas es necesario contar con el Bearer token del usuario de IBM.
Puedes recuperarlo utilizando la siguiente [URL](https://cloud.ibm.com/docs/key-protect?topic=key-protect-retrieve-access-token). Recuerda utilizar el de 2990184 Juan Nolazco

Una vez recuperado el Bearer token, tienes que cambiar el bearer token del codigo dentro del archivo pkg/utils.go y correr el proyecto.

Por el momento el api toma la imagen a partir del filesystem. Puedes actualizarlo con cambiar/agregar la imagen nueva dentro de la carpeta image y actualizar la configuracion dentro de pkg/utils.go
 
Nota: Se estara validado como eliminar este setteo manual del token. Asi como la obtencion de la imagen.

## Contributing

Pull requests son bienvenidos. Para cambios importantes, por favor abre primero un issue para discutir lo que te gustaría modificar.
