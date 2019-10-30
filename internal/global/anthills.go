package global

import pkg "github.com/gregmus2/ants-pkg"

type Anthills map[*User]map[pkg.Pos]*Anthill

func (ah Anthills) ByUser(user *User) map[pkg.Pos]*Anthill {
	return ah[user]
}

func (ah Anthills) ByPos(pos pkg.Pos) *Anthill {
	for _, anthills := range ah {
		if _, exist := anthills[pos]; exist {
			return anthills[pos]
		}
	}

	return nil
}

func (ah Anthills) FirstByUser(user *User) *Anthill {
	for _, anthill := range ah[user] {
		return anthill
	}

	return nil
}

func (ah Anthills) DeleteByPos(pos pkg.Pos) *Anthill {
	var obj *Anthill
	for i, anthills := range ah {
		if _, exist := anthills[pos]; exist {
			obj = ah[i][pos]
			delete(ah[i], pos)

			return obj
		}
	}

	return nil
}

func (ah Anthills) Add(user *User, pos pkg.Pos, anthill *Anthill) {
	if _, ok := ah[user]; !ok {
		ah[user] = make(map[pkg.Pos]*Anthill)
	}

	ah[user][pos] = anthill
}