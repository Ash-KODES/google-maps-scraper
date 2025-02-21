package gmaps

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"unicode"
)

type Image struct {
	Title string `json:"title"`
	Image string `json:"image"`
}

type LinkSource struct {
	Link   string `json:"link"`
	Source string `json:"source"`
}

type Owner struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Address struct {
	Borough    string `json:"borough"`
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	State      string `json:"state"`
	Country    string `json:"country"`
}

type Option struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type About struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Options []Option `json:"options"`
}

type Review struct {
	Name           string
	ProfilePicture string
	Rating         int
	Description    string
	Images         []string
	When           string
}

type WorkingHours struct {
	Day       string `json:"day" bson:"day"`
	OpenHours string `json:"openHours" bson:"openHours"`
	Open      bool   `json:"open" bson:"open"`
}

type Entry struct {
	Link         string         `json:"link"`
	Cid          string         `json:"cID"`
	Title        string         `json:"businessName"`
	Categories   []string       `json:"categories"`
	Category     string         `json:"category"`
	Address      string         `json:"address"`
	WorkingHours []WorkingHours `json:"WorkingHours"`
	// PopularTImes is a map with keys the days of the week
	// and value is a map with key the hour and value the traffic in that time
	PopularTimes     map[string]map[int]int `json:"popularTimes"`
	WebSite          string                 `json:"webSite"`
	Phone            string                 `json:"phone"`
	PlusCode         string                 `json:"plusCode"`
	ReviewCount      int                    `json:"reviewCount"`
	ReviewRating     float64                `json:"reviewRating"`
	ReviewsPerRating map[int]int            `json:"reviewsPerRating"`
	Latitude         float64                `json:"latitude"`
	Longtitude       float64                `json:"longtitude"`
	Status           string                 `json:"status"`
	Description      string                 `json:"description"`
	ReviewsLink      string                 `json:"reviewsLink"`
	Thumbnail        string                 `json:"thumbnail"`
	Timezone         string                 `json:"timeZone"`
	PriceRange       string                 `json:"priceRange"`
	DataID           string                 `json:"dataID"`
	Images           []Image                `json:"images"`
	Reservations     []LinkSource           `json:"reservations"`
	OrderOnline      []LinkSource           `json:"orderOnline"`
	Services         LinkSource             `json:"services"`
	Owner            Owner                  `json:"owner"`
	CompleteAddress  Address                `json:"completeAddress"`
	About            []About                `json:"about"`
	UserReviews      []Review               `json:"userReviews"`
	Emails           []string               `json:"emails"`
}

func (e *Entry) IsWebsiteValidForEmail() bool {
	if e.WebSite == "" {
		return false
	}

	needles := []string{
		"facebook",
		"instragram",
		"twitter",
	}

	for i := range needles {
		if strings.Contains(e.WebSite, needles[i]) {
			return false
		}
	}

	return true
}

func (e *Entry) Validate() error {
	if e.Title == "" {
		return fmt.Errorf("title is empty")
	}

	if e.Category == "" {
		return fmt.Errorf("category is empty")
	}

	return nil
}

func (e *Entry) CsvHeaders() []string {
	return []string{
		"link",
		"businessName",
		"category",
		"address",
		"workingHours",
		"popularTimes",
		"webSite",
		"phone",
		"plusCode",
		"reviewCount",
		"reviewRating",
		"reviewsPerRating",
		"latitude",
		"longitude",
		"cID",
		"status",
		"descriptions",
		"reviewsLink",
		"thumbnail",
		"timeZone",
		"priceRange",
		"dataID",
		"images",
		"reservations",
		"orderOnline",
		"services",
		"owner",
		"completeAddress",
		"about",
		"userReviews",
		"emails",
	}
}

func (e *Entry) CsvRow() []string {
	return []string{
		e.Link,
		e.Title,
		e.Category,
		e.Address,
		stringify(e.WorkingHours),
		stringify(e.PopularTimes),
		e.WebSite,
		e.Phone,
		e.PlusCode,
		stringify(e.ReviewCount),
		stringify(e.ReviewRating),
		stringify(e.ReviewsPerRating),
		stringify(e.Latitude),
		stringify(e.Longtitude),
		e.Cid,
		e.Status,
		e.Description,
		e.ReviewsLink,
		e.Thumbnail,
		e.Timezone,
		e.PriceRange,
		e.DataID,
		stringify(e.Images),
		stringify(e.Reservations),
		stringify(e.OrderOnline),
		stringify(e.Services),
		stringify(e.Owner),
		stringify(e.CompleteAddress),
		stringify(e.About),
		stringify(e.UserReviews),
		stringSliceToString(e.Emails),
	}
}

