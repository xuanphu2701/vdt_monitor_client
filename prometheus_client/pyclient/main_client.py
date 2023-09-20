import requests
import os

endpoint = f"http://localhost:7000/api/metrics/"
json = {"query":"prometheus_http_requests_total{handler=\"/graph\"}"}

response = requests.post(endpoint, json=json)
print(response.text)
print(response.status_code)
# print(response.json())