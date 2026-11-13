package command

type Type string

const (
	SET     Type = "SET"
	GET     Type = "GET"
	DEL     Type = "DEL"
	EXPIRE  Type = "EXPIRE"
	TTL     Type = "TTL"
	PERSIST Type = "PERSIST"
	QUIT    Type = "QUIT"
	KEYS    Type = "KEYS"
	EXISTS  Type = "EXISTS"
	PING    Type = "PING"
	INFO    Type = "INFO"
)

func (t Type) String() string {
	return string(t)
}

func All(write bool) []Type {

	if write {
		return []Type{
			SET,
			DEL,
			EXPIRE,
			PERSIST,
		}
	}
	return []Type{
		SET,
		GET,
		DEL,
		EXPIRE,
		TTL,
		PERSIST,
		QUIT,
		KEYS,
		EXISTS,
		PING,
		INFO,
	}
}

func (t Type) IsValid() bool {
	for _, cmd := range All(false) {
		if t == cmd {
			return true
		}
	}

	return false

	// return slices.Contains(All(false), t)
}

func (t Type) IsWriteCommand() bool {
	for _, cmd := range All(true) {
		if t == cmd {
			return true
		}
	}

	return false

	// return slices.Contains(All(true), t)
}
