from kit import Repository
from flask import Flask, request, jsonify

app = Flask(__name__)
repo = None

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

if __name__ == "__main__":
    # listens on port 5974
    app.run(host="0.0.0.0", port=5974)
