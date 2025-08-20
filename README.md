# tracklink
Herramienta para geolocalizacion por medio de ngrok y serveo - una con panel web y otra solo en consola. 

🌍 TrackLink - Herramienta de Geolocalización

TrackLink es una herramienta en Go que permite capturar la ubicación geográfica de un dispositivo a través del navegador, mostrando latitud, longitud, precisión, IP y agente de usuario.
Incluye un frontend minimalista con botones para obtener ubicación y ver datos guardados, además de la capacidad de exponer el servidor de manera pública con ngrok o serveo.

⚡ Características

📍 Obtención de coordenadas de geolocalización con precisión.

🌐 Detección de IP pública y User-Agent.

📊 Visualización de datos históricos en tabla.

🗺️ Enlaces directos a Google Maps para ver la ubicación.

🚀 Configuración automática de túneles públicos vía ngrok o serveo.

ServeoTrack - esto es el segundo script en go llamado serveotrack.go

📝 Descripción General

ServeoTrack es una aplicación en Go que levanta un servidor HTTP local con una interfaz falsa de Google.
Su objetivo es obtener datos de ubicación (latitud, longitud, precisión, IP, agente de usuario) de los visitantes mediante la API de geolocalización del navegador y almacenarlos en memoria.

Además:

Redirige búsquedas hacia Google real.

Crea un túnel público mediante Serveo (usando SSH reverso) para exponer el servidor en internet.

Incluye un panel de administración accesible en /admin para visualizar los datos recolectados.

Permite exportar la información en formato JSON desde /admin/data.

🔒 Servidor embebido en Go, sin dependencias externas adicionales.

📦 Requisitos

Go 1.18+
