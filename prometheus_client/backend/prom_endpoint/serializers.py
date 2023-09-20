from rest_framework import serializers
from .models import PromEndpoint

class PromEndpointSerializer(serializers.ModelSerializer):
    class Meta:
        model = PromEndpoint
        fields = '__all__'