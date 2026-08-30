FROM node:18

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .

EXPOSE 3000

# esto es para que corra bien dentro del contenedor sin preguntar nada
ENV CI=true
ENV HOST=0.0.0.0

CMD ["npm", "start"]
