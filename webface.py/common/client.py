import json

import requests
from django.http import JsonResponse, Http404


COMPONENT_HOST = {
    "gateway":  "http://127.0.0.1:8080",
    "vchain":   "http://127.0.0.1:8081",
    "consumer": "http://127.0.0.1:8082",
    "receiver": "http://127.0.0.1:8083",
}

def discover(component):
    # start a thread to run discover() with some interval
    # r = requests.get(COMPONENT_HOST["gateway"])
    # update_host(r)
    pass

def proxy(host, uri, *args):
    try:
        url = "%s%s" % (host, uri)
        if len(args) > 0:
            url = url % tuple(args)
        return JsonResponse(json.loads(requests.get(url)))
    except:
        raise Http404()


def serviceList(request, pid):
    return proxy(COMPONENT_HOST["consumer"], "/api/v1/%d/services", pid)


def serviceChainList(request, pid):
    return proxy(COMPONENT_HOST["consumer"], "/api/v1/%d/service-chain", pid)


def requestTypeList(request, pid):
    return proxy(COMPONENT_HOST["consumer"], "/api/v1/%d/request-types", pid)


def invokeChainList(request, pid):
    return proxy(COMPONENT_HOST["consumer"], "/api/v1/%d/invoke-chains", pid)


def requestLogList(request, pid):
    return proxy(COMPONENT_HOST["consumer"], "/api/v1/%d/request-logs", pid)
