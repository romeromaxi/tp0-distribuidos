# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker

## Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).

## Resolución
Para permitir que los cambios en los archivos de configuración (`config.ini` y `config.yaml`) sean efectivos sin necesidad de reconstruir las imágenes de Docker, se configuraron los `volumes` correspondientes dentro del archivo `docker-compose-dev.yaml`. Asegurando así que los cambios en estos archivos se reflejen en tiempo real en el contenedor sin necesidad de reconstruir la imagen.

Configuración para el `Servidor`:
```
volumes:
    - ./server/config.ini:/config.ini:ro
```

Mientras que para el `Cliente`:
```
volumes:
    - ./client/config.yaml:/config.yaml:ro
```
En ambos se utilizó la opción de `:ro` (read-only) para que sean tratados solamente como archivos de lectura, evitando posibles modificaciones.

Para realizar la prueba correspondiente se pueden ejecutar los contenedores (con el comando `make docker-compose-up`), luego modificar los archivos de configuración y ver como se actualizan los contenedores