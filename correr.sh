#!/bin/bash

echo "Iniciando Simulador de Red..."

docker compose down

docker compose up --build 

echo "========================================================"
echo "¡Proyecto desplegado con éxito!"
echo "Frontend: http://localhost:3000"
echo "Middleware: http://localhost:8100"
echo "Load Balancer: http://localhost:9100"
echo "Base de datos (PostgreSQL): localhost:5432"
echo "========================================================"
echo "Los servicios se están ejecutando en segundo plano."
echo "Para ver los logs en tiempo real, ejecuta: docker compose logs -f"
echo "Para detener el proyecto, ejecuta: docker compose down"
