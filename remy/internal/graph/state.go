package graph

import "sync"

type dedupKey struct {
	fromID, toID   uint64
	fromTag, toTag string
}

func dedupKeyOf(edge Edge) dedupKey {
	return dedupKey{
		fromID: edge.From.Key.ID(), fromTag: edge.From.Tag,
		toID: edge.To.Key.ID(), toTag: edge.To.Tag,
	}
}

type state struct {
	edgeSet map[dedupKey]struct{}
	edges   []Edge
	mutex   sync.RWMutex
}

func newState() *state {
	return &state{edgeSet: map[dedupKey]struct{}{}}
}

func (s *state) recordEdge(edge Edge) {
	key := dedupKeyOf(edge)

	s.mutex.RLock()
	_, seen := s.edgeSet[key]
	s.mutex.RUnlock()
	if seen {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, seen = s.edgeSet[key]; seen {
		return // another chain recorded it between the two locks
	}
	s.edgeSet[key] = struct{}{}
	s.edges = append(s.edges, edge)
}

func (s *state) snapshotEdges() []Edge {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	edges := make([]Edge, len(s.edges))
	copy(edges, s.edges)
	return edges
}
