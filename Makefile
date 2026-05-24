worker:
	cd ml-worker && PYTHONPATH=src .venv/bin/python3 -m worker.main
api:
	cd srv && go run main.go