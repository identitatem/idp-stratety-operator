// Copyright Red Hat

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	versioned "github.com/identitatem/idp-strategy-operator/api/client/clientset/versioned"
	internalinterfaces "github.com/identitatem/idp-strategy-operator/api/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/identitatem/idp-strategy-operator/api/client/listers/identitatem/v1alpha1"
	identitatemv1alpha1 "github.com/identitatem/idp-strategy-operator/api/identitatem/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// StrategyInformer provides access to a shared informer and lister for
// Strategies.
type StrategyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.StrategyLister
}

type strategyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewStrategyInformer constructs a new informer for Strategy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewStrategyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredStrategyInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredStrategyInformer constructs a new informer for Strategy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredStrategyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.IdentityconfigV1alpha1().Strategies(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.IdentityconfigV1alpha1().Strategies(namespace).Watch(context.TODO(), options)
			},
		},
		&identitatemv1alpha1.Strategy{},
		resyncPeriod,
		indexers,
	)
}

func (f *strategyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredStrategyInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *strategyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&identitatemv1alpha1.Strategy{}, f.defaultInformer)
}

func (f *strategyInformer) Lister() v1alpha1.StrategyLister {
	return v1alpha1.NewStrategyLister(f.Informer().GetIndexer())
}
