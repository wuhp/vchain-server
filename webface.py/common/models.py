from django.conf import settings
from django.db import models


class UserProfile(models.Model):
    ref = models.OneToOneField(settings.AUTH_USER_MODEL)


class Repository(models.Model):
    user = models.ForeignKey('UserProfile', null=False)
    name = models.CharField(max_length=100)
    key = models.CharField(max_length=255)
    created = models.DateTimeField()
