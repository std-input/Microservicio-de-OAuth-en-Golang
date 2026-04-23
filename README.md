# Microservicio de OAuth en Golang - Golang OAuth Microservice

Un microservicio de autenticación basado en Go que utiliza el framework Fiber v3, OAuth2 de Google, tokens JWT e integración con base de datos.

## Características

- Autenticación OAuth2 de Google
- Generación y refresco de tokens JWT
- Gestión de usuarios (operaciones CRUD)
- Integración con base de datos (PostgreSQL/MySQL/SQLite)
- Endpoints de API RESTful

## Prerrequisitos

- Go 1.19 o posterior
- Base de datos (PostgreSQL, MySQL o SQLite)
- Credenciales de OAuth2 de Google

## Instalación

1. Clona el repositorio:
   ```bash
   git clone <repository-url>
   cd auth_microservice
   ```

2. Instala las dependencias:
   ```bash
   go mod download
   ```

3. Crea el archivo de entorno:
   ```bash
   cp example.env .env
   ```
   Edita `.env` y completa con tus valores reales para las credenciales de la base de datos, ID/secreto del cliente de OAuth2 de Google y otras configuraciones.

4. Configura tu base de datos y actualiza la configuración en `configs/configs.go`.

5. Configura OAuth2 de Google:
   - Crea un proyecto en Google Cloud Console
   - Crea credenciales OAuth2
   - Actualiza `auth/google.go` con tu ID y secreto del cliente

6. Ejecuta la aplicación:
   ```bash
   go run main.go
   ```

## Endpoints de la API

- `GET api/auth` - Inicia el login con OAuth de Google
- `GET api/auth/google/callback` - Maneja el callback de OAuth
- `GET api/user/:id` - Obtiene usuario por ID
- `DELETE api/user` - Elimina el usuario actual
- `PUT api/user` - Actualiza el usuario actual
- `POST api/refresh` - Refresca el token JWT

## Contribuyendo

1. Haz un fork del repositorio
2. Crea una rama de características
3. Haz tus cambios
4. Envía un pull request

## Licencia

Este proyecto está licenciado bajo la Licencia MIT.
