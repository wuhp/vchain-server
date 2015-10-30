import uuid
import json
from datetime import datetime

from django.contrib.auth.models import User
from django.contrib.auth import login, logout, authenticate
from django.views.decorators.csrf import csrf_exempt
from django.http import JsonResponse, Http404

from . import models


def getUserProfile(user):
    if not user.is_authenticated():
        raise "Authentication Denied"
    try:
        return models.UserProfile.objects.get(ref=user)
    except:
        user = models.UserProfile(ref=user)
        user.save()
        return user


def checkUser(user0, user1):
    if user0.id != user1.id:
        raise Exception("Authorization Denied")
    return


def getInput(content):
    try:
        return json.loads(content)
    except:
        return {}


@csrf_exempt
def userRegister(request):
    try:
        obj = getInput(request.body)
        user = User(username=obj["username"], password=obj["password"])
        if User.objects.filter(username=user.username).count() > 0:
            return JsonResponse({"error": 1})
        user = User.objects.create_user(username=user.username, password=user.password)
        profile = getUserProfile(user)
        return JsonResponse({
            "error": 0,
            "id": profile.id,
            "username": user.username,
        })
    except Exception as e:
        print e
        raise Http404()


@csrf_exempt
def userLogin(request):
    try:
        obj = getInput(request.body)
        user = authenticate(username=obj["username"], password=obj["password"])
        if user is None:
            return JsonResponse({"error": 1})
        if not user.is_active:
            return JsonResponse({"error": 1})
        login(request, user)
        return JsonResponse({"error": 0})
    except Exception as e:
        print e
        raise Http404()


@csrf_exempt
def userLogout(request):
    try:
        logout(request)
        return JsonResponse({"error": 0})
    except Exception as e:
        print e
        raise Http404()


@csrf_exempt
def repositoryD0(request):
    if request.method in ["POST"]:
        return repositoryCreate(request)
    if request.method in ["GET"]:
        return repositoryList(request)
    raise Http404()


def repositoryCreate(request):
    try:
        user = getUserProfile(request.user)
        obj = getInput(request.body)
        # TODO: prevent malware or robot create repository infinitely
        if len(obj["name"]) == 0: raise
        repository = models.Repository(
            user=user,
            name=obj["name"],
            key=uuid.uuid4().hex,
            created=datetime.now()
        )
        if models.Repository.objects.filter(
            name=repository.name,
            user=user
        ).count() > 0:
            return JsonResponse({
                "error":   1,
                "message": "Repository Exists",
            })
        repository.save()
        return JsonResponse({
            "error": 0,
            "id":    repository.id,
            "name":  repository.name,
            "key":   repository.key,
        })
    except Exception as e:
        print e
        raise Http404()


def repositoryList(request):
    try:
        filter_cond = {}
        if "name" in request.GET:
            filter_cond["name"] = request.GET["name"]
        if "key" in request.GET:
            filter_cond["key"] = request.GET["key"]
        user = getUserProfile(request.user)
        repositories = models.Repository.objects.filter(
            user=user, **filter_cond
        ).order_by("name", "-created")
        return JsonResponse([
            {
                "id":   one.id,
                "name": one.name,
                "key":  one.key
            } for one in repositories
        ], safe=False)
    except Exception as e:
        print e
        raise Http404()


@csrf_exempt
def repositoryD1(request, id):
    if request.method in ["GET"]:
        return repositoryGet(request, id)
    if request.method in ["PUT", "PATCH"]:
        return repositoryUpdate(request, id)
    if request.method in ["DELETE"]:
        return repositoryDelete(request, id)
    raise Http404()


def repositoryGet(request, id):
    try:
        user = getUserProfile(request.user)
        repository = models.Repository.objects.get(id=id)
        checkUser(user, repository.user)
        return JsonResponse({
            "id":    repository.id,
            "name":  repository.name,
            "key":   repository.key,
        })
    except Exception as e:
        print e
        raise Http404()


def repositoryUpdate(request, id):
    try:
        user = getUserProfile(request.user)
        repository = models.Repository.objects.get(id=id)
        checkUser(user, repository.user)
        obj = getInput(request.body)
        changed = False
        if "name" in obj and len(obj["name"]) > 0:
            repository.name = obj["name"]
            changed = True
        if "key" in obj:
            repository.key = uuid.uuid4().hex
            changed = True
        if changed:
            repository.save()
        return JsonResponse({
            "id":    repository.id,
            "name":  repository.name,
            "key":   repository.key,
        })
    except Exception as e:
        print e
        raise Http404()


def repositoryDelete(request, id):
    try:
        user = getUserProfile(request.user)
        repository = models.Repository.objects.get(id=id)
        checkUser(user, repository.user)
        repository.delete()
        return JsonResponse({
            "name":  repository.name,
        })
    except Exception as e:
        print e
        raise Http404()
