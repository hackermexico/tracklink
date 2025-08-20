package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp string  `json:"timestamp"`
	IP        string  `json:"ip"`
	UserAgent string  `json:"userAgent"`
}

var locations []LocationData

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/location", locationHandler)
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/static/", staticHandler)

	fmt.Printf("Servidor iniciado en puerto %s\n", port)
	fmt.Println("Configurando túnel público...")

	// Intentar crear túnel con diferentes servicios
	go createTunnel(port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TrackLink - Herramienta de Geolocalización</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            text-align: center;
            color: #333;
            margin-bottom: 30px;
        }
        .location-info {
            background: #e8f5e8;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
            display: none;
        }
        .error {
            background: #ffe8e8;
            color: #d00;
            padding: 15px;
            border-radius: 8px;
            margin: 20px 0;
            display: none;
        }
        .btn {
            background: #007bff;
            color: white;
            padding: 12px 24px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            margin: 10px 5px;
        }
        .btn:hover {
            background: #0056b3;
        }
        .btn:disabled {
            background: #ccc;
            cursor: not-allowed;
        }
        .loading {
            display: none;
            text-align: center;
            margin: 20px 0;
        }
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #3498db;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 2s linear infinite;
            margin: 0 auto;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
        .data-table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
        }
        .data-table th, .data-table td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        .data-table th {
            background-color: #f2f2f2;
        }
        .map-link {
            color: #007bff;
            text-decoration: none;
        }
        .map-link:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🌍 TrackLink - Herramienta de Geolocalización</h1>
            <p>Obtén la ubicación actual del dispositivo</p>
        </div>

        <div style="text-align: center;">
            <button class="btn" onclick="getLocation()">📍 Obtener Mi Ubicación</button>
            <button class="btn" onclick="loadData()">📊 Ver Datos Guardados</button>
        </div>

        <div class="loading" id="loading">
            <div class="spinner"></div>
            <p>Obteniendo ubicación...</p>
        </div>

        <div class="error" id="error"></div>

        <div class="location-info" id="locationInfo">
            <h3>📍 Ubicación Detectada:</h3>
            <div id="locationDetails"></div>
        </div>

        <div id="dataContainer"></div>
    </div>

    <script>
        function getLocation() {
            if (!navigator.geolocation) {
                showError("La geolocalización no es soportada por este navegador.");
                return;
            }

            document.getElementById('loading').style.display = 'block';
            document.getElementById('error').style.display = 'none';
            document.getElementById('locationInfo').style.display = 'none';

            navigator.geolocation.getCurrentPosition(
                function(position) {
                    const locationData = {
                        latitude: position.coords.latitude,
                        longitude: position.coords.longitude,
                        accuracy: position.coords.accuracy,
                        timestamp: new Date().toISOString()
                    };

                    sendLocationData(locationData);
                },
                function(error) {
                    document.getElementById('loading').style.display = 'none';
                    let errorMsg = "Error desconocido";
                    switch(error.code) {
                        case error.PERMISSION_DENIED:
                            errorMsg = "Permiso de geolocalización denegado por el usuario.";
                            break;
                        case error.POSITION_UNAVAILABLE:
                            errorMsg = "Información de ubicación no disponible.";
                            break;
                        case error.TIMEOUT:
                            errorMsg = "Tiempo de espera agotado para obtener la ubicación.";
                            break;
                    }
                    showError(errorMsg);
                },
                {
                    enableHighAccuracy: true,
                    timeout: 10000,
                    maximumAge: 0
                }
            );
        }

        function sendLocationData(locationData) {
            fetch('/location', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(locationData)
            })
            .then(response => response.json())
            .then(data => {
                document.getElementById('loading').style.display = 'none';
                showLocationInfo(data);
            })
            .catch(error => {
                document.getElementById('loading').style.display = 'none';
                showError('Error al enviar datos: ' + error.message);
            });
        }

        function showLocationInfo(data) {
            var mapsUrl = 'https://www.google.com/maps?q=' + data.latitude + ',' + data.longitude;
            var details = '' +
                '<p><strong>Latitud:</strong> ' + data.latitude + '</p>' +
                '<p><strong>Longitud:</strong> ' + data.longitude + '</p>' +
                '<p><strong>Precisión:</strong> ' + Math.round(data.accuracy) + ' metros</p>' +
                '<p><strong>Timestamp:</strong> ' + new Date(data.timestamp).toLocaleString() + '</p>' +
                '<p><strong>IP:</strong> ' + data.ip + '</p>' +
                '<p><a href="' + mapsUrl + '" target="_blank" class="map-link">Ver en Google Maps</a></p>';
            
            document.getElementById('locationDetails').innerHTML = details;
            document.getElementById('locationInfo').style.display = 'block';
        }

        function showError(message) {
            document.getElementById('error').innerHTML = message;
            document.getElementById('error').style.display = 'block';
        }

        function loadData() {
            fetch('/data')
            .then(response => response.json())
            .then(data => {
                showDataTable(data);
            })
            .catch(error => {
                showError('Error al cargar datos: ' + error.message);
            });
        }

        function showDataTable(data) {
            if (data.length === 0) {
                document.getElementById('dataContainer').innerHTML = '<p>No hay datos guardados.</p>';
                return;
            }

            var table = '<h3>📊 Datos de Geolocalización Guardados:</h3>';
            table += '<table class="data-table">';
            table += '<tr><th>Timestamp</th><th>Latitud</th><th>Longitud</th><th>Precisión</th><th>IP</th><th>Mapa</th></tr>';
            
            data.forEach(function(item) {
                var mapsUrl = 'https://www.google.com/maps?q=' + item.latitude + ',' + item.longitude;
                var timestamp = new Date(item.timestamp).toLocaleString();
                table += '<tr>' +
                    '<td>' + timestamp + '</td>' +
                    '<td>' + item.latitude.toFixed(6) + '</td>' +
                    '<td>' + item.longitude.toFixed(6) + '</td>' +
                    '<td>' + Math.round(item.accuracy) + 'm</td>' +
                    '<td>' + item.ip + '</td>' +
                    '<td><a href="' + mapsUrl + '" target="_blank" class="map-link">Ver</a></td>' +
                '</tr>';
            });
            
            table += '</table>';
            document.getElementById('dataContainer').innerHTML = table;
        }
    </script>
