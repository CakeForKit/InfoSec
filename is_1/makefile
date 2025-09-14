GOCMD=go
SRC=src/*.go

run: enigma.exe
	./enigma.exe ./data/input.txt  ./data/outputEncr.txt  
	./enigma.exe ./data/outputEncr.txt  ./data/outputDecr.txt  

arch_run: enigma.exe
	./enigma.exe ./data2/input.rar  ./data2/outputEncr.rar
	./enigma.exe ./data2/outputEncr.rar  ./data2/outputDecr.rar  

enigma.exe: $(SRC)
	$(GOCMD) build -o enigma.exe $(SRC)

clean:
	rm -rf *.exe