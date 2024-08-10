package sitemap

import (
	"errors"
	"fmt"
	"github.com/Ja7ad/meilisitemap/config"
	"github.com/Ja7ad/meilisitemap/utils"
	"strings"
	"time"
	"unicode"
)

func uniqueToSlug(unique string) string {
	unique = strings.ToLower(unique)

	var builder strings.Builder

	for _, char := range unique {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			builder.WriteRune(char)
		} else if unicode.IsSpace(char) {
			builder.WriteRune('-')
		}
	}

	slug := builder.String()
	slug = strings.Trim(slug, "-")

	return slug
}

func getFileLoc(key string, doc map[string]interface{}) (string, error) {
	var val interface{}
	var ok bool

	// Split the key to check for any prefix, actual key, and suffix
	if strings.Contains(key, "|") {
		sp := strings.Split(key, "|")
		actualKey := sp[0]

		val := utils.PickByNestedKey(doc, actualKey)
		if val == nil {
			return "", fmt.Errorf("failed to get value loc for key: %s", actualKey)
		}

		// Construct the file location with the prefix and suffix
		switch len(sp) {
		case 2:
			if strings.HasSuffix(sp[1], "=") {
				return fmt.Sprintf("%s%v", sp[1], val), nil
			}
			return fmt.Sprintf("%s/%v", sp[1], val), nil
		case 3:
			if strings.HasSuffix(sp[1], "=") {
				return fmt.Sprintf("%s%v%s", sp[1], val, sp[2]), nil
			}
			return fmt.Sprintf("%s/%v%s", sp[1], val, sp[2]), nil
		default:
			return "", errors.New("invalid key format")
		}
	}

	val = utils.PickByNestedKey(doc, key)
	if val == nil {
		return "", fmt.Errorf("failed to get value loc for key: %s", key)
	}

	link, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("value is not a string for key: %s", key)
	}

	return link, nil
}

func imageFieldMapToSitemapImage(imgCfg *config.ImageConfig, doc map[string]interface{}) (*Image, error) {
	img := new(Image)

	loc, locErr := getFileLoc(imgCfg.Loc, doc)
	if locErr != nil {
		return nil, locErr
	}

	img.Loc = loc

	var (
		err error
	)

	if imgCfg.Title != "" {
		img.Title, err = getStringValueFromDoc(imgCfg.Title, doc)
		if err != nil {
			return nil, err
		}
	}

	if imgCfg.Caption != "" {
		img.Caption, err = getStringValueFromDoc(imgCfg.Caption, doc)
		if err != nil {
			return nil, err
		}
	}

	if imgCfg.GeoLocation != "" {
		img.GeoLocation, err = getStringValueFromDoc(imgCfg.GeoLocation, doc)
		if err != nil {
			return nil, err
		}
	}

	if imgCfg.License != "" {
		img.License, err = getStringValueFromDoc(imgCfg.License, doc)
		if err != nil {
			return nil, err
		}
	}

	return img, nil
}

