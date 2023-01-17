package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"github.com/teacatx/iknore/config"
	"github.com/teacatx/iknore/store"
	"gopkg.in/gographics/imagick.v3/imagick"
)

type ImageService struct {
	Config       *config.Config
	SizeAliases  map[string]map[string][2]int // { "user_avatar": { "small": [400, 400] } }
	Placeholders map[string][]byte            // { "user.png": [ 89, 64... ] }
	Formats      []string                     //
}

func NewImageService(conf *config.Config) *ImageService {
	svc := &ImageService{
		SizeAliases:  conf.InitSizeAliases(),
		Placeholders: conf.InitPlaceholders(),
		Config:       conf,
	}
	//
	imagick.Initialize()
	mw := imagick.NewMagickWand()

	// Load the available Imagick formats, so we can validate the requested formats.
	svc.Formats = lo.Map(mw.QueryFormats("*"), func(v string, _ int) string {
		return strings.ToLower(v)
	})
	return svc
}

func VariantPlaceholder(args *store.ImageArguments) string {
	path := fmt.Sprintf("%s-", args.Type)
	if args.Width != 0 {
		path += fmt.Sprintf("w_%d-", args.Width)
	}
	if args.Height != 0 {
		path += fmt.Sprintf("h_%d-", args.Height)
	}
	if args.CoverMode != store.CoverModeNone {
		path += fmt.Sprintf("c_%s-", args.CoverMode)
	}
	if args.BackgroundColor != "" {
		path += fmt.Sprintf("bc_%s-", args.BackgroundColor)
	}
	return strings.TrimSuffix(path, "-") + args.Extension
}

func (s *ImageService) AliasToSize(args *store.ImageArguments) (int, int) {
	typeMap, ok := s.SizeAliases[args.Type]
	if !ok {
		return 0, 0
	}
	v, ok := typeMap[args.Size]
	if !ok {
		return 0, 0
	}
	return v[0], v[1]
}

func (s *ImageService) CheckVaildSize(args *store.ImageArguments) bool {
	v, ok := s.Config.Types[args.Type]
	if !ok {
		return false
	}
	//
	if args.Width != 0 || args.Height != 0 {
		var str string
		if args.Width != 0 && args.Height == 0 {
			str = fmt.Sprintf("%dx", args.Width) // 500x
		} else if args.Width == 0 && args.Height != 0 {
			str = fmt.Sprintf("x%d", args.Height) // x500
		} else if args.Width != 0 && args.Height != 0 {
			str = fmt.Sprintf("%dx%d", args.Width, args.Height) // 500x500
		}
		return lo.Contains(lo.Map(v.Sizes, func(v string, _ int) string {
			return strings.Split(v, " ")[0] // 500x (small)
		}), str)
	}
	//
	if args.Size != "" {
		return lo.SomeBy(v.Sizes, func(v string) bool {
			return strings.Contains(v, fmt.Sprintf("(%s)", args.Size)) // (small)
		})
	}

	return true
}

func (s *ImageService) CheckValidCover(args *store.ImageArguments) bool {
	if args.CoverMode == store.CoverModeNone {
		return true
	}
	v, ok := s.Config.Types[args.Type]
	if !ok {
		return false
	}
	if lo.Contains(v.Covers, "*") {
		return true
	}
	return lo.Contains(v.Covers, string(args.CoverMode))
}

func (s *ImageService) CheckValidBackgroundColor(args *store.ImageArguments) bool {
	if args.BackgroundColor == "" {
		return true
	}
	v, ok := s.Config.Types[args.Type]
	if !ok {
		return false
	}
	if lo.Contains(v.BackgroundColors, "*") {
		return true
	}
	return lo.Contains(v.BackgroundColors, args.BackgroundColor)
}

func (s *ImageService) CheckValidFormat(args *store.ImageArguments) bool {
	v, ok := s.Config.Types[args.Type]
	if !ok {
		return false
	}
	if args.Format == v.Original.Format {
		return true
	}
	if lo.Contains(v.Formats, "*") {
		return true
	}
	return lo.Contains(v.Formats, args.Format)
}

func (s *ImageService) GetPlaceholder(args *store.ImageArguments) ([]byte, error) {
	placeholderName := VariantPlaceholder(args)

	// 找尋看看這個條件的 Placeholder 是不是有產出過了。
	v, ok := s.Placeholders[placeholderName]
	if ok {
		return v, nil
	}

	// 如果沒有產出過的話，就先從 Config 裡面抓出這個 Type 的 Placeholder。
	var defaultPlaceholder string
	for k, v := range s.Config.Types {
		if k != args.Type {
			continue
		}
		defaultPlaceholder = k + filepath.Ext(v.Placeholder)
	}
	// 如果找不到該種類 Placeholder，那就用全域的 Placeholder。
	if defaultPlaceholder == "" {
		defaultPlaceholder = "*" + filepath.Ext(s.Config.Placeholder)
	}

	// 把這個 Placeholder 存成暫存檔案，然後再做 Process，讓這個圖片塑造的像請求參數一樣。
	// 之後再存進記憶體的 Placeholder 供下次需要的時候直接拿取。
	tmpFile, err := os.CreateTemp("", "*"+filepath.Ext(defaultPlaceholder))
	if err != nil {
		return nil, err
	}
	if _, err := tmpFile.Write(s.Placeholders[defaultPlaceholder]); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	if err := s.MakeVariant(context.TODO(), tmpFile.Name(), args); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}
	s.Placeholders[placeholderName] = b

	return b, nil
}
