USERNAME = <割り振られた名前>

update_username:
	for f in 'api.yaml' 'ui.yaml'; do \
		sed -i '' 's/NAME/$(USERNAME)/' $$f; \
	done
