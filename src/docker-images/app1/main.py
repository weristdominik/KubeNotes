import os
import json
from functools import wraps
from flask import Flask, Blueprint, jsonify, redirect, session, request, url_for, Response
from keycloak import KeycloakOpenID
from keycloak.exceptions import KeycloakError
from jwcrypto.jwt import JWTExpired

PREFIX = os.getenv("PREFIX", "")
PORT = int(os.getenv("PORT", 5000))
DEBUG = os.getenv("DEBUG", "False") == "True"
KEYCLOAK_SERVER = os.getenv("KEYCLOAK_SERVER", "http://localhost:80")
REALM_NAME = os.getenv("REALM_NAME", "master")
CLIENT_ID = os.getenv("CLIENT_ID")
CLIENT_SECRET = os.getenv("CLIENT_SECRET")
REDIRECT_URI = os.getenv("REDIRECT_URI")

if PREFIX and not PREFIX.startswith("/"):
    PREFIX = "/" + PREFIX
PREFIX = PREFIX.rstrip("/")

keycloak_openid = KeycloakOpenID(
    server_url=KEYCLOAK_SERVER,
    client_id=CLIENT_ID,
    realm_name=REALM_NAME,
    client_secret_key=CLIENT_SECRET
)


def valid_session():
    """
    Returns:
        (is_valid, refreshed)
    In other routes use '@login_required'
    """
    id_token = session.get("id_token")
    refresh_token = session.get("refresh_token")

    if not id_token:
        return False, False

    try:
        keycloak_openid.decode_token(id_token)
        return True, False
    except KeycloakError:
        if not refresh_token:
            session.clear()
            return False, False

    except JWTExpired:
        if not refresh_token:
            session.clear()
            return False, False

        try:
            new_tokens = keycloak_openid.refresh_token(refresh_token)

            session["access_token"] = new_tokens["access_token"]
            session["refresh_token"] = new_tokens["refresh_token"]
            session["id_token"] = new_tokens["id_token"]

            return True, True

        except KeycloakError:
            session.clear()
            return False, False


def login_required(f):
    @wraps(f)
    def wrapper(*args, **kwargs):
        is_valid, _ = valid_session()
        if not is_valid:
            return redirect(url_for('app1.login'))
        return f(*args, **kwargs)

    return wrapper


bp = Blueprint("app1", __name__, url_prefix=PREFIX)


@bp.route("/")
def hello_world():
    return jsonify({"message": "Hello, from app1!"})


@bp.route("/login")
def login():
    auth_url = keycloak_openid.auth_url(
        redirect_uri=REDIRECT_URI,
        scope="openid profile email"
    )
    return redirect(auth_url)


@bp.route("/callback")
def callback():
    code = request.args.get("code")
    if not code:
        return jsonify({"error": "No code provided"}), 400
    try:
        token = keycloak_openid.token(
            grant_type='authorization_code',
            code=code,
            redirect_uri=os.environ['REDIRECT_URI']
        )
        session['id_token'] = token['id_token']
        session['access_token'] = token['access_token']
        session['refresh_token'] = token['refresh_token']
        return redirect(url_for('app1.home'))
    except Exception as e:
        return jsonify({"error": str(e)}), 500


@bp.route("/home")
@login_required
def home():
    id_token = session.get("id_token")
    access_token = session.get("access_token")
    refresh_token = session.get("refresh_token")

    data = {
        "status": "success",
        "message": "Dashboard (app1)",
        "session": {
            "id_token": {
                "raw": id_token,
                "claims": keycloak_openid.decode_token(id_token)
            },
            "access_token": {
                "raw": access_token,
                "claims": keycloak_openid.decode_token(access_token)
            },
            "refresh_token": refresh_token
        }
    }

    return Response(
        json.dumps(data, indent=4),
        mimetype='application/json',
        status=200
    )


@bp.route("/logout")
@login_required
def logout():
    session.clear()
    return jsonify({"message": "Logout successful"})


app = Flask(__name__)
app.secret_key = os.urandom(24)
app.register_blueprint(bp)

if __name__ == "__main__":
    print(f"Starting app on port {PORT}, debug={DEBUG}, prefix='{PREFIX}'")
    app.run(host="0.0.0.0", port=PORT, debug=DEBUG)
