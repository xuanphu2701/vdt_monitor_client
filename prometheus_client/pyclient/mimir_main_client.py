import requests
import os

endpoint = f"http://localhost:7000/api/metrics/"
json = {"query":"sum(prometheus_http_requests_total)"}
header = {"X-Scope-OrgID": "tenant1"}

response = requests.post(endpoint, json=json, headers = header)
print(response.text)
print(response.status_code)
print(response.json())