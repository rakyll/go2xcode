# go2xcode

go2xcode converts a Go package into a Xcode project that targets iOS 8.3.

Note: Not maintained, use [gomobile](godoc.org/golang.org/x/mobile/cmd/gomobile)

```
$ go get github.com/rakyll/go2xcode
```

## Usage

```
$ go get -d golang.org/x/mobile/example/basic/...
$ go2xcode golang.org/x/mobile/example/basic
$ cd xcode
$ xcodebuild # or open the Xcode project.
```

Note: go2xcode requires Go 1.5 darwin/arm and darwin/arm64 cross compilers.

In order to build these cross compilers, you need to Go 1.5 source.

Go 1.5 requires Go 1.4. Read more about this requirement at
http://golang.org/s/go15bootstrap.
Set GOROOT_BOOTSTRAP to the GOROOT of your existing 1.4 installation or
follow the steps below to checkout go1.4 from the source and build.

```
$ git clone https://go.googlesource.com/go $HOME/go1.4
$ cd $HOME/go1.4
$ git checkout go1.4.1
$ cd src && ./make.bash
```

Clone and build Go 1.5 and add Go 1.5 bin to your path.

```
$ git clone https://go.googlesource.com/go $HOME/go
$ cd $HOME/go/src && ./make.bash
$ export PATH=$PATH:$HOME/go/bin
```

And build the cross compilers.

```
$ GOARM=7 CGO_ENABLED=1 GOARCH=arm CC_FOR_TARGET=`pwd`/../misc/ios/clangwrap.sh \
  CXX_FOR_TARGET=`pwd`/../misc/ios/clangwrap.sh ./make.bash

$ CGO_ENABLED=1 GOARCH=arm64 CC_FOR_TARGET=`pwd`/../misc/ios/clangwrap.sh \
  CXX_FOR_TARGET=`pwd`/../misc/ios/clangwrap.sh ./make.bash
```
