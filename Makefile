BINARY = cocong
VERSION = v1.2.0

start-dev-db:
	docker-compose -f docker-compose.dev.yaml up -d
stop-dev-db:
	docker-compose -f docker-compose.dev.yaml down
build-linux:
	docker build -t cocong-rhel-builder -f Dockerfile.rhel .;
	docker build -t cocong-debian-builder -f Dockerfile.ubuntu .;

	docker run --platform=linux/amd64 --name cocong-rhel-builder -it -v ./:/build cocong-rhel-builder bash -c "source /root/.bashrc && go build -o cocong_rhel ./cmd/main.go";
	docker run --platform=linux/amd64 --name cocong-debian-builder  -it -v ./:/build cocong-debian-builder bash -c "export PATH=\$$PATH:/usr/local/go/bin && export GOPATH=\$$HOME/go && export PATH=\$$PATH:\$$GOPATH/bin && go build -o cocong_debian ./cmd/main.go";

	mkdir cocong-dist;

	mv cocong_* cocong-dist

	mkdir cocong-dist/cocong-admin
	cp -r cocong-admin/public cocong-dist/cocong-admin
	cp -r cocong-admin/src cocong-dist/cocong-admin
	cp cocong-admin/.eslintrc.json cocong-dist/cocong-admin
	cp cocong-admin/next-env.d.ts cocong-dist/cocong-admin
	cp cocong-admin/next.config.mjs cocong-dist/cocong-admin
	cp cocong-admin/package.json cocong-dist/cocong-admin
	cp cocong-admin/tsconfig.json cocong-dist/cocong-admin

	cp -r templates cocong-dist
	cp cocong.yaml cocong-dist
	cp install.sh cocong-dist
	cp uninstall.sh cocong-dist
	cp dependencies/* cocong-dist

	zip -r cocong-dist.zip cocong-dist

	rm -rf cocong-dist

	docker container rm cocong-rhel-builder
	docker container rm cocong-debian-builder

	docker rmi -f cocong-rhel-builder
	docker rmi -f cocong-debian-builder

build-rhel:
	docker build -t cocong-rhel-builder -f Dockerfile.rhel .;
	docker run --platform=linux/amd64 --name cocong-rhel-builder -it -v ./:/build cocong-rhel-builder bash -c "source /root/.bashrc && go build -o cocong_rhel ./cmd/main.go";
	docker container rm cocong-rhel-builder;