from rest_framework.decorators import api_view
from rest_framework.response import Response
import os 
import requests
from dotenv import load_dotenv
import json

load_dotenv()

@api_view(['POST', 'DELETE'])
def jwthandler(request, *args, **kwargs):
    http_method = request.method 
    if http_method == 'POST':
        headers = request.META
        tenantid = headers.get('HTTP_X_SCOPE_ORGID')
        tenantid_json = {"X-Scope-OrgID": tenantid}
        res = requests.post('http://localhost:8081/jwt', headers=tenantid_json)
        return Response(res)
    elif http_method == 'DELETE':
        headers = request.META
        authorization = headers.get('HTTP_AUTHORIZATION')
        authorization_json = {"Authorization": authorization}
        res = requests.delete('http://localhost:8081/jwt', headers=authorization_json)
        return Response(res)
    return Response({"message":"Method not Allowed"})