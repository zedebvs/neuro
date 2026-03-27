@echo off
for %%I in ("%CD%\Data\Lib") do set SHORT_DIR=%%~fsI
set CGO_CFLAGS=-I%SHORT_DIR%
set CGO_LDFLAGS=-L%SHORT_DIR% -lwhisper

go build -o app.exe ./app

set PATH=%CD%\Data\Lib;%PATH%

 .\app