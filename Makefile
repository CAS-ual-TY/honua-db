VERSION ?= 0.0.16
NAME ?= HONUA-DB

release:
	go mod tidy
	git add .
	git commit -m "[RELEASE] ${NAME}: changes for v${VERSION}"
	git tag v${VERSION}
	git push origin v${VERSION}
	git push