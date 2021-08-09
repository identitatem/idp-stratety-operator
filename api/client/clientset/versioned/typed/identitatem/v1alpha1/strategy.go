// Copyright Contributors to the Open Cluster Management project

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	scheme "github.com/identitatem/idp-strategy-operator/api/client/clientset/versioned/scheme"
	v1alpha1 "github.com/identitatem/idp-strategy-operator/api/identitatem/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// StrategiesGetter has a method to return a StrategyInterface.
// A group's client should implement this interface.
type StrategiesGetter interface {
	Strategies(namespace string) StrategyInterface
}

// StrategyInterface has methods to work with Strategy resources.
type StrategyInterface interface {
	Create(ctx context.Context, strategy *v1alpha1.Strategy, opts v1.CreateOptions) (*v1alpha1.Strategy, error)
	Update(ctx context.Context, strategy *v1alpha1.Strategy, opts v1.UpdateOptions) (*v1alpha1.Strategy, error)
	UpdateStatus(ctx context.Context, strategy *v1alpha1.Strategy, opts v1.UpdateOptions) (*v1alpha1.Strategy, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Strategy, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.StrategyList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Strategy, err error)
	StrategyExpansion
}

// strategies implements StrategyInterface
type strategies struct {
	client rest.Interface
	ns     string
}

// newStrategies returns a Strategies
func newStrategies(c *IdentityconfigV1alpha1Client, namespace string) *strategies {
	return &strategies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the strategy, and returns the corresponding strategy object, and an error if there is any.
func (c *strategies) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Strategy, err error) {
	result = &v1alpha1.Strategy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("strategies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Strategies that match those selectors.
func (c *strategies) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.StrategyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.StrategyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("strategies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested strategies.
func (c *strategies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("strategies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a strategy and creates it.  Returns the server's representation of the strategy, and an error, if there is any.
func (c *strategies) Create(ctx context.Context, strategy *v1alpha1.Strategy, opts v1.CreateOptions) (result *v1alpha1.Strategy, err error) {
	result = &v1alpha1.Strategy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("strategies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(strategy).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a strategy and updates it. Returns the server's representation of the strategy, and an error, if there is any.
func (c *strategies) Update(ctx context.Context, strategy *v1alpha1.Strategy, opts v1.UpdateOptions) (result *v1alpha1.Strategy, err error) {
	result = &v1alpha1.Strategy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("strategies").
		Name(strategy.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(strategy).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *strategies) UpdateStatus(ctx context.Context, strategy *v1alpha1.Strategy, opts v1.UpdateOptions) (result *v1alpha1.Strategy, err error) {
	result = &v1alpha1.Strategy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("strategies").
		Name(strategy.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(strategy).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the strategy and deletes it. Returns an error if one occurs.
func (c *strategies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("strategies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *strategies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("strategies").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched strategy.
func (c *strategies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Strategy, err error) {
	result = &v1alpha1.Strategy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("strategies").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
