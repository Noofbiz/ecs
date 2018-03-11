package ecs

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MySystemOneAble interface {
	Identifier
	MyComponent1Face
}
type NotSystemOneAble interface {
	NotSystem1Face
}
type MySystemOne struct{}

func (*MySystemOne) Priority() int { return 0 }
func (*MySystemOne) New(*World)    {}
func (*MySystemOne) Update(dt float32, entities []Identifier) {
	for _, e := range entities {
		entity := e.(MySystemOneAble)
		entity.GetMyComponent1().A++
	}
}
func (*MySystemOne) ComponentFilter(id Identifier) bool {
	if _, not := id.(NotSystemOneAble); not {
		return false
	}
	_, ok := id.(MySystemOneAble)
	return ok
}

type MySystemOneTwoAble interface {
	Identifier
	MyComponent1Face
	MyComponent2Face
}
type NotSystemOneTwoAble interface {
	NotSystem12Face
}
type MySystemOneTwo struct{}

func (*MySystemOneTwo) Priority() int { return 0 }
func (*MySystemOneTwo) New(*World)    {}
func (*MySystemOneTwo) Update(dt float32, entities []Identifier) {
	for _, e := range entities {
		entity := e.(MySystemOneTwoAble)
		entity.GetMyComponent1().B++
		entity.GetMyComponent2().D++
	}
}
func (*MySystemOneTwo) ComponentFilter(id Identifier) bool {
	if _, not := id.(NotSystemOneTwoAble); not {
		return false
	}
	_, ok := id.(MySystemOneTwoAble)
	return ok
}

type MySystemTwoAble interface {
	Identifier
	MyComponent2Face
}
type NotSystemTwoAble interface {
	NotSystem2Face
}
type MySystemTwo struct{}

func (*MySystemTwo) Priority() int { return 0 }
func (*MySystemTwo) New(*World)    {}
func (*MySystemTwo) Update(dt float32, entities []Identifier) {
	for _, e := range entities {
		entity := e.(MySystemTwoAble)
		entity.GetMyComponent2().C++
	}
}
func (*MySystemTwo) ComponentFilter(id Identifier) bool {
	if _, not := id.(NotSystemTwoAble); not {
		return false
	}
	_, ok := id.(MySystemTwoAble)
	return ok
}

type MyComponent1 struct {
	A, B int
}

func (c *MyComponent1) GetMyComponent1() *MyComponent1 {
	return c
}

type MyComponent1Face interface {
	GetMyComponent1() *MyComponent1
}

type MyComponent2 struct {
	C, D int
}

func (c *MyComponent2) GetMyComponent2() *MyComponent2 {
	return c
}

type MyComponent2Face interface {
	GetMyComponent2() *MyComponent2
}

type NotSystem1Component struct{}

func (c *NotSystem1Component) GetNotSystem1Component() *NotSystem1Component {
	return c
}

type NotSystem1Face interface {
	GetNotSystem1Component() *NotSystem1Component
}

type NotSystem12Component struct{}

func (c *NotSystem12Component) GetNotSystem12Component() *NotSystem12Component {
	return c
}

type NotSystem12Face interface {
	GetNotSystem12Component() *NotSystem12Component
}

type NotSystem2Component struct{}

func (c *NotSystem2Component) GetNotSystem2Component() *NotSystem2Component {
	return c
}

type NotSystem2Face interface {
	GetNotSystem2Component() *NotSystem2Component
}

// TestCreateEntity ensures IDs which are created, are unique
func TestCreateEntity(t *testing.T) {
	e1 := struct {
		BasicEntity
	}{
		NewBasic(),
	}

	e2 := struct {
		BasicEntity
	}{
		NewBasic(),
	}

	assert.NotEqual(t, e1.id, e2.id, "BasicEntity IDs should be unique")
}

// TestChangeableComponents ensures that Components which are being referenced, are changeable
func TestChangeableComponents(t *testing.T) {
	w := &World{}

	w.AddSystem(&MySystemOne{})

	e1 := struct {
		BasicEntity
		*MyComponent1
	}{
		NewBasic(),
		&MyComponent1{},
	}

	w.AddEntity(&e1)

	w.Update(0.125)
	assert.NotZero(t, e1.MyComponent1.A, "MySystemOne should have been able to change the value of MyComponent1.A")
}

// TestSystemEntityFiltering checks that entities go into the right systems and the flags are obeyed
func TestSystemEntityFiltering(t *testing.T) {
	w := &World{}

	w.AddSystem(&MySystemOne{})
	w.AddSystem(&MySystemTwo{})
	w.AddSystem(&MySystemOneTwo{})

	e1 := struct {
		BasicEntity
		*MyComponent1
	}{
		NewBasic(),
		&MyComponent1{},
	}
	w.AddEntity(&e1)

	e2 := struct {
		BasicEntity
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent2{},
	}
	w.AddEntity(&e2)

	e12 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
	}
	w.AddEntity(&e12)

	e12x1 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
		*NotSystem1Component
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
		&NotSystem1Component{},
	}
	w.AddEntity(&e12x1)

	e12x2 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
		*NotSystem2Component
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
		&NotSystem2Component{},
	}
	w.AddEntity(&e12x2)

	e12x1x2 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
		*NotSystem1Component
		*NotSystem2Component
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
		&NotSystem1Component{},
		&NotSystem2Component{},
	}
	w.AddEntity(&e12x1x2)

	w.Update(0.125)

	assert.Equal(t, 1, e1.A, "e1 was not updated by system 1")
	assert.Equal(t, 0, e1.B, "e1 was updated by system 12")

	assert.Equal(t, 1, e2.C, "e2 was not updated by system 2")
	assert.Equal(t, 0, e2.D, "e2 was updated by system 12")

	assert.Equal(t, 1, e12.A, "e12 was not updated by system 1")
	assert.Equal(t, 1, e12.B, "e12 was not updated by system 12")
	assert.Equal(t, 1, e12.C, "e12 was not updated by system 2")
	assert.Equal(t, 1, e12.D, "e12 was not updated by system 12")

	assert.Equal(t, 0, e12x1.A, "e12x1 was updated by system 1")
	assert.Equal(t, 1, e12x1.B, "e12x1 was not updated by system 12")
	assert.Equal(t, 1, e12x1.C, "e12x1 was not updated by system 2")
	assert.Equal(t, 1, e12x1.D, "e12x1 was not updated by system 12")

	assert.Equal(t, 1, e12x2.A, "e12x2 was not updated by system 1")
	assert.Equal(t, 1, e12x2.B, "e12x2 was not updated by system 12")
	assert.Equal(t, 0, e12x2.C, "e12x2 was updated by system 2")
	assert.Equal(t, 1, e12x2.D, "e12x2 was not updated by system 12")

	assert.Equal(t, 0, e12x1x2.A, "e12x1x2 was updated by system 1")
	assert.Equal(t, 1, e12x1x2.B, "e12x1x2 was not updated by system 12")
	assert.Equal(t, 0, e12x1x2.C, "e12x1x2 was updated by system 2")
	assert.Equal(t, 1, e12x1x2.D, "e12x1x2 was not updated by system 12")
}

