@echo off
set JAVA_HOME=C:\Users\Osher\.jdks\ms-17.0.16
set KC_BOOTSTRAP_ADMIN_USERNAME=admin
set KC_BOOTSTRAP_ADMIN_PASSWORD=admin
cd "C:\Keycloak\bin"
kc.bat start-dev --http-port 8180