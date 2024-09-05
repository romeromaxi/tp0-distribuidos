# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker
## Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

## Resolución
Para garantizar que tanto el cliente como el servidor terminen de manera _graceful_ al recibir una señal `SIGTERM`, se realizaron las siguientes modificaciones:

### Servidor
Se configuró el manejo de señales utilizando `signal.signal` en `Python`.Se estableció una función callback para manejar la señal `SIGTERM`, la cual se encarga de cerrar los recursos abiertos antes de que el servidor termine.

### Cliente
Se implementó una `goroutine` para manejar la recepción de señales del sistema, que se captura con `signal.Notify`. Cuando se recibe `SIGTERM`, la `goroutine` es la encargada de cerrar todos los recursos abiertos.

Por otro lado, se incrementó el valor del `timeout` dentro del `docker-compose-down`, para garantizar que Docker espere a que los contenedores terminen de manera ordenada antes de deternelos.

Para poder realizar la prueba correspondiente se deben seguir los siguientes pasos
1. Levanta el cliente y el servidor utilizando el comando:
    ```
    make docker-compose-up
    ```

2. Termina la ejecución anticipadamente con:
    ```
    make docker-compose-down
    ```
