.PHONY: grpc grpc_clean

generate_grpc: grpc_clean
	protoc --proto_path rpc/ --go_out=plugins=grpc:rpc rpc/account.proto

grpc_clean:
	rm rpc/*.go || true

remove_omitempty:
	ls rpc/*.pb.go | xargs -n1 -IX bash -c "sed -e 's/,omitempty//' X > X.tmp && mv X{.tmp,}"

grpc: grpc_clean generate_grpc remove_omitempty

migrate_account:
	cd services/account && tern migrate --migrations ./migrations

migrate_account_undo:
	cd services/account && tern migrate --migrations ./migrations --destination -1

migrate_order:
	cd services/order && tern migrate --migrations ./migrations

migrate_order_undo:
	cd services/order && tern migrate --migrations ./migrations --destination -1


local_dev:
	docker-compose up db

connect_db:
	docker-compose exec db psql -U postgres