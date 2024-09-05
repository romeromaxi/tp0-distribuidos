# TP0: Docker + Comunicaciones + Concurrencia

## Protocolo de Comunicación
En los ejercicios siguientes, se utiliza un protocolo de comunicación basado en una estructura común. Este protocolo establece que, antes de enviar el contenido de un mensaje, primero se debe comunicar el `Tipo de Mensaje`. Para ello, se emplean códigos predefinidos que, al ser enviados, deben ocupar exactamente `4 bytes`. Este formato uniforme asegura que todos los códigos tengan la misma longitud (en bytes), evitando problemas de lectura, ya que se espera siempre un tamaño fijo para el tipo de mensaje.

Una vez recibido el `Tipo de Mensaje`, se puede identificar la estructura del siguiente mensaje, si es necesario. Este segundo mensaje incluye un encabezado de `4 bytes` en formato `BigEndian`, que indica la longitud en bytes del payload (el contenido restante del mensaje). Esto permite saber cuántos bytes se deben leer para recibir el mensaje completo.

El payload se envía como una cadena de texto codificada en `UTF-8`. Si el mensaje contiene varios datos, estos se separan con el delimitador `|`. Este delimitador se eligió tras revisar los archivos de prueba en `./data/dataset.zip`, asegurando que el carácter `|` no aparece en ninguno de ellos.

### Tipos de Mensajes
Los códigos para los diferentes tipos de mensajes en el protocolo son los siguientes:

#### **Enviados por el Cliente**:
- `CONN`: Después de establecer la conexión, el cliente envía este mensaje para comunicar al servidor su *Id de Agencia*. Esto permite al servidor identificar el origen de cada mensaje sin necesidad de incluir esta información en otros (reduciendo su tamaño y procesamiento).
    - Payload: `Id de Agencia`.
- `BET`: para indicarle al servidor que quiere realizar una apuesta.
    - Payload: conformado por cada una de los campos necesarios para realizar una apuesta, separados por el caracter delimitador definido. El orden de los campos es `NOMBRE|APELLIDO|DOCUMENTO|NACIMIENTO|NUMERO`.
- `NBET`: le avisa al servidor que enviará un conjunto de apuestas en un solo mensaje (procesamiento en _batchs_).
    - Payload: el primer campo corresponde a la cantidad de apuestas que tiene el mensaje, seguidos por todos los campos de cada apuesta (como para el mensaje `BET`). Por ejemplo, en el caso de `2 (dos)` apuestas: `2|NOM_1|APE_1|DOC_1|NAC_1|NRO_1|NOM_2|APE_2|DOC_2|NAC_2|NRO_2`.
- `END`: se utiliza cuando el cliente (agencia de apuestas) ha terminado de enviar todas las apuestas correspondientes.
    - Payload: no aplica.
- `GWIN`: le solicita al servidor los resultados del sorteo para obtener los ganadores
    - Payload: no aplica.

#### **Enviados por el Servidor**:
- `OK`: le informa al cliente en cuestión que la acción solicitada se completó con éxito. Este mensaje puede referirse a la realización de una apuesta, al procesamiento de varias apuestas en batch, o a la finalización de todas las apuestas del cliente.
    - Payload: no aplica.
- `NOK`: le indica al cliente que ocurrió un error durante la ejecución de la acción solicitada. Es la respuesta opuesta al mensaje `OK`.
    - Payload: no aplica.
- `NEND`: para comunicarle al cliente que debe esperar para obtener los ganadores del sorteo, ya que algunas agencias aún no han completado la carga de apuestas. Esto es una respuesta al mensaje `GWIN`, en el caso que el sorteo no haya terminado.
    - Payload: no aplica.
- `RWIN`: le notifica al cliente que puede obtener los resultados del sorteo, los cuales se comunicarán en el siguiente mensaje (respuesta a `GWIN` cuando el sorteo si ha  finalizado).
    - Payload: el primer campo corresponde a la cantidad de ganadores, seguido por cada uno de los documentos. Para el caso de `3 (tres)` ganadores sería: `3|DOC_1|DOC_2|DOC_3`.

## Parte 2: Repaso de Comunicaciones

## Ejercicio N°7:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

## Resolución
Para este ejercicio, se mantienen los `5 (cinco)` clientes definidos en el archivo DockerCompose del ejercicio anterior. Además, se ha añadido una nueva variable de entorno en el servidor llamada`NUMBER_AGENCIES`, que indica la cantidad de agencias por las que se debe esperar para poder finalizar el sorteo.


En este caso se utilizaron solamente los mensajes de 
- `CONN`
- `NBET`
- `OK`
- `NOK`
- `END`
- `GWIN`
- `NEND`
- `RWIN`

Aunque el servidor podría seguir reconociendo un paquete del tipo `BET`.

**Pasos**

En este punto el flujo es el siguiente:
1. El `Cliente` envía el mensaje de conexión `CONN`, y luego su número de agencia (`Id`).

2. A continuación, el `Cliente` envía el mensaje `NBET` con una lote de apuestas, en un proceso _batch_.

3. El `Servidor` recibe las apuestas, las decodifica y las almacena. Si no se produce ningún error en este proceso, responde al `Cliente` con el mensaje de `OK`. Si se produjo un error, le enviará el correspondiente `NOK`. 

4. Si el `Cliente`  tiene apuestas pendientes por enviar, vuelve al _Paso 2_.

5. Una vez que el `Cliente` ha enviado todas las apuestas, comunica la finalización de la carga con el mensaje `END`.

6. A lo que el `Servidor` registra que dicha agencia terminó, y devuelve el mensaje de estado correspondiente (`OK` si se realizó con éxito, o `NOK` si ocurrió algún imprevisto).

7. Inmediatamente después de recibir la confirmación de éxito, el `Cliente`, le solicita al `Servidor` por los ganadores del sorteo.
    - Si el sorteo ha finalizado, el `Servidor`consulta por todos los ganadores dentro de esta agencia, y responde con el mensaje `RWIN`.
    - Si el sorteo aún no ha terminado, el `Cliente` se desconecta por unos instantes, luego vuelve a establecer comunicación con el `Servidor`, enviando nuevamente el mensaje `CONN`, y repitiendo el _Paso 7_.


De forma gráfica se vería

<div align="center">
    <img src="assets/Ej7.png" alt="Flujo de mensajes - Ej7" width="600">
</div>


### Ejecución
Para una correcta ejecución se deben seguir estos pasos:
1. Extraer los archivos de las apuestas de `./data/dataset.zip`, y colocarlos dentro de `./data/dataset/`.

2. Ejecutar el siguiente comando:
    ```
    make docker-compose-up
    ```