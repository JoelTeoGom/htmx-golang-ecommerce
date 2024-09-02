# Partim d'una imatge de Go oficial
FROM golang:1.20-alpine

# Definim el directori de treball dins del contenidor
WORKDIR /app

# Copiem el codi font al contenidor
COPY . .

# Compilem l'aplicació
RUN go build -o main .

# Exposem el port en el qual l'aplicació escoltarà
EXPOSE 8080

# Definim el punt d'entrada
CMD ["./main"]
