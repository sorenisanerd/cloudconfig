package cloudconfig

import (
	yaml "gopkg.in/yaml.v3"
)

type CloudConfig struct {
	Password  *string   `yaml:",omitempty"`
	Packages  []string  `yaml:",omitempty"`
	SshPwAuth *bool     `yaml:"ssh_pwauth,omitempty"`
	ChPasswd  *ChPasswd `yaml:",omitempty"`
	Commands  []Command `yaml:"runcmd,flow,omitempty"`
	Mounts    []Mount   `yaml:",omitempty,flow"`
}

type Mount []string
type Command []string

type ChPasswd struct {
	Expire bool
}

type CloudConfigValue interface{}

func New() *CloudConfig {
	return &CloudConfig{}
}

func (cc *CloudConfig) AddPackage(pkgname string) {
	cc.Packages = append(cc.Packages, pkgname)
}

// AddMount adds a mount entry (like a line in /etc/fstab).
// First argument is `fs_spec`. It describes the block special device
// or remote filesystem to be mounted.
// Second argument is `fs_file`. It describes the mount point.
// An optional third argument is filesystem type.
// An optional fourth argument is mount options (noatime and such).
// An optional fifth argument tells dump(8) is mount options (noatime and such).
// An optional sixth argument is passno. It determines the order in which
// fsck is run at boot.
func (cc *CloudConfig) AddMount(device, mountPoint string, moreInfo ...string) {
	cc.Mounts = append(cc.Mounts, append([]string{device, mountPoint}, moreInfo...))
}

func (cc *CloudConfig) AddRunCmd(cmdArgs ...string) {
	cc.Commands = append(cc.Commands, cmdArgs)
}

func (cc *CloudConfig) SetPassword(password string) {
	cc.Password = &password
}

func (cc *CloudConfig) SetSshPwAuth(enable bool) {
	cc.SshPwAuth = &enable
}

func (cc *CloudConfig) SetChpasswdExpire(enable bool) {
	cc.ChPasswd = &ChPasswd{Expire: enable}
}

func (cc *CloudConfig) AddBashScript(script string) {
	cc.AddRunCmd("/bin/bash", "-c", script)
}

func (cc *CloudConfig) Generate() ([]byte, error) {
	out, _ := yaml.Marshal(cc)
	return append([]byte("#cloud-config\n"), out...), nil
}
