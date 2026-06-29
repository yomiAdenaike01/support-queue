package images

import (
	"fmt"
	"strings"
)

type Image struct {
	Name        string
	Tag         string
	Shortname   string
	Environment map[string]string
}

func (i Image) GetPath() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

type ImagesList []Image

func (i ImagesList) String() string {
	builder := strings.Builder{}
	for _, image := range i {
		builder.WriteString(fmt.Sprintf("%s,", image.GetPath()))
	}
	return builder.String()
}

func getAppName(appName string) string {
	imgBasename := "adenaikeyomi/support-queue"
	return fmt.Sprintf("%s-%s", imgBasename, appName)
}

func GetImagesList(version string) ImagesList {
	return ImagesList{
		{
			Name:      "pgvector/pgvector",
			Tag:       "pg17",
			Shortname: "db",
		},
		{
			Name:      "redis",
			Tag:       "latest",
			Shortname: "redis",
		},
		{
			Name:      getAppName("ml-worker"),
			Tag:       version,
			Shortname: "ml-worker",
		},
		{
			Name:      getAppName("frontend"),
			Tag:       version,
			Shortname: "frontend",
		},
		{
			Name:      getAppName("srv"),
			Tag:       version,
			Shortname: "srv",
		},
	}

}
