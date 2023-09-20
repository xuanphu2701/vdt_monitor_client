from rest_framework.decorators import api_view
from rest_framework.response import Response

from prom_endpoint.models import PromEndpoint
from prom_endpoint.serializers import PromEndpointSerializer

# Create your views here.
@api_view(["GET"])
def api_home(request, *args, **kwargs):
    instances = PromEndpoint.objects.all()
    data = {}
    if instances:
        for i in range(len(instances)):
            data[i] = PromEndpointSerializer(instances[i]).data
    return Response(data)

