import json
import os
from datetime import datetime
from http.server import HTTPServer, BaseHTTPRequestHandler
import psycopg2
import psycopg2.extras

class DeviceService(BaseHTTPRequestHandler):
    
    def __init__(self, *args, **kwargs):
        # Подключение к базе данных
        self.db_conn = psycopg2.connect(
            dbname=os.getenv('DB_NAME', 'warmhouse'),
            user=os.getenv('DB_USER', 'postgres'),
            password=os.getenv('DB_PASSWORD', 'postgres'),
            host=os.getenv('DB_HOST', 'localhost'),
            port=os.getenv('DB_PORT', '5432')
        )
        self.init_database()
        super().__init__(*args, **kwargs)
    
    def init_database(self):
        """Инициализация таблиц если их нет"""
        try:
            with self.db_conn.cursor() as cursor:
                cursor.execute("""
                    CREATE TABLE IF NOT EXISTS devices (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(100) NOT NULL,
                        location VARCHAR(100) NOT NULL,
                        type VARCHAR(50) NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
                    )
                """)
                
                # Проверяем есть ли данные
                cursor.execute("SELECT COUNT(*) FROM devices")
                count = cursor.fetchone()[0]
                
                if count == 0:
                    cursor.execute("""
                        INSERT INTO devices (name, location, type) VALUES
                        ('Living Room Sensor', 'Living Room', 'temperature'),
                        ('Bedroom Sensor', 'Bedroom', 'temperature'),
                        ('Kitchen Switch', 'Kitchen', 'switch')
                    """)
                
                self.db_conn.commit()
        except Exception as e:
            print(f"Database init error: {e}")
    
    def do_GET(self):
        """Обработка GET запросов"""
        if self.path == '/health':
            self.send_health()
        elif self.path == '/devices':
            self.get_devices()
        elif self.path.startswith('/devices/'):
            device_id = self.path.split('/')[-1]
            if device_id.isdigit():
                self.get_device(int(device_id))
            else:
                self.send_error(400, "Invalid device ID")
        else:
            self.send_error(404, "Endpoint not found")
    
    def do_POST(self):
        """Обработка POST запросов"""
        if self.path == '/devices':
            self.create_device()
        else:
            self.send_error(404, "Endpoint not found")
    
    def send_health(self):
        """Health check endpoint"""
        try:
            with self.db_conn.cursor() as cursor:
                cursor.execute("SELECT 1")
            db_status = "connected"
        except:
            db_status = "disconnected"
        
        response = {
            "status": "ok",
            "service": "device-service-python",
            "database": db_status,
            "timestamp": datetime.now().isoformat()
        }
        
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps(response).encode())
    
    def get_devices(self):
        """Получение списка устройств"""
        try:
            with self.db_conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as cursor:
                cursor.execute("SELECT id, name, location, type, created_at FROM devices ORDER BY id")
                devices = cursor.fetchall()
                
                result = []
                for device in devices:
                    result.append({
                        "id": device["id"],
                        "name": device["name"],
                        "location": device["location"],
                        "type": device["type"],
                        "created_at": device["created_at"].isoformat()
                    })
                
                self.send_response(200)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(result).encode())
        except Exception as e:
            self.send_error(500, f"Database error: {str(e)}")
    
    def get_device(self, device_id):
        """Получение устройства по ID"""
        try:
            with self.db_conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as cursor:
                cursor.execute(
                    "SELECT id, name, location, type, created_at FROM devices WHERE id = %s",
                    (device_id,)
                )
                device = cursor.fetchone()
                
                if device:
                    response = {
                        "id": device["id"],
                        "name": device["name"],
                        "location": device["location"],
                        "type": device["type"],
                        "created_at": device["created_at"].isoformat()
                    }
                    
                    self.send_response(200)
                    self.send_header('Content-Type', 'application/json')
                    self.end_headers()
                    self.wfile.write(json.dumps(response).encode())
                else:
                    self.send_error(404, "Device not found")
        except Exception as e:
            self.send_error(500, f"Database error: {str(e)}")
    
    def create_device(self):
        """Создание нового устройства"""
        try:
            content_length = int(self.headers['Content-Length'])
            post_data = self.rfile.read(content_length)
            data = json.loads(post_data)
            
            # Проверка обязательных полей
            required_fields = ['name', 'location', 'type']
            for field in required_fields:
                if field not in data:
                    self.send_error(400, f"Missing field: {field}")
                    return
            
            with self.db_conn.cursor() as cursor:
                cursor.execute(
                    "INSERT INTO devices (name, location, type) VALUES (%s, %s, %s) RETURNING id, created_at",
                    (data['name'], data['location'], data['type'])
                )
                result = cursor.fetchone()
                self.db_conn.commit()
                
                response = {
                    "id": result[0],
                    "name": data['name'],
                    "location": data['location'],
                    "type": data['type'],
                    "created_at": result[1].isoformat()
                }
                
                self.send_response(201)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(response).encode())
        except json.JSONDecodeError:
            self.send_error(400, "Invalid JSON")
        except Exception as e:
            self.send_error(500, f"Database error: {str(e)}")
    
    def log_message(self, format, *args):
        """Кастомное логирование"""
        print(f"[{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}] {args[0]}")

def run_server():
    """Запуск сервера"""
    port = int(os.getenv('PORT', 8082))
    server_address = ('0.0.0.0', port)
    
    print(f"Starting Device Service (Python) on port {port}")
    print(f"Database: {os.getenv('DB_HOST', 'localhost')}:{os.getenv('DB_PORT', '5432')}/{os.getenv('DB_NAME', 'warmhouse')}")
    
    httpd = HTTPServer(server_address, DeviceService)
    httpd.serve_forever()

if __name__ == '__main__':
    run_server()
