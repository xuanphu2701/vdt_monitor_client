from django.urls import path
from . import views

urlpatterns = [
    path('jwt/', views.jwthandler),
]