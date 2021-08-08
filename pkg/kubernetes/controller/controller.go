package controller

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	crdv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
	clientsetv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1/apis/clientset/versioned"
	informersv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1/apis/informers/externalversions"
	"github.com/gotway/gotway/pkg/log"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Options struct {
	Namespace    string
	ResyncPeriod time.Duration
}

type IngressMatcher = func(*crdv1alpha1.IngressHTTP) bool

var (
	ErrIngressNotFound = errors.New("ingress not found")
)

type Controller struct {
	options Options

	ingresshttpInformer cache.SharedIndexInformer
	ingressMux          sync.RWMutex

	queue  workqueue.RateLimitingInterface
	logger log.Logger
}

func (c *Controller) Run(ctx context.Context) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("starting controller")

	c.logger.Info("starting informer")
	go c.ingresshttpInformer.Run(ctx.Done())

	c.logger.Info("waiting for informer caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), c.ingresshttpInformer.HasSynced) {
		err := errors.New("failed to wait for informers caches to sync")
		utilruntime.HandleError(err)
		return err
	}
	c.logger.Info("controller ready")

	<-ctx.Done()
	c.logger.Info("stopping controller")

	return nil
}

func (c *Controller) ListIngresses() ([]crdv1alpha1.IngressHTTP, error) {
	c.ingressMux.RLock()
	defer c.ingressMux.RUnlock()

	var ingresses []crdv1alpha1.IngressHTTP
	for _, obj := range c.ingresshttpInformer.GetIndexer().List() {
		if ingress, ok := obj.(*crdv1alpha1.IngressHTTP); ok {
			ingresses = append(ingresses, *ingress)
			continue
		}
		return nil, fmt.Errorf("unexpected object %v", obj)
	}
	return ingresses, nil
}

func (c *Controller) FindIngress(matchFn IngressMatcher) (crdv1alpha1.IngressHTTP, error) {
	c.ingressMux.RLock()
	defer c.ingressMux.RUnlock()

	for _, obj := range c.ingresshttpInformer.GetIndexer().List() {
		if ingress, ok := obj.(*crdv1alpha1.IngressHTTP); ok {
			if matchFn(ingress) {
				return *ingress, nil
			}
			continue
		}
		c.logger.Error(fmt.Sprintf("unexpected object %v", obj))
	}
	return crdv1alpha1.IngressHTTP{}, ErrIngressNotFound
}

func (c *Controller) UpdateIngressStatus(
	ctx context.Context,
	ingress crdv1alpha1.IngressHTTP,
	healthy bool,
) error {
	if ingress.Status.IsServiceHealthy == healthy {
		return nil
	}
	c.ingressMux.Lock()
	defer c.ingressMux.Unlock()

	key, err := cache.MetaNamespaceKeyFunc(&ingress)
	if err != nil {
		return err
	}
	obj, exists, err := c.ingresshttpInformer.GetIndexer().GetByKey(key)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("ingress %v not found", ingress.Name)
	}

	i, ok := obj.(*crdv1alpha1.IngressHTTP)
	if !ok {
		return fmt.Errorf("unexpected object %v", obj)
	}
	i.Status.IsServiceHealthy = healthy

	return nil
}

func New(
	options Options,
	ingresshttpClientSet clientsetv1alpha1.Interface,
	logger log.Logger,
) *Controller {

	informerFactory := informersv1alpha1.NewSharedInformerFactory(
		ingresshttpClientSet,
		options.ResyncPeriod,
	)
	ingresshttpInformer := informerFactory.Gotway().V1alpha1().IngressHTTPs().Informer()

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	return &Controller{
		options:             options,
		ingresshttpInformer: ingresshttpInformer,
		queue:               queue,
		logger:              logger,
	}
}
