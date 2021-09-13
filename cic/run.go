package cic

import (
	"github.com/sirupsen/logrus"
	"runcic/cic/capabilities"
	"runcic/containerimage"
	"runcic/containerimage/common"
)

func mergeCap(CapAdd, CapDrop []string) {
	hashCap := make(map[string]bool)
	for _, c := range capabilities.DefaultCapabilities {
		hashCap[c] = true
	}
	for _, a := range CapAdd {
		hashCap[a] = true
	}
	for _, d := range CapDrop {
		delete(hashCap, d)
	}
	defaultCap := make([]string, 0)
	for c, _ := range hashCap {
		defaultCap = append(defaultCap, c)
	}
	capabilities.DefaultCapabilities = defaultCap
}
func Run(cfg CicConfig) (err error) {
	run := &Runcic{
		Volume:          cfg.Volume,
		CopyEnv:         cfg.CopyParentEnv,
		Name:            cfg.Name,
		Command:         cfg.Cmd,
		CicVolume:       cfg.CicVolume,
		ImagePullPolicy: cfg.ImagePullPolicy,
	}
	for i := 0; i < len(cfg.Images); i++ {
		run.Images = append(run.Images, &common.Image{
			Image: cfg.Images[i],
		})
	}
	mergeCap(cfg.CapAdd, cfg.CapDrop)
	switch run.ImagePullPolicy {
	case imagePullPolicyAlways:
		for i := 0; i < len(run.Images); i++ {
			err = pullimage(run.Images[i].Image, cfg.Authfile)
			if err != nil {
				return
			}
		}
	case ImagePullPolicyfNotPresent:
		fallthrough
	default:
		for i := 0; i < len(run.Images); i++ {
			logrus.Infof("runcic imagedriver spec image %s", run.Images[i].Image)
			imagespec := containerimage.Driver().Spec(run.Images[i].Image)
			if imagespec == nil {
				logrus.Warnf("runcic imagedriver not found image %s", run.Images[i].Image)
				pullimage(run.Images[i].Image, cfg.Authfile)
			}
		}
	}
	for i := 0; i < len(run.Images); i++ {
		imgi := containerimage.Driver().Spec(run.Images[i].Image)
		if imgi == nil {
			logrus.Errorf("runcic imagedriver spec image is nil,your image=%s", run.Images[i].Image)
			return
		}
		run.Images[i] = imgi
		logrus.Infof("runcic imagedriver spec image %+v", imgi)
	}

	//todo 创建之前，需要检测是否已存在
	run.mergeEnv(cfg.Env)
	run.mergeCmd()
	if run.Name == "" {
		run.ContainerID = newID()
		run.Name = newName()
	}

	if err = run.Create(); err != nil {
		logrus.Errorf("create cic by images %+v fail,error: %+v", run.ImageArray(), err.Error())
		return
	}

	if err = run.Start(); err != nil {
		logrus.Errorf("start image %+v %+v fail,error: %+v", run.ImageArray(), run.Command, err.Error())
		return
	}
	return
}
