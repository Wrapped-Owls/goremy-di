package graph

import (
	"errors"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type view struct {
	root  *Injector
	binds types.BindEnumerator
}

func (v *view) Edges() []Edge {
	return v.root.state.snapshotEdges()
}

func (v *view) ResolveAll() ([]FailedNode, error) {
	if v.binds == nil {
		return nil, nil
	}

	var failures []FailedNode
	v.binds.ForEachBind(func(tag string, value any) bool {
		keyedBind, isBind := value.(types.AnyBind)
		if !isBind {
			return true // raw instances were already resolved at registration
		}

		node := Node{Key: keyedBind.ElementKey(), Tag: tag}
		if genErr := tryGenerate(keyedBind, v.root.rootedAt(node)); genErr != nil {
			failures = append(failures, FailedNode{Node: node, Err: genErr})
		}
		return true
	})

	reasons := make([]error, len(failures))
	for index, failure := range failures {
		reasons[index] = failure.Err
	}

	return failures, errors.Join(reasons...)
}

func tryGenerate(bind types.AnyBind, retriever types.DependencyRetriever) (err error) {
	defer func() {
		if recovered := remyErrs.FromRecovered(recover()); recovered != nil {
			err = recovered
		}
	}()

	_, err = bind.GenAsAny(retriever)
	return err
}
