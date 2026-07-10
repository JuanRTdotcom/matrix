# Etapa 1: Build
FROM node:20-alpine AS builder

WORKDIR /app

COPY package.json ./
RUN npm install

COPY . .
RUN npm run build

# Etapa 2: Imagen final
FROM node:20-alpine

WORKDIR /app

COPY package.json ./
RUN npm install --omit=dev

COPY --from=builder /app/build ./build

# El adaptador de Node escucha en el puerto definido por PORT (por defecto 3000).
ENV PORT=5173
EXPOSE 5173

CMD ["node", "build/index.js"]
