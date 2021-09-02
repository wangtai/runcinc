package cic

import (
	"github.com/sirupsen/logrus"
	"os"
	"runcic/cic/fs"
	"runcic/utils"
)

func (r *Runcic) SetEnv(copyParentEnv bool) {
	if copyParentEnv {
		r.Envs = append(r.Envs, os.Environ()...)
	}
	os.Clearenv()
	for k, v := range utils.ParseEnvs(r.Envs) {
		os.Setenv(k, v)
	}
}
func (r *Runcic) Start() (err error) {
	if err = r.cicvolume(); err != nil {
		return
	}
	if err = r.rootfspath(); err != nil {
		return
	}
	if err = r.mountoverlay(); err != nil {
		return
	}

	if r.ParentRootfs, err = realChroot(r.Roorfs()); err != nil {
		logrus.Errorf("chroot failed %s", err.Error())
		return
	}

	if err = fs.Mount(); err != nil {
		logrus.Errorf("fs mount failed %s", err.Error())
	}
	if err = fs.Link(); err != nil {
		logrus.Errorf("fs link failed %s", err.Error())
	}
	r.SetEnv(false)
	logrus.Infof("%+v %+v", r.Command, r.Envs)

	err = Execv(r.Command[0], r.Command[1:], r.Envs)
	if err != nil {
		logrus.Errorf("exec exited %v", err.Error())
	}
	return err
}
