import requests

endpoint = "http://localhost:8000/api/endpoint/1/"

response = requests.get(endpoint)
print(response.text)
print(response.status_code)