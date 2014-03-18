CC=gcc
RM=rm

CFLAGS=-Wall -Wextra -Werror -fmax-errors=10 -std=c99
LDFLAGS=
RMFLAGS=-rfv

NAME=a
EXT=.exe

# No editing should be required beyond this point

BIN=$(NAME)$(EXT)
SRC=out.c

$(BIN): $(SRC)
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $<

clean:
	$(RM) $(RMFLAGS) $(BIN) *.o $(SRC)
