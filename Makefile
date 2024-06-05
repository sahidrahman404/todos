build:
	go build -o=/tmp/bin/api ./cmd/api

run: build
	/tmp/bin/api

run/live:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "/tmp/bin/api" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"

test:
	go test -v ./...

atlas/migrate/diff:
	atlas migrate diff ${name} \
	--env local \
  --format '{{ sql . "  " }}'

atlas/migrate/apply:
		atlas migrate apply --env local
