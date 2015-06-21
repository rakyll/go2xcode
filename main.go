// Copyright 2015 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main contains a program that converts a Go package
// into an XCode project that targets iOS.
// Usage:
//
//    go2xcode <pkg> [-o]
//
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var o = flag.String("o", "", "directory to output the Xcode project; by default pwd")

func main() {
	flag.Parse()
	pkg := "."
	if args := flag.Args(); len(args) > 0 {
		pkg = args[0]
	}
	// main.xcodeproj/project.pbxproj
	// main (TODO: build the fat binary)
	// Info.plist
	// Images.xcassets/AppIcon.appiconset/Contents.json
	dir, err := ioutil.TempDir("", "go2xcode")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	// TODO(jbd): Build and add binary.
	layout := map[string][]byte{
		filepath.Join(dir, "main.xcodeproj/project.pbxproj"):                        projPbxproj,
		filepath.Join(dir, "main/Info.plist"):                                       infoPlist,
		filepath.Join(dir, "main/Images.xcassets/AppIcon.appiconset/Contents.json"): contentsJSON,
	}

	for dst, v := range layout {
		if err := touchAndWrite(dst, v); err != nil {
			log.Fatalf("failed to touch and write to %v; err = %v", dst, err)
		}
	}

	// TODO(jbd): Build darwin/arm and darwin/arm64 cross compilers.
	cmd := exec.Command("go", "build", "-o", filepath.Join(dir, "main-arm"), pkg)
	cmd.Env = []string{
		"GOOS=darwin",
		"GOARCH=arm",
		"CGO_ENABLED=1",
		"GOPATH=" + goEnv("GOPATH"),
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command("go", "build", "-o", filepath.Join(dir, "main-arm64"), pkg)
	cmd.Env = []string{
		"GOOS=darwin",
		"GOARCH=arm64",
		"CGO_ENABLED=1",
		"GOPATH=" + goEnv("GOPATH"),
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"xcrun",
		"lipo",
		"-create",
		filepath.Join(dir, "main-arm"),
		filepath.Join(dir, "main-arm64"),
		"-o", filepath.Join(dir, "main/main"),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command("rm", filepath.Join(dir, "main-arm"), filepath.Join(dir, "main-arm64"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// TODO(jbd): Copy one by one if Rename fails due to temp dir
	// being on a different device.
	if *o == "" {
		*o = "."
	}
	dst := filepath.Join(*o, "xcode")
	if err := os.RemoveAll(dst); err != nil {
		panic(err)
	}
	if err := os.Rename(dir, dst); err != nil {
		panic(err)
	}
}

func touchAndWrite(dst string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0775|os.ModeDir); err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, 0644)
}

func goEnv(name string) string {
	if val := os.Getenv(name); val != "" {
		return val
	}
	val, err := exec.Command("go", "env", name).Output()
	if err != nil {
		panic(err) // the Go tool was tested to work earlier
	}
	return strings.TrimSpace(string(val))
}

var infoPlist = []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleDevelopmentRegion</key>
  <string>en</string>
  <key>CFBundleExecutable</key>
  <string>main</string>
  <key>CFBundleIdentifier</key>
  <string>org.golang.todo.$(PRODUCT_NAME:rfc1034identifier)</string>
  <key>CFBundleInfoDictionaryVersion</key>
  <string>6.0</string>
  <key>CFBundleName</key>
  <string>$(PRODUCT_NAME)</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleShortVersionString</key>
  <string>1.0</string>
  <key>CFBundleSignature</key>
  <string>????</string>
  <key>CFBundleVersion</key>
  <string>1</string>
  <key>LSRequiresIPhoneOS</key>
  <true/>
  <key>UILaunchStoryboardName</key>
  <string>LaunchScreen</string>
  <key>UIRequiredDeviceCapabilities</key>
  <array>
    <string>armv7</string>
  </array>
  <key>UISupportedInterfaceOrientations</key>
  <array>
    <string>UIInterfaceOrientationPortrait</string>
    <string>UIInterfaceOrientationLandscapeLeft</string>
    <string>UIInterfaceOrientationLandscapeRight</string>
  </array>
  <key>UISupportedInterfaceOrientations~ipad</key>
  <array>
    <string>UIInterfaceOrientationPortrait</string>
    <string>UIInterfaceOrientationPortraitUpsideDown</string>
    <string>UIInterfaceOrientationLandscapeLeft</string>
    <string>UIInterfaceOrientationLandscapeRight</string>
  </array>
</dict>
</plist>
`)

var projPbxproj = []byte(`// !$*UTF8*$!
{
  archiveVersion = 1;
  classes = {
  };
  objectVersion = 46;
  objects = {

/* Begin PBXBuildFile section */
    254BB84F1B1FD08900C56DE9 /* Images.xcassets in Resources */ = {isa = PBXBuildFile; fileRef = 254BB84E1B1FD08900C56DE9 /* Images.xcassets */; };
    254BB8681B1FD16500C56DE9 /* main in Resources */ = {isa = PBXBuildFile; fileRef = 254BB8671B1FD16500C56DE9 /* main */; };
/* End PBXBuildFile section */

/* Begin PBXFileReference section */
    254BB83E1B1FD08900C56DE9 /* main.app */ = {isa = PBXFileReference; explicitFileType = wrapper.application; includeInIndex = 0; path = main.app; sourceTree = BUILT_PRODUCTS_DIR; };
    254BB8421B1FD08900C56DE9 /* Info.plist */ = {isa = PBXFileReference; lastKnownFileType = text.plist.xml; path = Info.plist; sourceTree = "<group>"; };
    254BB84E1B1FD08900C56DE9 /* Images.xcassets */ = {isa = PBXFileReference; lastKnownFileType = folder.assetcatalog; path = Images.xcassets; sourceTree = "<group>"; };
    254BB8671B1FD16500C56DE9 /* main */ = {isa = PBXFileReference; lastKnownFileType = "compiled.mach-o.executable"; path = main; sourceTree = "<group>"; };
/* End PBXFileReference section */

/* Begin PBXGroup section */
    254BB8351B1FD08900C56DE9 = {
      isa = PBXGroup;
      children = (
        254BB8401B1FD08900C56DE9 /* main */,
        254BB83F1B1FD08900C56DE9 /* Products */,
      );
      sourceTree = "<group>";
      usesTabs = 0;
    };
    254BB83F1B1FD08900C56DE9 /* Products */ = {
      isa = PBXGroup;
      children = (
        254BB83E1B1FD08900C56DE9 /* main.app */,
      );
      name = Products;
      sourceTree = "<group>";
    };
    254BB8401B1FD08900C56DE9 /* main */ = {
      isa = PBXGroup;
      children = (
        254BB8671B1FD16500C56DE9 /* main */,
        254BB84E1B1FD08900C56DE9 /* Images.xcassets */,
        254BB8411B1FD08900C56DE9 /* Supporting Files */,
      );
      path = main;
      sourceTree = "<group>";
    };
    254BB8411B1FD08900C56DE9 /* Supporting Files */ = {
      isa = PBXGroup;
      children = (
        254BB8421B1FD08900C56DE9 /* Info.plist */,
      );
      name = "Supporting Files";
      sourceTree = "<group>";
    };
/* End PBXGroup section */

/* Begin PBXNativeTarget section */
    254BB83D1B1FD08900C56DE9 /* main */ = {
      isa = PBXNativeTarget;
      buildConfigurationList = 254BB8611B1FD08900C56DE9 /* Build configuration list for PBXNativeTarget "main" */;
      buildPhases = (
        254BB83C1B1FD08900C56DE9 /* Resources */,
      );
      buildRules = (
      );
      dependencies = (
      );
      name = main;
      productName = main;
      productReference = 254BB83E1B1FD08900C56DE9 /* main.app */;
      productType = "com.apple.product-type.application";
    };
/* End PBXNativeTarget section */

/* Begin PBXProject section */
    254BB8361B1FD08900C56DE9 /* Project object */ = {
      isa = PBXProject;
      attributes = {
        LastUpgradeCheck = 0630;
        ORGANIZATIONNAME = Developer;
        TargetAttributes = {
          254BB83D1B1FD08900C56DE9 = {
            CreatedOnToolsVersion = 6.3.1;
          };
        };
      };
      buildConfigurationList = 254BB8391B1FD08900C56DE9 /* Build configuration list for PBXProject "main" */;
      compatibilityVersion = "Xcode 3.2";
      developmentRegion = English;
      hasScannedForEncodings = 0;
      knownRegions = (
        en,
        Base,
      );
      mainGroup = 254BB8351B1FD08900C56DE9;
      productRefGroup = 254BB83F1B1FD08900C56DE9 /* Products */;
      projectDirPath = "";
      projectRoot = "";
      targets = (
        254BB83D1B1FD08900C56DE9 /* main */,
      );
    };
/* End PBXProject section */

/* Begin PBXResourcesBuildPhase section */
    254BB83C1B1FD08900C56DE9 /* Resources */ = {
      isa = PBXResourcesBuildPhase;
      buildActionMask = 2147483647;
      files = (
        254BB8681B1FD16500C56DE9 /* main in Resources */,
        254BB84F1B1FD08900C56DE9 /* Images.xcassets in Resources */,
      );
      runOnlyForDeploymentPostprocessing = 0;
    };
/* End PBXResourcesBuildPhase section */

/* Begin XCBuildConfiguration section */
    254BB85F1B1FD08900C56DE9 /* Debug */ = {
      isa = XCBuildConfiguration;
      buildSettings = {
        ALWAYS_SEARCH_USER_PATHS = NO;
        CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
        CLANG_CXX_LIBRARY = "libc++";
        CLANG_ENABLE_MODULES = YES;
        CLANG_ENABLE_OBJC_ARC = YES;
        CLANG_WARN_BOOL_CONVERSION = YES;
        CLANG_WARN_CONSTANT_CONVERSION = YES;
        CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
        CLANG_WARN_EMPTY_BODY = YES;
        CLANG_WARN_ENUM_CONVERSION = YES;
        CLANG_WARN_INT_CONVERSION = YES;
        CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
        CLANG_WARN_UNREACHABLE_CODE = YES;
        CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
        "CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
        COPY_PHASE_STRIP = NO;
        DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
        ENABLE_STRICT_OBJC_MSGSEND = YES;
        GCC_C_LANGUAGE_STANDARD = gnu99;
        GCC_DYNAMIC_NO_PIC = NO;
        GCC_NO_COMMON_BLOCKS = YES;
        GCC_OPTIMIZATION_LEVEL = 0;
        GCC_PREPROCESSOR_DEFINITIONS = (
          "DEBUG=1",
          "$(inherited)",
        );
        GCC_SYMBOLS_PRIVATE_EXTERN = NO;
        GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
        GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
        GCC_WARN_UNDECLARED_SELECTOR = YES;
        GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
        GCC_WARN_UNUSED_FUNCTION = YES;
        GCC_WARN_UNUSED_VARIABLE = YES;
        IPHONEOS_DEPLOYMENT_TARGET = 8.3;
        MTL_ENABLE_DEBUG_INFO = YES;
        ONLY_ACTIVE_ARCH = YES;
        SDKROOT = iphoneos;
        TARGETED_DEVICE_FAMILY = "1,2";
      };
      name = Debug;
    };
    254BB8601B1FD08900C56DE9 /* Release */ = {
      isa = XCBuildConfiguration;
      buildSettings = {
        ALWAYS_SEARCH_USER_PATHS = NO;
        CLANG_CXX_LANGUAGE_STANDARD = "gnu++0x";
        CLANG_CXX_LIBRARY = "libc++";
        CLANG_ENABLE_MODULES = YES;
        CLANG_ENABLE_OBJC_ARC = YES;
        CLANG_WARN_BOOL_CONVERSION = YES;
        CLANG_WARN_CONSTANT_CONVERSION = YES;
        CLANG_WARN_DIRECT_OBJC_ISA_USAGE = YES_ERROR;
        CLANG_WARN_EMPTY_BODY = YES;
        CLANG_WARN_ENUM_CONVERSION = YES;
        CLANG_WARN_INT_CONVERSION = YES;
        CLANG_WARN_OBJC_ROOT_CLASS = YES_ERROR;
        CLANG_WARN_UNREACHABLE_CODE = YES;
        CLANG_WARN__DUPLICATE_METHOD_MATCH = YES;
        "CODE_SIGN_IDENTITY[sdk=iphoneos*]" = "iPhone Developer";
        COPY_PHASE_STRIP = NO;
        DEBUG_INFORMATION_FORMAT = "dwarf-with-dsym";
        ENABLE_NS_ASSERTIONS = NO;
        ENABLE_STRICT_OBJC_MSGSEND = YES;
        GCC_C_LANGUAGE_STANDARD = gnu99;
        GCC_NO_COMMON_BLOCKS = YES;
        GCC_WARN_64_TO_32_BIT_CONVERSION = YES;
        GCC_WARN_ABOUT_RETURN_TYPE = YES_ERROR;
        GCC_WARN_UNDECLARED_SELECTOR = YES;
        GCC_WARN_UNINITIALIZED_AUTOS = YES_AGGRESSIVE;
        GCC_WARN_UNUSED_FUNCTION = YES;
        GCC_WARN_UNUSED_VARIABLE = YES;
        IPHONEOS_DEPLOYMENT_TARGET = 8.3;
        MTL_ENABLE_DEBUG_INFO = NO;
        SDKROOT = iphoneos;
        TARGETED_DEVICE_FAMILY = "1,2";
        VALIDATE_PRODUCT = YES;
      };
      name = Release;
    };
    254BB8621B1FD08900C56DE9 /* Debug */ = {
      isa = XCBuildConfiguration;
      buildSettings = {
        ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
        INFOPLIST_FILE = main/Info.plist;
        LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
        PRODUCT_NAME = "$(TARGET_NAME)";
      };
      name = Debug;
    };
    254BB8631B1FD08900C56DE9 /* Release */ = {
      isa = XCBuildConfiguration;
      buildSettings = {
        ASSETCATALOG_COMPILER_APPICON_NAME = AppIcon;
        INFOPLIST_FILE = main/Info.plist;
        LD_RUNPATH_SEARCH_PATHS = "$(inherited) @executable_path/Frameworks";
        PRODUCT_NAME = "$(TARGET_NAME)";
      };
      name = Release;
    };
/* End XCBuildConfiguration section */

/* Begin XCConfigurationList section */
    254BB8391B1FD08900C56DE9 /* Build configuration list for PBXProject "main" */ = {
      isa = XCConfigurationList;
      buildConfigurations = (
        254BB85F1B1FD08900C56DE9 /* Debug */,
        254BB8601B1FD08900C56DE9 /* Release */,
      );
      defaultConfigurationIsVisible = 0;
      defaultConfigurationName = Release;
    };
    254BB8611B1FD08900C56DE9 /* Build configuration list for PBXNativeTarget "main" */ = {
      isa = XCConfigurationList;
      buildConfigurations = (
        254BB8621B1FD08900C56DE9 /* Debug */,
        254BB8631B1FD08900C56DE9 /* Release */,
      );
      defaultConfigurationIsVisible = 0;
      defaultConfigurationName = Release;
    };
/* End XCConfigurationList section */
  };
  rootObject = 254BB8361B1FD08900C56DE9 /* Project object */;
}
`)

var contentsJSON = []byte(`{
  "images" : [
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "29x29",
      "scale" : "3x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "40x40",
      "scale" : "3x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "2x"
    },
    {
      "idiom" : "iphone",
      "size" : "60x60",
      "scale" : "3x"
    },
    {
      "idiom" : "ipad",
      "size" : "29x29",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "29x29",
      "scale" : "2x"
    },
    {
      "idiom" : "ipad",
      "size" : "40x40",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "40x40",
      "scale" : "2x"
    },
    {
      "idiom" : "ipad",
      "size" : "76x76",
      "scale" : "1x"
    },
    {
      "idiom" : "ipad",
      "size" : "76x76",
      "scale" : "2x"
    }
  ],
  "info" : {
    "version" : 1,
    "author" : "xcode"
  }
}
`)
