package types

type Repo struct {
	Owner string
	Repo  string
}

func (r *Repo)GetOwner() string {
	return r.Owner
}

func (r *Repo)GetRepo() string {
	return r.Repo
}

