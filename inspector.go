package inspector

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Total       int           `json:"total"`
	TotalPage   int           `json:"total_page"`
	CurrentPage int           `json:"current_page"`
	PerPage     int           `json:"per_page"`
	HasNext     bool          `json:"has_next"`
	HasPrev     bool          `json:"has_prev"`
	NextPageUrl string        `json:"next_page_url"`
	PrevPageUrl string        `json:"prev_page_url"`
	Data        []RequestStat `json:"data"`
}

type RequestStat struct {
	RequestedAt   time.Time   `json:"requested_at"`
	RequestUrl    string      `json:"request_url"`
	HttpMethod    string      `json:"http_method"`
	HttpStatus    int         `json:"http_status"`
	ContentType   string      `json:"content_type"`
	GetParams     interface{} `json:"get_params"`
	PostParams    interface{} `json:"post_params"`
	PostMultipart interface{} `json:"post_multipart"`
	ClientIP      string      `json:"client_ip"`
	Cookies       interface{} `json:"cookies"`
	Headers       interface{} `json:"headers"`
	Body          string      `json:"body"`
}

type AllRequests struct {
	Requets []RequestStat `json:"requests"`
}

var allRequests = AllRequests{}
var pagination = Pagination{}

func GetPaginator() Pagination {
	return pagination
}

func LoadHtml(router *gin.Engine) (err error) {
	router.Delims("{{", "}}")
	tmpl, err := template.New(HtmlName).Delims("{{", "}}").Funcs(template.FuncMap{
		"formatDate": formatDate,
	}).Parse(Html)
	fmt.Printf("%v", tmpl)
	router.SetHTMLTemplate(tmpl)
	return
}

func formatDate(t time.Time) string {
	return t.Format(time.RFC822)
}

func Frontend(c *gin.Context) {
	// c.JSON(200, inspector.GetPaginator())
	c.HTML(http.StatusOK, HtmlName, map[string]interface{}{
		"title":      "Gin Inspector",
		"pagination": GetPaginator(),
	})
}

func JsonFrontend(c *gin.Context) {
	c.JSON(200, GetPaginator())
}

const HtmlName = "inspector.html"

func InspectorStats() gin.HandlerFunc {
	return func(c *gin.Context) {

		urlPath := c.Request.URL.Path

		if strings.Contains(urlPath, "/_inspector") {

			page, _ := strconv.ParseFloat(c.DefaultQuery("page", "1"), 64)
			perPage, _ := strconv.ParseFloat(c.DefaultQuery("per_page", "20"), 64)
			total := float64(len(allRequests.Requets))
			totalPage := math.Ceil(total / perPage)
			offset := (page - 1) * perPage

			if offset < 0 {
				offset = 0
			}

			pagination.HasPrev = false
			pagination.HasNext = false
			pagination.CurrentPage = int(page)
			pagination.PerPage = int(perPage)
			pagination.TotalPage = int(totalPage)
			pagination.Total = int(total)
			pagination.Data = paginate(allRequests.Requets, int(offset), int(perPage))

			if pagination.CurrentPage > 1 {
				pagination.HasPrev = true
				pagination.PrevPageUrl = urlPath + "?page=" + strconv.Itoa(pagination.CurrentPage-1) + "&per_page=" + strconv.Itoa(pagination.PerPage)
			}

			if pagination.CurrentPage < pagination.TotalPage {
				pagination.HasNext = true
				pagination.NextPageUrl = urlPath + "?page=" + strconv.Itoa(pagination.CurrentPage+1) + "&per_page=" + strconv.Itoa(pagination.PerPage)
			}

		} else {

			start := time.Now()

			c.Request.ParseForm()
			c.Request.ParseMultipartForm(10000)

			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Println(err)
			}
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			request := RequestStat{
				RequestedAt:   start,
				RequestUrl:    urlPath,
				HttpMethod:    c.Request.Method,
				HttpStatus:    c.Writer.Status(),
				ContentType:   c.ContentType(),
				Headers:       c.Request.Header,
				Cookies:       c.Request.Cookies(),
				GetParams:     c.Request.URL.Query(),
				PostParams:    c.Request.PostForm,
				PostMultipart: c.Request.MultipartForm,
				ClientIP:      c.ClientIP(),
				Body:          string(body),
			}

			allRequests.Requets = append([]RequestStat{request}, allRequests.Requets...)

		}

		c.Next()

	}
}

func paginate(s []RequestStat, offset, length int) []RequestStat {
	end := offset + length
	if end < len(s) {
		return s[offset:end]
	}

	return s[offset:]
}

