package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {

	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.*/
	w.Header().Add("Access-Control-Allow-Origin", "*")

	/*	- Get the `url` query string parameter value from the request.
		If not supplied, respond with an http.StatusBadRequest error.*/
	query := r.URL.Query()
	if len(query) == 0 {
		http.Error(w, "Missing Query Parameters", http.StatusBadRequest)
		return
	}

	/*	- Call fetchHTML() to fetch the requested URL. See comments in that
		function for more details.*/
	fetchedHTML, err := fetchHTML(query.Get("url"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	/*	- Call extractSummary() to extract the page summary meta-data,
		as directed in the assignment. See comments in that function
		for more details*/
	summary, err := extractSummary(query.Get("url"), fetchedHTML)
	if err != nil {
		log.Fatal(err)
	}

	/*	- Close the response HTML stream so that you don't leak resources. */
	fetchedHTML.Close()

	/*	- Finally, respond with a JSON-encoded version of the PageSummary
			struct. That way the client can easily parse the JSON back into
			an object. Remember to tell the client that the response content
			type is JSON.

		Helpful Links:
		https://golang.org/pkg/net/http/#Request.FormValue
		https://golang.org/pkg/net/http/#Error
		https://golang.org/pkg/encoding/json/#NewEncoder */
	encoder := json.NewEncoder(w)
	err = encoder.Encode(summary)
	if err != nil {
		log.Fatal(err)
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.*/
	getResponse, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if getResponse.StatusCode >= 400 {
		err = errors.New("Status code above 400")
		return nil, err
	} else if !strings.HasPrefix(getResponse.Header.Get("Content-type"), "text/html") { //thank you logan
		err = errors.New("Content type incorrect")
		return nil, err
	}
	return getResponse.Body, nil

	/*	To test your implementation of this function, run the TestFetchHTML
		test in summary_test.go. You can do that directly in Visual Studio Code,
		or at the command line by running:
			go test -run TestFetchHTML

		Helpful Links:
		https://golang.org/pkg/net/http/#Get
	*/
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.*/
	tokenizer := html.NewTokenizer(htmlStream)

	//Creating Struct early
	summary := &PageSummary{}
	var imagesSlice []*PreviewImage

	for {
		tokenType := tokenizer.Next()
		//tricky error handling provided by the reading
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				//end of the file, break out of the loop
				break
			}
		}
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				break // found the end of the head tag so all meta attributes must be done
			}
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken { //Found a token that is an actual tag
			token := tokenizer.Token()

			//Check that it is a meta tag
			if token.Data == "meta" {
				//===========================================og:image tag handling===========================================
				URLImplicit, Success := gimmeMetaPropertyContent(token, "og:image")
				if Success {
					/* Found an og:image tag - Cases
					1. This is the very first image tag found, start of new image
					2. This is the second image tag found, still start of a new image
					*/

					imagesSlice = append(imagesSlice, &PreviewImage{})

					//initial og:image tag could not contain the actual URL
					//this would cause the returned value to be an empty string
					//even though the above method found the right tag.
					if URLImplicit != "" {
						//testing if it is a relative url
						if strings.Contains(URLImplicit, "://") {
							imagesSlice[len(imagesSlice)-1].URL = URLImplicit
						} else {
							imagesSlice[len(imagesSlice)-1].URL = urlRebaser(pageURL, URLImplicit)
						}
					}
				}
				//We have already detected a new og:image (or maybe we havent)
				//now we need to check if it has a secure URL also given
				//I'm going to assume that secure URL should take precedence over the standard URL
				URLExplicit, Success := gimmeMetaPropertyContent(token, "og:image:url")
				if Success {
					//testing if it is a relative url
					if strings.Contains(URLExplicit, "://") {
						imagesSlice[len(imagesSlice)-1].URL = URLExplicit
					} else {
						imagesSlice[len(imagesSlice)-1].URL = urlRebaser(pageURL, URLExplicit)
					}
				}
				SecureURL, Success := gimmeMetaPropertyContent(token, "og:image:secure_url")
				if Success {
					imagesSlice[len(imagesSlice)-1].SecureURL = SecureURL
				}
				imgType, Success := gimmeMetaPropertyContent(token, "og:image:type")
				if Success {
					imagesSlice[len(imagesSlice)-1].Type = imgType
				}
				Width, Success := gimmeMetaPropertyContent(token, "og:image:width")
				if Success && len(Width) != 0 {
					widthInt, err := strconv.Atoi(Width)
					if err != nil {
						log.Fatal(err)
					}
					imagesSlice[len(imagesSlice)-1].Width = widthInt
				}
				Height, Success := gimmeMetaPropertyContent(token, "og:image:height")
				if Success && len(Height) != 0 {
					heightInt, err := strconv.Atoi(Height)
					if err != nil {
						log.Fatal(err)
					}
					imagesSlice[len(imagesSlice)-1].Height = heightInt
				}
				Alt, Success := gimmeMetaPropertyContent(token, "og:image:alt")
				if Success {
					imagesSlice[len(imagesSlice)-1].Alt = Alt
				}
				//=============================================================================================================

				//Rest of standard meta tags
				//try and fetch the attribute
				//if successful in finding it, assign it to the struct
				Type, Success := gimmeMetaPropertyContent(token, "og:type")
				if Success {
					summary.Type = Type
				}
				URL, Success := gimmeMetaPropertyContent(token, "og:url")
				if Success {
					summary.URL = URL
				}
				Title, Success := gimmeMetaPropertyContent(token, "og:title")
				if Success {
					summary.Title = Title
				}
				SiteName, Success := gimmeMetaPropertyContent(token, "og:site_name")
				if Success {
					summary.SiteName = SiteName
				}
				OGDescription, Success := gimmeMetaPropertyContent(token, "og:description")
				if Success {
					summary.Description = OGDescription
				} else {
					Description, Success := gimmeMetaPropertyContent(token, "description")
					if Success && summary.Description == "" {
						summary.Description = Description
					}
				}
				Author, Success := gimmeMetaPropertyContent(token, "author")
				if Success {
					summary.Author = Author
				}
				Keywords, Success := gimmeMetaPropertyContent(token, "keywords")
				if Success {
					KeywordsSlice := strings.Split(Keywords, ",")
					var TrimmedKeywordsSlice []string
					for _, word := range KeywordsSlice {
						TrimmedKeywordsSlice = append(TrimmedKeywordsSlice, strings.TrimSpace(word))
					}
					summary.Keywords = TrimmedKeywordsSlice
				}
			}

			//Looking for the icon image
			if token.Data == "link" {
				for _, attribute := range token.Attr {
					if attribute.Key == "rel" && attribute.Val == "icon" {
						iconImgPreview := getIcon(token, pageURL)
						summary.Icon = iconImgPreview
					}
				}

			}
			//looking for the title if it isnt present in opengraph form
			if token.Data == "title" {
				tokenType = tokenizer.Next()
				//only grab this if it wasnt already provided by OG
				if tokenType == html.TextToken && summary.Title == "" {
					summary.Title = tokenizer.Token().Data
				}
			}
		}
	}
	if len(imagesSlice) != 0 {
		summary.Images = imagesSlice
	}
	return summary, nil
	/*To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/
}

func gimmeMetaPropertyContent(token html.Token, propertyName string) (string, bool) {
	foundAttribute := false
	foundContent := ""
	for _, attribute := range token.Attr { //for each loop to find what we need
		//Make sure that the attribute is in the provided tag
		if (attribute.Key == "property" || attribute.Key == "name") && attribute.Val == propertyName {
			foundAttribute = true
		}
		//Grabs the content of the provided tag
		if attribute.Key == "content" {
			foundContent = attribute.Val
		}
	}
	return foundContent, foundAttribute
}

//rebases the relative url given
func urlRebaser(pageURL string, relativeURL string) string {
	stringURL := strings.Split(pageURL, "/")
	stringURL[len(stringURL)-1] = strings.Trim(relativeURL, "/")
	newURL := strings.Join(stringURL, "/")
	return newURL
}

//tests if the given token is an image. Can test if it's part of an image or the beginning of an OG image
func testIfImg(token html.Token, flag string) bool {
	for _, attribute := range token.Attr {
		if attribute.Key == "property" && attribute.Val == "og:image" && flag == "beginning" {
			return true
		}
		if attribute.Key == "property" && strings.Contains(attribute.Val, "og:image") && flag == "still" {
			return true
		}
	}

	return false
}

func getIcon(token html.Token, pageURL string) *PreviewImage {
	iconImgPreview := new(PreviewImage)
	for _, attribute := range token.Attr {
		if attribute.Key == "href" {
			if strings.Contains(attribute.Val, "://") {
				iconImgPreview.URL = attribute.Val
			} else {
				iconImgPreview.URL = urlRebaser(pageURL, attribute.Val)
			}
		}
		if attribute.Key == "type" {
			iconImgPreview.Type = attribute.Val
		}
		if attribute.Key == "sizes" { //gotta  chop this up cause its a messy string
			if attribute.Val != "any" && len(attribute.Val) != 0 {
				heightXwidth := strings.Split(attribute.Val, "x")
				heightInt, err := strconv.Atoi(heightXwidth[0])
				if err != nil {
					log.Fatal(err)
				}
				widthInt, err := strconv.Atoi(heightXwidth[1])
				if err != nil {
					log.Fatal(err)
				}
				iconImgPreview.Height = heightInt
				iconImgPreview.Width = widthInt
			}
		}
	}
	return iconImgPreview
}
