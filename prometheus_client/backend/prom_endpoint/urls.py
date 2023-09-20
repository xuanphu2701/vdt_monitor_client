# urls.py
from django.urls import path, include
from rest_framework.routers import DefaultRouter 
from . import views

# Create a router and register our viewset with it.
# router = DefaultRouter()
# router.register(r'prometheus-endpoints', views.PromEndpointViewSet)

# The API URLs are now determined automatically by the router.
urlpatterns = [
    # path('', include(router.urls)),
    # path('<int:pk>/', views.PromEndpointDetailAPIView.as_view())
    path('normal/', views.get_prom_metrics),
    path('normal/range/', views.get_prom_range_metrics),
    path('custom/', views.get_mimir_metrics),
]