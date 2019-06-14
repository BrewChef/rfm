package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wilriker/librfm"
	"github.com/wilriker/rfm"
)

// FileinfoOptions holds the specific parameters for fileinfo requests
type FileinfoOptions struct {
	*BaseOptions
	path          string
	humanReadable bool
}

// Check checks all parameters for valid values
func (f *FileinfoOptions) Check() {
	f.BaseOptions.Check()

	if f.path == "" {
		log.Fatal("-path is mandatory")
	}
	f.path = rfm.CleanRemotePath(f.path)
}

// InitFileinfoOptions inializes a FileinfoOptions instance from command-line parameters
func InitFileinfoOptions(arguments []string) *FileinfoOptions {
	f := FileinfoOptions{BaseOptions: &BaseOptions{}}

	fs := f.GetFlagSet()
	fs.BoolVar(&f.humanReadable, "h", false, "Display size in human readable units")
	fs.Parse(arguments)

	if fs.NArg() > 0 {
		f.path = fs.Arg(0)
	}

	f.Check()

	f.Connect()

	return &f
}

// DoFileinfo is a convenience function to run a download from command-line parameters
func DoFileinfo(arguments []string) error {
	fo := InitFileinfoOptions(arguments)
	return NewFileinfo(fo).Fileinfo(fo.path)
}

// Fileinfo provides a singl method to fetch infos on a file
type Fileinfo interface {
	Fileinfo(path string) error
}

// fileinfo implement the Fileinfo interface
type fileinfo struct {
	o *FileinfoOptions
}

// NewFileinfo creates a new instance of the Fileinfo interface
func NewFileinfo(fo *FileinfoOptions) Fileinfo {
	return &fileinfo{
		o: fo,
	}
}

// Fileinfo fetches information on a file and prints it to sdtout
func (f *fileinfo) Fileinfo(path string) error {
	fi, err := f.o.Rfm.Fileinfo(path)
	if err != nil {
		return err
	}
	f.print(path, fi)
	return nil
}

func (f *fileinfo) print(path string, fi *librfm.Fileinfo) {
	fmt.Println(path)
	fmt.Printf("Size:               %s\n", f.getSize(fi.Size))
	fmt.Printf("Last modified:      %s\n", fi.LastModified().Format(librfm.TimeFormat))
	if fi.Height > 0 {
		fmt.Printf("Height:             %.2fmm\n", fi.Height)
	}
	if fi.FirstLayerHeight > 0 {
		fmt.Printf("First layer height: %.2fmm\n", fi.FirstLayerHeight)
	}
	if fi.LayerHeight > 0 {
		fmt.Printf("Layer height:       %.2fmm\n", fi.LayerHeight)
	}
	if fi.PrintTime > 0 {
		fmt.Printf("Print time:         %s\n", f.getPrintTime(fi.PrintTime))
	}
	if len(fi.Filament) > 0 {
		fmt.Printf("Filament usage:     %s\n", f.getFilamentUsage(fi.Filament))
	}
	if fi.GeneratedBy != "" {
		fmt.Printf("Generated by:       %s\n", fi.GeneratedBy)
	}
}

func (f *fileinfo) getPrintTime(seconds uint64) string {
	if f.o.humanReadable {
		d := time.Duration(time.Duration(seconds) * time.Second)
		return d.String()
	}
	return fmt.Sprintf("%ds", seconds)
}

func (f *fileinfo) getFilamentUsage(filaments []float64) string {
	if len(filaments) == 1 {
		return fmt.Sprintf("%.1fmm", filaments[0])
	}
	var b strings.Builder
	for _, fil := range filaments {
		b.WriteString(fmt.Sprintf("%.1fmm", fil))
		b.WriteString(", ")
	}
	return strings.TrimSuffix(b.String(), ", ")
}

func (f *fileinfo) getSize(size uint64) string {
	if f.o.humanReadable {
		return rfm.HumanReadableSize(size)
	}
	return fmt.Sprintf("%d", size)
}
