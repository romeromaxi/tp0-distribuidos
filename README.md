# TP0: Docker + Comunicaciones + Concurrencia

## Parte 2: Repaso de Comunicaciones

## Protocolo de Comunicación
En los ejercicios siguientes, se utiliza un protocolo de comunicación basado en una estructura común. Este protocolo establece que, antes de enviar el contenido de un mensaje, primero se debe comunicar el `Tipo de Mensaje`. Para ello, se emplean códigos predefinidos que, al ser enviados, deben ocupar exactamente `4 bytes`. Este formato uniforme asegura que todos los códigos tengan la misma longitud (en bytes), evitando problemas de lectura, ya que se espera siempre un tamaño fijo para el tipo de mensaje.

Una vez recibido el `Tipo de Mensaje`, se puede identificar la estructura del siguiente mensaje, si es necesario. Este segundo mensaje incluye un encabezado de `4 bytes` en formato `BigEndian`, que indica la longitud en bytes del payload (el contenido restante del mensaje). Esto permite saber cuántos bytes se deben leer para recibir el mensaje completo.

El payload se envía como una cadena de texto codificada en `UTF-8`. Si el mensaje contiene varios datos, estos se separan con el delimitador `|`. Este delimitador se eligió tras revisar los archivos de prueba en `./data/dataset.zip`, asegurando que el carácter `|` no aparece en ninguno de ellos.

### Tipos de Mensajes
Los códigos para los diferentes tipos de mensajes en el protocolo son los siguientes:

#### **Enviados por el Cliente**:
- `CONN`: Después de establecer la conexión, el cliente envía este mensaje para comunicar al servidor su *Id de Agencia*. Esto permite al servidor identificar el origen de cada mensaje sin necesidad de incluir esta información en otros (reduciendo su tamaño y procesamiento)
    - Payload: `Id de Agencia`
- `BET`: para indicarle al servidor que quiere realizar una apuesta.
    - Payload: conformado por cada una de los campos necesarios para realizar una apuesta, separados por el caracter delimitador definido. El orden de los campos es `NOMBRE|APELLIDO|DOCUMENTO|NACIMIENTO|NUMERO`
- `NBET`: le avisa al servidor que enviará un conjunto de apuestas en un solo mensaje (procesamiento en _batchs_)
    - Payload: el primer campo corresponde a la cantidad de apuestas que tiene el mensaje, seguidos por todos los campos de cada apuesta (como para el mensaje `BET`). Por ejemplo, en el caso de `2 (dos)` apuestas: `2|NOM_1|APE_1|DOC_1|NAC_1|NRO_1|NOM_2|APE_2|DOC_2|NAC_2|NRO_2`
- `END`: se utiliza cuando el cliente (agencia de apuestas) ha terminado de enviar todas las apuestas correspondientes.
    - Payload: no aplica
- `GWIN`: le solicita al servidor los resultados del sorteo para obtener los ganadores
    - Payload: no aplica

#### **Enviados por el Servidor**:
- `OK`: le informa al cliente en cuestión que la acción solicitada se completó con éxito. Este mensaje puede referirse a la realización de una apuesta, al procesamiento de varias apuestas en batch, o a la finalización de todas las apuestas del cliente.
    - Payload: no aplica
- `NOK`: le indica al cliente que ocurrió un error durante la ejecución de la acción solicitada. Es la respuesta opuesta al mensaje `OK`
    - Payload: no aplica
- `NEND`: para comunicarle al cliente que debe esperar para obtener los ganadores del sorteo, ya que algunas agencias aún no han completado la carga de apuestas. Esto es una respuesta al mensaje `GWIN`, en el caso que el sorteo no haya terminado
    - Payload: no aplica
- `RWIN`: le notifica al cliente que puede obtener los resultados del sorteo, los cuales se comunicarán en el siguiente mensaje (respuesta a `GWIN` cuando el sorteo si ha  finalizado)
    - Payload: el primer campo corresponde a la cantidad de ganadores, seguido por cada uno de los documentos. Para el caso de `3 (tres)` ganadores sería: `3|DOC_1|DOC_2|DOC_3`

## Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).

## Resolución
Para este ejercicio, se añadieron `5 (cinco)` clientes al archivo DockerCompose. Cada cliente tiene definidas las siguientes variables de entorno, que proporcionan los datos necesarios para realizar una apuesta:
- `CLI_BET_NAME`: nombre de la persona que realiza la apuesta
- `CLI_BET_SURNAME`: apellido de la persona que realiza la apuesta
- `CLI_BET_DNI`: documento de la persona que realiza la apuesta
- `CLI_BET_BIRTH`: fecha de nacimiento de la persona que realiza la apuesta
- `CLI_BET_NUMBER`: número de la apuesta

En este caso se utilizaron solamente los mensajes de 
- `CONN`
- `BET`
- `OK`

**Pasos**

En este punto el flujo es el siguiente:
1. El `Cliente` envía el mensaje de conexión `CONN`, y luego su número de agencia (`Id`)

2. Inmediatamente después, el `Cliente` envía el mensaje `BET` y luego los datos necesarios para realizar una apuesta

3. El `Servidor`, al recibir la apuesta, la almacena y le responde de forma satisfactoria al `Cliente` con el mensaje de `OK`


De forma gráfica se vería

<div align="center">
    <img src="assets/Ej5.png" alt="Flujo de mensajes" width="600">
</div>


### Ejecución
Para poder correr este ejercicio basta con ejecutar el siguiente comando:
```
make docker-compose-up
```