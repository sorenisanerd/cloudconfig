package cloudconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCloudConfig(t *testing.T) {
	cc := New()
	cc.SetPassword("passw0rd")
	cc.SetSshPwAuth(true)
	cc.AddPackage("nfs-client")
	cc.AddPackage("nfs-common")
	cc.AddMount("1.2.3.4:/foobar", "/shared", "nfs")
	cc.AddMount("/dev/sda1", "/data", "ext4", "noatime,noexec")
	cc.AddRunCmd("mount", "-a")
	cc.AddRunCmd("echo", `"'`)
	cc.AddBashScript(`#!/bin/bash
echo "this"
echo 'that'
exit 2
`)
	cc.SetChpasswdExpire(false)
	expected := `#cloud-config
password: passw0rd
packages:
    - nfs-client
    - nfs-common
ssh_pwauth: true
chpasswd:
    expire: false
commands: [[mount, -a], [echo, '"'''], [/bin/bash, -c, "#!/bin/bash\necho \"this\"\necho 'that'\nexit 2\n"]]
mounts: [['1.2.3.4:/foobar', /shared, nfs], [/dev/sda1, /data, ext4, 'noatime,noexec']]
`
	check(cc, expected, t)
}

func check(cc *CloudConfig, expected string, t *testing.T) {
	got, _ := cc.Generate()
	s := string(got)
	if diff := cmp.Diff(expected, s); diff != "" {
		t.Errorf("cc.Generated() mismatch (-want +got):\n%s", diff)
	}
}

func TestCloudConfigEmpty(t *testing.T) {
	cc := New()
	check(cc, "#cloud-config\n{}\n", t)
}
