#!/bin/bash

set -e 

echo "🚀 Iniciando la instalación de la API de Canciones..."

# 1️⃣ 
echo "🔄 Actualizando paquetes del sistema..."
sudo apt update -y && sudo apt upgrade -y

# 2️⃣ 
echo "📦 Instalando dependencias..."
sudo apt install -y curl unzip

# 3️⃣ 
if ! command -v docker &> /dev/null
then
    echo "🐳 Instalando Docker..."
    curl -fsSL https://get.docker.com | sudo bash
    sudo systemctl enable docker
    sudo systemctl start docker
else
    echo "✅ Docker ya está instalado."
fi

# 4️⃣ 
if ! command -v docker-compose &> /dev/null
then
    echo "🛠️ Instalando Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
else
    echo "✅ Docker Compose ya está instalado."
fi

# 5️⃣ 
echo "📂 Descomprimiendo archivos..."
unzip canciones_go.zip -d canciones_go

# 6️⃣ 
cd canciones_go

# 7️⃣ 
echo "🐳 Cargando imagen de la API..."
docker load -i canciones_go-app.tar

# 9️⃣ 
echo "🚀 Iniciando contenedores con Docker Compose..."
docker-compose up -d

# 🔟 
echo "🔍 Verificando servicios..."
docker ps

echo "✅ Instalación completada. La API está corriendo."

echo "🌎 Puedes probar la API en:"
echo "📌 A través de Nginx: http://localhost/login"
echo "📌 Directamente en Go: http://localhost:3000/login"

