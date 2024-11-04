@echo off
go mod tidy

SET CGO_ENABLED=0

FOR %%P IN (windows linux darwin) DO (
    FOR %%A IN (amd64 386 arm) DO (
        SET GOOS=%%P
        SET GOARCH=%%A
        IF "%%P"=="windows" (
            go build -o pupload_%%P_%%A.exe
        ) ELSE (
            go build -o pupload_%%P_%%A
        )
    )
)

SET GOOS=windows
SET GOARCH=amd64