package pulet

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os/exec"
	"os/user"
	"time"
)

type ImportSpec struct {
	SourceHost string `json:"src_host"`
	SourcePort int    `json:"src_port"`
	SourceId   string `json:"src_id"`
}

type ImageSpec struct {
	Repo string `json:"repo"`
	Tag  string `json:"tag"`
}

func randomUniqString() string {
	var d [8]byte
	io.ReadFull(rand.Reader, d[:])
	str := hex.EncodeToString(d[:])
	return fmt.Sprintf("%x-%v", time.Now().Unix(), str)
}

type Pulet struct {
}

func NewPulet() *Pulet {
	u, err := user.Current()
	if err != nil {
		panic(err.Error())
	}
	if u.Username != "root" {
		panic("must be root!")
	}
	return &Pulet{}
}

func (self *Pulet) Import(spec *ImportSpec) (*ImageSpec, error) {
	importUrl := fmt.Sprintf("http://%v:%v/export/%v", spec.SourceHost, spec.SourcePort, spec.SourceId)
	imgSpec := &ImageSpec{
		Repo: randomUniqString(),
		Tag:  randomUniqString(),
	}
	alias := fmt.Sprintf("%v:%v", imgSpec.Repo, imgSpec.Tag)
	log.Printf("docker import %v %v", importUrl, alias)
	cmd := exec.Command("docker", "import", importUrl, alias)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return imgSpec, nil
}

// runArgs: arguments sends after docker run. Something like -p -v, etc.
// cmdList: command runs in the container
func (self *Pulet) RunImage(img *ImageSpec, runArgs []string, cmdList []string) (string, error) {
	alias := fmt.Sprintf("%v:%v", img.Repo, img.Tag)
	args := make([]string, 0, len(runArgs)+len(cmdList)+3)
	args = append(args, "run")
	args = append(args, "-d")
	args = append(args, runArgs...)
	args = append(args, alias)
	args = append(args, cmdList...)

	log.Printf("docker %+v", args)
	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	log.Printf("output of docker: %v\n", string(output))
	return string(output), nil
}
