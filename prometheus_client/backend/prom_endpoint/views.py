from rest_framework.decorators import api_view
from rest_framework.response import Response
import os 
from dotenv import load_dotenv
from prometheus_api_client import PrometheusConnect
from .models import PromEndpoint, MimirEndpoint
import json
from datetime import datetime

load_dotenv()

@api_view(['POST'])
def get_prom_metrics(request, *args, **kwargs):
    url = request.data['url']
    prom_instance = PromEndpoint.objects.filter(url=url)
    if not prom_instance.exists():
        return({"msg":"Your url is invalid"})
    prom = PrometheusConnect(url = url, disable_ssl=True)
    query = request.data['query']
    print(query)
    res = prom.custom_query(query=query)
    return Response(res)

@api_view(['POST'])
def get_prom_range_metrics(request, *args, **kwargs):
    url = request.data['url']
    prom_instance = PromEndpoint.objects.filter(url=url)
    if not prom_instance.exists():
        return({"msg":"Your url is invalid"})
    prom = PrometheusConnect(url = url, disable_ssl=True)
    metric_name = request.data['metric_name']
    label_config = json.loads(request.data['label_config'])
    try:
        start_time = datetime.fromtimestamp(int(request.data['start_time']))
        end_time = datetime.fromtimestamp(int(request.data['end_time']))
        res = prom.get_metric_range_data(metric_name=metric_name, label_config=label_config, 
                                        start_time=start_time, end_time=end_time)

    except:
        res = prom.get_metric_range_data(metric_name=metric_name, label_config=label_config)
    return Response(res)

@api_view(['POST'])
def get_mimir_metrics(request, *args, **kwargs):
    tenantId = request.META.get('HTTP_X_SCOPE_ORGID')
    request_url = request.data['url']
    mimir_instance = MimirEndpoint.objects.filter(tenantId=tenantId).first()
    url = mimir_instance.url

    if not mimir_instance or (mimir_instance and url != request_url) :
        return
    
    prom = PrometheusConnect(url = url, headers={"X-Scope-OrgID":tenantId}, disable_ssl=True)
    query = request.data['query']
    print(query)
    res = prom.custom_query(query=query)
    return Response(res)



