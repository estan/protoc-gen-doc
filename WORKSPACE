workspace(
    name = "protoc_gen_doc",
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7f1aa43d986df189f7cf30e81dd2dc9d8ed7c74e356341a17267f6d7e5748382",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.1/rules_go-v0.24.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.1/rules_go-v0.24.1.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "d4113967ab451dd4d2d767c3ca5f927fec4b30f3b2c6f8135a2033b9c05a5687",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.0/bazel-gazelle-v0.22.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.0/bazel-gazelle-v0.22.0.tar.gz",
    ],
)

http_archive(
    name = "com_google_protobuf",
    sha256 = "0e8e32d44c9d4572975f43591b51cd3c77392661e4ded17fdfab81e8460344e8",
    strip_prefix = "protobuf-31ebe2ac71400344a5db91ffc13c4ddfb7589f92/",
    urls = [
        "https://github.com/google/protobuf/archive/31ebe2ac71400344a5db91ffc13c4ddfb7589f92.tar.gz",
    ],
)

load(
    "@com_google_protobuf//:protobuf_deps.bzl",
    "protobuf_deps",
)

protobuf_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//defs:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains()

gazelle_dependencies()
