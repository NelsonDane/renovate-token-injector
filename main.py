# Nelson Dane
# Modified From: https://github.com/orgs/community/discussions/48186

import os
import sys
import time
from pathlib import Path

import jwt
import requests

# Get environment variables
if any([not os.getenv("PRIVATE_KEY"), not os.getenv("CLIENT_ID"), not os.getenv("INSTALL_ID")]):
    print("Missing PRIVATE_KEY, CLIENT_ID or INSTALL_ID")
    sys.exit(1)
pem = os.environ["PRIVATE_KEY"]
client_id = os.environ["CLIENT_ID"]
install_id = os.environ["INSTALL_ID"]
out_path = os.getenv("OUT_PATH", "/app/out/config.env")

# Open PEM
signing_key = jwt.jwk_from_pem(bytes(pem, "utf-8"))
payload = {
    # Issued at time
    "iat": int(time.time()),
    # JWT expiration time (10 minutes maximum)
    "exp": int(time.time()) + 600,
    # GitHub App's identifier
    "iss": client_id,
}

# Create JWT
jwt_instance = jwt.JWT()
encoded_jwt = jwt_instance.encode(payload, signing_key, alg="RS256")
print("Created JWT")

# JWT -> GitHub Token
response = requests.post(
    url=f"https://api.github.com/app/installations/{install_id}/access_tokens",
    headers={
        "Authorization": f"Bearer {encoded_jwt}",
        "Accept": "application/vnd.github+json",
    },
    timeout=3,
)
if not response.ok:
    print(f"Bad response: {response.text}")
    sys.exit(1)

# Write result to file
with Path(out_path).open("w", encoding="utf-8") as file:
    file.write(f"GH_TOKEN={response.json()['token']}\n")
print(f"Saved GitHub token to {out_path}")
