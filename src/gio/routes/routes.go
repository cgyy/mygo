package routes

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"
)

type Route struct {
	method  string
	regex   *regexp.Regexp
	params  map[int]string
	handler http.HandlerFunc
}

type Router struct {
	routes  []*Route
}

func NewRouter() *Router {
	this := Router{}
	return &this
}

// Adds a new Route to the Handler
func (this *Router) AddRoute(method string, pattern string, handler http.HandlerFunc) *Route {

	//split the url into sections
	parts := strings.Split(pattern, "/")

	//find params that start with ":"
	//replace with regular expressions
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’ 
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}

	//recreate the url pattern, with parameters replaced
	//by regular expressions. then compile the regex
	pattern = strings.Join(parts, "/")
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
		return nil
	}

	//now create the Route
	route := &Route{}
	route.method = method
	route.regex = regex
	route.handler = handler
	route.params = params

	//and finally append to the list of Routes
	this.routes = append(this.routes, route)

	return route
}

// Adds a new Route for GET requests
func (this *Router) Get(pattern string, handler http.HandlerFunc) *Route {
	return this.AddRoute(GET, pattern, handler)
}

// Adds a new Route for PUT requests
func (this *Router) Put(pattern string, handler http.HandlerFunc) *Route {
	return this.AddRoute(PUT, pattern, handler)
}

// Adds a new Route for DELETE requests
func (this *Router) Del(pattern string, handler http.HandlerFunc) *Route {
	return this.AddRoute(DELETE, pattern, handler)
}

// Adds a new Route for PATCH requests
func (this *Router) Patch(pattern string, handler http.HandlerFunc) *Route {
	return this.AddRoute(PATCH, pattern, handler)
}

// Adds a new Route for POST requests
func (this *Router) Post(pattern string, handler http.HandlerFunc) *Route {
	return this.AddRoute(POST, pattern, handler)
}

// Required by http.Handler interface. This method is invoked by the
// http server and will handle all page routing
func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestPath := r.URL.Path

	//find a matching Route
    var route *Route = nil

    for _, _route := range this.routes {

		//if the methods don't match, skip this handler
		//i.e if request.Method is 'PUT' Route.Method must be 'PUT'
		if r.Method != _route.method {
			continue
		}

		//check if Route pattern matches url
		if !_route.regex.MatchString(requestPath) {
			continue
		}

		//get submatches (params)
		matches := _route.regex.FindStringSubmatch(requestPath)

		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(requestPath) {
			continue
		}

		//add url parameters to the query param map
		values := r.URL.Query()
		for i, match := range matches[1:] {
			values.Add(_route.params[i], match)
		}

		//reassemble query params and add to RawQuery
		r.URL.RawQuery = url.Values(values).Encode()

        route = _route
		break
	}

    if route != nil {
		//Invoke the request handler
		route.handler(w, r)
    } else {
		http.NotFound(w, r)
	}
}