//nolint:gomnd // it's ok, I need the indexes
func EntryFromJSON(raw []byte) (entry Entry, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered from panic: %v stack: %s", r, debug.Stack())

			return
		}
	}()

	var jd []any
	if err := json.Unmarshal(raw, &jd); err != nil {
		return entry, err
	}

	if len(jd) < 7 {
		return entry, fmt.Errorf("invalid json")
	}

	darray, ok := jd[6].([]any)
	if !ok {
		return entry, fmt.Errorf("invalid json")
	}

	entry.Link = getNthElementAndCast[string](darray, 27)
	entry.Title = getNthElementAndCast[string](darray, 11)

	categoriesI := getNthElementAndCast[[]any](darray, 13)

	entry.Categories = make([]string, len(categoriesI))
	for i := range categoriesI {
		entry.Categories[i], _ = categoriesI[i].(string)
	}

	if len(entry.Categories) > 0 {
		entry.Category = entry.Categories[0]
	}

	entry.Address = strings.TrimSpace(
		strings.TrimPrefix(getNthElementAndCast[string](darray, 18), entry.Title+","),
	)
	entry.WorkingHours = getHours(darray)
	entry.PopularTimes = getPopularTimes(darray)
	entry.WebSite = getNthElementAndCast[string](darray, 7, 0)
	entry.Phone = getNthElementAndCast[string](darray, 178, 0, 0)
	entry.PlusCode = getNthElementAndCast[string](darray, 183, 2, 2, 0)
	entry.ReviewCount = int(getNthElementAndCast[float64](darray, 4, 8))
	entry.ReviewRating = getNthElementAndCast[float64](darray, 4, 7)
	entry.Latitude = getNthElementAndCast[float64](darray, 9, 2)
	entry.Longtitude = getNthElementAndCast[float64](darray, 9, 3)
	entry.Cid = getNthElementAndCast[string](jd, 25, 3, 0, 13, 0, 0, 1)
	entry.Status = getNthElementAndCast[string](darray, 34, 4, 4)
	entry.Description = getNthElementAndCast[string](darray, 32, 1, 1)
	entry.ReviewsLink = getNthElementAndCast[string](darray, 4, 3, 0)
	entry.Thumbnail = getNthElementAndCast[string](darray, 72, 0, 1, 6, 0)
	entry.Timezone = getNthElementAndCast[string](darray, 30)
	entry.PriceRange = getNthElementAndCast[string](darray, 4, 2)
	entry.DataID = getNthElementAndCast[string](darray, 10)

	items := getLinkSource(getLinkSourceParams{
		arr:    getNthElementAndCast[[]any](darray, 171, 0),
		link:   []int{3, 0, 6, 0},
		source: []int{2},
	})

	entry.Images = make([]Image, len(items))

	for i := range items {
		entry.Images[i] = Image{
			Title: items[i].Source,
			Image: items[i].Link,
		}
	}

	entry.Reservations = getLinkSource(getLinkSourceParams{
		arr:    getNthElementAndCast[[]any](darray, 46),
		link:   []int{0},
		source: []int{1},
	})

	orderOnlineI := getNthElementAndCast[[]any](darray, 75, 0, 1, 2)

	if len(orderOnlineI) == 0 {
		orderOnlineI = getNthElementAndCast[[]any](darray, 75, 0, 0, 2)
	}

	entry.OrderOnline = getLinkSource(getLinkSourceParams{
		arr:    orderOnlineI,
		link:   []int{1, 2, 0},
		source: []int{0, 0},
	})

	entry.Services = LinkSource{
		Link:   getNthElementAndCast[string](darray, 38, 0),
		Source: getNthElementAndCast[string](darray, 38, 1),
	}

	entry.Owner = Owner{
		ID:   getNthElementAndCast[string](darray, 57, 2),
		Name: getNthElementAndCast[string](darray, 57, 1),
	}

	if entry.Owner.ID != "" {
		entry.Owner.Link = fmt.Sprintf("https://www.google.com/maps/contrib/%s", entry.Owner.ID)
	}

	entry.CompleteAddress = Address{
		Borough:    getNthElementAndCast[string](darray, 183, 1, 0),
		Street:     getNthElementAndCast[string](darray, 183, 1, 1),
		City:       getNthElementAndCast[string](darray, 183, 1, 3),
		PostalCode: getNthElementAndCast[string](darray, 183, 1, 4),
		State:      getNthElementAndCast[string](darray, 183, 1, 5),
		Country:    getNthElementAndCast[string](darray, 183, 1, 6),
	}

	aboutI := getNthElementAndCast[[]any](darray, 100, 1)

	for i := range aboutI {
		el := getNthElementAndCast[[]any](aboutI, i)
		about := About{
			ID:   getNthElementAndCast[string](el, 0),
			Name: getNthElementAndCast[string](el, 1),
		}

		optsI := getNthElementAndCast[[]any](el, 2)
		for j := range optsI {
			opt := Option{
				Enabled: getNthElementAndCast[int](optsI, j, 2, 1, 0, 0) == 1,
				Name:    getNthElementAndCast[string](optsI, j, 1),
			}

			if opt.Name != "" {
				about.Options = append(about.Options, opt)
			}
		}

		entry.About = append(entry.About, about)
	}

	entry.ReviewsPerRating = map[int]int{
		1: int(getNthElementAndCast[float64](darray, 52, 3, 0)),
		2: int(getNthElementAndCast[float64](darray, 52, 3, 1)),
		3: int(getNthElementAndCast[float64](darray, 52, 3, 2)),
		4: int(getNthElementAndCast[float64](darray, 52, 3, 3)),
		5: int(getNthElementAndCast[float64](darray, 52, 3, 4)),
	}

	reviewsI := getNthElementAndCast[[]any](darray, 52, 0)

	for i := range reviewsI {
		el := getNthElementAndCast[[]any](reviewsI, i)
		review := Review{
			Name:           getNthElementAndCast[string](el, 0, 1),
			ProfilePicture: getNthElementAndCast[string](el, 0, 2),
			When:           getNthElementAndCast[string](el, 1),
			Rating:         int(getNthElementAndCast[float64](el, 4)),
			Description:    getNthElementAndCast[string](el, 3),
		}

		if review.Name == "" {
			continue
		}

		optsI := getNthElementAndCast[[]any](el, 14)

		for j := range optsI {
			val := getNthElementAndCast[string](optsI, j, 6, 0)
			if val != "" {
				review.Images = append(review.Images, val)
			}
		}

		entry.UserReviews = append(entry.UserReviews, review)
	}

	return entry, nil
}

