@echo off

:: Percorso del file di configurazione
set "CONFIG_FILE=startconfig.txt"

:: Controlla se il file di configurazione esiste
if not exist "%CONFIG_FILE%" (
    echo [Config] > "%CONFIG_FILE%"
    echo path=%~dp0 >> "%CONFIG_FILE%"
    echo Configuration completed successfully
) else (
    echo Found: starting session
)

:: Leggi il percorso configurato nel file
for /f "tokens=2 delims==" %%A in ('findstr "path" "%CONFIG_FILE%"') do set "PROJECT_PATH=%%A"

:: Rimuove eventuali spazi bianchi finali dal percorso
set "PROJECT_PATH=%PROJECT_PATH: =%"

:: Cambia directory al percorso configurato
cd /d "%PROJECT_PATH%"

:: Verifica se Go è installato nel sistema
where go >nul 2>nul
if errorlevel 1 (
    echo Go = FALSE.
    pause
    exit /b
)

:: Esegui il comando 'go run .'
go run .

:: Mostra un messaggio per evitare che la finestra del terminale si chiuda subito
pause
