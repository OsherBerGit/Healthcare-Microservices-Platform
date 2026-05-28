@echo off
title Apache Kafka (KRaft Mode)
set JAVA_HOME=C:\Users\Osher\.jdks\ms-17.0.16

cd /d "C:\Kafka"
echo Starting Kafka in KRaft mode...
.\bin\windows\kafka-server-start.bat C:\Kafka\config\server.properties

pause