type getLinkSourceParams struct {
	arr    []any
	source []int
	link   []int
}

func getLinkSource(params getLinkSourceParams) []LinkSource {
	var result []LinkSource

	for i := range params.arr {
		item := getNthElementAndCast[[]any](params.arr, i)

		el := LinkSource{
			Source: getNthElementAndCast[string](item, params.source...),
			Link:   getNthElementAndCast[string](item, params.link...),
		}
		if el.Link != "" && el.Source != "" {
			result = append(result, el)
		}
	}

	return result
}

//nolint:gomnd // it's ok, I need the indexes
func getHours(darray []any) []WorkingHours {
	items := getNthElementAndCast[[]any](darray, 34, 1)
	var workingHours []WorkingHours

	for _, item := range items {
		day := getNthElementAndCast[string](item.([]any), 0)
		timesI := getNthElementAndCast[[]any](item.([]any), 1)
		times := make([]string, len(timesI))

		for i := range timesI {
			times[i], _ = timesI[i].(string)
		}
		timesInString := strings.Join(times, "") //""10 am–10 pm""
		// Check if there are any time slots, and set the Open field accordingly
		// startTime := ""
		// endTime := ""
		// Remove non-digit characters and split into start and end parts
		cleanedTime := strings.ReplaceAll(timesInString, " ", "")
		superCleanedTime := strings.ReplaceAll(cleanedTime, "\"", "")
		// matches := regexp.MustCompile(`(\d{1,2}[APMapm]+)-(\d{1,2}[APMapm]+)`).FindStringSubmatch(superCleanedTime)
		open := false
		for _, char := range superCleanedTime {
			if unicode.IsDigit(char) {
				open = true
			}
		}
		println(superCleanedTime)
		// Create a WorkingHour instance and append it to the result slice
		workingHour := WorkingHours{
			Day:       day,
			OpenHours: superCleanedTime,
			Open:      open,
		}

		workingHours = append(workingHours, workingHour)
	}

	return workingHours
}

func getPopularTimes(darray []any) map[string]map[int]int {
	items := getNthElementAndCast[[]any](darray, 84, 0) //nolint:gomnd // it's ok, I need the indexes
	popularTimes := make(map[string]map[int]int, len(items))

	dayOfWeek := map[int]string{
		1: "Monday",
		2: "Tuesday",
		3: "Wednesday",
		4: "Thursday",
		5: "Friday",
		6: "Saturday",
		7: "Sunday",
	}

	for ii := range items {
		item, ok := items[ii].([]any)
		if !ok {
			return nil
		}

		day := int(getNthElementAndCast[float64](item, 0))

		timesI := getNthElementAndCast[[]any](item, 1)

		times := make(map[int]int, len(timesI))

		for i := range timesI {
			t, ok := timesI[i].([]any)
			if !ok {
				return nil
			}

			v, ok := t[1].(float64)
			if !ok {
				return nil
			}

			h, ok := t[0].(float64)
			if !ok {
				return nil
			}

			times[int(h)] = int(v)
		}

		popularTimes[dayOfWeek[day]] = times
	}

	return popularTimes
}

func getNthElementAndCast[T any](arr []any, indexes ...int) T {
	var (
		defaultVal T
		idx        int
	)

	if len(indexes) == 0 {
		return defaultVal
	}

	for len(indexes) > 1 {
		idx, indexes = indexes[0], indexes[1:]

		if idx >= len(arr) {
			return defaultVal
		}

		next := arr[idx]

		if next == nil {
			return defaultVal
		}

		var ok bool

		arr, ok = next.([]any)
		if !ok {
			return defaultVal
		}
	}

	if len(indexes) == 0 || len(arr) == 0 {
		return defaultVal
	}

	ans, ok := arr[indexes[0]].(T)
	if !ok {
		return defaultVal
	}

	return ans
}

func stringSliceToString(s []string) string {
	return strings.Join(s, ", ")
}

func stringify(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%f", val)
	case nil:
		return ""
	default:
		d, _ := json.Marshal(v)
		return string(d)
	}
}
