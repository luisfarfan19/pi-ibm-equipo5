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
Simplemente ejecuta el api de POST /watson/validate/shelve 
Donde el body debe recibir:
```
{
"image64": "{imagen en base64}"
"imageType": "{mime type de la imagen}"
}
```
Cada ejecucion al api de watson, internamente se llama al api para obtener el bearer token. Esto habria que validar si lo dejamos asi o buscamos otra alternativa ya que tarda bastante el api en responder

## Contributing

Pull requests son bienvenidos. Para cambios importantes, por favor abre primero un issue para discutir lo que te gustaría modificar.
