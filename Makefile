PROG=g001_simplest

all: $(PROG)

$(PROG): $(PROG).go
	go run $(PROG).go

run: all

clean:
	rm -f $(PROG) $(PROG).exe

.PHONY: clean run
