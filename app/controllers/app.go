package controllers

import "io/ioutil"
import "net/http"
import "encoding/json"
import "html/template"
import "github.com/robfig/revel"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	var url, post_title, post_content, post_url string
	var id string = c.Params.Get("id")

	if id != "" {
		url = "http://api.tumblr.com/v2/blog/tumblr.pierrebeaucamp.com/posts?api_key=SUPERTUMBLRPRIVATEKEY&limit=1&id=" + id
	} else {
		url = "http://api.tumblr.com/v2/blog/tumblr.pierrebeaucamp.com/posts?api_key=SUPERTUMBLRPRIVATEKEY&limit=1"
	}

	response, err := http.Get(url)
	if err != nil {
		post_title = "Error"
		post_content = err.Error()
	} else {
		defer response.Body.Close()
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			post_title = "Error"
			post_content = err.Error()
		} else {
			var message map[string]interface{}
			err := json.Unmarshal([]byte(content), &message)
			if err != nil {
				post_title = "Error"
				post_content = err.Error()
			} else {
				post_title = message["response"].(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["title"].(string)
				post_content = message["response"].(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["body"].(string)
				post_url = message["response"].(map[string]interface{})["posts"].([]interface{})[0].(map[string]interface{})["post_url"].(string)
			}
		}
	}
	post_content_html := template.HTML(post_content)
	return c.Render(post_title, post_content_html, post_url)
}
