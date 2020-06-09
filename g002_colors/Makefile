PROG=g002_colors

all: $(PROG)

$(PROG): $(PROG).go
	go run $(PROG).go

run: all

clean:
	rm -f $(PROG) $(PROG).exe

.PHONY: clean run
