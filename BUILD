load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")

# gazelle:prefix github.com/groundfoundation/gotabgo
gazelle(name = "gazelle")

go_library(
    name = "gotabgo",
    srcs = [
        "client.go",
        "http.go",
        "httpclient.go",
        "model.go",
        "tabapi.go",
        "types.go",
    ],
    importpath = "github.com/groundfoundation/gotabgo",
    visibility = ["//visibility:public"],
    deps = ["@com_github_sirupsen_logrus//:logrus"],
)
