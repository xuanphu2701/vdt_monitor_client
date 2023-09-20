from django.db import models

# Create your models here.
class PromEndpoint(models.Model):
    url = models.URLField(max_length=100)
    description = models.TextField(blank=True, null=True)
    
    def __str__(self):
        return self.url
    
class MimirEndpoint(models.Model):
    url = models.URLField(max_length=100)
    description = models.TextField(blank=True, null=True)
    tenantId = models.CharField(max_length=100)
    def __str__(self):
        return self.url