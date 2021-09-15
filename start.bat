@echo off
TITLE PsstRobot
:: Enables virtual env mode and then starts PsstRobot
env\scripts\activate.bat && py -m bot.py
