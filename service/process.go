package service

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"github.com/teacatx/iknore/store"
	"gopkg.in/gographics/imagick.v3/imagick"
)

func Compress(c context.Context, path string) (int, int, error) {
	imagick.Initialize()
	//
	mw := imagick.NewMagickWand()
	if err := mw.ReadImage(path); err != nil {
		return 0, 0, err
	}
	if err := mw.SetCompressionQuality(60); err != nil {
		return 0, 0, err
	}
	if err := mw.SetImageCompressionQuality(60); err != nil {
		return 0, 0, err
	}
	if err := mw.SetImageFormat("avif"); err != nil {
		return 0, 0, err
	}
	if err := mw.StripImage(); err != nil {
		return 0, 0, err
	}
	width, height := int(mw.GetImageWidth()), int(mw.GetImageHeight())
	//
	if err := mw.WriteImage(path); err != nil {
		return 0, 0, err
	}
	//
	mw.Destroy()
	imagick.Terminate()

	//
	return width, height, nil
}

func CoverTypeToImagickGravity(v store.CoverMode) imagick.GravityType {
	switch v {
	case store.CoverModeNorthWest:
		return imagick.GRAVITY_NORTH_WEST
	case store.CoverModeNorth:
		return imagick.GRAVITY_NORTH
	case store.CoverModeNorthEast:
		return imagick.GRAVITY_NORTH_EAST
	case store.CoverModeWest:
		return imagick.GRAVITY_WEST
	case store.CoverModeCenter:
		return imagick.GRAVITY_CENTER
	case store.CoverModeEast:
		return imagick.GRAVITY_EAST
	case store.CoverModeSouthWest:
		return imagick.GRAVITY_SOUTH_WEST
	case store.CoverModeSouth:
		return imagick.GRAVITY_SOUTH
	case store.CoverModeSouthEast:
		return imagick.GRAVITY_SOUTH_EAST
	default:
		return imagick.GRAVITY_UNDEFINED
	}
}

func (s *ImageService) SuffixToFormat(v string) string {
	switch v {
	case ".png":
		return "png"
	case ".bmp":
		return "bmp"
	case ".jpg":
		return "jpeg" // jpeg?
	case ".webp":
		return "webp"
	case ".avif":
		return "avif"
	case ".gif":
		return "gif"
	default:
		return ""
	}
}

func (s *ImageService) IsSupportedFormat(v string) error {
	if !lo.Contains(s.Formats, v) {
		return errors.New("unsupported target format")
	}
	return nil
}

func (s *ImageService) GetImageSize(mw *imagick.MagickWand) (int, int) {
	return int(mw.GetImageWidth()), int(mw.GetImageHeight())
}

func (s *ImageService) GetImageRatio(w, h int) float32 {
	return float32(w) / float32(h)
}

func (s *ImageService) CalculateCoverXY(origW, origH int, args *store.ImageArguments) (x int, y int) {
	switch args.CoverMode {
	case store.CoverModeNorthWest:
		return
	case store.CoverModeNorth:
		x = (origW / 2) - (args.Width / 2)
		return
	case store.CoverModeNorthEast:
		x = origW - args.Width
		return
	case store.CoverModeWest:
		y = (origH / 2) - (args.Height / 2)
		return
	case store.CoverModeCenter:
		x = (origW / 2) - (args.Width / 2)
		y = (origH / 2) - (args.Height / 2)
		return
	case store.CoverModeContain:
		x = (origW / 2) - (args.Width / 2)
		y = (origH / 2) - (args.Height / 2)
		return
	case store.CoverModeEast:
		x = origW - args.Width
		y = (origH / 2) - (args.Height / 2)
		return
	case store.CoverModeSouthWest:
		y = origH - args.Height
		return
	case store.CoverModeSouth:
		x = (origW / 2) - (args.Width / 2)
		y = origH - args.Height
		return
	case store.CoverModeSouthEast:
		x = origW - args.Width
		y = origH - args.Height
		return
	}
	return
}