const Html = `
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="Mark Otto, Jacob Thornton, and Bootstrap contributors">
    <meta name="generator" content="Jekyll v3.8.5">
    <title>{{.title}}</title>

    <!-- Bootstrap core CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/css/bootstrap.min.css" integrity="sha384-GJzZqFGwb1QTTN6wy59ffF1BuGJpLSa9DkKMp0DgiMDm4iYMj70gZWKYbI706tWS" crossorigin="anonymous">
    <link href="https://fonts.googleapis.com/css?family=Poppins:300,500,700" rel="stylesheet">


    <style>
      body {
        font-family: 'Poppins', sans-serif;
        font-size: 15px;
        background: #f0f5fb!important;
        -webkit-font-smoothing: antialiased;
        text-rendering: optimizeLegibility;
        -moz-osx-font-smoothing: grayscale;
        font-weight: 300;
      }
      .navbar {
        background-color: #37474f;
      }
      .table thead th {
        text-transform: uppercase;
        font-size: 13px;
        font-weight: 500;
        color: #607d8b;
    
      }
      .table td, .table th {
        border-top: none;
      }
      .table tr th {
        font-weight: 500;
      }
      .table thead th, .table tr {
        border-bottom: 1px solid #eee;
      }
      .table tr:last-child {
        border-bottom: none;
      }
      .badge {
        padding: 6px 12px;
        font-size: 13px;
        font-weight: 500;
      }



      .btn-detail {
        font-size: 13px;
        font-weight: 500;
        padding: 5px 20px;
        background: #00bcd4;
        border:none;
      }
      .btn-detail:hover,.btn-detail:active {
        background-color: #4dd0e1!important;
     
      }
      .badge-200 {
          color: #fff!important;
          background-color: #4ac35e!important;
      }

      .badge-500 {
          color: #fff!important;
          background-color: #f44336;
      }

      .bd-placeholder-img {
        font-size: 1.125rem;
        text-anchor: middle;
      }
      .shadow-sm {
        box-shadow: 0 .125rem 1.25rem rgba(0,0,0,.02) !important;
      }
      .page-link {
        border-color: #01bcd4;
        color: #01bcd4;
      }
      .page-link:hover {
        background-color: #00bcd4;
        color: #fff;
      }

      .nav-tabs {
        border: none;
      }
      .nav-tabs .nav-link {
        border-radius: 40px;
        background-color: #eee;
        color:#999; 
        border:none;
        font-size: 15px;
        font-weight: 500;
        padding: 10px 30px;
        margin-right:10px;
      }
      .nav-tabs .nav-link.active {
        background-color: #00bcd4;
        color: #fff;
      }

      @media (min-width: 768px) {
        .bd-placeholder-img-lg {
          font-size: 3.5rem;
        }
      }

      @media (min-width: 576px) {
        .modal-dialog {
          max-width: none !important;
        }
      }
    </style>
    <!-- Custom styles for this template -->
  </head>
  <body class="bg-light">
      <nav class="navbar navbar-dark">
          <a class="navbar-brand mx-auto" href="#">
            <img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" width="30"  class="d-inline-block align-top" alt="">
            <div style="display: inline-block;margin-top: 10px;margin-left: 10px;">Inspector</div>
          </a>
        </nav>


<main role="main" class="container" >

  

  <div class="my-3 p-4 bg-white rounded shadow-sm" > 

    <table class="table m-0">
            <thead>
              <tr>
                <th scope="col">Method</th>
                <th scope="col">Url</th>
                <th scope="col">Status</th>
                <th scope="col">Date</th>
                <th scope="col">Inspect</th>
              </tr>
            </thead>
            <tbody>
            {{ range $i,$value :=.pagination.Data }}
              <tr>
                <th scope="row">{{ $value.HttpMethod }}</th>
                <td><span>{{$value.RequestUrl}}</span></td>
                <td>
                    <span class="badge badge-secondary badge-{{$value.HttpStatus}}">{{ $value.HttpStatus }}</span>
                  
                </td>
                <td>{{$value.RequestedAt | formatDate }}</td>
                <td>
                  
                    <button type="button" class="btn btn-primary btn-detail" data-toggle="modal" data-target="#modal-{{$i}}">
                        Detail
                    </button>
                
                    <div class="modal fade" id="modal-{{$i}}" tabindex="-1" role="dialog" aria-labelledby="exampleModalLongTitle" aria-hidden="true">
                      <div class="container">
                        <div class="modal-dialog" role="document">
                          <div class="modal-content">
                            <div class="modal-header pl-4 pr-4">
                              <h5 class="modal-title" id="exampleModalLongTitle">{{.RequestUrl}}  <span class="badge badge-secondary badge-{{.HttpStatus}}">{{.HttpStatus}}</span> </h5>
                              <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                              </button>
                            </div>
                            <div class="modal-body p-4">

                                <ul class="nav nav-tabs mb-4">
                                    <li class="nav-item">
                                      <a class="nav-link active" data-toggle="tab" href="#modal-{{$i}}-tab1">Parameters</a>
                                    </li>
                                    <li class="nav-item">
                                      <a class="nav-link" data-toggle="tab" href="#modal-{{$i}}-tab2">Headers</a>
                                    </li>
                                    
                                  </ul>
                                
                                  <!-- Tab panes -->
                                  <div class="tab-content">

                                    <div id="modal-{{$i}}-tab1" class="tab-pane active">
                                      {{ if .GetParams }}
                                      <h3>Query Parameters</h3>
                                      <table class="table table-hover">
                                          <thead>
                                            <tr>
                                              <th scope="col">Key</th>
                                              <th scope="col">Value</th>
                                        
                                            </tr>
                                          </thead>
                                          <tbody>
                                            
                                            {{ range $key, $value :=  .GetParams}}
                                            <tr>
                                              <th scope="row">{{$key}}</th>
                                              <td>{{$value}}</td>
                                          
                                            </tr>
                                            {{end}}
                                            
                                          </tbody>
                                        </table>
                                        {{end}}

                                        {{ if .PostParams }}
                                        <h3>Post Parameters</h3>
                                          <table class="table table-hover">
                                            <thead>
                                              <tr>
                                                <th scope="col">Key</th>
                                                <th scope="col">Value</th>
                                          
                                              </tr>
                                            </thead>
                                            <tbody>
                                              {{ range $key, $value :=  .PostParams}}
                                              <tr>
                                                <th scope="row">{{$key}}</th>
                                                <td>{{$value}}</td>
                                            
                                              </tr>
                                              {{end}}
                                              
                                            </tbody>
                                          </table>
                                        {{end}}
                                        {{ if .Body }}
                                        <h3>Post Parameters</h3>
                                          <table class="table table-hover">
                                            <thead>
                                              <tr>
                                                <th scope="col">Key</th>
                                                <th scope="col">Value</th>
                                              </tr>
                                            </thead>
                                            <tbody>
                                              <tr>
                                                <th scope="row">Body</th>
                                                <td>{{.Body}}</td>
                                            
                                              </tr>
                                              
                                            </tbody>
                                          </table>
                                        {{end}}

                                        {{ if .PostMultipart }}

                                        

                                            {{if .PostMultipart.File}}
                                            <h3>Post Multipart Files</h3>
                                              <table class="table table-hover">
                                                <thead>
                                                  <tr>
                                                    <th scope="col">Key</th>
                                                    <th scope="col">Value</th>
                                              
                                                  </tr>
                                                </thead>
                                                <tbody>
                                                  {{ range $key, $value :=  .PostMultipart.File}}
                                                  <tr>
                                                    <th scope="row">{{$key}}</th>
                                                    <td>{{$value}}</td>
                                              
                                                  </tr>
                                                  {{end}}
                                                  
                                                </tbody>
                                              </table>
                                            {{end}}

                                            
                                        {{ end }}

                                        
                                        
                                       

                                    </div>
                                    <div id="modal-{{$i}}-tab2" class="tab-pane fade">
                                      
                                      <h3>Headers</h3>
                                      <table class="table table-hover">
                                          <thead>
                                            <tr>
                                              <th scope="col">Key</th>
                                              <th scope="col">Value</th>
                                        
                                            </tr>
                                          </thead>
                                          <tbody>
                                            {{ range $key, $value :=  .Headers}}
                                            <tr>
                                              <th scope="row">{{$key}}</th>
                                              <td>{{$value}}</td>
                                        
                                            </tr>
                                            {{end}}
                                            
                                          </tbody>
                                        </table>

                                    </div>
                                    
                                  </div>

                                             
                            </div>
                      
                          </div>
                        </div>
                      </div>  
                    </div>
                
                </td>
                

              </tr>

              {{ end }}  
             
            </tbody>
          </table>


          <nav aria-label="...">
              <ul class="pagination mt-3 mb-0">
                {{ if .pagination.HasPrev }}
                <li class="page-item">
                  <a class="page-link" href="{{.pagination.PrevPageUrl}}" tabindex="-1">Previous</a>
                </li>
                {{ end }}
                {{ if .pagination.HasNext }}
                <li class="page-item">
                  <a class="page-link" href="{{.pagination.NextPageUrl}}">Next</a>
                </li>
                {{ end }}
              </ul>
            </nav>
  </div>

  </div>
</main>
<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.6/umd/popper.min.js" integrity="sha384-wHAiFfRlMFy6i5SRaxvfOCifBUQy1xHdJ/yoi7FRNXMRBu5WHdZYu1hA6ZOblgut" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.2.1/js/bootstrap.min.js" integrity="sha384-B0UglyR+jN6CkvvICOB2joaf5I4l3gm9GU6Hc1og6Ls7i6U/mkkaduKaBhlAXv9k" crossorigin="anonymous"></script>
<script>
</script>
</html>
`
