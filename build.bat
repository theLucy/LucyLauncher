@echo off
IF NOT "%~1"=="startas" GOTO START

set "app_name=go_launcher_app"

cls
windres -i "resources\versioninfo.coff" -o "src\versioninfo.syso"

cd src

go run gen.go

IF "%~2"=="" GOTO DEBUG
IF "%~2"=="run" GOTO DEBUG
IF "%~2"=="release" GOTO RELEASE
IF "%~2"=="rrun" GOTO RELEASE


:DEBUG
title Build [DEBUG]
set "exe_name=%app_name%_debug"
echo Building executable in DEBUG mode..
go build -ldflags "-X main.build_type=dev"
move %app_name%.exe %exe_name%.exe >nul
echo Finished DEBUG Build.
GOTO POST

:RELEASE
title Build [RELEASE]
set "exe_name=%app_name%"

echo Running test and benchmarks..
go test -bench=. -benchmem

echo Building executable in RELEASE mode..
go build -ldflags "-s -w -H windowsgui"
upx --ultra-brute %app_name%.exe
echo Finished RELEASE Build.

:POST
del ..\target\* /Q /S
move %exe_name%.exe ..\target\ >nul
IF "%~2"=="run" GOTO RUN
IF "%~2"=="rrun" GOTO RUN

:EXIT
echo Press ENTER to exit..
pause >nul
exit

:RUN
cd ..\target
echo Running app..
.\%exe_name%.exe
GOTO EXIT

:START
cd ..
start build.bat startas %~1