func (s *ImageService) MakeVariant(c context.Context, orig string, args *store.ImageArguments) error {
	imagick.Initialize()
	//
	mw := imagick.NewMagickWand()
	//
	if err := s.IsSupportedFormat(args.Format); err != nil {
		return err
	}
	//
	if err := mw.ReadImage(orig); err != nil {
		return err
	}
	//
	origW, origH := s.GetImageSize(mw)
	//
	ratio := s.GetImageRatio(origW, origH)

	//
	//coverMode := CoverTypeToImagickGravity(cover) // default covermode?
	//if coverMode != imagick.GRAVITY_UNDEFINED {
	//	if err := mw.SetGravity(coverMode); err != nil {
	//		return err
	//	}
	//}

	if args.BackgroundColor != "" {
		pmw := imagick.NewPixelWand()
		pmw.SetColor(args.BackgroundColor)

		if err := mw.SetImageBackgroundColor(pmw); err != nil {
			return err
		}
	}

	// 如果只有指定一邊的寬或是高，那就依照等比例自動算出缺少的那一邊。
	// 然後縮放圖片（放大或縮小）。
	if (args.Width != 0 && args.Height == 0) || (args.Width == 0 && args.Height != 0) {
		newW := args.Width
		newH := args.Height

		if newW == 0 {
			if args.IgnoreAspectRatio {
				newW = origW
			} else {
				newW = int(float32(newH) * ratio)
			}

		}
		if newH == 0 {
			if args.IgnoreAspectRatio {
				newH = origH
			} else {
				newH = int(float32(newW) / ratio)
			}
		}

		if err := mw.ResizeImage(uint(newW), uint(newH), imagick.FILTER_LANCZOS); err != nil {
			return err
		}
	}

	if args.Width != 0 && args.Height != 0 {

		if args.CoverMode == store.CoverModeContain {
			imgW := 0
			imgH := 0

			if origW > origH {
				imgW = args.Width
				imgH = int((float32(imgW) / ratio))

				if args.Height < imgH {
					imgH = args.Height
					imgW = int((float32(imgH) * ratio))
				}

			} else {
				imgH = args.Height
				imgW = int((float32(imgH) * ratio))

				if args.Width < imgW {
					imgW = args.Width
					imgH = int((float32(imgW) / ratio))
				}
			}

			if err := mw.ResizeImage(uint(imgW), uint(imgH), imagick.FILTER_LANCZOS); err != nil {
				return err
			}

			x, y := s.CalculateCoverXY(imgW, imgH, args)

			if err := mw.ExtentImage(uint(args.Width), uint(args.Height), x, y); err != nil {
				return err
			}
		} else if args.CoverMode != store.CoverModeNone {

			imgW := 0
			imgH := 0

			if origW > origH {
				imgH = args.Height
				imgW = int((float32(imgH) * ratio))

				if args.Width > imgW {
					imgW = args.Width
					imgH = int((float32(imgW) / ratio))
				}

			} else {
				imgW = args.Width
				imgH = int((float32(imgW) / ratio))

				if args.Height > imgH {
					imgH = args.Height
					imgW = int((float32(imgH) * ratio))
				}
			}

			if err := mw.ResizeImage(uint(imgW), uint(imgH), imagick.FILTER_LANCZOS); err != nil {
				return err
			}

			x, y := s.CalculateCoverXY(imgW, imgH, args)

			if err := mw.ExtentImage(uint(args.Width), uint(args.Height), x, y); err != nil {
				return err
			}
		}

	}

	if err := mw.SetCompressionQuality(60); err != nil {
		return err
	}
	if err := mw.SetImageCompressionQuality(60); err != nil {
		return err
	}

	if err := mw.SetImageFormat(args.Format); err != nil {
		return err
	}
	if err := mw.StripImage(); err != nil {
		return err
	}

	// mw.IdentifyImage()
	// case ".bmp", ".heic", ".heif", ".jfif", ".jpeg", ".jpg", ".png", ".avif", ".gif", ".webp":

	//
	if err := mw.WriteImage(orig); err != nil {
		return err
	}
	//
	mw.Destroy()
	imagick.Terminate()

	//
	return nil
}
