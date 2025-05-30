from kit import Repository
from flask import Flask, request, jsonify

app = Flask(__name__)
repo = None

@app.route("/")
def index():
    return jsonify({"Result": "Heartbeat"}), 200

@app.route("/init", methods=["POST"])
def init():
    data = request.get_json()
    print(data)
    path = data.get("Data").get("path")
    if path is None:
        return jsonify({"error": "path is required"}), 400
    global repo
    repo = Repository(path)
    print(f"Repository initialized at {path}")
    return jsonify({"Result": {"path": path}, "Status": 200}), 200

@app.route("/file_tree", methods=["POST"])
def file_tree():
    if repo is None:
        return jsonify({"error": "repository not initialized"}), 400
    return jsonify({"Result": repo.get_file_tree(), "Status": 200}), 200

@app.route("/find_symbols", methods=["POST"])
def find_symbols():
    data = request.get_json()
    result = repo.find_symbol_usages(
        symbol_name=data.get("Data").get("symbol"),
        symbol_type=data.get("Data").get("symbol_type")
    )
    if result is None:
        result = jsonify({"error": "no results found"})
    return jsonify({"Result": result, "Status": 200}), 200


@app.route("/get_file_content", methods=["POST"])
def get_file_content():
    data = request.get_json()
    result = repo.get_file_content(data.get("Data").get("path"))
    if result is None:
        result = jsonify({"error": "no results found"})
    return jsonify({"Result": result, "Status": 200}), 200

if __name__ == "__main__":
    # listens on port 5974
    print("FLASK SERVER READY")
    app.run(host="0.0.0.0", port=5974, debug=False)
