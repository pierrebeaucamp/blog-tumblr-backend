package controllers

import "io/ioutil"
import "strings"
import "net/http"
import "encoding/json"
import "html/template"
import "github.com/robfig/revel"

const baseurl string = "http://api.tumblr.com/v2/blog/tumblr.pierrebeaucamp.com/posts?api_key=Grsk0wqWNw9IjnCTmKGfxodp1vsk1jneVUhAfK8F5SKyrdEpcy"

type Blog struct{ *revel.Controller }

func simplify_message(m map[string]interface{}, c string, i int) string {
	switch c {
	case "title":
		return m["response"].(map[string]interface{})["posts"].([]interface{})[i].(map[string]interface{})["title"].(string)
	case "content":
		return m["response"].(map[string]interface{})["posts"].([]interface{})[i].(map[string]interface{})["body"].(string)
	case "url":
		return strings.TrimPrefix(m["response"].(map[string]interface{})["posts"].([]interface{})[i].(map[string]interface{})["post_url"].(string), "http://tumblr.pierrebeaucamp.com/post")
	default:
		panic("Don't call this function with other cases then those aboth")
	}
}

func get_json(url string) map[string]interface{} {
	var message map[string]interface{}
	response, _ := http.Get(url)
	defer response.Body.Close()
	content, _ := ioutil.ReadAll(response.Body)

	err := json.Unmarshal([]byte(content), &message)
	if err != nil {
		panic("Couldn't unmarshal json")
	}

	return message
}

func get_post(url string) map[string]string {
	message := get_json(url)
	data := map[string]string{
		"title":   simplify_message(message, "title", 0),
		"content": simplify_message(message, "content", 0),
		"url":     simplify_message(message, "url", 0),
	}
	return data
}

func get_menu() map[string]string {
	menu := make(map[string]string)
	message := get_json(baseurl)
	for i := range message["response"].(map[string]interface{})["posts"].([]interface{}) {
		title := simplify_message(message, "title", i)
		url := simplify_message(message, "url", i)
		menu[title] = url
	}
	return menu
}

func (c Blog) Permalink(id string) revel.Result {
	var post, menu map[string]string
	post = get_post(baseurl + "&limit=1&id=" + id)
	menu = get_menu()
	title := post["title"]
	content := template.HTML(post["content"])
	return c.Render(title, content, menu)
}

func (c Blog) Index() revel.Result {
	var post map[string]string
	post = get_post(baseurl + "&limit=1")
	return c.Redirect(post["url"])
}

func (c Blog) Ajax(id string) revel.Result {
	var post map[string]string
	post = get_post(baseurl + "&limit=1&id=" + id)
	return c.RenderJson(post)
}
