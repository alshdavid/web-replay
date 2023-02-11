package har

type Model struct {
	Log Log `json:"log"`
}

type Log struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	ResourceType string   `json:"_resourceType"`
	Request      Request  `json:"request"`
	Response     Response `json:"response"`
}

type Request struct {
	Method   string          `json:"method"`
	Url      string          `json:"url"`
	Headers  []Header        `json:"headers"`
	PostData RequestPostData `json:"postData"`
}

type RequestPostData struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
	Encoding string `json:"encoding"`
}

type Response struct {
	Status  int             `json:"status"`
	Headers []Header        `json:"headers"`
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
	Encoding string `json:"encoding"`
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
