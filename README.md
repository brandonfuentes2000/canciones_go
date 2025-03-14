# Proyecto 1

Realiza busqueda de canciones por medio de itunes. Este proyecto utiliza *GOLANG*, *MongoDB* y esta diseñado para ejecutarse en un entorno con *Docker* y *NGINX* 

- [Proyecto 1](#proyecto-1)
  - [**Introducción**](#introducción)
  - [**Características**](#características)
  - [**Requisitos Previos**](#requisitos-previos)
  - [**Instalación**](#instalación)
    - [**1. Clonar el repositorio**](#1-clonar-el-repositorio)
    - [**2. Configurar variables de entorno**](#2-configurar-variables-de-entorno)
    - [**3. Instalar usando el script de instalación**](#3-instalar-usando-el-script-de-instalación)
    - [**4. Ejecutar manualmente con Docker Compose**](#4-ejecutar-manualmente-con-docker-compose)
      - [Construir y ejecutar contenedores](#construir-y-ejecutar-contenedores)
  - [**Arquitectura del Proyecto**](#arquitectura-del-proyecto)
  - [**Consideraciones y Recomendaciones**](#consideraciones-y-recomendaciones)
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
   - Linux limpio (Ubuntu, Debian, etc.)
   - Docker y Docker Compose instalados (si no los tienes, el script de intalación lo hará).
   - Git instalado.

2. **Conexión a Internet**:
   - Para descargar imágenes de contenedores y dependencias.

---

## **Instalación**

### **1. Clonar el repositorio**
```bash
git clone https://github.com/brandonfuentes2000/canciones_go.git
```

### **2. Configurar variables de entorno**

DB_HOST=mongo
DB_PORT=21017
DB_NAME=songs_db
APP_PORT=3000
JWT_SECRET=mysecretkey

### **3. Instalar usando el script de instalación**

Si has recibido un archivo ZIP (canciones_go.zip), sigue estos pasos:

Descomprimir el archivo ZIP:

unzip canciones_go.zip -d canciones_go
cd canciones_go

Dar permisos de ejecución al script de instalación:

chmod +x install.sh

Ejecutar el script para instalar y desplegar la API:

./install.sh

Este script: instalará Docker, extraerá la imagen de la API, configurará MongoDB y Nginx, y levantará los contenedores.

### **4. Ejecutar manualmente con Docker Compose**

#### Construir y ejecutar contenedores
- 1. Puedes utilizar este comando para levantar el contenedor y visualizar los logs
docker-compose up --build 

- 2. Puedes utilizar este comando para levantar el contenedor pero se ejecutan en segundo plano y asi la terminal no te queda mostrando los logs. 
docker-compose up --build -d

Esto iniciará tres servicios:

app: Servicio API disponible en http://localhost:3000.
mongo: Base de datos accesible en el puerto 27017.
nginx: Servidor web para exponer el API en htpp://localhost.

- Despues de eso puedes usar el siguiente comando para ver los contenedores que se estan ejecutando
docker ps

## **Arquitectura del Proyecto**
```canciones/
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
├── nginx/                  # Configuración servidor nginx
   └── conf.d/default.conf  # Configuración del servidor y proxy. Aquí es donde definimos cómo Nginx redirige las peticiones a go
├── .env                    # Archivo de configuración de variables de entorno
├── .gitignore              # Configuración de exclusiones de Git
├── docker-compose.yml      # Configuración para orquestación de contenedores
├── Dockerfile              # Definición de imagen para el contenedor de la app
├── go.mod                  # Dependencias del proyecto en Go
├── go.sum                  # Registro de versiones de las dependencias
├── nginx.conf              # Este archivo contiene solo configuraciones generales de Nginx
└── README.md               # Documentación de instalación, configuración y ejecución del proyecto.
└── API_DOCUMENTACION.md    # Documentación de uso para consumir la API (endpoints, autenticación)```

## **Consideraciones y Recomendaciones**
Configura correctamente las variables de entorno (.env).
Asegúrate de exponer los puertos necesarios para los servicios (80, 3000, 27017).
Utiliza un cliente HTTP como Insomnia para validar la API.