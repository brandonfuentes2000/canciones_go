# Documentación del Servicio - Canciones API 🎵

Este documento detalla cómo consumir el servicio **Canciones API**, que proporciona búsqueda de canciones desde iTunes Y ChartLyrics.
---

## **Base URL**
La API se encuentra disponible en:
```
http://localhost
```
Si deseas acceder directamente a la API sin pasar por Nginx:
```
http://localhost:3000
```

---

## **Autenticación**
La API requiere autenticación mediante **JSON Web Token (JWT)**. Primero, el usuario debe autenticarse enviando una solicitud `POST` al endpoint `/login`.

### **Endpoint: Iniciar sesión**
```
http://localhost:3000/login
```

### **Ejemplo de solicitud**
```json
{
  "username": "brandon",
  "password": "contra2024"
}
```

### **Respuesta esperada**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImJyYW5kb24iLCJleHAiOjE2ODMwMDAwMDB9.BK3YV_6sMEo"
}
```

El token JWT tiene una duración de **2 horas**. 

Para consumir la API, el token debe enviarse en el encabezado `Authorization` en cada solicitud:
```
Authorization: Bearer {token}
```

---

## **Búsqueda de canciones**
Este endpoint permite buscar canciones por nombre, artista y album.

### **Endpoint: Buscar canciones**
```
GET /api/search?name={nombre}&artist={artista}
```

### **Ejemplo de solicitud**
```
GET /api/search?name=Shape+of+You&artist=Ed+Sheeran
```

---

## **Errores Comunes y Solución**
### **1. Token inválido**
**Solución:** Asegúrate de que el token JWT esté presente en el encabezado `Authorization` y no haya expirado.

### **2. Petición incorrecta**
**Solución:** Verifica que los parámetros de la URL estén bien formateados.

### **3. Problema del servidor**
**Solución:** Revisa los logs de la API y asegúrate de que MongoDB y Nginx estén corriendo correctamente.

---

## **Conclusión**
Esta documentación proporciona todo lo necesario para consumir la API. Asegúrate de manejar correctamente la autenticación y errores para una mejor experiencia de usuario. 

