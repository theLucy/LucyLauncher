@echo off
IF NOT "%~1"=="startas" GOTO START

set "app_name=launcher_service"

cls
windres -i "resources\versioninfo_service.coff" -o "service\versioninfo.syso"

cd service

REM go run gen.go

IF "%~2"=="" GOTO DEBUG
IF "%~2"=="run" GOTO DEBUG
IF "%~2"=="release" GOTO RELEASE
IF "%~2"=="rrun" GOTO RELEASE


:DEBUG
title Build [DEBUG]
set "exe_name=%app_name%_debug.exe"
echo Building executable in DEBUG mode..
go build -ldflags "-X main.build_type=dev" -o %exe_name%
echo Finished DEBUG Build.
GOTO POST

:RELEASE
title Build [RELEASE]
set "exe_name=%app_name%.exe"

echo Running test and benchmarks..
go test -bench=. -benchmem

echo Building executable in RELEASE mode..
go build -ldflags "-s -w -H windowsgui" -o %exe_name%
upx --ultra-brute %app_name%.exe
echo Finished RELEASE Build.

:POST
move %exe_name% ..\target\ >nul
IF "%~2"=="run" GOTO RUN
IF "%~2"=="rrun" GOTO RUN

:EXIT
echo Press ENTER to exit..
pause >nul
exit

:RUN
cd ..\target
echo Running app..
.\%exe_name%
GOTO EXIT

:START
cd ..
start build_service.bat startas %~1