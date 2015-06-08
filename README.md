# go2xcode

go2xcode converts a Go package into a Xcode project that targets iOS 8.3.

```
$ go get -d golang.org/x/mobile/example/basic
$ go2xcode golang.org/x/mobile/example/basic
$ cd xcode
$ xcodebuild # or open the Xcode project.
```

Note: go2xcode requires darwin/arm and darwin/arm64 cross compilers.