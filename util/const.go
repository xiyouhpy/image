package util

// LogoDir 压缩的水印图片目录
const LogoDir = "./config/logo/"

// TtfDir 压缩的水印ttf目录
const TtfDir = "./config/ttf/"

// ImgWmDir 水印的结果图片目录
const ImgWmDir = "./data/image/watermark/"

// TmpWmDir 临时文件目录
const TmpWmDir = "/tmp/image/watermark/"

// ImgCutDir 裁剪的结果图片目录
const ImgCutDir = "./data/image/cut/"

// ArrDirs 检查、创建、清理的目录列表
var ArrDirs = []string{
	ImgWmDir,
	TmpWmDir,
	ImgCutDir,
}
