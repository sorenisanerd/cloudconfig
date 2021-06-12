package cloudconfig

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
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
packages:
    - nfs-client
    - nfs-common
ssh_pwauth: true
password: passw0rd
chpasswd:
    expire: false
runcmd: 
  - [mount, -a]
  - [echo, '"''']
  - [/bin/bash, -c, "#!/bin/bash\necho \"this\"\necho 'that'\nexit 2\n"]
mounts: 
  - ['1.2.3.4:/foobar', /shared, nfs]
  - [/dev/sda1, /data, ext4, 'noatime,noexec']
`
	check(cc, expected, t)
}

func check(cc *CloudConfig, expected string, t *testing.T) {
	gotMap := make(map[string]interface{})
	expectedMap := make(map[string]interface{})

	got, _ := cc.Generate()

	err := yaml.Unmarshal(got, gotMap)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(expected), expectedMap)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(expectedMap, gotMap); diff != "" {
		t.Errorf("cc.Generated() mismatch (-want +got):\n%s", diff)
	}
}

func TestCloudConfigEmpty(t *testing.T) {
	cc := New()
	check(cc, "#cloud-config\n{}\n", t)
}
