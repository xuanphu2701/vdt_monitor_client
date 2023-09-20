from prometheus_api_client import PrometheusConnect
import os
from dotenv import load_dotenv

load_dotenv()

prom = PrometheusConnect(url = "http://localhost:9090", disable_ssl=True)
metric_list = prom.all_metrics()
print(metric_list)
query = 'prometheus_http_requests_total'.replace('"', "'")
print(prom.custom_query(query=query))