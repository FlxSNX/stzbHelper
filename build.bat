@echo off

cd web
echo 当前工作目录: %cd%
echo 开始打包前端

call npm run build

cd ../
echo 当前工作目录: %cd%
echo 开始编译后端

rem 创建 dist 目录，如果目录已存在则忽略错误
mkdir dist 2>nul

echo 开始编译windows 64位
rem 编译为 Windows 64 位可执行文件
set GOOS=windows
set GOARCH=amd64
go build -tags="nomsgpack" -ldflags="-s -w" -o dist\stzbHelper-windows-amd64.exe stzbHelper

echo 编译完成，可执行文件已输出到 dist 目录。
pause