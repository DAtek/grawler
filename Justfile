packages := ". \
./cache \
./page_loader
"

coverfile := ".coverage"


test *options:
    go test {{ options }} {{ packages }}


test-cover *options:
    go test {{ options }} -coverprofile {{ coverfile }} {{ packages }}


show-coverage:
    go tool cover -html {{ coverfile }}
