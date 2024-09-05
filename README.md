# TP0: Docker + Comunicaciones + Concurrencia

## Parte 3: Repaso de Concurrencia

## Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

En caso de que el alumno implemente el servidor Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

## Resolución
En este ejercicio, se decidió implementar el servidor utilizando _multithreading_. Al aceptar una conexión de un cliente, se crea un nuevo thread para manejar las solicitudes y mensajes de dicho cliente, lo que permite procesar múltiples conexiones en paralelo. Aunque el `Global Interpreter Lock (GIL)` de Python limita la ejecución paralela de operaciones en CPU, en este caso no genera inconvenientes, ya que la mayoría de las operaciones que realiza el servidor son de I/O (como la lectura y escritura en archivos), las cuales no están bloqueadas por el `GIL`. Esto permite que el servidor acepte conexiones y procese mensajes de manera concurrente de forma eficiente.

Para garantizar la correcta sincronización, se emplearon dos mecanismos de locks:
- Un _lock_ sobre el archivo de apuestas, utilizado tanto para las operaciones de escritura (`store_bets`) como de lectura (`load_bets`), asegurando que el acceso concurrente al archivo no cause inconsistencias.
- Otro _lock_ que se utiliza para almacenar las agencias que confirmaron haber finalizado sus apuestas, evitando problemas de concurrencia al actualizar esta información en paralelo.

### Ejecución
1. Extraer los archivos de las apuestas de `./data/dataset.zip`, y colocarlos dentro de `./data/dataset/`

2. Ejecutar el siguiente comando
    ```
    make docker-compose-up
    ```