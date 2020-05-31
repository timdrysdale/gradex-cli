package pagedata

import (
	"errors"
)

type Link struct {
	First    string
	Last     string
	Sequence []string
	IsLinked bool
}

// A non-nil error means there is a broken sequence on at least one page
// the Linkmap has the details ....
func GetLinkMap(pageDataMap map[int]PageData) (map[int]Link, error) {

	allLinksOk := true
	linkMap := make(map[int]Link)

	for pageNumber, pdPage := range pageDataMap {

		pds := []PageDetail{}

		for _, pd := range pdPage.Previous {

			pds = append(pds, pd)
		}

		pds = append(pds, pdPage.Current)

		pageLinkMap := make(map[string]string)

		for _, pd := range pds {
			pageLinkMap[pd.UUID] = pd.Follows
		}

		last := pdPage.Current.UUID
		s := []string{last}

		thisPageLinksOk := true

		// avoid infinite cycles by limiting iterations
		for i := 0; i < len(pageLinkMap); i++ {
			if next, ok := pageLinkMap[last]; ok {
				if next == "" {
					break
				}
				s = append(s, next)
				last = next
			} else {
				break
			}
		}

		if len(s) < len(pageLinkMap) {
			thisPageLinksOk = false
			allLinksOk = false
		}

		//https://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go/19239850
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}

		linkMap[pageNumber] = Link{
			First:    s[0],
			Last:     s[len(s)-1],
			Sequence: s,
			IsLinked: thisPageLinksOk,
		}

	}

	if allLinksOk {
		return linkMap, nil
	} else {
		return linkMap, errors.New("Links are not in sequence")
	}
}
