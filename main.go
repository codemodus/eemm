package main

func main() {
	cs := newComs()
	trip := tripFn(cs)
	log(cs)

	cs.ic <- "hello"
	// TODO: add sub-command for migration
	cs.ic <- "starting migration tool"

	// TODO: config = slice of dstCnf/srcCnf pairs
	dstConf := sessionConfig{
		server:   "mail.host.invalid",
		port:     "993",
		account:  "dst@example.com",
		password: "invalid",
	}

	srcConf := sessionConfig{
		server:   "mail.host.invalid",
		port:     "993",
		account:  "srv@example.com",
		password: "invalid",
	}

	bs, err := newBondedSession(cs, 11, dstConf, srcConf)
	trip(err)
	defer bs.close()

	trip(bs.sync())

	cs.ic <- "goodbye"

	// TODO: add flag to control concurrency
	// TODO: add flag(s) to restrict message handling to span (i.e. "after", "before")

	// TODO: add sub-command for duplicate removal
	// TODO: config = slice of dstCnf
}
