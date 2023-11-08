all: clean

clean:
	docker-compose down
	docker volume rm $$(docker volume list -q)

.PHONY: clean
	
