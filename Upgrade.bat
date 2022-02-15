@echo off
git pull && powershell -command "Stop-service -Force -name "PsstRobot" -ErrorAction SilentlyContinue; go mod tidy; go build; Start-service -name "PsstRobot""
:: Hail Hydra