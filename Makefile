# Run templ generation in watch mode
live/templ:
	templ generate --watch --proxy="http://localhost:42068" --proxybind="localhost" --proxyport="42069" --open-browser=true -v

# Run air to detect any go file changes to re-buid an restart the server
live/server:
	go run github.com/cosmtrek/air@v1.52.0 \
		--build.cmd 'go build -o ./dist/main ./src/main.go && rsync -ar ./static ./dist --exclude "*.css"' \
		--build.bin "cd dist && ENV=development ./main" \
		--build.delay "100" \
		--build.exclude_dir "assets,tmp,vendor,testdata,dist,cmd,static,node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "true"

# Run tailwindcss to generate the styles bundle in watch mode
live/css:
	pnpm build:css --watch

# Watch for any js or css changes in the asset folder, then reload the browser via templ proxy 
live/sync_assets:
	go run github.com/cosmtrek/air@v1.52.0 \
		--build.cmd "templ generate --notify-proxy --proxybind='localhost' --proxyport='42069'"  \
		--build.bin "true" \
		--build.delay "100" \
		--build.exclude_dir "" \
		--build.exclude_regex "*.css" \
		--build.include_dir "dist/static"  


live:
	make clean
	mkdir -p dist
	go run ./cmd/app/
	make -j4 live/templ live/server live/css live/sync_assets

prod:
	make clean
	mkdir -p dist 
	go run ./cmd/app/
	templ generate
	make -j2 prod/css prod/build
	rsync -ar ./static ./dist --exclude "*.css"

prod/css:
	pnpm build:css --minify

prod/build:
	go build -o ./dist/main ./src/main.go 


.PHONY: clean
clean:
	find ./src -name '*templ.go' -type f -delete
	find ./src -name '*app.go' -type f -delete
	rm -rf ./dist
