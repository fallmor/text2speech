build:
	GOOS=linux go build -o app
	docker build -t gcr.io/kubernetes-285020/say .
	rm -f app

push: 
	gcloud docker -- push gcr.io/kubernetes-285020/say