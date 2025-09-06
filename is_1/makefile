GOCMD=go
SRC=src/*.go

run: enigma.exe
	./enigma.exe ./data/byte.txt  ./data/outputEncr.txt  
	./enigma.exe ./data/outputEncr.txt  ./data/outputDecr.txt  

enigma.exe:
	$(GOCMD) build -o enigma.exe $(SRC)

clean:
	rm -rf *.exe