GOCMD=go
SRC=src/*.go

run: enigma.exe
	./enigma.exe ./data/output.txt  ./data/input.txt  

enigma.exe:
	$(GOCMD) build -o enigma.exe $(SRC)



clean:
	rm -rf *.exe