# TP0: Docker + Comunicaciones + Concurrencia

## Parte 1: Introducción a Docker
## Ejercicio N°1:
Además, definir un script de bash `generar-compose.sh` que permita crear una definición de DockerCompose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

## Resolución
Para este ejercicio, se desarrolló el script `generar-compose.sh`, el cual recibe dos parámetros:
- `<output_filename>`: nombre del archivo de salida
- `<clients_number>`: la cantidad de clientes que se deben incluir en la definición de Docker Compose

El modo de ejecución puede ser tanto de forma directa `./` (si se tienen permisos de ejecución), como utilizando `sh`
```
./generar-compose.sh <output_filename> <clients_number>
```
```
sh generar-compose.sh <output_filename> <clients_number>
```
Por ejemplo, si se desea que el archivo de salida se llame `docker-compose-dev-test.yaml` y configurar únicamente `2 (dos)` clientes, se debería usar el siguiente comando:
```
./generar-compose.sh docker-compose-dev-test.yaml 2
```
Cuyo resultado por pantalla sería 
```
Generating configuration file docker-compose-dev-test.yaml with 2 client(s)...
The file docker-compose-dev-test.yaml with 2 client(s) was generated successfully.
```

En caso de que no se proporcionen los parámetros necesarios, o si los parámetros proporcionados son insuficientes, el script mostrará el siguiente mensaje indicando cómo debe ejecutarse
```
Usage: ./generar-compose.sh <output_filename> <clients_number>
```