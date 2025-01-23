# proyecto 1

Realiza busqueda de canciones por medio de itunes. Este proyecto utiliza *GOLANG*, *MongoDB* y esta diseñado para ejecutarse en un entorno con *Docker* y *NGINX* 

## **Tabla de Contenidos**

1. [Introducción](#introducción)
2. [Características](#características)
3. [Requisitos Previos](#requisitos-previos)
4. [Instalación](#instalación)
    - [1. Clonar el repositorio](#1-clonar-el-repositorio)
    - [2. Configurar las variables de entorno](#2-configurar-las-variables-de-entorno)
    - [3. Ejecutar con Docker Compose](#3-ejecutar-con-docker-compose)
5. [Uso](#uso)
6. [Arquitectura del Proyecto](#arquitectura-del-proyecto)
7. [Consideraciones y Recomendaciones](#consideraciones-y-recomendaciones)
---

## **Introducción**

Canciones API centraliza las búsquedas de canciones desde Itunes. Proporciona una única interfaz para consultar canciones por:
- Nombre
- Artista
- Álbum

Este servicio está diseñado para:
- Reducir la complejidad al interactuar con múltiples APIs es decir que solo manejen un endpoint y la persona elija con que parametro buscar.
- Almacenar los resultados en una base de datos para su reutilización.
- Responder con datos estandarizados y bien formateados.

---

## **Características**

- **Autenticación**: Proporciona un sistema seguro para consumir el servicio.
- **Persistencia**: Almacena respuestas en MongoDB para minimizar llamadas repetitivas.
- **Contenerización**: Desplegable fácilmente con Docker y Docker Compose.
- **Extensibilidad**: Es posible agregar más fuentes de datos en el futuro.

---

## **Requisitos Previos**

Asegúrate de que tu sistema cumpla con los siguientes requisitos antes de comenzar:

1. **Herramientas necesarias**:
   - Docker y Docker Compose instalados.
   - Git instalado.
   - Editor de texto compatible con Markdown (por ejemplo, VS Code).

2. **Conexión a Internet**:
   - Para descargar imágenes de contenedores y dependencias.

---

## **Instalación**

### **1. Clonar el repositorio**
```bash
git clone https://github.com/brandonfuentes2000/canciones_go.git
```

## **2. Configurar variables de entorno**

# Variables de entorno
DB_HOST=mongo
DB_PORT=21017
DB_NAME=songs_db
APP_PORT=3000
JWT_SECRET=mysecretkey

## **3. Ejecutar con Docker

# Construir y ejecutar contenedores
- 1. Puedes utilizar este comando para levantar el contenedor y visualizar los logs
docker-compose up --build 

- 2. Puedes utilizar este comando para levantar el contenedor pero se ejecutan en segundo plano y asi la terminal no te queda mostrando los logs. 
docker-compose up --build -d

Esto iniciará tres servicios:

app: Servicio API disponible en http://localhost:3000.
mongo: Base de datos accesible en el puerto 27017.
nginx: Servidor web para exponer el API.

- Despues de eso puedes usar el siguiente comando para ver los contenedores que se estan ejecutando
docker ps

## **Uso**
Utilice Insomnia o Postman
Realize una solicitud Post al siguiente endpoint y con la siguiente estructura:
http://localhost:3000/login

{
	"username" : "brandon", 
	"password" : "contra2024"
}

Se va a generar un token el cual tiene un tiempo de duracion de 2 horas

Con ese token debe abrir otro request en insomnia o postman dirigirse a Headers agregar uno nuevo.
En header colocar 
 ## **Authorization**
En el value colocar el token que les genero con esto incluido
## **Bearer {Aqui va el token}**
Ahora realiza una solicitud GET al siguiente endpoint 
http://localhost:3000/api/search?name=bad+bunny&artist=drake

## **Arquitectura del Proyecto**
canciones/
├── cmd/                    # Punto de entrada principal
│   └── main.go             # Archivo principal para iniciar la aplicación
├── internal/               # Lógica interna del proyecto
│   ├── handlers/           # Controladores para manejar las rutas de la API
│   │   ├── auth.go         # Manejador de autenticación
│   │   ├── health.go       # Manejador para verificar el estado del servicio
│   │   └── search.go       # Manejador para búsqueda de canciones
│   ├── middleware/         # Funcionalidades intermedias, como autenticación
│   │   └── jwt.go          # Middleware para autenticación con JWT
│   ├── models/             # Modelos de datos usados en la aplicación
│   │   ├── song.go         # Modelo para canciones
│   │   └── user.go         # Modelo para usuarios
│   └── storage/            # Capa de acceso y persistencia de datos
│       ├── db.go           # Conexión y configuración de MongoDB
│       └── external.go     # Interacciones con APIs externas
├── .env                    # Archivo de configuración de variables de entorno
├── .gitignore              # Configuración de exclusiones de Git
├── docker-compose.yml      # Configuración para orquestación de contenedores
├── Dockerfile              # Definición de imagen para el contenedor de la app
├── go.mod                  # Dependencias del proyecto en Go
├── go.sum                  # Registro de versiones de las dependencias
├── nginx.conf              # Configuración de NGINX 
└── README.md               # Documentación del proyecto


## **Consideraciones y Recomendaciones**
Configura correctamente las variables de entorno.
Asegúrate de exponer los puertos necesarios para los servicios.
Utiliza un cliente HTTP como Insomnia para validar la API.