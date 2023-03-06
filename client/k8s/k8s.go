package k8s

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	scalerv1 "scaler/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type K8sClient struct {
	clientSet *kubernetes.Clientset
}

func New() (*K8sClient, error) {
	conf, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, err
	}
	return &K8sClient{clientSet: clientSet}, nil
}

func (k *K8sClient) CreateDeployment(app scalerv1.Application) error {
	ports := make([]corev1.ContainerPort, 0)
	for _, p := range app.Ports {
		ports = append(ports, corev1.ContainerPort{
			Name:          p.Name,
			ContainerPort: p.Port,
		})
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.NameSpace,
			Labels:    app.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &app.Replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": app.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": app.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            app.Name,
							Image:           app.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports:           ports,
						},
					},
				},
			},
		},
	}
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := k.clientSet.AppsV1().Deployments(app.NameSpace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if errors.IsAlreadyExists(err) {
			// 如果 Service 已经存在，则更新 Service
			existingDeploy, err := k.clientSet.AppsV1().Deployments(app.NameSpace).Get(context.Background(), app.Name, metav1.GetOptions{})
			if err != nil {
				log.Log.Error(err, "get deployment error")
				return err
			}
			existingDeploy.Spec = deployment.Spec
			_, err = k.clientSet.AppsV1().Deployments(app.NameSpace).Update(context.Background(), existingDeploy, metav1.UpdateOptions{})
			if err != nil {
				log.Log.Error(err, "update deployment error")
				return err
			}
			return nil
		}
		if err != nil {
			log.Log.Error(err, "create deployment error")
		}
		return err
	})
	if err != nil {
		return err
	}

	//等待 deployment 创建成功
	time.Sleep(5 * time.Second)

	log.Log.Info("deployment created")

	return nil

}

func (k *K8sClient) CreateService(app scalerv1.Application) error {
	servicePorts := make([]corev1.ServicePort, 0)
	for _, p := range app.Ports {
		servicePorts = append(servicePorts, corev1.ServicePort{
			TargetPort: intstr.FromInt(int(p.Port)),
			Protocol:   corev1.ProtocolTCP,
		})
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.NameSpace,
		},
		Spec: corev1.ServiceSpec{
			Ports: servicePorts,
			Selector: map[string]string{
				"app": app.Name,
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		_, err := k.clientSet.CoreV1().Services(app.NameSpace).Create(context.Background(), service, metav1.CreateOptions{})
		if errors.IsAlreadyExists(err) {
			// 如果 Service 已经存在，则更新 Service
			existingService, err := k.clientSet.CoreV1().Services(app.NameSpace).Get(context.Background(), app.Name, metav1.GetOptions{})
			if err != nil {
				log.Log.Error(err, "get service error")
				return err
			}
			existingService.Spec = service.Spec
			_, err = k.clientSet.CoreV1().Services(app.NameSpace).Update(context.Background(), existingService, metav1.UpdateOptions{})
			if err != nil {
				log.Log.Error(err, "update service error")
				return err
			}
			return nil
		}
		if err != nil {
			log.Log.Error(err, "create service error")
		}
		return err
	})
	if err != nil {
		return err
	}

	// 等待 service 创建成功
	time.Sleep(5 * time.Second)

	log.Log.Info("service created")
	return nil
}