func videoFieldMapToSitemapVideo(vidCfg *config.VideoConfig, doc map[string]interface{}) (*Video, error) {
	vid := new(Video)

	if vidCfg.ThumbnailLoc != "" {
		loc, locErr := getFileLoc(vidCfg.ThumbnailLoc, doc)
		if locErr != nil {
			return nil, locErr
		}
		vid.ThumbnailLoc = loc
	}

	if vidCfg.ContentLoc != "" {
		loc, locErr := getFileLoc(vidCfg.ContentLoc, doc)
		if locErr != nil {
			return nil, locErr
		}
		vid.ContentLoc = loc
	}

	var (
		err error
	)

	if vidCfg.PlayerAutoPlay != "" {
		val := utils.PickByNestedKey(doc, vidCfg.PlayerAutoPlay)
		if val == nil {
			return nil, fmt.Errorf("failed to get value loc for key: %s", vidCfg.PlayerAutoPlay)
		}

		vid.PlayerLocAutoplay, err = getBoolValueToZeroOrOne(vidCfg.PlayerAutoPlay, val)
		if err != nil {
			return nil, err
		}

		vid.PlayerLocAutoplay = "ap=" + vid.PlayerLocAutoplay
	}

	if vidCfg.Title != "" {
		vid.Title, err = getStringValueFromDoc(vidCfg.Title, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.Duration != "" {
		vid.Duration, err = getStringValueFromDoc(vidCfg.Duration, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.Description != "" {
		vid.Description, err = getStringValueFromDoc(vidCfg.Description, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.ExpirationDate != "" {
		val := utils.PickByNestedKey(doc, vidCfg.ExpirationDate)
		if val == nil {
			return nil, fmt.Errorf("failed to get value for key: %s", vidCfg.ExpirationDate)
		}
		vid.ExpirationDate, err = getDateTimeFromDoc(val)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.FamilyFriendly != "" {
		val := utils.PickByNestedKey(doc, vidCfg.FamilyFriendly)
		if val == nil {
			return nil, fmt.Errorf("failed to get value loc for key: %s", vidCfg.FamilyFriendly)
		}

		vid.FamilyFriendly, err = getBoolValueToYesOrNo(vidCfg.FamilyFriendly, val)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.Live != "" {
		val := utils.PickByNestedKey(doc, vidCfg.Live)
		if val == nil {
			return nil, fmt.Errorf("failed to get value loc for key: %s", vidCfg.Live)
		}

		vid.Live, err = getBoolValueToYesOrNo(vidCfg.Live, val)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.PlayerLoc != "" {
		vid.PlayerLoc, err = getStringValueFromDoc(vidCfg.PlayerLoc, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.PublicationDate != "" {
		val := utils.PickByNestedKey(doc, vidCfg.PublicationDate)
		if val == nil {
			return nil, fmt.Errorf("failed to get value for key: %s", vidCfg.PublicationDate)
		}
		vid.PublicationDate, err = getDateTimeFromDoc(val)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.Rating != "" {
		vid.Rating, err = getStringValueFromDoc(vidCfg.Rating, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.RequiresSubscription != "" {
		vid.RequiresSubscription, err = getStringValueFromDoc(vidCfg.RequiresSubscription, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.Restriction != "" {
		vid.Restriction, err = getStringValueFromDoc(vidCfg.Restriction, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.RestrictionRelationship != "" {
		vid.RestrictionRelationship, err = getStringValueFromDoc(vidCfg.RestrictionRelationship, doc)
		if err != nil {
			return nil, err
		}
	}

	if vidCfg.ViewCount != "" {
		vid.ViewCount, err = getStringValueFromDoc(vidCfg.ViewCount, doc)
		if err != nil {
			return nil, err
		}
	}

	return vid, nil
}

func newsFieldMapToSitemapNews(newsCfg *config.NewsConfig, doc map[string]interface{}) (*News, error) {
	news := new(News)

	var err error

	if newsCfg.Title != "" {
		news.Title, err = getStringValueFromDoc(newsCfg.Title, doc)
		if err != nil {
			return nil, err
		}
	}

	if newsCfg.Description != "" {
		news.Description, err = getStringValueFromDoc(newsCfg.Description, doc)
		if err != nil {
			return nil, err
		}
	}

	if newsCfg.PubDate != "" {
		val := utils.PickByNestedKey(doc, newsCfg.PubDate)
		if val == nil {
			return nil, fmt.Errorf("failed to get value for key: %s", newsCfg.PubDate)
		}
		news.PubDate, err = getDateTimeFromDoc(val)
		if err != nil {
			return nil, err
		}
	}

	if newsCfg.Keywords != "" {
		val := utils.PickByNestedKey(doc, newsCfg.Keywords)
		if val == nil {
			return nil, fmt.Errorf("failed to get value for key: %s", newsCfg.Keywords)
		}
		news.Keywords, err = getArrayFromDoc(val)
		if err != nil {
			return nil, err
		}
	}

	if newsCfg.Publication != nil {
		news.Publication = new(NewsPublication)
		if newsCfg.Publication.Name != "" {
			news.Publication.Name, err = getStringValueFromDoc(newsCfg.Publication.Name, doc)
			if err != nil {
				return nil, err
			}
		}

		if newsCfg.Publication.Language != "" {
			news.Publication.Language, err = getStringValueFromDoc(newsCfg.Publication.Language, doc)
			if err != nil {
				return nil, err
			}
		}
	}

	return news, nil
}

func getStringValueFromDoc(key string, doc map[string]interface{}) (string, error) {
	if strings.Contains(key, "|") {
		sp := strings.Split(key, "|")
		vals := make([]string, 0)

		for _, s := range sp {
			val := utils.PickByNestedKey(doc, s)
			if val == nil {
				return "", fmt.Errorf("failed to get value loc for key: %s", s)
			}

			vals = append(vals, fmt.Sprintf("%v", val))
		}

		return strings.Join(vals, " "), nil
	}

	val := utils.PickByNestedKey(doc, key)
	if val == nil {
		return "", fmt.Errorf("failed to get value loc for key: %s", key)
	}

	return fmt.Sprintf("%v", val), nil
}

func getBoolValueToYesOrNo(key string, val interface{}) (string, error) {
	b, ok := val.(bool)
	if !ok {
		return "", fmt.Errorf("value for key %s is not a boolean", key)
	}

	if b {
		return "yes", nil
	}

	return "no", nil
}

func getBoolValueToZeroOrOne(key string, val interface{}) (string, error) {
	b, ok := val.(bool)
	if !ok {
		return "", fmt.Errorf("value for key %s is not a boolean", key)
	}

	if b {
		return "1", nil
	}

	return "0", nil
}

func getDateTimeFromDoc(val interface{}) (string, error) {
	switch val.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, val.(string))
		if err != nil {
			return "", err
		}
		return t.Format(_datetimeLayout), nil
	case time.Time:
		return val.(time.Time).Format(_datetimeLayout), nil
	case int64:
		return time.Unix(val.(int64), 0).Format(_datetimeLayout), nil
	default:
		return "", fmt.Errorf("unsupported datetime format")
	}
}

func getArrayFromDoc(val interface{}) (string, error) {
	vArry, ok := val.([]string)
	if !ok {
		return "", fmt.Errorf("unsupported array format")
	}

	return strings.Join(vArry, ", "), nil
}
