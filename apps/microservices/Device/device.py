from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)
telemetry_data = []

@app.route('/telemetry', methods=['POST'])
def receive_telemetry():
    data = request.get_json()
    data['timestamp'] = datetime.utcnow().isoformat()
    telemetry_data.append(data)
    return jsonify({
        "message": "Телеметрия получена",
        "data": data,
        "count": len(telemetry_data)
    }), 201

@app.route('/telemetry', methods=['GET'])
def get_telemetry():
    return jsonify({
        "data": telemetry_data,
        "total": len(telemetry_data)
    })

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)