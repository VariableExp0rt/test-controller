package main

import (
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/labels"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	ourAnnotation = "controller-example/myannotation"
	myNamespace   = "default"
)

//Controller is a low-level controller that is parameterized by a Config and used in sharedIndexInformer.
type Controller struct {
	deploymentGetter appsv1.DeploymentsGetter
	deploymentLister appslisters.DeploymentLister
	deploymentSynced cache.InformerSynced
	queue            workqueue.RateLimitingInterface
}

//NewController returns a new sample controller - if your Controller stuct is named differently,
//it would be New<whatever the name is>()
func NewController(kubeclient *kubernetes.Clientset, deploymentInformer appsinformers.DeploymentInformer) *Controller {

	c := &Controller{

		deploymentGetter: kubeclient.AppsV1(),
		deploymentLister: deploymentInformer.Lister(),
		deploymentSynced: deploymentInformer.Informer().HasSynced,
		queue:            workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "deploymentlister"),
	}

	//The objects here are the full JSON blobs that are submitted to the API server, perhaps we could have
	//a function that gets the deployments name and only prints if it is a deployment we want to know about.
	deploymentInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				log.Print("deployment added")
				c.getDeploymentByAnnotation(myNamespace)
			},

			UpdateFunc: func(oldObj, newObj interface{}) {
				log.Print("deployment updated")
				c.getDeploymentByAnnotation(myNamespace)
			},

			DeleteFunc: func(obj interface{}) {
				log.Print("deployment deleted")
			},
		},
	)
	return c
}

//Run will run our new controller loop
func (c *Controller) Run(stop <-chan struct{}) {
	var wg sync.WaitGroup

	defer func() {
		// make sure the work queue is shut down which will trigger workers to end
		log.Print("shutting down queue")
		c.queue.ShutDown()

		// wait on the workers
		log.Print("shutting down workers")
		wg.Wait()

		log.Print("workers are all done")
	}()

	log.Print("waiting for cache sync")
	if !cache.WaitForCacheSync(
		stop, c.deploymentSynced,
	) {
		log.Print("timed out waiting for cache sync")
		return
	}
	log.Print("caches are synced")

	// wait until we're told to stop
	log.Print("waiting for stop signal")
	<-stop
	log.Print("received stop signal")
}

func (c *Controller) getDeploymentByAnnotation(ns string) error {
	rawDeployments, err := c.deploymentLister.Deployments(ns).List(labels.Everything())

	if err != nil {
		return err
	}

	for _, dep := range rawDeployments {
		if _, ok := dep.Annotations[ourAnnotation]; ok {
			log.Printf("seen annotation")
		}
	}
	return err
}
