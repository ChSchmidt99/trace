package raytracer

import (
	"github/chschmidt99/pt/pkg/pt"
	"time"
)

type Path interface {
	transformationAtTime(t float64) pt.Matrix4
	duration() time.Duration
}

type Animation struct {
	node         *pt.SceneNode
	path         Path
	startingTime time.Time
	isRepeating  bool
}

func NewAnimation(node *pt.SceneNode, path Path, start time.Time, repeating bool) *Animation {
	return &Animation{
		node:         node,
		path:         path,
		startingTime: start,
		isRepeating:  repeating,
	}
}

func (a Animation) t(t time.Time) float64 {
	return t.Sub(a.startingTime).Seconds() / a.path.duration().Seconds()
}

type animationScheduler struct {
	queue animationQueue
	t     time.Time
}

func newAnimationScheduler() *animationScheduler {
	return &animationScheduler{
		t:     time.Now(),
		queue: animationQueue{},
	}
}

func (context *animationScheduler) queueAnimation(animation *Animation) {
	context.queue.add(animation)
}

func (context *animationScheduler) animationTick(deltaTime time.Duration) {
	context.t = context.t.Add(deltaTime)
	removed := context.queue.removeInActive(context.t)
	for _, a := range removed {
		if a.isRepeating {
			a.startingTime = context.t
			context.queue.add(a)
		}
	}
	active := context.queue.getActive(context.t)
	for _, a := range active {
		a.node.SetTransformation(a.path.transformationAtTime(a.t(context.t)))
	}
}

type waypoint struct {
	coordinate pt.Vector3
	offset     time.Duration // How much time is spent to reach this waypoint from the former one
}

type Line struct {
	origin      waypoint
	destination waypoint
}

func NewLinePath(origin pt.Vector3, destination pt.Vector3, offset time.Duration, duration time.Duration) Line {
	org := waypoint{
		coordinate: origin,
		offset:     offset,
	}
	dst := waypoint{
		coordinate: destination,
		offset:     duration,
	}
	return Line{
		origin:      org,
		destination: dst,
	}
}

func (l Line) transformationAtTime(t float64) pt.Matrix4 {
	direction := l.destination.coordinate.Sub(l.origin.coordinate)
	magnitude := direction.Mul(t)
	position := l.origin.coordinate.Add(magnitude)
	return pt.Translate(position.X, position.Y, position.Z)
}

func (l Line) duration() time.Duration {
	return l.origin.offset + l.destination.offset
}

type Sequence struct {
	waypoints []waypoint
	d         time.Duration
}

func NewUniformSequence(wps []pt.Vector3, offset time.Duration, duration time.Duration) Sequence {
	var totalDistance float64 = 0
	for i := 0; i < len(wps)-1; i++ {
		ab := wps[i+1].Sub(wps[i])
		totalDistance += ab.Length()
	}
	waypoints := make([]waypoint, len(wps))
	waypoints[0] = waypoint{
		coordinate: wps[0],
		offset:     offset,
	}
	for i := 1; i < len(wps); i++ {
		distance := wps[i].Sub(wps[i-1]).Length()
		off := (distance / totalDistance) * float64(duration)
		waypoints[i] = waypoint{
			coordinate: wps[i],
			offset:     time.Duration(off),
		}
	}
	return Sequence{
		waypoints: waypoints,
		d:         duration,
	}
}

func (s Sequence) transformationAtTime(t float64) pt.Matrix4 {
	delta := t * float64(s.d)
	for i := 0; i < len(s.waypoints)-1; i++ {
		if delta > float64(s.waypoints[i+1].offset) {
			delta -= float64(s.waypoints[i+1].offset)
			continue
		}
		tSegment := delta / float64(s.waypoints[i+1].offset)
		wp1 := s.waypoints[i]
		wp2 := s.waypoints[i+1]
		direction := wp2.coordinate.Sub(wp1.coordinate)
		magnitude := direction.Mul(tSegment)
		position := wp1.coordinate.Add(magnitude)
		return pt.Translate(position.X, position.Y, position.Z)
	}
	dst := s.waypoints[len(s.waypoints)-1].coordinate
	return pt.Translate(dst.X, dst.Y, dst.Z)
}

func (s Sequence) duration() time.Duration {
	return s.d
}

type animationQueue struct {
	head *animationQueueNode
}

func (q *animationQueue) add(a *Animation) {
	node := &animationQueueNode{
		a: a,
	}
	if q.head == nil {
		q.head = node
		return
	}
	q.head = q.head.insert(node)
}

func (q *animationQueue) getActive(t time.Time) []*Animation {
	if q.head == nil {
		return nil
	}
	acc := make([]*Animation, 0)
	q.head.collectActive(t, &acc)
	return acc
}

func (q *animationQueue) removeInActive(t time.Time) []*Animation {
	if q.head != nil {
		removed := make([]*Animation, 0)
		q.head = q.head.removeIfInactive(t, &removed)
		return removed
	}
	return nil
}

type animationQueueNode struct {
	a    *Animation
	next *animationQueueNode
}

func (node *animationQueueNode) insert(in *animationQueueNode) *animationQueueNode {
	if in.a.startingTime.Before(node.a.startingTime) {
		in.next = node
		return in
	}
	if node.next == nil {
		node.next = in
	} else {
		node.next = node.next.insert(in)
	}
	return node
}

// Removes this and all following nodes, if their finish time was before t
// Returns the new first node in the chain
func (node *animationQueueNode) removeIfInactive(t time.Time, removedAcc *[]*Animation) *animationQueueNode {
	endTime := node.a.startingTime.Add(node.a.path.duration())
	if endTime.Before(t) {
		*removedAcc = append(*removedAcc, node.a)
		if node.next == nil {
			return nil
		}
		return node.next.removeIfInactive(t, removedAcc)
	}
	if node.next == nil {
		return node
	}
	node.next = node.next.removeIfInactive(t, removedAcc)
	return node
}

func (node *animationQueueNode) collectActive(t time.Time, acc *[]*Animation) {
	if node.a.startingTime.After(t) {
		return
	}
	*acc = append(*acc, node.a)
	if node.next != nil {
		node.next.collectActive(t, acc)
	}
}
