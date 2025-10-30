@echo off
REM Go Koperasi Management System
REM Windows Batch file for development tasks

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="build" goto build
if "%1"=="run" goto run
if "%1"=="test" goto test
if "%1"=="migrate" goto migrate
if "%1"=="seed" goto seed
if "%1"=="migrate-fresh" goto migrate-fresh
if "%1"=="migrate-drop" goto migrate-drop
if "%1"=="dev-setup" goto dev-setup
if "%1"=="clean" goto clean
if "%1"=="fmt" goto fmt
if "%1"=="lint" goto lint
if "%1"=="dev" goto dev
if "%1"=="docs" goto docs
if "%1"=="security" goto security
if "%1"=="mod-tidy" goto mod-tidy
if "%1"=="mod-vendor" goto mod-vendor
if "%1"=="install-tools" goto install-tools
if "%1"=="quick-start" goto quick-start
if "%1"=="env-info" goto env-info

echo Unknown command: %1
echo Run "make.bat help" to see available commands
goto end

:help
echo Available commands:
echo.
echo   help          - Show this help message
echo   build         - Build the application
echo   run           - Run the application
echo   test          - Run tests
echo   migrate       - Run database migrations
echo   seed          - Run database seeders
echo   migrate-fresh - Drop all tables, migrate, and seed
echo   migrate-drop  - Drop all tables and migrate
echo   dev-setup     - Setup development environment
echo   clean         - Clean build artifacts
echo   fmt           - Format code
echo   lint          - Lint code (requires golangci-lint)
echo   dev           - Start development server with hot reload
echo   docs          - Generate documentation
echo   security      - Run security scan
echo   mod-tidy      - Tidy modules
echo   mod-vendor    - Vendor dependencies
echo   install-tools - Install development tools
echo   quick-start   - Complete setup for new developers
echo   env-info      - Show environment information
echo.
echo Examples:
echo   make.bat run           - Start the application
echo   make.bat migrate-fresh - Fresh database setup
echo   make.bat dev           - Development with hot reload
goto end

:build
echo Building application...
go build -o bin\main.exe cmd\main.go
if %errorlevel% neq 0 (
    echo Build failed!
    exit /b 1
)
echo Build completed successfully!
goto end

:run
echo Running application...
go run cmd\main.go
goto end

:test
echo Running tests...
go test .\... -v
goto end

:migrate
echo Running database migrations...
go run cmd\migrate\main.go
goto end

:seed
echo Running database seeders...
go run cmd\seeder\*.go
goto end

:migrate-fresh
echo Running fresh migration...
go run cmd\migrate\main.go -fresh
goto end

:migrate-drop
echo Dropping tables and migrating...
go run cmd\migrate\main.go -drop
goto end

:dev-setup
echo Setting up development environment...
echo Installing dependencies...
go mod download
if %errorlevel% neq 0 (
    echo Failed to download dependencies!
    exit /b 1
)
echo Running fresh migration...
go run cmd\migrate\main.go -fresh
if %errorlevel% neq 0 (
    echo Migration failed!
    exit /b 1
)
echo Development setup complete!
goto end

:clean
echo Cleaning build artifacts...
if exist bin rmdir /s /q bin
go clean
echo Clean completed!
goto end

:fmt
echo Formatting code...
go fmt .\...
echo Code formatting completed!
goto end

:lint
echo Linting code...
where golangci-lint >nul 2>nul
if %errorlevel% neq 0 (
    echo golangci-lint not installed. Install with:
    echo go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    goto end
)
golangci-lint run
goto end

:dev
echo Starting development server with hot reload...
where air >nul 2>nul
if %errorlevel% neq 0 (
    echo Air not installed. Installing...
    go install github.com/cosmtrek/air@latest
    if %errorlevel% neq 0 (
        echo Failed to install air!
        exit /b 1
    )
)
air
goto end

:docs
echo Generating documentation...
where godoc >nul 2>nul
if %errorlevel% neq 0 (
    echo godoc not installed. Install with:
    echo go install golang.org/x/tools/cmd/godoc@latest
    goto end
)
echo Starting godoc server at http://localhost:6060
start http://localhost:6060
godoc -http=:6060
goto end

:security
echo Running security scan...
where gosec >nul 2>nul
if %errorlevel% neq 0 (
    echo gosec not installed. Install with:
    echo go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
    goto end
)
gosec .\...
goto end

:mod-tidy
echo Tidying modules...
go mod tidy
echo Module tidy completed!
goto end

:mod-vendor
echo Vendoring dependencies...
go mod vendor
echo Vendoring completed!
goto end

:install-tools
echo Installing development tools...
echo Installing air...
go install github.com/cosmtrek/air@latest
echo Installing golangci-lint...
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
echo Installing godoc...
go install golang.org/x/tools/cmd/godoc@latest
echo Installing gosec...
go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
echo Development tools installed successfully!
goto end

:quick-start
echo Starting quick setup for new developers...
call :install-tools
if %errorlevel% neq 0 (
    echo Tool installation failed!
    exit /b 1
)
call :dev-setup
if %errorlevel% neq 0 (
    echo Dev setup failed!
    exit /b 1
)
echo.
echo 🎉 Quick start complete!
echo.
echo Next steps:
echo 1. Configure your .env file
echo 2. Run "make.bat run" to start the application
echo 3. Run "make.bat dev" for hot reload development
echo.
goto end

:env-info
echo Environment Information:
echo.
go version
echo GOOS: %GOOS%
echo GOARCH: %GOARCH%
echo GOPATH: %GOPATH%
echo Working directory: %CD%
echo.
goto end

:end