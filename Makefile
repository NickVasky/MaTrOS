mail_listener_run:
	cd cmd/listener && go run main.go

count_lines:
	bash -O globstar -c 'wc -l **/*.go'

count_lines_tests:
	bash -O globstar -c 'wc -l **/*_test.go'