</body>
</html>`

	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var locationData LocationData
	if err := json.NewDecoder(r.Body).Decode(&locationData); err != nil {
		http.Error(w, "Error al decodificar JSON", http.StatusBadRequest)
		return
	}

	// Agregar información adicional
	locationData.IP = getClientIP(r)
	locationData.UserAgent = r.UserAgent()
	if locationData.Timestamp == "" {
		locationData.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Guardar en memoria
	locations = append(locations, locationData)

	// Log para el servidor
	fmt.Printf("Nueva ubicación recibida: Lat: %.6f, Lng: %.6f, IP: %s\n",
		locationData.Latitude, locationData.Longitude, locationData.IP)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locationData)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func getClientIP(r *http.Request) string {
	// Verificar headers de proxy
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	return r.RemoteAddr
}

func createTunnel(port string) {
	time.Sleep(2 * time.Second) // Esperar a que el servidor inicie

	fmt.Println("\n=== Configurando Túnel Público ===")

	// Intentar con ngrok primero
	if tryNgrok(port) {
		return
	}

	// Si ngrok no funciona, intentar con serveo
	if tryServeo(port) {
		return
	}

	// Si ninguno funciona, mostrar instrucciones manuales
	showManualInstructions(port)
}

func tryNgrok(port string) bool {
	fmt.Println("Intentando crear túnel con ngrok...")

	cmd := exec.Command("ngrok", "http", port)
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error al iniciar ngrok: %v\n", err)
		return false
	}

	time.Sleep(3 * time.Second)

	// Intentar obtener la URL pública de ngrok
	resp, err := http.Get("http://localhost:4040/api/tunnels")
	if err != nil {
		fmt.Println("No se pudo conectar a la API de ngrok")
		return false
	}
	defer resp.Body.Close()

	var tunnels struct {
		Tunnels []struct {
			PublicURL string `json:"public_url"`
		} `json:"tunnels"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tunnels); err != nil {
		fmt.Println("Error al decodificar respuesta de ngrok")
		return false
	}

	if len(tunnels.Tunnels) > 0 {
		fmt.Printf("✅ Túnel ngrok creado exitosamente!\n")
		fmt.Printf("🌐 URL pública: %s\n", tunnels.Tunnels[0].PublicURL)
		fmt.Printf("📱 Comparte este link para obtener geolocalizaciones remotas\n\n")
		return true
	}

	return false
}

func tryServeo(port string) bool {
	fmt.Println("Intentando crear túnel con serveo...")

	cmd := exec.Command("ssh", "-R", "80:localhost:"+port, "serveo.net")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error al iniciar serveo: %v\n", err)
		return false
	}

	fmt.Println("Túnel serveo iniciado. Presiona Ctrl+C para detener.")
	fmt.Println("Para usar el túnel, accede a: https://serveo.net")
	fmt.Println("El servidor está escuchando en localhost:" + port)
	fmt.Println("Presiona Ctrl+C para salir.")

	// Mantener el túnel activo
	cmd.Wait()
	return true
}

func showManualInstructions(port string) {
	fmt.Println("❌ No se pudieron crear túneles automáticamente.")
	fmt.Println("Instrucciones manuales:")
	fmt.Println("1. Instala ngrok: https://ngrok.com/download")
	fmt.Println("2. Ejecuta: ngrok http " + port)
	fmt.Println("3. Copia la URL pública generada por ngrok")
	fmt.Println("4. Comparte esa URL con otros usuarios")
}
