package ecs

import (
	"reflect"
	"sort"
)

// World contains a bunch of Entities, and a bunch of Systems. It is the
// recommended way to run ecs.
type World struct {
	systems     systems
	entities    []Identifier
	newSystems  []System
	newEntities []Identifier
	sysEntities map[reflect.Type][]Identifier
}

// AddSystem adds the given System to the World, sorted by priority.
func (w *World) AddSystem(system System) {
	if initializer, ok := system.(Initializer); ok {
		initializer.New(w)
	}

	w.systems = append(w.systems, system)
	w.newSystems = append(w.newSystems, system)
	sort.Sort(w.systems)
}

// Systems returns the list of Systems managed by the World.
func (w *World) Systems() []System {
	return w.systems
}

// Update updates each System managed by the World. It is invoked by the engine
// once every frame, with dt being the duration since the previous update.
func (w *World) Update(dt float32) {
	w.updateEntSys()
	for _, system := range w.Systems() {
		system.Update(dt, w.sysEntities[reflect.TypeOf(system)])
	}
}

// AddEntity adds the given entity to the World
func (w *World) AddEntity(id Identifier) {
	w.entities = append(w.entities, id)
	w.newEntities = append(w.newEntities, id)
}

// RemoveEntity removes the entity across all systems.
func (w *World) RemoveEntity(id Identifier) {
	if w.sysEntities == nil {
		w.sysEntities = make(map[reflect.Type][]Identifier)
	}
	idx := -1
	for i, n := range w.entities {
		if n.ID() == id.ID() {
			idx = i
			break
		}
	}
	if idx >= 0 {
		w.entities = append(w.entities[:idx], w.entities[idx+1:]...)
	}
	for _, sys := range w.systems {
		idx = -1
		for i, n := range w.sysEntities[reflect.TypeOf(sys)] {
			if n.ID() == id.ID() {
				idx = i
				break
			}
		}
		if idx >= 0 {
			w.sysEntities[reflect.TypeOf(sys)] = append(w.sysEntities[reflect.TypeOf(sys)][:idx], w.sysEntities[reflect.TypeOf(sys)][idx+1:]...)
		}
	}
}

func (w *World) updateEntSys() {
	if w.sysEntities == nil {
		w.sysEntities = make(map[reflect.Type][]Identifier)
	}
	done := make(map[reflect.Type]struct{})
	for _, sys := range w.newSystems {
		for _, ent := range w.entities {
			if sys.ComponentFilter(ent) {
				w.sysEntities[reflect.TypeOf(sys)] = append(w.sysEntities[reflect.TypeOf(sys)], ent)
				done[reflect.TypeOf(sys)] = struct{}{}
			}
		}
	}
	for _, ent := range w.newEntities {
		for _, sys := range w.systems {
			if _, ok := done[reflect.TypeOf(sys)]; ok {
				continue
			}
			if sys.ComponentFilter(ent) {
				w.sysEntities[reflect.TypeOf(sys)] = append(w.sysEntities[reflect.TypeOf(sys)], ent)
			}
		}
	}
	w.newEntities = make([]Identifier, 0)
	w.newSystems = make([]System, 0)
}
