dice:  wordlist.go
	go build .

wordlist.go: eff_large_wordlist.txt
	/bin/echo -ne 'package main\n\nvar eff_wordlist = `' > $@
	cat $< >> $@
	/bin/echo -e '`\n' >> $@
