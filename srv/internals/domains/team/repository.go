package team

import "sync"

type Repository struct {
	mu     sync.RWMutex
	teams  map[string]Team
	nextId int
}

func NewRepository() *Repository {
	return &Repository{
		teams:  make(map[string]Team),
		nextId: 1,
	}
}

func (r *Repository) List() []Team {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Team, 0, len(r.teams))
	for _, team := range r.teams {
		result = append(result, team)
	}
	return result
}

func (r *Repository) Get(slug string) (Team, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	team, ok := r.teams[slug]
	return team, ok
}

func (r *Repository) Create(team Team) (Team, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.teams[team.Slug]; ok {
		return Team{}, false
	}
	team = r.prepareForSave(team)
	r.teams[team.Slug] = team
	return team, true
}

func (r *Repository) Update(slug string, team Team) (Team, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.teams[slug]
	if !ok {
		return Team{}, false
	}
	team.Slug = slug
	if team.TeamId == 0 {
		team.TeamId = existing.TeamId
	}
	team = r.prepareForSave(team)
	r.teams[slug] = team
	return team, true
}

func (r *Repository) Delete(slug string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.teams[slug]; !ok {
		return false
	}
	delete(r.teams, slug)
	return true
}

func (r *Repository) prepareForSave(team Team) Team {
	if team.TeamId == 0 {
		team.TeamId = r.nextId
	}
	if team.TeamId >= r.nextId {
		r.nextId = team.TeamId + 1
	}
	if team.Hours == nil {
		team.Hours = []int{0, 24}
	}
	return team
}
