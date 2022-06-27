CC = go build

FLAGS =  -ldflags "-w -s"

TARGET = crack

FILE = main.go

all: $(TARGET)

build_all: ${OBJ}
	$(CC) $(FLAGS) -o $(TARGET) $(FILE)

$(TARGET): build_all

clean:
	rm -f $(TARGET)

fclean: clean

re: fclean all

.PHONY: all re clean fclean
