package webface;

import (
  "fmt"
  "io/ioutil"
);

var __webfaceCache map[string]string = make(map[string]string);

func LoadStaticContents (relativePath) string {
  return LoadStaticContents0(relativePath, 0);
}

func LoadStaticContents0 (relativePath string, force int) string {
  if (!force && len(__webfaceCache[relativePath]) > 0) {
    return __webfaceCache[relativePath];
  }
  var content []byte;
  content, _ := ioutil.ReadFile(relativePath);
  // TODO: error handle
  __webfaceCache[relativePath] = string(content);
  return __webfaceCache[relativePath];
}
