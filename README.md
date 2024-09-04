# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker


## Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`).

## Solución:
Para este ejercicio, se desarrolló el script `validar-echo-server.sh` que verifica el funcionamiento del servidor utilizando `netcat`. Se puede ejecutar mediante el comando `sh` o de forma directa `./` (si se tienen permisos de ejecución).

Antes de ejecutar el script, se debe asegurar de que el servidor esté activo (`make docker-compose-up`)

```
./validar-echo-server.sh
```

A su vez se le pueden pasar dos parámetros por si se requiere configurar el puerto del servidor o cambiar el mensaje que se envia
```
./validar-echo-server.sh <port> <message>
```
Por ejemplo, si quere enviarse el mensaje `Test echo-server` pero se quiere seguir usando el puerto preestablecido (`12345`), puede ejecutarse el siguiente comando
```
./validar-echo-server.sh - "Test echo-server"
```