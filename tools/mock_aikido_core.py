from flask import Flask, request, jsonify
import sys

app = Flask(__name__)

responses = {
    "/config": {},
    "/api/runtime/config": {},
    "/api/runtime/events": {},
}

events = []


def update_response(response_key, response_value):
    responses[response_key] = response_value


@app.route('/config', methods=['GET'])
def get_config():
    return jsonify(responses["/config"])


@app.route('/api/runtime/config', methods=['GET'])
def get_runtime_config():
    return jsonify(responses["/api/runtime/config"])


@app.route('/api/runtime/events', methods=['POST'])
def post_events():
    print(request.get_json())
    events.append(request.get_json())
    return jsonify(responses["/api/runtime/events"])


@app.route('/mock/set/config', methods=['POST'])
def mock_set_config():
    update_response("/config", request.get_json())
    return jsonify({})


@app.route('/mock/set/api/runtime/config', methods=['POST'])
def mock_set_runtime_config():
    update_response("/api/runtime/config", request.get_json())
    return jsonify({})


@app.route('/mock/set/api/runtime/events', methods=['POST'])
def mock_set_events():
    update_response("/api/runtime/events", request.get_json())
    return jsonify({})


@app.route('/mock/get/events', methods=['POST'])
def mock_get_events():
    return jsonify(events)


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("Usage: python mock_server.py <port>")
        sys.exit(1)
    
    port = int(sys.argv[1])
    app.run(host='127.0.0.1', port=port)
