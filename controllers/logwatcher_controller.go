package controllers

import (
	"bufio"
	"context"
	"regexp"
	"strings"
	"time"

	monitoringv1alpha1 "github.com/yourname/log-watcher/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type LogWatcherReconciler struct {
	client.Client
}

// +kubebuilder:rbac:groups=monitoring.example.com,resources=logwatchers,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;delete
// +kubebuilder:rbac:groups=core,resources=pods/log,verbs=get

func (r *LogWatcherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var lw monitoringv1alpha1.LogWatcher
	if err := r.Get(ctx, req.NamespacedName, &lw); err != nil {
		if kerrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	podList := &corev1.PodList{}
	sel := labels.SelectorFromSet(lw.Spec.PodLabelSelector)

	if err := r.List(ctx, podList, client.InNamespace(lw.Spec.PodNamespace), client.MatchingLabelsSelector{Selector: sel}); err != nil {
		return ctrl.Result{}, err
	}

	pattern, err := regexp.Compile(lw.Spec.MatchPattern)
	if err != nil {
		logger.Error(err, "Invalid regex pattern")
		return ctrl.Result{}, nil
	}

	for _, pod := range podList.Items {
		req := client.ObjectKey{Name: pod.Name, Namespace: pod.Namespace}
		p := &corev1.Pod{}
		if err := r.Get(ctx, req, p); err != nil {
			continue
		}

		logs, err := r.getRecentLogs(ctx, p)
		if err != nil {
			logger.Error(err, "Could not fetch logs", "pod", pod.Name)
			continue
		}

		if pattern.MatchString(logs) {
			logger.Info("Pattern matched, deleting pod", "pod", pod.Name)
			if err := r.Delete(ctx, p); err != nil {
				logger.Error(err, "Failed to delete pod", "pod", pod.Name)
			}
		}
	}

	return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
}

func (r *LogWatcherReconciler) getRecentLogs(ctx context.Context, pod *corev1.Pod) (string, error) {
	podLogOpts := &corev1.PodLogOptions{
		TailLines: int64Ptr(100),
	}
	req := r.Client.RESTClient().Get().
		Namespace(pod.Namespace).
		Name(pod.Name).
		Resource("pods").
		SubResource("log").
		VersionedParams(podLogOpts, client.ParameterCodec)

	stream, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	var logs []string
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		logs = append(logs, scanner.Text())
	}
	return strings.Join(logs, "\n"), scanner.Err()
}

func int64Ptr(i int64) *int64 { return &i }

func (r *LogWatcherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.LogWatcher{}).
		Complete(r)
}
