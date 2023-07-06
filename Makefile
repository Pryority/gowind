
install:
	cd web && npm install

dev:
	go run cmd/server/main.go

start:
	docker build -t gowind .
	docker run -d --name gowind-app -p 8080:8080 gowind /app/goserver
	cd web && npm run dev

stop:
	docker stop gowind-app
	docker rm gowind-app
