run: 
	docker build -t go-concurrency-demo .
	docker run -dit --name go_concurrency_demo go-concurrency-demo
	docker attach go_concurrency_demo