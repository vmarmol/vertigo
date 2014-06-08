package pulet

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
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
	return &Pulet{}
}

func (self *Pulet) Import(spec *ImportSpec) (*ImageSpec, error) {
	importUrl := fmt.Sprintf("http://%v:%v/export/%v", spec.SourceHost, spec.SourcePort, spec.SourceId)
	imgSpec := &ImageSpec{
		Repo: randomUniqString(),
		Tag:  randomUniqString(),
	}
	alias := fmt.Sprintf("%v:%v", imgSpec.Repo, imgSpec.Tag)
	fmt.Printf("docker import %v %v", importUrl, alias)
	cmd := exec.Command("docker", "import", importUrl, alias)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return imgSpec, nil
}

// runArgs: arguments sends after docker run. Something like -p -v, etc.
// cmdList: command runs in the container
func (self *Pulet) RunImage(img *ImageSpec, runArgs []string, cmdList []string) error {
	alias := fmt.Sprintf("%v:%v", img.Repo, img.Tag)
	args := make([]string, 0, len(runArgs)+len(cmdList)+2)
	args = append(args, "run")
	for _, arg := range runArgs {
		args = append(args, arg)
	}
	args = append(args, alias)
	for _, arg := range cmdList {
		args = append(args, arg)
	}

	fmt.Printf("docker %+v", args)
	cmd := exec.Command("docker", args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
