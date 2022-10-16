package utils

import (
	"github.com/prisma/prisma-client-go/binaries"
	"github.com/prisma/prisma-client-go/binaries/platform"
)

func GetEnginePath(dir string, name string) string {
	binaryName := platform.BinaryPlatformName()
	enginePath := binaries.GetEnginePath(dir, name, binaryName) // 构建引擎路径
	return enginePath
}

func GetQueryEnginePath(dir string) string {
	return GetEnginePath(dir, "query-engine")
}

func GetIntrospectionEnginePath(dir string) string {
	return GetEnginePath(dir, "introspection-engine")
}
