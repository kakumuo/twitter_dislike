from golang:latest

WORKDIR ./backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/*.go ./
COPY backend/dislikes ./dislikes

RUN go build -o twitter-plus-backend

CMD ["twitter-plus-backend"]