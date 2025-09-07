GOCMD=go
SRC=src/*.go

run: enigma.exe
	./enigma.exe ./data/input.txt  ./data/outputEncr.txt  
	./enigma.exe ./data/outputEncr.txt  ./data/outputDecr.txt  

enigma.exe: $(SRC)
	$(GOCMD) build -o enigma.exe $(SRC)

clean:
	rm -rf *.exe