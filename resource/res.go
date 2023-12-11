package main

import (
	"log"
	"os"

	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
)

func main() {
	// First create an empty resource set
	rs := winres.ResourceSet{}
	/*
		// Make an icon group from a png file
		f, err := os.Open("icon.png")
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatalln(err)
		}
		f.Close()
		icon, _ := winres.NewIconFromResizedImage(img, nil)

		// Add the icon to the resource set, as "APPICON"
		rs.SetIcon(winres.Name("APPICON"), icon)
	*/
	// Make a VersionInfo structure
	vi := version.Info{
		FileVersion:    [4]uint16{0, 0, 28, 0},
		ProductVersion: [4]uint16{0, 0, 1, 0},
	}
	vi.Set(0, version.ProductName, "WinGUI example")
	vi.Set(0, version.ProductVersion, "v0.0.1")
	vi.Set(0, version.FileVersion, "v0.0.28.0")

	// Add the VersionInfo to the resource set
	rs.SetVersionInfo(vi)
	/*
		// Add a manifest
		rs.SetManifest(winres.AppManifest{
			ExecutionLevel:      RequireAdministrator,
			DPIAwareness:        DPIPerMonitorV2,
			UseCommonControlsV6: true,
		})
	*/
	// Create an object file for amd64
	out, err := os.Create("rsrc_windows_amd64.syso")
	defer out.Close()
	if err != nil {
		log.Fatalln(err)
	}
	err = rs.WriteObject(out, winres.ArchAMD64)
	if err != nil {
		log.Fatalln(err)
	}
}