// TestDelete tests a commonly used method for removing an entity from the list of entities
func TestDelete(t *testing.T) {
	const maxEntities = 10

	for j := 1; j < maxEntities; j++ {
		w := &World{}

		w.AddSystem(&MySystemOneTwo{})

		var entities []BasicEntity

		// Add all of them
		for i := 0; i < maxEntities; i++ {
			e := struct{ BasicEntity }{NewBasic()}
			w.AddEntity(&e)
			entities = append(entities, e.BasicEntity) // in order to remove it without having a reference to e
		}

		before := len(w.entities)

		// Attempt to remove j
		w.RemoveEntity(entities[j])

		assert.Len(t, w.entities, before-1, "World should now have exactly one less Entity")
	}
}

// TestAddRemoveUpdate tests if the add and remove work after an update
func TestAddRemoveUpdate(t *testing.T) {
	w := &World{}

	w.AddSystem(&MySystemOne{})

	e1 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
	}
	w.AddEntity(&e1)

	e2 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
	}
	w.AddEntity(&e2)

	w.Update(0.125)

	w.AddSystem(&MySystemTwo{})

	e3 := struct {
		BasicEntity
		*MyComponent1
		*MyComponent2
	}{
		NewBasic(),
		&MyComponent1{},
		&MyComponent2{},
	}
	w.AddEntity(&e3)

	w.Update(0.125)

	w.RemoveEntity(&e2)

	w.Update(0.125)

	assert.Equal(t, 3, e1.A, "e1 was not updated by system 1 3 times")
	assert.Equal(t, 2, e1.C, "e12 was not updated by system 2 2 times")

	assert.Equal(t, 2, e2.A, "e2 was not updated by system 1 2 times")
	assert.Equal(t, 1, e2.C, "e2 was not updated by system 2")

	assert.Equal(t, 2, e3.A, "e3 was not updated by system 1 2 times")
	assert.Equal(t, 2, e3.C, "e3 was not updated by system 2")
}

type MyEntity struct {
	BasicEntity
}

// TestIdentifierInterface makes sure that my entity can be stored as an Identifier interface
func TestIdentifierInterface(t *testing.T) {
	e1 := MyEntity{}
	e1.BasicEntity = NewBasic()

	var slice = []Identifier{e1}

	_, ok := slice[0].(MyEntity)
	assert.True(t, ok, "MyEntity should have been recoverable from the Identifier interface")
}

func TestSortableIdentifierSlice(t *testing.T) {
	e1 := MyEntity{}
	e1.BasicEntity = NewBasic()
	e2 := MyEntity{}
	e2.BasicEntity = NewBasic()

	var entities IdentifierSlice = []Identifier{e2, e1}
	sort.Sort(entities)
	assert.ObjectsAreEqual(e1, entities[0])
	assert.ObjectsAreEqual(e2, entities[1])
}

func BenchmarkIdiomatic(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&MySystemOne{})

		e1 := struct {
			BasicEntity
			*MyComponent1
		}{
			NewBasic(),
			&MyComponent1{},
		}
		w.AddEntity(&e1)
	}

	Bench(b, preload, setup)
}

func BenchmarkIdiomaticDouble(b *testing.B) {
	preload := func() {}
	setup := func(w *World) {
		w.AddSystem(&MySystemOneTwo{})

		e12 := struct {
			BasicEntity
			*MyComponent1
			*MyComponent2
		}{
			NewBasic(),
			&MyComponent1{},
			&MyComponent2{},
		}
		e12.BasicEntity = NewBasic()

		w.AddEntity(&e12)
	}

	Bench(b, preload, setup)
}

// Bench is a helper-function to easily benchmark one frame, given a preload / setup function
func Bench(b *testing.B, preload func(), setup func(w *World)) {
	w := &World{}

	preload()
	setup(w)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w.Update(1 / 120) // 120 fps
	}
}

func BenchmarkNewBasic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasic()
	}
}

func BenchmarkNewBasics1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(1)
	}
}

func BenchmarkNewBasic10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			NewBasic()
		}
	}
}

func BenchmarkNewBasics10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(10)
	}
}

func BenchmarkNewBasic100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			NewBasic()
		}
	}
}

func BenchmarkNewBasics100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(100)
	}
}

func BenchmarkNewBasic1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			NewBasic()
		}
	}
}

func BenchmarkNewBasics1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBasics(1000)
	}
}
