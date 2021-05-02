import hmac
import hashlib
import json
import time

API_KEY = "API_KEY"
SECRET_KEY = "SECRET_KEY"

req = {
  "id": 11,
  "method": "private/get-order-detail",
  "api_key": API_KEY,
  "params": {
    "order_id": "337843775021233500",
  },
  "nonce": int(time.time() * 1000)
};

# First ensure the params are alphabetically sorted by key
paramString = ""

if "params" in req:
  for key in sorted(req['params']):
    paramString += key
    paramString += str(req['params'][key])

sigPayload = req['method'] + str(req['id']) + req['api_key'] + paramString + str(req['nonce'])

req['sig'] = hmac.new(
  bytes(str(SECRET_KEY), 'utf-8'),
  msg=bytes(sigPayload, 'utf-8'),
  digestmod=hashlib.sha256
).hexdigest()

print(req['nonce']) # 1619956517732
print(req['sig']) # '8c17b4cfbb7073a5453e348ecb1b20a1e709af910a7d1fe4b4569bcf29736e58'