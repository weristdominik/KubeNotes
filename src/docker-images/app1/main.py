import os
from flask import Flask, Blueprint, jsonify

PREFIX = os.getenv("PREFIX", "")
PORT = int(os.getenv("PORT", 5000))
DEBUG = os.getenv("DEBUG", "False") == "True"

if PREFIX and not PREFIX.startswith("/"):
    PREFIX = "/" + PREFIX
PREFIX = PREFIX.rstrip("/")

bp = Blueprint("app1", __name__, url_prefix=PREFIX)


@bp.route("/")
def hello_world():
    return jsonify({"message": "Hello, from app1!"})


app = Flask(__name__)
app.register_blueprint(bp)

if __name__ == "__main__":
    print(f"Starting app on port {PORT}, debug={DEBUG}, prefix='{PREFIX}'")
    app.run(host="0.0.0.0", port=PORT, debug=DEBUG)