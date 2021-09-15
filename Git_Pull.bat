@echo off
TITLE Pulling updates from git...

:: Print the branch cause ..oooooo fancy!
echo Pulling from branch:
git branch
echo.
git pull
Title Running a pre-commit check
pre-commit run --all-files
cls
echo All done! check if any errors exist!
sleep 5
exit
