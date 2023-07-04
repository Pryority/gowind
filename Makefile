
install:
	cd web && npm install

start:
	docker build -t gowind .
	docker run -d --name gowind-app -p 8080:8080 gowind /app/goserver
	cd web && npm run dev

stop:
	docker stop gowind-app
	docker rm gowind-app
