package helm_client

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

type HelmClient struct {
	actionConfig *action.Configuration
}
type InstallRequest struct {
	ReleaseName string
	Namespace   string
	RepoName    string
	RepoURl     string
	ValuesSet   chartutil.Values
}
type UnInstallRequest struct {
	ReleaseName string
	Namespace   string
	Timeout     time.Duration
}

func New() (*HelmClient, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), "", "", func(format string, v ...interface{}) {
		format = fmt.Sprintf("[helm debug] %s\n", format)
		log.Log.Info("helm", v...)
	})
	if err != nil {
		log.Log.Error(err, "helm New")
		return nil, err
	}
	return &HelmClient{actionConfig: actionConfig}, nil
}
func (h *HelmClient) CheckExistence(releaseName string) (bool, error) {
	client := action.NewList(h.actionConfig)
	results, err := client.Run()
	if err != nil {
		log.Log.Error(err, "helm list error")
		return false, err
	}
	found := false

	for _, release := range results {
		log.Log.Info("release", "name", release.Name, "namespace", release.Namespace)
		if release.Name == releaseName {
			found = true
			break
		}
	}
	return found, nil
}
func (h *HelmClient) UnInstall(req UnInstallRequest) error {
	client := action.NewUninstall(h.actionConfig)

	client.Timeout = req.Timeout
	rel, err := client.Run(req.ReleaseName)
	if err != nil {
		log.Log.Error(err, "helm uninstall")
		return err
	}
	log.Log.Info("helm uninstall", "release", rel.Release.Name, "namespace", rel.Release.Namespace, "info", rel.Info)
	return nil
}
func (h *HelmClient) Install(req InstallRequest) error {
	err := addHelmRepository(req.RepoName, req.RepoURl)
	if err != nil {
		log.Log.Error(err, "helm install")
		return err
	}

	// 创建 Helm 客户端
	client := action.NewInstall(h.actionConfig)

	// 设置 Helm 参数
	client.ReleaseName = req.ReleaseName
	client.Namespace = req.Namespace
	chart, err := loader.Load(fmt.Sprintf("%s/%s", req.RepoName, req.ReleaseName))
	if err != nil {
		log.Log.Error(err, "helm install")
		return err
	}

	// 安装 Helm Chart
	rel, err := client.Run(chart, req.ValuesSet)
	if err != nil {
		log.Log.Error(err, "helm install")
		return err
	}
	log.Log.Info("helm install", "release", rel.Name, "namespace", rel.Namespace, "info", rel.Info)
	return nil
}
func addHelmRepository(name, url string) error {
	// 创建 Helm 仓库
	r, err := repo.NewChartRepository(&repo.Entry{
		Name: name,
		URL:  url,
	}, getter.All(&cli.EnvSettings{}))
	if err != nil {
		log.Log.Error(err, "add helm repository failed")
		return err
	}

	// 更新 Helm 仓库索引
	if _, err := r.DownloadIndexFile(); err != nil {
		log.Log.Error(err, "add helm repository failed")
		return err
	}
	return nil
}
