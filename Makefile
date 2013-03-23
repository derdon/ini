all:
	@echo "no default target"

cov:
	gocov test ini | gocov report

countTests:
	grep -c "func Test*" *_test.go
