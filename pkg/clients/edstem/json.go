package edstem

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (c *Client) JsonParse() ([]Announcement, error) {
	announcement := []Announcement{}
	res, err := c.GetJson(http.MethodGet, "https://edstem.org/api/user")
	if err != nil {
		return nil, err
	}

	courses := Courses{}

	err = json.Unmarshal(res, &courses)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(courses.Courses); i++ {
		res, err = c.GetJson(http.MethodGet, "https://edstem.org/api/courses/"+strconv.Itoa(courses.Courses[0].Course.Id)+"/lessons")
		if err != nil {
			return nil, err
		}

		lessons := Lessons{}
		err = json.Unmarshal(res, &lessons)
		if err != nil {
			return nil, err
		}

		for _, les := range lessons.Lessons {
			if les.Status == "attempted" {
				res, err = c.GetJson(http.MethodGet, "https://edstem.org/api/lessons/slides/"+strconv.Itoa(les.LastViewedSlideId)+"/questions")
				if err != nil {
					return nil, err
				}
				if string(res) == ErNotfound || string(res) == ErQuize {
					continue
				} else {
					announcement = append(announcement, Announcement{Title: les.Title, SubjectName: courses.Courses[i].Course.Name, Deadline: nil})
				}
			}
		}
	}
	return announcement, nil
}
